package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// アイテムからイベントをトリガーする
func ItemTrigger(creator *ecs.Entity, item ecs.Entity, targets Targets, world w.World) {
	EventTrigger(creator, item, targets, world)

	gameComponents := world.Components.Game.(*gc.Components)
	_, ok := gameComponents.Consumable.Get(item).(*gc.Consumable)
	if ok {
		world.Manager.DeleteEntity(item)
	}
}

// イベントをトリガーする
func EventTrigger(creator *ecs.Entity, item ecs.Entity, targets Targets, world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	healing, ok := gameComponents.ProvidesHealing.Get(item).(*gc.ProvidesHealing)
	if ok {
		AddEffect(creator, Healing{Amount: healing.Amount}, targets)
	}

	damage, ok := gameComponents.InflictsDamage.Get(item).(*gc.InflictsDamage)
	if ok {
		AddEffect(creator, Damage{Amount: damage.Amount}, targets)
	}
}
