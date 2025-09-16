package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/worldhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestItemCollisionDetection(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを配置
	_, err = worldhelper.SpawnPlayer(world, 3, 3, "セレスティン")
	require.NoError(t, err)

	var playerEntity ecs.Entity
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = entity
	}))

	// アイテムを配置
	item, err := worldhelper.SpawnFieldItem(world, "回復薬", gc.Tile(3), gc.Tile(3))
	require.NoError(t, err)

	// GridElementベースの衝突判定をテスト
	playerGrid := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)
	itemGrid := world.Components.GridElement.Get(item).(*gc.GridElement)

	// 同じタイルにいるかテスト
	isSameTile := playerGrid.X == itemGrid.X && playerGrid.Y == itemGrid.Y
	assert.True(t, isSameTile, "プレイヤーとアイテムは同じタイルにいるべき")
}

func TestCollectFieldItem(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// フィールドアイテムを配置
	item, err := worldhelper.SpawnFieldItem(world, "回復薬", gc.Tile(5), gc.Tile(5))
	require.NoError(t, err)

	// 収集前の状態確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドにあるべき")
	assert.True(t, item.HasComponent(world.Components.GridElement), "アイテムはグリッド要素を持つべき")

	// アイテムを収集（パニックしないことを確認）
	require.NotPanics(t, func() {
		collectFieldItem(world, item)
	}, "collectFieldItem関数はパニックしてはいけない")

	// 収集後の状態確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationInBackpack), "アイテムはバックパックにあるべき")
	assert.False(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドにないべき")
	assert.False(t, item.HasComponent(world.Components.GridElement), "アイテムはグリッド要素を持たないべき")
	assert.False(t, item.HasComponent(world.Components.SpriteRender), "アイテムはスプライトを持たないべき")
}

func TestHandleItemCollectionInput(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーとアイテムを同じ位置に配置
	_, err = worldhelper.SpawnPlayer(world, 5, 5, "セレスティン")
	require.NoError(t, err)
	item, err := worldhelper.SpawnFieldItem(world, "回復薬", gc.Tile(5), gc.Tile(5))
	require.NoError(t, err)

	// 手動収集実行前の状態確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドにあるべき")

	// HandleItemCollectionInputを実行（Enterキー入力をシミュレート）
	require.NotPanics(t, func() {
		HandleItemCollectionInput(world)
	}, "HandleItemCollectionInputはパニックしてはいけない")

	// 手動収集実行後の状態確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationInBackpack), "アイテムはバックパックにあるべき")
	assert.False(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドにないべき")
}
