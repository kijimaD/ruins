package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ActionAPI は全てのアクションを統一的に管理するAPI
type ActionAPI struct {
	manager *ActivityManager
	logger  *logger.Logger
}

// NewActionAPI は新しいActionAPIを作成する
func NewActionAPI() *ActionAPI {
	actionLogger := logger.New(logger.CategoryAction)
	return &ActionAPI{
		manager: NewActivityManager(actionLogger),
		logger:  actionLogger,
	}
}

// Execute は指定されたアクション（アクティビティ）を実行する
// 即座実行アクション（移動、攻撃等）も継続アクション（休息等）も統一的に処理
func (api *ActionAPI) Execute(activityType ActivityType, params ActionParams, world w.World) (*ActionResult, error) {
	api.logger.Debug("アクション実行開始",
		"type", activityType.String(),
		"actor", params.Actor)

	// アクティビティを作成
	activity, err := api.createActivity(activityType, params, world)
	if err != nil {
		return &ActionResult{
			Success:      false,
			ActivityType: activityType,
			Message:      err.Error(),
		}, err
	}

	// アクティビティを開始
	if err := api.manager.StartActivity(activity, world); err != nil {
		return &ActionResult{
			Success:      false,
			ActivityType: activityType,
			Message:      err.Error(),
		}, err
	}

	// 即座実行アクション（1ターン）の場合は即座に処理
	if activity.TurnsTotal == 1 {
		// ターン処理実行
		api.manager.ProcessTurn(world)

		// ターン管理システムに移動コストを通知
		api.consumeMoveCost(world, activityType, params.Actor)

		// 結果を確認
		currentActivity := api.manager.GetCurrentActivity(params.Actor)
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

// ProcessTurn は全アクティビティの1ターン分の処理を実行する
func (api *ActionAPI) ProcessTurn(world w.World) {
	api.manager.ProcessTurn(world)
}

// InterruptActivity は指定されたエンティティのアクティビティを中断する
func (api *ActionAPI) InterruptActivity(entity ecs.Entity, reason string) error {
	return api.manager.InterruptActivity(entity, reason)
}

// ResumeActivity は指定されたエンティティのアクティビティを再開する
func (api *ActionAPI) ResumeActivity(entity ecs.Entity, world w.World) error {
	return api.manager.ResumeActivity(entity, world)
}

// CancelActivity は指定されたエンティティのアクティビティをキャンセルする
func (api *ActionAPI) CancelActivity(entity ecs.Entity, reason string, world w.World) {
	api.manager.CancelActivity(entity, reason, world)
}

// GetCurrentActivity は指定されたエンティティの現在のアクティビティを取得する
func (api *ActionAPI) GetCurrentActivity(entity ecs.Entity) *Activity {
	return api.manager.GetCurrentActivity(entity)
}

// HasActivity は指定されたエンティティがアクティビティを実行中かを返す
func (api *ActionAPI) HasActivity(entity ecs.Entity) bool {
	return api.manager.HasActivity(entity)
}

// GetActivitySummary は全アクティビティの要約情報を取得する
func (api *ActionAPI) GetActivitySummary() map[string]interface{} {
	return api.manager.GetActivitySummary()
}

// createActivity はActivityTypeとパラメータからアクティビティを作成する
func (api *ActionAPI) createActivity(activityType ActivityType, params ActionParams, world w.World) (*Activity, error) {
	switch activityType {
	case ActivityMove:
		if params.Destination == nil {
			return nil, fmt.Errorf("移動先が指定されていません")
		}
		// AP積み上げ方式で必要ターン数を計算
		characterAP := api.getEntityMaxAP(params.Actor, world)
		duration := CalculateRequiredTurns(ActivityMove, characterAP)
		activity := NewActivity(ActivityMove, params.Actor, duration)
		activity.Position = params.Destination
		return activity, nil

	case ActivityAttack:
		if params.Target == nil {
			return nil, fmt.Errorf("攻撃対象が指定されていません")
		}
		characterAP := api.getEntityMaxAP(params.Actor, world)
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
			// AP積み上げ方式で必要ターン数を計算
			characterAP := api.getEntityMaxAP(params.Actor, world)
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
			// AP積み上げ方式で必要ターン数を計算
			characterAP := api.getEntityMaxAP(params.Actor, world)
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
			// AP積み上げ方式で必要ターン数を計算
			characterAP := api.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityRead, characterAP)
		}
		activity := NewActivity(ActivityRead, params.Actor, duration)
		activity.Target = params.Target
		return activity, nil

	case ActivityCraft:
		duration := params.Duration
		if duration <= 0 {
			// AP積み上げ方式で必要ターン数を計算
			characterAP := api.getEntityMaxAP(params.Actor, world)
			duration = CalculateRequiredTurns(ActivityCraft, characterAP)
		}
		activity := NewActivity(ActivityCraft, params.Actor, duration)
		activity.Target = params.Target
		return activity, nil

	default:
		return nil, fmt.Errorf("未対応のアクティビティタイプ: %v", activityType)
	}
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

// QuickMove は移動アクションのショートカット
func (api *ActionAPI) QuickMove(actor ecs.Entity, dest gc.Position, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor:       actor,
		Destination: &dest,
	}
	return api.Execute(ActivityMove, params, world)
}

// QuickAttack は攻撃アクションのショートカット
func (api *ActionAPI) QuickAttack(actor ecs.Entity, target ecs.Entity, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor:  actor,
		Target: &target,
	}
	return api.Execute(ActivityAttack, params, world)
}

// QuickPickup はアイテム拾得アクションのショートカット
func (api *ActionAPI) QuickPickup(actor ecs.Entity, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor: actor,
	}
	return api.Execute(ActivityPickup, params, world)
}

// QuickWarp はワープアクションのショートカット
func (api *ActionAPI) QuickWarp(actor ecs.Entity, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor: actor,
	}
	return api.Execute(ActivityWarp, params, world)
}

// StartRest は休息アクティビティのショートカット
func (api *ActionAPI) StartRest(actor ecs.Entity, duration int, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor:    actor,
		Duration: duration,
	}
	return api.Execute(ActivityRest, params, world)
}

// StartWait は待機アクティビティのショートカット
func (api *ActionAPI) StartWait(actor ecs.Entity, duration int, reason string, world w.World) (*ActionResult, error) {
	params := ActionParams{
		Actor:    actor,
		Duration: duration,
		Reason:   reason,
	}
	return api.Execute(ActivityWait, params, world)
}

// consumeMoveCost はターン管理システムに移動コストを通知する
func (api *ActionAPI) consumeMoveCost(world w.World, activityType ActivityType, actor ecs.Entity) {
	// ターン管理リソースを取得
	if world.Resources.TurnManager == nil {
		api.logger.Warn("TurnManagerリソースが見つかりません")
		return
	}

	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	// アクティビティ種別に応じた移動コストを計算
	cost, actionName := GetActivityCost(activityType)

	// エンティティタイプに応じて適切なコスト消費メソッドを呼ぶ
	if actor.HasComponent(world.Components.Player) {
		// プレイヤーの場合は従来のPlayerMoves消費
		turnManager.ConsumePlayerMoves(actionName, cost)
	} else {
		// AIエンティティの場合はActionPoints消費
		success := turnManager.ConsumeActionPoints(world, actor, actionName, cost)
		if !success {
			api.logger.Debug("AI移動コスト消費失敗", "actor", actor, "cost", cost)
		}
	}

	api.logger.Debug("移動コスト消費",
		"activity", activityType.String(),
		"cost", cost,
		"actor", actor,
		"isPlayer", actor.HasComponent(world.Components.Player))
}

// GetActivityCost はアクティビティタイプに応じたコストと名前を返す
//
// activityInfosテーブルから情報を取得し、統一的なコスト管理を提供する。
// キャラクターの初期APが100であることを基準として、相対的なコストを定義している。
//
// 返り値:
//   - int: アクションポイントコスト（キャラクター初期AP100を基準とした相対値）
//   - string: アクション名（ログ用）
func GetActivityCost(activityType ActivityType) (int, string) {
	info := GetActivityInfo(activityType)
	return info.ActionPointCost, info.Name
}

// getEntityMaxAP はエンティティの最大AP値を取得する
func (api *ActionAPI) getEntityMaxAP(entity ecs.Entity, world w.World) int {
	// TurnBasedコンポーネントからAP値を取得
	if turnBasedComponent := world.Components.TurnBased.Get(entity); turnBasedComponent != nil {
		turnBased := turnBasedComponent.(*gc.TurnBased)
		return turnBased.AP.Max
	}

	// TurnBasedコンポーネントがない場合（プレイヤーなど）はデフォルト値を返す
	api.logger.Debug("TurnBasedコンポーネントが見つからない", "entity", entity)
	return 100 // デフォルトAP値
}
