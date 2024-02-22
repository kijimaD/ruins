package craft

import (
	"math/rand"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/materialhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 所持数と必要数を比較してクラフト可能か判定する
func CanCraft(world w.World, name string) bool {
	canCraft := true
	for _, recipe := range RequiredMaterials(world, name) {
		if !(materialhelper.GetAmount(recipe.Name, world) >= recipe.Amount) {
			canCraft = false
			break
		}
	}

	return canCraft
}

// 指定したレシピに必要な素材一覧
func RequiredMaterials(world w.World, goal string) []components.RecipeInput {
	required := []components.RecipeInput{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Recipe,
		gameComponents.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		if name.Name == goal {
			recipe := gameComponents.Recipe.Get(entity).(*gc.Recipe)
			for _, r := range recipe.Inputs {
				required = append(required, r)
			}
		}
	}))

	return required
}

func Randomize(world w.World, entity ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	if entity.HasComponent(gameComponents.Weapon) {
		weapon := gameComponents.Weapon.Get(entity).(*gc.Weapon)

		weapon.Accuracy += (-10 + rand.Intn(20))        // -10 ~ +9
		weapon.BaseDamage += (-5 + rand.Intn(15))       // -5  ~ +9
		weapon.EnergyConsumption += (-1 + rand.Intn(3)) // -1  ~ +1
	}
}
