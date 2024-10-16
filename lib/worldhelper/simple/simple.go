package simple

import (
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 意味がないのでこれらのGet系ヘルパーは削除する。
func GetCard(world w.World, target ecs.Entity) *components.Card {
	var result *components.Card
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Card).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Card) {
			result = gameComponents.Card.Get(entity).(*gc.Card)
		}
	}))

	return result
}

func GetAttack(world w.World, target ecs.Entity) *components.Attack {
	var result *components.Attack
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Attack).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Attack) {
			result = gameComponents.Attack.Get(entity).(*gc.Attack)
		}
	}))

	return result
}

func GetWearable(world w.World, target ecs.Entity) *components.Wearable {
	var result *components.Wearable
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Wearable).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Wearable) {
			result = gameComponents.Wearable.Get(entity).(*gc.Wearable)
		}
	}))

	return result
}

func GetMaterial(world w.World, target ecs.Entity) *components.Material {
	var result *components.Material
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Material).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Material) {
			result = gameComponents.Material.Get(entity).(*gc.Material)
		}
	}))

	return result
}

func GetDescription(world w.World, target ecs.Entity) components.Description {
	gameComponents := world.Components.Game.(*gc.Components)
	description := gameComponents.Description.Get(target).(*gc.Description)

	return *description
}

// 所持中の素材
// TODO: worldを先に置く
func OwnedMaterial(f func(entity ecs.Entity), world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Material,
		gameComponents.ItemLocationInBackpack,
	).Visit(ecs.Visit(f))
}

// パーティメンバー
func InPartyMember(world w.World, f func(entity ecs.Entity)) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(f))
}
