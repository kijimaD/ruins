package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/maingame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCraftMenuSortIntegration(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	state := &CraftMenuState{}

	// テスト用レシピエンティティを作成（名前順ではない順序で作成）
	// 道具レシピ
	recipe1 := world.Manager.NewEntity()
	recipe1.AddComponent(world.Components.Name, &gc.Name{Name: "Potion Recipe"})
	recipe1.AddComponent(world.Components.Recipe, &gc.Recipe{})
	recipe1.AddComponent(world.Components.Consumable, &gc.Consumable{})

	recipe2 := world.Manager.NewEntity()
	recipe2.AddComponent(world.Components.Name, &gc.Name{Name: "Antidote Recipe"})
	recipe2.AddComponent(world.Components.Recipe, &gc.Recipe{})
	recipe2.AddComponent(world.Components.Consumable, &gc.Consumable{})

	recipe3 := world.Manager.NewEntity()
	recipe3.AddComponent(world.Components.Name, &gc.Name{Name: "Elixir Recipe"})
	recipe3.AddComponent(world.Components.Recipe, &gc.Recipe{})
	recipe3.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// queryMenuConsumableのテスト（道具タブ）
	consumables := state.queryMenuConsumable(world)
	require.Len(t, consumables, 3, "道具レシピが3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(consumables[0]).(*gc.Name)
	name2 := world.Components.Name.Get(consumables[1]).(*gc.Name)
	name3 := world.Components.Name.Get(consumables[2]).(*gc.Name)

	assert.Equal(t, "Antidote Recipe", name1.Name, "1番目のレシピ名が正しくない")
	assert.Equal(t, "Elixir Recipe", name2.Name, "2番目のレシピ名が正しくない")
	assert.Equal(t, "Potion Recipe", name3.Name, "3番目のレシピ名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(recipe1)
	world.Manager.DeleteEntity(recipe2)
	world.Manager.DeleteEntity(recipe3)
}

func TestCraftMenuCardSortIntegration(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	state := &CraftMenuState{}

	// カードレシピを作成
	card1 := world.Manager.NewEntity()
	card1.AddComponent(world.Components.Name, &gc.Name{Name: "Thunder Card"})
	card1.AddComponent(world.Components.Recipe, &gc.Recipe{})
	card1.AddComponent(world.Components.Card, &gc.Card{})

	card2 := world.Manager.NewEntity()
	card2.AddComponent(world.Components.Name, &gc.Name{Name: "Fire Card"})
	card2.AddComponent(world.Components.Recipe, &gc.Recipe{})
	card2.AddComponent(world.Components.Card, &gc.Card{})

	card3 := world.Manager.NewEntity()
	card3.AddComponent(world.Components.Name, &gc.Name{Name: "Ice Card"})
	card3.AddComponent(world.Components.Recipe, &gc.Recipe{})
	card3.AddComponent(world.Components.Card, &gc.Card{})

	// queryMenuCardのテスト（手札タブ）
	cards := state.queryMenuCard(world)
	require.Len(t, cards, 3, "カードレシピが3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(cards[0]).(*gc.Name)
	name2 := world.Components.Name.Get(cards[1]).(*gc.Name)
	name3 := world.Components.Name.Get(cards[2]).(*gc.Name)

	assert.Equal(t, "Fire Card", name1.Name, "1番目のカード名が正しくない")
	assert.Equal(t, "Ice Card", name2.Name, "2番目のカード名が正しくない")
	assert.Equal(t, "Thunder Card", name3.Name, "3番目のカード名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(card1)
	world.Manager.DeleteEntity(card2)
	world.Manager.DeleteEntity(card3)
}

func TestCraftMenuWearableSortIntegration(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	state := &CraftMenuState{}

	// 装備レシピを作成
	wearable1 := world.Manager.NewEntity()
	wearable1.AddComponent(world.Components.Name, &gc.Name{Name: "Shield Recipe"})
	wearable1.AddComponent(world.Components.Recipe, &gc.Recipe{})
	wearable1.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable2 := world.Manager.NewEntity()
	wearable2.AddComponent(world.Components.Name, &gc.Name{Name: "Armor Recipe"})
	wearable2.AddComponent(world.Components.Recipe, &gc.Recipe{})
	wearable2.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable3 := world.Manager.NewEntity()
	wearable3.AddComponent(world.Components.Name, &gc.Name{Name: "Helmet Recipe"})
	wearable3.AddComponent(world.Components.Recipe, &gc.Recipe{})
	wearable3.AddComponent(world.Components.Wearable, &gc.Wearable{})

	// queryMenuWearableのテスト（装備タブ）
	wearables := state.queryMenuWearable(world)
	require.Len(t, wearables, 3, "装備レシピが3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(wearables[0]).(*gc.Name)
	name2 := world.Components.Name.Get(wearables[1]).(*gc.Name)
	name3 := world.Components.Name.Get(wearables[2]).(*gc.Name)

	assert.Equal(t, "Armor Recipe", name1.Name, "1番目の装備名が正しくない")
	assert.Equal(t, "Helmet Recipe", name2.Name, "2番目の装備名が正しくない")
	assert.Equal(t, "Shield Recipe", name3.Name, "3番目の装備名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(wearable1)
	world.Manager.DeleteEntity(wearable2)
	world.Manager.DeleteEntity(wearable3)
}
