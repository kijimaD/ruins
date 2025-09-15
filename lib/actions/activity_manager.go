package actions

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/logger"
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
func NewActivityManager() *ActivityManager {
	return &ActivityManager{
		currentActivities: make(map[ecs.Entity]*Activity),
		logger:            logger.New(logger.CategoryAction),
	}
}

// StartActivity は新しいアクティビティを開始する
func (am *ActivityManager) StartActivity(activity *Activity, world w.World) error {
	if activity == nil {
		return fmt.Errorf("アクティビティがnilです")
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
		return fmt.Errorf("アクティビティアクターが見つかりません: %s", activity.Type.String())
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
		return fmt.Errorf("アクティビティが見つかりません")
	}

	return activity.Interrupt(reason)
}

// ResumeActivity は指定されたエンティティのアクティビティを再開する
func (am *ActivityManager) ResumeActivity(entity ecs.Entity, world w.World) error {
	activity := am.GetCurrentActivity(entity)
	if activity == nil {
		return fmt.Errorf("再開するアクティビティが見つかりません")
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

// GetAllActiveActivities は全てのアクティブなアクティビティを取得する
func (am *ActivityManager) GetAllActiveActivities() map[ecs.Entity]*Activity {
	result := make(map[ecs.Entity]*Activity)
	for entity, activity := range am.currentActivities {
		if activity.IsActive() {
			result[entity] = activity
		}
	}
	return result
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
		return fmt.Errorf("アクティビティがnilです")
	}

	if activity.Actor == 0 {
		return fmt.Errorf("アクターが設定されていません")
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

// CleanupCompletedActivities は完了したアクティビティをクリーンアップする
func (am *ActivityManager) CleanupCompletedActivities() {
	for entity, activity := range am.currentActivities {
		if activity.IsCompleted() || activity.IsCanceled() {
			delete(am.currentActivities, entity)
		}
	}
}
