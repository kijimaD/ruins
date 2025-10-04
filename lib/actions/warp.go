package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// WarpActivity はActivityInterfaceの実装
type WarpActivity struct{}

func init() {
	RegisterActivityActor(ActivityWarp, &WarpActivity{})
}

// Info はActivityInterfaceの実装
func (wa *WarpActivity) Info() ActivityInfo {
	return ActivityInfo{
		Type:             ActivityWarp,
		Name:             "ワープ",
		Description:      "ワープホールを使用する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeTime,
		ActionPointCost:  0,
		TotalRequiredAP:  0,
		RequiresTarget:   false,
		RequiresPosition: false,
	}
}

// Validate はワープアクティビティの検証を行う
// Validate はActivityInterfaceの実装
func (wa *WarpActivity) Validate(act *Activity, world w.World) error {
	// プレイヤーの現在位置のワープホールをチェック
	warp := wa.getPlayerWarp(act, world)
	if warp == nil {
		return fmt.Errorf("ワープホールが見つかりません")
	}

	// ワープモードが有効かチェック
	switch warp.Mode {
	case gc.WarpModeNext, gc.WarpModeEscape:
		// 有効なワープモード
	default:
		return fmt.Errorf("不明なワープタイプ: %v", warp.Mode)
	}

	return nil
}

// Start はワープ開始時の処理を実行する
// Start はActivityInterfaceの実装
func (wa *WarpActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("ワープ開始", "actor", act.Actor)
	return nil
}

// DoTurn はワープアクティビティの1ターン分の処理を実行する
// DoTurn はActivityInterfaceの実装
func (wa *WarpActivity) DoTurn(act *Activity, world w.World) error {
	// ワープ可能性をチェック
	warp := wa.getPlayerWarp(act, world)
	if warp == nil {
		act.Cancel("ワープホールが見つかりません")
		return fmt.Errorf("ワープホールがありません")
	}

	// ワープ実行
	if err := wa.performWarp(act, world, warp); err != nil {
		act.Cancel(fmt.Sprintf("ワープエラー: %s", err.Error()))
		return err
	}

	// ワープ処理完了
	act.Complete()
	return nil
}

// Finish はワープ完了時の処理を実行する
// Finish はActivityInterfaceの実装
func (wa *WarpActivity) Finish(act *Activity, world w.World) error {
	act.Logger.Debug("ワープアクティビティ完了", "actor", act.Actor)

	// プレイヤーの場合のみワープ完了メッセージを表示
	if isPlayerActivity(act, world) {
		gamelog.New(gamelog.FieldLog).
			Magic("空間移動した。").
			Log()
	}

	return nil
}

// Canceled はワープキャンセル時の処理を実行する
// Canceled はActivityInterfaceの実装
func (wa *WarpActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("ワープキャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// performWarp は実際のワープ処理を実行する
func (wa *WarpActivity) performWarp(act *Activity, world w.World, warp *gc.Warp) error {

	switch warp.Mode {
	case gc.WarpModeNext:
		world.Resources.Dungeon.SetStateEvent(resources.StateEventWarpNext)
		act.Logger.Debug("次の階へワープ", "actor", act.Actor)
		return nil

	case gc.WarpModeEscape:
		world.Resources.Dungeon.SetStateEvent(resources.StateEventWarpEscape)
		act.Logger.Debug("脱出ワープ", "actor", act.Actor)
		return nil

	default:
		return fmt.Errorf("不明なワープタイプ: %v", warp.Mode)
	}
}

// getPlayerWarp はプレイヤーの現在位置のワープホールを取得する
func (wa *WarpActivity) getPlayerWarp(_ *Activity, world w.World) *gc.Warp {
	// プレイヤーエンティティを探す
	var playerEntity ecs.Entity
	world.Manager.Join(world.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = entity
	}))

	if playerEntity == 0 {
		return nil
	}

	gridElement := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)

	// プレイヤーと同じ座標にあるWarpコンポーネントを探す
	var warp *gc.Warp
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Warp,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		ge := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if ge.X == gridElement.X && ge.Y == gridElement.Y {
			warp = world.Components.Warp.Get(entity).(*gc.Warp)
		}
	}))

	return warp
}
