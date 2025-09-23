package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/maingame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestMergeMaterialIntoInventoryWithMaterial(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// 既存のmaterialをバックパックに配置（初期数量5）
	existingMaterial, err := SpawnMaterial(world, "鉄くず", 5, gc.ItemLocationInBackpack)
	require.NoError(t, err)

	// 新しいmaterialを作成（数量3）
	newMaterial, err := SpawnMaterial(world, "鉄くず", 3, gc.ItemLocationInBackpack)
	require.NoError(t, err)

	// MergeMaterialIntoInventoryを実行
	err = MergeMaterialIntoInventory(world, newMaterial, "鉄くず")
	require.NoError(t, err)

	// 既存のmaterialの数量が統合されていることを確認
	updatedMat := world.Components.Material.Get(existingMaterial).(*gc.Material)
	assert.Equal(t, 8, updatedMat.Amount, "数量が正しく統合されていない")

	// 新しいmaterialエンティティが削除されていることを確認（コンポーネントが存在しない）
	assert.False(t, newMaterial.HasComponent(world.Components.Material), "新しいmaterialエンティティが削除されていない")
}

func TestMergeMaterialIntoInventoryWithNewMaterial(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// 新しいmaterialを作成（既存のものはなし）
	newMaterial, err := SpawnMaterial(world, "緑ハーブ", 2, gc.ItemLocationInBackpack)
	require.NoError(t, err)

	// バックパック内のmaterial数をカウント（統合前）
	materialCountBefore := 0
	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		materialCountBefore++
	}))

	// MergeMaterialIntoInventoryを実行
	err = MergeMaterialIntoInventory(world, newMaterial, "緑ハーブ")
	require.NoError(t, err)

	// バックパック内のmaterial数をカウント（統合後）
	materialCountAfter := 0
	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		materialCountAfter++
	}))

	// 新しいmaterialとして追加されていることを確認
	assert.Equal(t, materialCountBefore, materialCountAfter, "新しいmaterialが正しく追加されていない")
	assert.True(t, newMaterial.HasComponent(world.Components.Material), "新しいmaterialエンティティが生きているべき")

	// 数量が維持されていることを確認
	updatedMat := world.Components.Material.Get(newMaterial).(*gc.Material)
	assert.Equal(t, 2, updatedMat.Amount, "数量が維持されていない")
}

func TestMergeMaterialIntoInventoryWithNonMaterial(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// 既存のアイテム（material以外）をバックパックに配置
	existingItem, err := SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	require.NoError(t, err)

	// 新しい同じアイテムを作成
	newItem, err := SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	require.NoError(t, err)

	// バックパック内のアイテム数をカウント（統合前）
	itemCountBefore := 0
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		itemCountBefore++
	}))

	// MergeMaterialIntoInventoryを実行
	err = MergeMaterialIntoInventory(world, newItem, "回復薬")
	require.NoError(t, err)

	// バックパック内のアイテム数をカウント（統合後）
	itemCountAfter := 0
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		itemCountAfter++
	}))

	// material以外は統合されず、2つのアイテムが存在することを確認
	assert.Equal(t, itemCountBefore, itemCountAfter, "material以外は統合されないべき")
	assert.True(t, existingItem.HasComponent(world.Components.Item), "既存のアイテムエンティティが生きているべき")
	assert.True(t, newItem.HasComponent(world.Components.Item), "新しいアイテムエンティティが生きているべき")
}

func TestMergeMaterialIntoInventoryWithoutItemOrMaterialComponent(t *testing.T) {
	t.Parallel()
	world, err := maingame.InitWorld(960, 720)
	require.NoError(t, err)

	// ItemもMaterialコンポーネントも持たないエンティティを作成
	componentList := entities.ComponentList[gc.EntitySpec]{}
	componentList.Entities = append(componentList.Entities, gc.EntitySpec{
		Name: &gc.Name{Name: "テスト"},
	})
	entities := entities.AddEntities(world, componentList)
	nonItemEntity := entities[0]

	// MergeMaterialIntoInventoryを実行してエラーが発生することを確認
	err = MergeMaterialIntoInventory(world, nonItemEntity, "テスト")
	require.Error(t, err, "ItemもMaterialコンポーネントも持たないエンティティに対してエラーを返すべき")
	assert.Contains(t, err.Error(), "does not have Item or Material component", "エラーメッセージが適切であるべき")
}
