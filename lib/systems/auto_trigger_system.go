package systems

import (
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AutoTriggerSystem はプレイヤーが自動実行のTriggerに接触した際に自動実行する
func AutoTriggerSystem(world w.World) error {
	// プレイヤーエンティティを取得
	playerEntity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return err
	}

	// プレイヤーの位置を取得
	if !playerEntity.HasComponent(world.Components.GridElement) {
		return nil
	}
	playerGrid := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)

	// プレイヤーの範囲内にあるTriggerを検索
	var triggersToProcess []ecs.Entity
	world.Manager.Join(
		world.Components.Trigger,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		trigger := world.Components.Trigger.Get(entity).(*gc.Trigger)
		triggerGrid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// ActivationRangeに応じた範囲チェック
		if !isInActivationRange(playerGrid, triggerGrid, trigger.ActivationRange) {
			return
		}

		triggersToProcess = append(triggersToProcess, entity)
	}))

	// 自動実行トリガーを処理する
	for _, triggerEntity := range triggersToProcess {
		trigger := world.Components.Trigger.Get(triggerEntity).(*gc.Trigger)

		// 手動実行はスルー
		if trigger.ActivationMode != gc.ActivationModeAuto {
			continue
		}

		// 自動実行のTriggerを実行する
		activity := &actions.TriggerActivateActivity{
			TriggerEntity: triggerEntity,
		}
		params := actions.ActionParams{
			Actor: playerEntity,
		}
		manager := actions.NewActivityManager(logger.New(logger.CategoryAction))
		_, err := manager.Execute(activity, params, world)
		if err != nil {
			return err
		}
	}

	return nil
}

// isInActivationRange はプレイヤーがトリガーの発動範囲内にいるかを判定する
func isInActivationRange(playerGrid, triggerGrid *gc.GridElement, activationRange gc.ActivationRange) bool {
	switch activationRange {
	case gc.ActivationRangeSameTile:
		// 直上（同じタイル）
		return playerGrid.X == triggerGrid.X && playerGrid.Y == triggerGrid.Y
	case gc.ActivationRangeAdjacent:
		// 隣接タイル（近傍8タイル）
		diffX := int(playerGrid.X) - int(triggerGrid.X)
		diffY := int(playerGrid.Y) - int(triggerGrid.Y)
		dx := max(diffX, -diffX)
		dy := max(diffY, -diffY)
		return dx <= 1 && dy <= 1
	default:
		return false
	}
}
