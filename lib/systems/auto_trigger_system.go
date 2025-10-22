package systems

import (
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AutoTriggerSystem はプレイヤーがAutoExecute=trueのTriggerに接触した際に自動実行する
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

	// プレイヤーと同じ位置にあるTriggerを検索
	var triggersToProcess []ecs.Entity
	world.Manager.Join(
		world.Components.Trigger,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		triggerGrid := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if triggerGrid.X == playerGrid.X && triggerGrid.Y == playerGrid.Y {
			triggersToProcess = append(triggersToProcess, entity)
		}
	}))

	// AutoExecute=trueのTriggerのみ自動実行
	for _, triggerEntity := range triggersToProcess {
		trigger := world.Components.Trigger.Get(triggerEntity).(*gc.Trigger)

		// AutoExecute=falseのTriggerは手動実行のみ（Enterキー）
		if !trigger.AutoExecute {
			continue
		}

		// AutoExecute=trueのTriggerを自動実行
		activity := &actions.TriggerActivateActivity{
			TriggerEntity: triggerEntity,
		}
		params := actions.ActionParams{
			Actor: playerEntity,
		}
		manager := actions.NewActivityManager(logger.New(logger.CategoryAction))
		_, _ = manager.Execute(activity, params, world)
	}

	return nil
}
