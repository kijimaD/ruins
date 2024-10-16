package simple

import (
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 意味がないのでこれらのGet系ヘルパーは削除する。
func GetCard(world w.World, entity ecs.Entity) *components.Card {
	gameComponents := world.Components.Game.(*gc.Components)
	card := gameComponents.Card.Get(entity).(*gc.Card)

	return card
}

func GetAttack(world w.World, entity ecs.Entity) *components.Attack {
	gameComponents := world.Components.Game.(*gc.Components)
	attack := gameComponents.Attack.Get(entity).(*gc.Attack)

	return attack
}

func GetWearable(world w.World, entity ecs.Entity) *components.Wearable {
	gameComponents := world.Components.Game.(*gc.Components)
	wearable := gameComponents.Wearable.Get(entity).(*gc.Wearable)

	return wearable
}

func GetMaterial(world w.World, entity ecs.Entity) *components.Material {
	gameComponents := world.Components.Game.(*gc.Components)
	material := gameComponents.Material.Get(entity).(*gc.Material)

	return material
}

func GetDescription(world w.World, entity ecs.Entity) components.Description {
	gameComponents := world.Components.Game.(*gc.Components)
	description := gameComponents.Description.Get(entity).(*gc.Description)

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
