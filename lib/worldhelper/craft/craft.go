package craft

import (
	"fmt"
	"math/rand"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/spawner"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func Craft(world w.World, name string) (*ecs.Entity, error) {
	if !CanCraft(world, name) {
		return nil, fmt.Errorf("必要素材が足りない")
	}

	resultEntity := spawner.SpawnItem(world, name, raw.SpawnInBackpack)
	randomize(world, resultEntity)
	consumeMaterials(world, name)

	return &resultEntity, nil
}

// 所持数と必要数を比較してクラフト可能か判定する
func CanCraft(world w.World, name string) bool {
	canCraft := true
	for _, recipeInput := range requiredMaterials(world, name) {
		if !(material.GetAmount(recipeInput.Name, world) >= recipeInput.Amount) {
			canCraft = false
			break
		}
	}

	return canCraft
}

// アイテム合成に必要な素材を消費する
func consumeMaterials(world w.World, goal string) {
	for _, recipeInput := range requiredMaterials(world, goal) {
		material.MinusAmount(recipeInput.Name, recipeInput.Amount, world)
	}
}

// 指定したレシピに必要な素材一覧
func requiredMaterials(world w.World, goal string) []components.RecipeInput {
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

func randomize(world w.World, entity ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	if entity.HasComponent(gameComponents.Weapon) {
		weapon := gameComponents.Weapon.Get(entity).(*gc.Weapon)

		weapon.Accuracy += (-10 + rand.Intn(20))        // -10 ~ +9
		weapon.BaseDamage += (-5 + rand.Intn(15))       // -5  ~ +9
		weapon.EnergyConsumption += (-1 + rand.Intn(3)) // -1  ~ +1
	}
	if entity.HasComponent(gameComponents.Wearable) {
		wearable := gameComponents.Wearable.Get(entity).(*gc.Wearable)

		wearable.BaseDefense += (-4 + rand.Intn(20)) // -4 ~ +9
	}
}
