package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestAddToInventoryWithMaterial(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 既存のmaterialをバックパックに配置（初期数量5）
	existingMaterial := SpawnMaterial(world, "鉄くず", 5, gc.ItemLocationInBackpack)

	// 新しいmaterialを作成（数量3）
	newMaterial := SpawnMaterial(world, "鉄くず", 3, gc.ItemLocationInBackpack)

	// AddToInventoryを実行
	AddToInventory(world, newMaterial, "鉄くず")

	// 既存のmaterialの数量が統合されていることを確認
	updatedMat := world.Components.Material.Get(existingMaterial).(*gc.Material)
	assert.Equal(t, 8, updatedMat.Amount, "数量が正しく統合されていない")

	// 新しいmaterialエンティティが削除されていることを確認（コンポーネントが存在しない）
	assert.False(t, newMaterial.HasComponent(world.Components.Material), "新しいmaterialエンティティが削除されていない")
}

func TestAddToInventoryWithNewMaterial(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 新しいmaterialを作成（既存のものはなし）
	newMaterial := SpawnMaterial(world, "緑ハーブ", 2, gc.ItemLocationInBackpack)

	// バックパック内のmaterial数をカウント（統合前）
	materialCountBefore := 0
	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		materialCountBefore++
	}))

	// AddToInventoryを実行
	AddToInventory(world, newMaterial, "緑ハーブ")

	// バックパック内のmaterial数をカウント（統合後）
	materialCountAfter := 0
	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		materialCountAfter++
	}))

	// 新しいmaterialとして追加されていることを確認
	assert.Equal(t, materialCountBefore, materialCountAfter, "新しいmaterialが正しく追加されていない")
	assert.True(t, newMaterial.HasComponent(world.Components.Material), "新しいmaterialエンティティが生きているべき")

	// 数量が維持されていることを確認
	updatedMat := world.Components.Material.Get(newMaterial).(*gc.Material)
	assert.Equal(t, 2, updatedMat.Amount, "数量が維持されていない")
}

func TestAddToInventoryWithNonMaterial(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 既存のアイテム（material以外）をバックパックに配置
	existingItem := SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)

	// 新しい同じアイテムを作成
	newItem := SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)

	// バックパック内のアイテム数をカウント（統合前）
	itemCountBefore := 0
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		itemCountBefore++
	}))

	// AddToInventoryを実行
	AddToInventory(world, newItem, "回復薬")

	// バックパック内のアイテム数をカウント（統合後）
	itemCountAfter := 0
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		itemCountAfter++
	}))

	// material以外は統合されず、2つのアイテムが存在することを確認
	assert.Equal(t, itemCountBefore, itemCountAfter, "material以外は統合されないべき")
	assert.True(t, existingItem.HasComponent(world.Components.Item), "既存のアイテムエンティティが生きているべき")
	assert.True(t, newItem.HasComponent(world.Components.Item), "新しいアイテムエンティティが生きているべき")
}

func TestAddToInventoryWithoutItemOrMaterialComponent(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// ItemもMaterialコンポーネントも持たないエンティティを作成
	componentList := entities.ComponentList{}
	componentList.Game = append(componentList.Game, gc.GameComponentList{
		Name: &gc.Name{Name: "テスト"},
	})
	entities := entities.AddEntities(world, componentList)
	nonItemEntity := entities[0]

	// AddToInventoryを実行してパニックが発生することを確認
	require.Panics(t, func() {
		AddToInventory(world, nonItemEntity, "テスト")
	}, "ItemもMaterialコンポーネントも持たないエンティティに対してパニックするべき")
}
