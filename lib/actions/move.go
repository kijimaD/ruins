package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/movement"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// MoveActivity は移動アクティビティの実装
type MoveActivity struct{}

func init() {
	// 移動アクティビティをレジストリに登録
	RegisterActivityActor(ActivityMove, &MoveActivity{})
}

// Info は移動アクティビティの情報を返す
func (ma *MoveActivity) Info() ActivityInfo {
	return ActivityInfo{
		Type:             ActivityMove,
		Name:             "移動",
		Description:      "隣接するタイルに移動する",
		Interruptible:    false,
		Resumable:        false,
		TimingMode:       TimingModeSpeed,
		ActionPointCost:  100, // 初期AP相当（基本アクション）
		TotalRequiredAP:  100,
		RequiresTarget:   false,
		RequiresPosition: true, // 移動先が必要
	}
}

// Validate は移動アクティビティの検証を行う
func (ma *MoveActivity) Validate(act *Activity, world w.World) error {
	// 移動先の確認
	if act.Position == nil {
		return ErrMoveTargetNotSet
	}

	// 目的地が有効な座標範囲内かチェック
	destX, destY := int(act.Position.X), int(act.Position.Y)
	if destX < 0 || destY < 0 {
		return ErrMoveTargetCoordInvalid
	}

	// GridElementコンポーネントの存在チェック
	gridElement := world.Components.GridElement.Get(act.Actor)
	if gridElement == nil {
		return ErrMoveNoGridElement
	}

	// 移動可能性をチェック
	if !movement.CanMoveTo(world, int(act.Position.X), int(act.Position.Y), act.Actor) {
		return ErrMoveTargetInvalid
	}

	return nil
}

// Start は移動開始時の処理を実行する
func (ma *MoveActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("移動開始", "actor", act.Actor, "destination", *act.Position)
	return nil
}

// DoTurn は移動アクティビティの1ターン分の処理を実行する
func (ma *MoveActivity) DoTurn(act *Activity, world w.World) error {
	// 移動先の確認
	if act.Position == nil {
		act.Cancel("移動先が設定されていません")
		return ErrMoveTargetNotSet
	}

	// 移動可能性をチェック
	if !ma.canMove(act, world) {
		act.Cancel("移動できません")
		return ErrMoveTargetInvalid
	}

	// 移動実行
	if err := ma.performMove(act, world); err != nil {
		act.Cancel(fmt.Sprintf("移動エラー: %s", err.Error()))
		return err
	}

	// 移動処理完了
	act.Complete()
	return nil
}

// Finish は移動完了時の処理を実行する
func (ma *MoveActivity) Finish(act *Activity, _ w.World) error {
	act.Logger.Debug("移動アクティビティ完了", "actor", act.Actor)
	// 移動完了のログは通常は出力しない（頻繁すぎるため）
	return nil
}

// Canceled は移動キャンセル時の処理を実行する
func (ma *MoveActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("移動キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// performMove は実際の移動処理を実行する
func (ma *MoveActivity) performMove(act *Activity, world w.World) error {
	// GridElementを取得
	gridElement := world.Components.GridElement.Get(act.Actor)
	if gridElement == nil {
		return ErrGridElementNotFound
	}

	grid := gridElement.(*gc.GridElement)
	oldX, oldY := int(grid.X), int(grid.Y)

	// 座標を更新
	grid.X = gc.Tile(act.Position.X)
	grid.Y = gc.Tile(act.Position.Y)

	// TODO: 移動だけでなく、ターンを消費するすべての操作で空腹度を上げる必要がする気もする
	ma.increasePlayerHunger(act.Actor, world)

	act.Logger.Debug("移動完了",
		"actor", act.Actor,
		"from", fmt.Sprintf("(%d,%d)", oldX, oldY),
		"to", fmt.Sprintf("(%.1f,%.1f)", act.Position.X, act.Position.Y))

	return nil
}

// increasePlayerHunger はプレイヤーの場合は空腹度を増加させる
func (ma *MoveActivity) increasePlayerHunger(entity ecs.Entity, world w.World) {
	// プレイヤーエンティティのみが対象
	if !entity.HasComponent(world.Components.Player) {
		return
	}

	if hungerComponent := world.Components.Hunger.Get(entity); hungerComponent != nil {
		hunger := hungerComponent.(*gc.Hunger)
		hunger.Increase(1) // 移動1回につき空腹度+1
	}
}

// canMove は移動可能かをチェックする
func (ma *MoveActivity) canMove(act *Activity, world w.World) bool {
	// GridElementコンポーネントの存在チェック
	gridElement := world.Components.GridElement.Get(act.Actor)
	if gridElement == nil {
		return false
	}

	// 移動先の確認
	if act.Position == nil {
		return false
	}

	// movement.CanMoveToを使用して移動可能性をチェック
	return movement.CanMoveTo(world, int(act.Position.X), int(act.Position.Y), act.Actor)
}
