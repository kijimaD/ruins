package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAmount(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(world.Components.Stackable, &gc.Stackable{Count: 10})
	materialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "鉄"})

	// 素材の数量を取得
	entity, found := FindStackableInInventory(world, "鉄")
	require.True(t, found, "素材が見つからない")
	stackable := world.Components.Stackable.Get(entity).(*gc.Stackable)
	assert.Equal(t, 10, stackable.Count, "素材の数量が正しく取得できない")

	// 存在しない素材の数量を取得
	_, found = FindStackableInInventory(world, "存在しない素材")
	assert.False(t, found, "存在しない素材が見つかってはいけない")

	// クリーンアップ
	world.Manager.DeleteEntity(materialEntity)
}

func TestPlusMinusAmount(t *testing.T) {
	t.Parallel()
	world := testutil.InitTestWorld(t)

	// テスト用素材エンティティを作成
	materialEntity := world.Manager.NewEntity()
	materialEntity.AddComponent(world.Components.Stackable, &gc.Stackable{Count: 10})
	materialEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	materialEntity.AddComponent(world.Components.Name, &gc.Name{Name: "鉄"})

	// 数量を増加
	err := AddStackableCount(world, "鉄", 5)
	require.NoError(t, err)
	entity, found := FindStackableInInventory(world, "鉄")
	require.True(t, found)
	stackable := world.Components.Stackable.Get(entity).(*gc.Stackable)
	assert.Equal(t, 15, stackable.Count, "数量増加が正しく動作しない")

	// 数量を減少
	err = RemoveStackableCount(world, "鉄", 3)
	require.NoError(t, err)
	entity, found = FindStackableInInventory(world, "鉄")
	require.True(t, found)
	stackable = world.Components.Stackable.Get(entity).(*gc.Stackable)
	assert.Equal(t, 12, stackable.Count, "数量減少が正しく動作しない")

	// 大量追加テスト（制限なし）
	err = AddStackableCount(world, "鉄", 1000)
	require.NoError(t, err)
	entity, found = FindStackableInInventory(world, "鉄")
	require.True(t, found)
	stackable = world.Components.Stackable.Get(entity).(*gc.Stackable)
	assert.Equal(t, 1012, stackable.Count, "数量が正しく加算されない")

	// 0以下になるとエンティティが削除される
	err = RemoveStackableCount(world, "鉄", 1500)
	require.NoError(t, err)
	_, found = FindStackableInInventory(world, "鉄")
	assert.False(t, found, "0以下になったらエンティティが削除されるべき")
}
