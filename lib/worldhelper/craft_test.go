package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
)

func TestCanCraft(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game.(*gc.Components)

	// レシピエンティティを作成
	recipe := world.Manager.NewEntity()
	recipe.AddComponent(gameComponents.Recipe, &gc.Recipe{
		Inputs: []gc.RecipeInput{
			{Name: "鉄", Amount: 2},
			{Name: "木", Amount: 1},
		},
	})
	recipe.AddComponent(gameComponents.Name, &gc.Name{Name: "鉄剣"})

	// 必要な素材を作成（十分な量）
	ironMaterial := world.Manager.NewEntity()
	ironMaterial.AddComponent(gameComponents.Material, &gc.Material{Amount: 5})
	ironMaterial.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	ironMaterial.AddComponent(gameComponents.Name, &gc.Name{Name: "鉄"})

	woodMaterial := world.Manager.NewEntity()
	woodMaterial.AddComponent(gameComponents.Material, &gc.Material{Amount: 2})
	woodMaterial.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	woodMaterial.AddComponent(gameComponents.Name, &gc.Name{Name: "木"})

	// クラフト可能かテスト
	canCraft, err := CanCraft(world, "鉄剣")
	assert.True(t, canCraft, "十分な素材があるときはクラフト可能であるべき")
	assert.NoError(t, err, "十分な素材があるときはエラーが発生してはいけない")

	// 素材が不足している場合のテスト
	woodMaterialComp := gameComponents.Material.Get(woodMaterial).(*gc.Material)
	woodMaterialComp.Amount = 0 // 木の量を0にする

	canCraft, err = CanCraft(world, "鉄剣")
	assert.False(t, canCraft, "素材が不足しているときはクラフト不可能であるべき")
	assert.NoError(t, err, "素材が不足してもエラーは発生しないべき")

	// 存在しないレシピのテスト
	canCraft, err = CanCraft(world, "存在しない武器")
	assert.False(t, canCraft, "存在しないレシピはクラフト不可能であるべき")
	assert.Error(t, err, "存在しないレシピでエラーが発生するべき")
	assert.Contains(t, err.Error(), "レシピが存在しません", "エラーメッセージにレシピ不存在の内容が含まれるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(recipe)
	world.Manager.DeleteEntity(ironMaterial)
	world.Manager.DeleteEntity(woodMaterial)
}

func TestCraft(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game.(*gc.Components)

	// 存在しないレシピでのクラフト試行
	result, err := Craft(world, "存在しない武器")
	assert.Nil(t, result, "存在しないレシピでは結果がnilであるべき")
	assert.Error(t, err, "存在しないレシピでエラーが返されるべき")
	assert.Contains(t, err.Error(), "レシピが存在しません", "エラーメッセージにレシピ不存在の内容が含まれるべき")

	// レシピエンティティを作成
	recipe := world.Manager.NewEntity()
	recipe.AddComponent(gameComponents.Recipe, &gc.Recipe{
		Inputs: []gc.RecipeInput{
			{Name: "鉄", Amount: 1},
		},
	})
	recipe.AddComponent(gameComponents.Name, &gc.Name{Name: "簡単な剣"})

	// 素材不足でのクラフト試行
	result, err = Craft(world, "簡単な剣")
	assert.Nil(t, result, "素材不足では結果がnilであるべき")
	assert.Error(t, err, "素材不足でエラーが返されるべき")
	assert.Contains(t, err.Error(), "必要素材が足りません", "エラーメッセージに素材不足の内容が含まれるべき")

	// クリーンアップ
	world.Manager.DeleteEntity(recipe)
}
