package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ActivityManager はアクティビティの管理を行う
type ActivityManager struct {
	// 現在実行中の全アクティビティ(全エンティティごと)
	// 1エンティティで最大1アクティビティ
	currentActivities map[ecs.Entity]*Activity
	logger            *logger.Logger
}

// NewActivityManager は新しいActivityManagerを作成する
func NewActivityManager(logger *logger.Logger) *ActivityManager {
	return &ActivityManager{
		currentActivities: make(map[ecs.Entity]*Activity),
		logger:            logger,
	}
}

// Execute は指定されたアクション（アクティビティ）を実行する
// 即座実行アクション（移動、攻撃等）も継続アクション（休息等）も統一的に処理
func (am *ActivityManager) Execute(activityType ActivityType, params ActionParams, world w.World) (*ActionResult, error) {
	am.logger.Debug("アクション実行開始",
		"type", activityType.String(),
		"actor", params.Actor)

	// アクティビティを作成
	activity, err := am.createActivity(activityType, params, world)
	if err != nil {
		return &ActionResult{
			Success:      false,
			ActivityType: activityType,
			Message:      err.Error(),
		}, err
	}

	// アクティビティを開始
	if err := am.StartActivity(activity, world); err != nil {
		return &ActionResult{
			Success:      false,
			ActivityType: activityType,
			Message:      err.Error(),
		}, err
	}

	// 即座実行アクション（1ターン）の場合は即座に処理
	if activity.TurnsTotal == 1 {
		// ターン処理実行
		am.ProcessTurn(world)

		// ターン管理システムに移動コストを通知
		am.consumeMoveCost(world, activityType, params.Actor)

		// 結果を確認
		currentActivity := am.GetCurrentActivity(params.Actor)
		if currentActivity == nil || currentActivity.IsCompleted() {
			return &ActionResult{
				Success:      true,
				ActivityType: activityType,
				Message:      "アクション完了",
			}, nil
		} else if currentActivity.IsCanceled() {
			return &ActionResult{
				Success:      false,
				ActivityType: activityType,
				Message:      currentActivity.CancelReason,
			}, fmt.Errorf("アクション失敗: %s", currentActivity.CancelReason)
		}
	}

	// 継続アクションの場合は開始成功を返す
	return &ActionResult{
		Success:      true,
		ActivityType: activityType,
		Message:      "アクション開始",
	}, nil
}

// StartActivity は新しいアクティビティを開始する
func (am *ActivityManager) StartActivity(activity *Activity, world w.World) error {
	if activity == nil {
		return ErrActivityNil
	}

	// 既存のアクティビティがある場合は中断
	if currentActivity := am.GetCurrentActivity(activity.Actor); currentActivity != nil {
		if err := am.InterruptActivity(activity.Actor, "新しいアクティビティを開始"); err != nil {
			am.logger.Warn("既存アクティビティの中断に失敗", "entity", activity.Actor, "error", err.Error())
		}
	}

	// アクティビティアクターを取得
	activityActor := GetActivityActor(activity.Type)
	if activityActor == nil {
		return fmt.Errorf("%w: %s", ErrActivityActorNotFound, activity.Type.String())
	}

	// アクティビティアクターでの検証
	if err := activityActor.Validate(activity, world); err != nil {
		return fmt.Errorf("アクティビティ検証失敗: %w", err)
	}

	// 基本的な必須項目チェック
	if err := am.validateBasicRequirements(activity); err != nil {
		return fmt.Errorf("基本要件検証失敗: %w", err)
	}

	// アクティビティを登録
	am.currentActivities[activity.Actor] = activity
	activity.State = ActivityStateRunning

	// アクティビティアクターのStart処理を実行
	if err := activityActor.Start(activity, world); err != nil {
		// 開始に失敗した場合はクリーンアップ
		delete(am.currentActivities, activity.Actor)
		return fmt.Errorf("アクティビティ開始失敗: %w", err)
	}

	am.logger.Debug("アクティビティ開始",
		"entity", activity.Actor,
		"type", activity.Type.String(),
		"duration", activity.TurnsTotal)

	return nil
}

// GetCurrentActivity は指定されたエンティティの現在のアクティビティを取得する
func (am *ActivityManager) GetCurrentActivity(entity ecs.Entity) *Activity {
	return am.currentActivities[entity]
}

// HasActivity は指定されたエンティティがアクティビティを実行中かを返す
func (am *ActivityManager) HasActivity(entity ecs.Entity) bool {
	activity := am.GetCurrentActivity(entity)
	return activity != nil && activity.IsActive()
}

// InterruptActivity は指定されたエンティティのアクティビティを中断する
func (am *ActivityManager) InterruptActivity(entity ecs.Entity, reason string) error {
	activity := am.GetCurrentActivity(entity)
	if activity == nil {
		return ErrActivityNotFound
	}

	return activity.Interrupt(reason)
}

// ResumeActivity は指定されたエンティティのアクティビティを再開する
func (am *ActivityManager) ResumeActivity(entity ecs.Entity, world w.World) error {
	activity := am.GetCurrentActivity(entity)
	if activity == nil {
		return ErrActivityNotFound
	}

	// 再開条件をチェック
	if err := am.validateResume(activity, world); err != nil {
		return fmt.Errorf("アクティビティ再開検証失敗: %w", err)
	}

	return activity.Resume()
}

// CancelActivity は指定されたエンティティのアクティビティをキャンセルする
func (am *ActivityManager) CancelActivity(entity ecs.Entity, reason string, world w.World) {
	activity := am.GetCurrentActivity(entity)
	if activity == nil {
		return
	}

	// アクティビティアクターを取得してCanceled処理を実行
	activityActor := GetActivityActor(activity.Type)
	if activityActor != nil {
		if err := activityActor.Canceled(activity, world); err != nil {
			am.logger.Warn("アクティビティキャンセル処理エラー",
				"entity", entity,
				"error", err.Error())
		}
	}

	// アクティビティ自体をキャンセル状態に
	activity.Cancel(reason)
	delete(am.currentActivities, entity)

	am.logger.Debug("アクティビティキャンセル",
		"entity", entity,
		"type", activity.Type.String(),
		"reason", reason)
}

// ProcessTurn は全てのアクティブなアクティビティの1ターン分の処理を実行する
func (am *ActivityManager) ProcessTurn(world w.World) {
	am.logger.Debug("アクティビティターン処理開始", "count", len(am.currentActivities))

	// 完了・キャンセルされたアクティビティを削除するためのリスト
	var toRemove []ecs.Entity

	for entity, activity := range am.currentActivities {
		// アクティブなアクティビティのみ処理
		if !activity.IsActive() {
			if activity.IsCompleted() || activity.IsCanceled() {
				toRemove = append(toRemove, entity)
			}
			continue
		}

		// アクティビティアクターのDoTurn処理を実行
		activityActor := GetActivityActor(activity.Type)
		if activityActor == nil {
			am.logger.Error("アクティビティアクターが見つかりません",
				"entity", entity,
				"type", activity.Type.String())
			activity.Cancel("アクティビティアクターが見つかりません")
			toRemove = append(toRemove, entity)
			continue
		}

		// ターン処理を実行
		if err := activityActor.DoTurn(activity, world); err != nil {
			am.logger.Error("アクティビティターン処理エラー",
				"entity", entity,
				"type", activity.Type.String(),
				"error", err.Error())

			// エラーが発生した場合はキャンセル
			am.CancelActivity(entity, fmt.Sprintf("エラー: %s", err.Error()), world)
			toRemove = append(toRemove, entity)
			continue
		}

		// 完了したアクティビティの処理
		if activity.IsCompleted() {
			// Finish処理を実行
			if err := activityActor.Finish(activity, world); err != nil {
				am.logger.Error("アクティビティ完了処理エラー",
					"entity", entity,
					"type", activity.Type.String(),
					"error", err.Error())
			}

			am.logger.Debug("アクティビティ完了",
				"entity", entity,
				"type", activity.Type.String())
			toRemove = append(toRemove, entity)
		}
	}

	// 完了・キャンセルされたアクティビティを削除
	for _, entity := range toRemove {
		delete(am.currentActivities, entity)
	}

	am.logger.Debug("アクティビティターン処理完了", "removed", len(toRemove))
}

// GetActivitySummary はアクティビティの要約情報を取得する
func (am *ActivityManager) GetActivitySummary() map[string]interface{} {
	summary := make(map[string]interface{})

	var activeCount, pausedCount, totalCount int
	for _, activity := range am.currentActivities {
		totalCount++
		switch activity.State {
		case ActivityStateRunning:
			activeCount++
		case ActivityStatePaused:
			pausedCount++
		}
	}

	summary["total"] = totalCount
	summary["active"] = activeCount
	summary["paused"] = pausedCount

	return summary
}

// validateBasicRequirements はアクティビティの基本要件を検証する
// 詳細な検証は各アクティビティのValidateメソッドで行う
func (am *ActivityManager) validateBasicRequirements(activity *Activity) error {
	// 基本的なnilチェックのみ実行
	if activity == nil {
		return ErrActivityNil
	}

	if activity.Actor == 0 {
		return ErrActorNotSet
	}

	return nil
}

// validateResume はアクティビティの再開可能性を検証する
func (am *ActivityManager) validateResume(activity *Activity, world w.World) error {
	if !activity.CanResume() {
		return fmt.Errorf("アクティビティ '%s' は再開できません", activity.GetDisplayName())
	}

	// アクティビティアクターでの検証を再実行
	activityActor := GetActivityActor(activity.Type)
	if activityActor != nil {
		if err := activityActor.Validate(activity, world); err != nil {
			return fmt.Errorf("再開時検証失敗: %w", err)
		}
	}

	// 基本要件を再チェック
	return am.validateBasicRequirements(activity)
}

// createActivity はActivityTypeとパラメータからアクティビティを作成する
func (am *ActivityManager) createActivity(activityType ActivityType, params ActionParams, world w.World) (*Activity, error) {
	switch activityType {
	case ActivityMove:
		if params.Destination == nil {
			return nil, fmt.Errorf("移動先が指定されていません")
		}
		characterAP := am.getEntityMaxAP(params.Actor, world)
		duration := CalculateRequiredTurns(ActivityMove, characterAP)
		activity := NewActivity(ActivityMove, params.Actor, duration)
		activity.Position = params.Destination
		return activity, nil

	case ActivityAttack:
		if params.Target == nil {
			return nil, fmt.Errorf("攻撃対象が指定されていません")
		}
		characterAP := am.getEntityMaxAP(params.Actor, world)
		duration := CalculateRequiredTurns(ActivityAttack, characterAP)
		activity := NewActivity(ActivityAttack, params.Actor, duration)
		activity.Target = params.Target
		return activity, nil

	case ActivityPickup:
		return &Activity{
			Actor:      params.Actor,
			Type:       ActivityPickup,
			State:      ActivityStateRunning,
			TurnsLeft:  1,
			TurnsTotal: 1,
			Message:    "アイテムを拾得中...",
			Logger:     logger.New(logger.CategoryAction),
		}, nil

	case ActivityWarp:
		return &Activity{
			Actor:      params.Actor,
			Type:       ActivityWarp,
			State:      ActivityStateRunning,
			TurnsLeft:  1,
			TurnsTotal: 1,
			Message:    "ワープ中...",
			Logger:     logger.New(logger.CategoryAction),
		}, nil

	case ActivityRest:
		duration := params.Duration
		if duration <= 0 {
			characterAP := am.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityRest, characterAP)
		}
		return &Activity{
			Actor:      params.Actor,
			Type:       ActivityRest,
			State:      ActivityStateRunning,
			TurnsLeft:  duration,
			TurnsTotal: duration,
			Message:    "休息中...",
			Logger:     logger.New(logger.CategoryAction),
		}, nil

	case ActivityWait:
		duration := params.Duration
		if duration <= 0 {
			characterAP := am.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityWait, characterAP)
		}
		return &Activity{
			Actor:      params.Actor,
			Type:       ActivityWait,
			State:      ActivityStateRunning,
			TurnsLeft:  duration,
			TurnsTotal: duration,
			Message:    "待機中...",
			Logger:     logger.New(logger.CategoryAction),
		}, nil

	case ActivityRead:
		duration := params.Duration
		if duration <= 0 {
			characterAP := am.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityRead, characterAP)
		}
		activity := NewActivity(ActivityRead, params.Actor, duration)
		activity.Target = params.Target
		return activity, nil

	case ActivityCraft:
		duration := params.Duration
		if duration <= 0 {
			characterAP := am.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityCraft, characterAP)
		}
		activity := NewActivity(ActivityCraft, params.Actor, duration)
		activity.Target = params.Target
		return activity, nil

	default:
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedActivity, activityType)
	}
}

// consumeMoveCost はターン管理システムに移動コストを通知する
func (am *ActivityManager) consumeMoveCost(world w.World, activityType ActivityType, actor ecs.Entity) {
	if world.Resources.TurnManager == nil {
		am.logger.Warn("TurnManagerリソースが見つかりません")
		return
	}

	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	cost, actionName := GetActivityCost(activityType)

	if actor.HasComponent(world.Components.Player) {
		turnManager.ConsumePlayerMoves(actionName, cost)
	} else {
		success := turnManager.ConsumeActionPoints(world, actor, actionName, cost)
		if !success {
			am.logger.Debug("AI移動コスト消費失敗", "actor", actor, "cost", cost)
		}
	}

	am.logger.Debug("移動コスト消費",
		"activity", activityType.String(),
		"cost", cost,
		"actor", actor,
		"isPlayer", actor.HasComponent(world.Components.Player))
}

// getEntityMaxAP はエンティティの最大AP値を取得する
func (am *ActivityManager) getEntityMaxAP(entity ecs.Entity, world w.World) int {
	if turnBasedComponent := world.Components.TurnBased.Get(entity); turnBasedComponent != nil {
		turnBased := turnBasedComponent.(*gc.TurnBased)
		return turnBased.AP.Max
	}
	am.logger.Debug("TurnBasedコンポーネントが見つからない", "entity", entity)
	return 100 // デフォルトAP値
}

// ActionParams はアクション実行時のパラメータを表す
type ActionParams struct {
	Actor       ecs.Entity   // アクションを実行するエンティティ
	Target      *ecs.Entity  // 対象エンティティ（攻撃等で使用）
	Destination *gc.Position // 対象位置（移動等で使用）
	Duration    int          // 継続時間（休息、待機等で使用）
	Reason      string       // 理由（待機等で使用）
}

// ActionResult はアクション実行結果を表す
type ActionResult struct {
	Success      bool         // 実行成功/失敗
	ActivityType ActivityType // 実行されたアクティビティタイプ
	Message      string       // 結果メッセージ
}

// GetActivityCost はアクティビティタイプに応じたコストと名前を返す
func GetActivityCost(activityType ActivityType) (int, string) {
	info := GetActivityInfo(activityType)
	return info.ActionPointCost, info.Name
}
