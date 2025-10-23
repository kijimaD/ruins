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

		if !worldhelper.IsInActivationRange(playerGrid, triggerGrid, trigger.Data.Config().ActivationRange) {
			return
		}

		triggersToProcess = append(triggersToProcess, entity)
	}))

	// 検索した自動実行トリガーを処理する
	for _, triggerEntity := range triggersToProcess {
		trigger := world.Components.Trigger.Get(triggerEntity).(*gc.Trigger)
		config := trigger.Data.Config()

		// Triggerの設定が有効かチェック
		if err := config.ActivationRange.Valid(); err != nil {
			logger.New(logger.CategoryAction).Warn("無効なActivationRangeを持つトリガーをスキップ",
				"entity", triggerEntity,
				"range", config.ActivationRange,
				"error", err)
			continue
		}
		if err := config.ActivationWay.Valid(); err != nil {
			logger.New(logger.CategoryAction).Warn("無効なActivationWayを持つトリガーをスキップ",
				"entity", triggerEntity,
				"way", config.ActivationWay,
				"error", err)
			continue
		}

		if config.ActivationWay != gc.ActivationWayAuto {
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
