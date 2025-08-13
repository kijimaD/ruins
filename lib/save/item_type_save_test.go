package save

import (
	"os"
	"path/filepath"
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSaveLoadItemTypes(t *testing.T) {
	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "save_test_")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// ワールドを作成
	w, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 防具アイテムを作成
	wearableItem := w.Manager.NewEntity()
	wearableItem.AddComponent(w.Components.Item, &gc.Item{})
	wearableItem.AddComponent(w.Components.Name, &gc.Name{Name: "テスト防具"})
	wearableItem.AddComponent(w.Components.Description, &gc.Description{Description: "テスト用の防具です"})
	wearableItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	wearableItem.AddComponent(w.Components.Wearable, &gc.Wearable{
		Defense:           10,
		EquipmentCategory: gc.EquipmentTorso,
		EquipBonus: gc.EquipBonus{
			Vitality: 2,
		},
	})

	// カードアイテムを作成
	cardItem := w.Manager.NewEntity()
	cardItem.AddComponent(w.Components.Item, &gc.Item{})
	cardItem.AddComponent(w.Components.Name, &gc.Name{Name: "テストカード"})
	cardItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	cardItem.AddComponent(w.Components.Card, &gc.Card{
		TargetType: gc.TargetType{
			TargetGroup: gc.TargetGroupAlly,
			TargetNum:   gc.TargetSingle,
		},
		Cost: 3,
	})

	// 素材アイテムを作成
	materialItem := w.Manager.NewEntity()
	materialItem.AddComponent(w.Components.Name, &gc.Name{Name: "テスト素材"})
	materialItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	materialItem.AddComponent(w.Components.Material, &gc.Material{
		Amount: 5,
	})

	// 消費アイテムを作成
	consumableItem := w.Manager.NewEntity()
	consumableItem.AddComponent(w.Components.Item, &gc.Item{})
	consumableItem.AddComponent(w.Components.Name, &gc.Name{Name: "テスト消費アイテム"})
	consumableItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	consumableItem.AddComponent(w.Components.Consumable, &gc.Consumable{
		TargetType: gc.TargetType{
			TargetGroup: gc.TargetGroupAlly,
			TargetNum:   gc.TargetAll,
		},
	})

	// 攻撃属性を持つアイテムを作成
	attackItem := w.Manager.NewEntity()
	attackItem.AddComponent(w.Components.Item, &gc.Item{})
	attackItem.AddComponent(w.Components.Name, &gc.Name{Name: "テスト武器"})
	attackItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	attackItem.AddComponent(w.Components.Attack, &gc.Attack{
		Damage:      50,
		Accuracy:    90,
		AttackCount: 1,
		Element:     gc.ElementTypeNone,
	})

	// レシピを持つアイテムを作成
	recipeItem := w.Manager.NewEntity()
	recipeItem.AddComponent(w.Components.Item, &gc.Item{})
	recipeItem.AddComponent(w.Components.Name, &gc.Name{Name: "合成可能アイテム"})
	recipeItem.AddComponent(w.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	recipeItem.AddComponent(w.Components.Recipe, &gc.Recipe{
		Inputs: []gc.RecipeInput{
			{Name: "素材A", Amount: 2},
			{Name: "素材B", Amount: 1},
		},
	})

	// セーブマネージャーを作成
	sm := NewSerializationManager(tempDir)

	// ワールドを保存
	err = sm.SaveWorld(w, "test_slot")
	require.NoError(t, err)

	// セーブファイルが存在することを確認
	saveFile := filepath.Join(tempDir, "test_slot.json")
	_, err = os.Stat(saveFile)
	require.NoError(t, err)

	// 新しいワールドを作成してロード
	newWorld, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	err = sm.LoadWorld(newWorld, "test_slot")
	require.NoError(t, err)

	// 防具アイテムが正しくロードされたか確認
	wearableCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Item,
		newWorld.Components.Wearable,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "テスト防具", name.Name)

		desc := newWorld.Components.Description.Get(entity).(*gc.Description)
		assert.Equal(t, "テスト用の防具です", desc.Description)

		wearable := newWorld.Components.Wearable.Get(entity).(*gc.Wearable)
		assert.Equal(t, 10, wearable.Defense)
		assert.Equal(t, gc.EquipmentTorso, wearable.EquipmentCategory)
		assert.Equal(t, 2, wearable.EquipBonus.Vitality)

		wearableCount++
	}))
	assert.Equal(t, 1, wearableCount, "防具アイテムが正しくロードされていない")

	// カードアイテムが正しくロードされたか確認
	cardCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Item,
		newWorld.Components.Card,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "テストカード", name.Name)

		card := newWorld.Components.Card.Get(entity).(*gc.Card)
		assert.Equal(t, gc.TargetGroupAlly, card.TargetType.TargetGroup)
		assert.Equal(t, gc.TargetSingle, card.TargetType.TargetNum)
		assert.Equal(t, 3, card.Cost)

		cardCount++
	}))
	assert.Equal(t, 1, cardCount, "カードアイテムが正しくロードされていない")

	// 素材アイテムが正しくロードされたか確認
	materialCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Material,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "テスト素材", name.Name)

		material := newWorld.Components.Material.Get(entity).(*gc.Material)
		assert.Equal(t, 5, material.Amount)

		materialCount++
	}))
	assert.Equal(t, 1, materialCount, "素材アイテムが正しくロードされていない")

	// 消費アイテムが正しくロードされたか確認
	consumableCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Item,
		newWorld.Components.Consumable,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "テスト消費アイテム", name.Name)

		consumable := newWorld.Components.Consumable.Get(entity).(*gc.Consumable)
		assert.Equal(t, gc.TargetGroupAlly, consumable.TargetType.TargetGroup)
		assert.Equal(t, gc.TargetAll, consumable.TargetType.TargetNum)

		consumableCount++
	}))
	assert.Equal(t, 1, consumableCount, "消費アイテムが正しくロードされていない")

	// 攻撃属性を持つアイテムが正しくロードされたか確認
	attackCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Item,
		newWorld.Components.Attack,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "テスト武器", name.Name)

		attack := newWorld.Components.Attack.Get(entity).(*gc.Attack)
		assert.Equal(t, 50, attack.Damage)
		assert.Equal(t, 90, attack.Accuracy)
		assert.Equal(t, 1, attack.AttackCount)
		assert.Equal(t, gc.ElementTypeNone, attack.Element)

		attackCount++
	}))
	assert.Equal(t, 1, attackCount, "攻撃アイテムが正しくロードされていない")

	// レシピを持つアイテムが正しくロードされたか確認
	recipeCount := 0
	newWorld.Manager.Join(
		newWorld.Components.Item,
		newWorld.Components.Recipe,
		newWorld.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := newWorld.Components.Name.Get(entity).(*gc.Name)
		assert.Equal(t, "合成可能アイテム", name.Name)

		recipe := newWorld.Components.Recipe.Get(entity).(*gc.Recipe)
		require.Len(t, recipe.Inputs, 2)
		assert.Equal(t, "素材A", recipe.Inputs[0].Name)
		assert.Equal(t, 2, recipe.Inputs[0].Amount)
		assert.Equal(t, "素材B", recipe.Inputs[1].Name)
		assert.Equal(t, 1, recipe.Inputs[1].Amount)

		recipeCount++
	}))
	assert.Equal(t, 1, recipeCount, "レシピアイテムが正しくロードされていない")
}
