package states

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEquipMenuWearSortIntegration(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	state := &EquipMenuState{}

	// テスト用防具エンティティを作成（名前順ではない順序で作成）
	wearable1 := world.Manager.NewEntity()
	wearable1.AddComponent(world.Components.Name, &gc.Name{Name: "Shield"})
	wearable1.AddComponent(world.Components.Item, &gc.Item{})
	wearable1.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable1.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable2 := world.Manager.NewEntity()
	wearable2.AddComponent(world.Components.Name, &gc.Name{Name: "Armor"})
	wearable2.AddComponent(world.Components.Item, &gc.Item{})
	wearable2.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable2.AddComponent(world.Components.Wearable, &gc.Wearable{})

	wearable3 := world.Manager.NewEntity()
	wearable3.AddComponent(world.Components.Name, &gc.Name{Name: "Helmet"})
	wearable3.AddComponent(world.Components.Item, &gc.Item{})
	wearable3.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	wearable3.AddComponent(world.Components.Wearable, &gc.Wearable{})

	// queryMenuWearのテスト
	wearables := state.queryMenuWear(world)
	require.Len(t, wearables, 3, "防具が3つ見つかるべき")

	// ソート順を確認（名前順）
	name1 := world.Components.Name.Get(wearables[0]).(*gc.Name)
	name2 := world.Components.Name.Get(wearables[1]).(*gc.Name)
	name3 := world.Components.Name.Get(wearables[2]).(*gc.Name)

	assert.Equal(t, "Armor", name1.Name, "1番目の防具名が正しくない")
	assert.Equal(t, "Helmet", name2.Name, "2番目の防具名が正しくない")
	assert.Equal(t, "Shield", name3.Name, "3番目の防具名が正しくない")

	// クリーンアップ
	world.Manager.DeleteEntity(wearable1)
	world.Manager.DeleteEntity(wearable2)
	world.Manager.DeleteEntity(wearable3)
}

func TestEquipMenuCardSortIntegration(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	state := &EquipMenuState{}

	// テスト用カードエンティティを作成（名前順ではない順序で作成）
	card1 := world.Manager.NewEntity()
	card1.AddComponent(world.Components.Name, &gc.Name{Name: "Thunder Card"})
	card1.AddComponent(world.Components.Item, &gc.Item{})
	card1.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	card1.AddComponent(world.Components.Card, &gc.Card{})

	card2 := world.Manager.NewEntity()
	card2.AddComponent(world.Components.Name, &gc.Name{Name: "Fire Card"})
	card2.AddComponent(world.Components.Item, &gc.Item{})
	card2.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	card2.AddComponent(world.Components.Card, &gc.Card{})

	card3 := world.Manager.NewEntity()
	card3.AddComponent(world.Components.Name, &gc.Name{Name: "Ice Card"})
	card3.AddComponent(world.Components.Item, &gc.Item{})
	card3.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	card3.AddComponent(world.Components.Card, &gc.Card{})

	// queryMenuCardのテスト
	cards := state.queryMenuCard(world)
	require.Len(t, cards, 3, "カードが3つ見つかるべき")

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
