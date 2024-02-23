package simple

import (
	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func GetWeapon(world w.World, target ecs.Entity) components.Weapon {
	result := components.Weapon{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Weapon).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Weapon) {
			weapon := gameComponents.Weapon.Get(entity).(*gc.Weapon)
			result = *weapon
		}
	}))

	return result
}

func GetMaterial(world w.World, target ecs.Entity) components.Material {
	result := components.Material{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Material).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Material) {
			material := gameComponents.Material.Get(entity).(*gc.Material)
			result = *material
		}
	}))

	return result
}

// 所持中の素材
func OwnedMaterial(f func(entity ecs.Entity), world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Material,
		gameComponents.InBackpack,
	).Visit(ecs.Visit(f))
}

func GetDescription(world w.World, target ecs.Entity) components.Description {
	result := components.Description{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == target && entity.HasComponent(gameComponents.Description) {
			description := gameComponents.Description.Get(entity).(*gc.Description)
			result = *description
		}
	}))

	return result
}
