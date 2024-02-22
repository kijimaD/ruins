package craft

import (
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func Randmize(entity ecs.Entity, world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	if entity.HasComponent(gameComponents.Weapon) {
		weapon := gameComponents.Weapon.Get(entity).(*gc.Weapon)

		weapon.Accuracy += (-10 + rand.Intn(20))        // -10 ~ +9
		weapon.BaseDamage += (-5 + rand.Intn(15))       // -5 ~ +9
		weapon.EnergyConsumption += (-1 + rand.Intn(3)) // -1 ~ +1
	}
}
