package worldhelper

import (
	"fmt"
	"math/rand/v2"

	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

// Craft はアイテムをクラフトする
func Craft(world w.World, name string) (*ecs.Entity, error) {
	canCraft, err := CanCraft(world, name)
	if err != nil {
		// レシピが存在しない場合
		return nil, err
	}
	if !canCraft {
		// 素材不足の場合
		return nil, fmt.Errorf("必要素材が足りません")
	}

	resultEntity := SpawnItem(world, name, gc.ItemLocationInBackpack)
	randomize(world, resultEntity)
	consumeMaterials(world, name)

	return &resultEntity, nil
}

// CanCraft は所持数と必要数を比較してクラフト可能か判定する
func CanCraft(world w.World, name string) (bool, error) {
	required := requiredMaterials(world, name)
	// レシピが存在しない場合はエラー
	if len(required) == 0 {
		return false, fmt.Errorf("レシピが存在しません: %s", name)
	}

	// 素材不足をチェック（素材不足はエラーではなくfalseを返す）
	for _, recipeInput := range required {
		currentAmount := GetAmount(recipeInput.Name, world)
		if currentAmount < recipeInput.Amount {
			return false, nil // 素材不足はエラーではない
		}
	}

	return true, nil
}

// consumeMaterials はアイテム合成に必要な素材を消費する
func consumeMaterials(world w.World, goal string) {
	for _, recipeInput := range requiredMaterials(world, goal) {
		MinusAmount(recipeInput.Name, recipeInput.Amount, world)
	}
}

// requiredMaterials は指定したレシピに必要な素材一覧
func requiredMaterials(world w.World, goal string) []gc.RecipeInput {
	required := []gc.RecipeInput{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Recipe,
		gameComponents.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		if name.Name == goal {
			recipe := gameComponents.Recipe.Get(entity).(*gc.Recipe)
			required = append(required, recipe.Inputs...)
		}
	}))

	return required
}

// randomize はアイテムにランダム値を設定する
func randomize(world w.World, entity ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	if entity.HasComponent(gameComponents.Attack) {
		attack := gameComponents.Attack.Get(entity).(*gc.Attack)

		attack.Accuracy += (-10 + rand.IntN(20)) // -10 ~ +9
		attack.Damage += (-5 + rand.IntN(15))    // -5  ~ +9
	}
	if entity.HasComponent(gameComponents.Wearable) {
		wearable := gameComponents.Wearable.Get(entity).(*gc.Wearable)

		wearable.Defense += (-4 + rand.IntN(20)) // -4 ~ +9
	}
}
