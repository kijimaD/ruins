package items

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
