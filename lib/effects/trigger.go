package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AddEffectのラッパー群。アイテムからトリガーする

// アイテムについたComponentからEffectを登録する
// 消費アイテムはなくなる
func ItemTrigger(creator *ecs.Entity, item ecs.Entity, targets Targets, world w.World) {
	eventTrigger(creator, item, targets, world)

	gameComponents := world.Components.Game.(*gc.Components)
	_, ok := gameComponents.Consumable.Get(item).(*gc.Consumable)
	if ok {
		world.Manager.DeleteEntity(item)
	}
}

// TODO: 地雷など、フィールド上で一度しか動作しないギミックを動作させるのに使う予定
// Consumable Itemと同じように、使ったあとに消す処理を入れるだろう
// func Trigger() {
//	...
// }

// アイテムからコンポーネントを取り出し、対応したEffectをトリガーする
func eventTrigger(creator *ecs.Entity, entity ecs.Entity, targets Targets, world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	healing, ok := gameComponents.ProvidesHealing.Get(entity).(*gc.ProvidesHealing)
	if ok {
		AddEffect(creator, Healing{Amount: healing.Amount}, targets)
	}

	damage, ok := gameComponents.InflictsDamage.Get(entity).(*gc.InflictsDamage)
	if ok {
		AddEffect(creator, Damage{Amount: damage.Amount}, targets)
	}
}
