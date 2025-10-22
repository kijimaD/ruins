package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TriggerActivateActivity はTriggerを発動するActivity
type TriggerActivateActivity struct {
	TriggerEntity ecs.Entity
}

// Info はActivityInterfaceの実装
func (ta *TriggerActivateActivity) Info() ActivityInfo {
	return ActivityInfo{
		Name:            "トリガー発動",
		Description:     "トリガーを発動する",
		Interruptible:   false,
		Resumable:       false,
		ActionPointCost: 0,
		TotalRequiredAP: 0,
	}
}

// String はActivityInterfaceの実装
func (ta *TriggerActivateActivity) String() string {
	return "TriggerActivate"
}

// Validate はトリガー発動アクティビティの検証を行う
func (ta *TriggerActivateActivity) Validate(_ *Activity, world w.World) error {
	// TriggerEntityが存在するかチェック
	if !ta.TriggerEntity.HasComponent(world.Components.Trigger) {
		return fmt.Errorf("指定されたエンティティはTriggerを持っていません")
	}
	return nil
}

// Start はトリガー発動開始時の処理を実行する
func (ta *TriggerActivateActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("トリガー発動開始", "actor", act.Actor, "trigger", ta.TriggerEntity)
	return nil
}

// DoTurn はトリガー発動アクティビティの1ターン分の処理を実行する
func (ta *TriggerActivateActivity) DoTurn(act *Activity, world w.World) error {
	trigger := world.Components.Trigger.Get(ta.TriggerEntity).(*gc.Trigger)

	// Triggerの種類に応じた処理を実行
	switch content := trigger.Data.(type) {
	case gc.WarpNextTrigger:
		ta.executeWarpNext(act, world, content)
	case gc.WarpEscapeTrigger:
		ta.executeWarpEscape(act, world, content)
	default:
		err := fmt.Errorf("未知のトリガータイプ: %T", trigger)
		act.Cancel(fmt.Sprintf("トリガー発動エラー: %s", err.Error()))
		return err
	}

	// Consumableコンポーネントがある場合はエンティティを削除
	if ta.TriggerEntity.HasComponent(world.Components.Consumable) {
		world.Manager.DeleteEntity(ta.TriggerEntity)
	}

	act.Complete()
	return nil
}

// Finish はトリガー発動完了時の処理を実行する
func (ta *TriggerActivateActivity) Finish(act *Activity, _ w.World) error {
	act.Logger.Debug("トリガー発動アクティビティ完了", "actor", act.Actor)
	return nil
}

// Canceled はトリガー発動キャンセル時の処理を実行する
func (ta *TriggerActivateActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("トリガー発動キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// executeWarpNext は次の階へワープするトリガーを実行する
func (ta *TriggerActivateActivity) executeWarpNext(act *Activity, world w.World, _ gc.WarpNextTrigger) {
	world.Resources.Dungeon.SetStateEvent(resources.WarpNextEvent{})
	act.Logger.Debug("次の階へワープ", "actor", act.Actor)
	if isPlayerActivity(act, world) {
		gamelog.New(gamelog.FieldLog).
			Magic("空間移動した。").
			Log()
	}
}

// executeWarpEscape は脱出ワープするトリガーを実行する
func (ta *TriggerActivateActivity) executeWarpEscape(act *Activity, world w.World, _ gc.WarpEscapeTrigger) {
	world.Resources.Dungeon.SetStateEvent(resources.WarpEscapeEvent{})
	act.Logger.Debug("脱出ワープ", "actor", act.Actor)
	if isPlayerActivity(act, world) {
		gamelog.New(gamelog.FieldLog).
			Magic("脱出した。").
			Log()
	}
}
