package simple

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// メニュー上でいうアイテムを取得する。エンティティ上はCardもItemを持っているので、排除する
func QueryMenuItem(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.InBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity.HasComponent(gameComponents.Card) {
			return
		}

		items = append(items, entity)
	}))

	return items
}
