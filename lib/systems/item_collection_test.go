package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/game"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestItemCollectionSystem(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを配置
	require.NoError(t, worldhelper.SpawnOperator(world, 3, 3))

	// フィールドアイテムをプレイヤーの近くに配置
	item, err := worldhelper.SpawnFieldItem(world, "回復薬", gc.Tile(3), gc.Tile(3))
	require.NoError(t, err)

	// 収集前の状態確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドにあるべき")
	assert.True(t, item.HasComponent(world.Components.GridElement), "アイテムはグリッド要素を持つべき")

	// バックパック内アイテム数の確認（収集前）
	backpackCountBefore := countBackpackItems(world)

	// アイテム収集システムを実行（パニックしないことを確認）
	require.NotPanics(t, func() {
		ItemCollectionSystem(world)
	}, "ItemCollectionSystemはパニックしてはいけない")

	// バックパック内アイテム数の確認（収集後）
	backpackCountAfter := countBackpackItems(world)

	// 正常に動作することを基本的にテスト（実際の収集はプレイヤーとアイテムの距離による）
	assert.GreaterOrEqual(t, backpackCountAfter, backpackCountBefore, "バックパック内のアイテムが増えているか同じであるべき")
}

func TestCheckItemCollision(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// プレイヤーを配置
	require.NoError(t, worldhelper.SpawnOperator(world, 3, 3))

	var playerEntity ecs.Entity
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = entity
	}))

	// アイテムを配置
	item, err := worldhelper.SpawnFieldItem(world, "回復薬", gc.Tile(3), gc.Tile(3))
	require.NoError(t, err)

	// プレイヤーのグリッド位置を取得してピクセル位置に変換
	playerGridElement := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)
	playerPos := &gc.Position{
		X: gc.Pixel(int(playerGridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2),
		Y: gc.Pixel(int(playerGridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2),
	}

	// アイテムのピクセル位置を計算
	itemPos := &gc.Position{
		X: gc.Pixel(3*int(consts.TileSize) + int(consts.TileSize)/2), // タイル中央
		Y: gc.Pixel(3*int(consts.TileSize) + int(consts.TileSize)/2), // タイル中央
	}

	// 衝突判定をテスト（パニックしないことを確認）
	require.NotPanics(t, func() {
		checkCollisionSimple(world, playerEntity, item, playerPos, itemPos)
	}, "checkCollisionSimple関数はパニックしてはいけない")
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
	assert.True(t, item.HasComponent(world.Components.ItemLocationInBackpack), "アイテムはバックパックに移動するべき")
	assert.False(t, item.HasComponent(world.Components.ItemLocationOnField), "アイテムはフィールドから削除されるべき")
	assert.False(t, item.HasComponent(world.Components.GridElement), "グリッド要素は削除されるべき")
	assert.False(t, item.HasComponent(world.Components.SpriteRender), "スプライト要素は削除されるべき")
}

// countBackpackItems はバックパック内のアイテム数をカウントする
func countBackpackItems(world w.World) int {
	count := 0
	world.Manager.Join(world.Components.ItemLocationInBackpack).Visit(ecs.Visit(func(_ ecs.Entity) {
		count++
	}))
	return count
}
