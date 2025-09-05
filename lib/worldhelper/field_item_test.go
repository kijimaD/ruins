package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSpawnFieldItem(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// フィールドアイテムを生成
	item, err := SpawnFieldItem(world, "回復薬", gc.Row(5), gc.Col(10))
	require.NoError(t, err, "SpawnFieldItem should not return error")
	require.NotNil(t, item, "アイテムエンティティが生成されるべき")

	// Nameコンポーネントの確認
	require.True(t, item.HasComponent(world.Components.Name), "Nameコンポーネントが必要")
	name := world.Components.Name.Get(item).(*gc.Name)
	assert.Equal(t, "回復薬", name.Name, "アイテム名が正しくない")

	// GridElementコンポーネントの確認
	require.True(t, item.HasComponent(world.Components.GridElement), "GridElementコンポーネントが必要")
	gridElement := world.Components.GridElement.Get(item).(*gc.GridElement)
	assert.Equal(t, gc.Row(5), gridElement.Row, "行位置が正しくない")
	assert.Equal(t, gc.Col(10), gridElement.Col, "列位置が正しくない")

	// SpriteRenderコンポーネントの確認
	require.True(t, item.HasComponent(world.Components.SpriteRender), "SpriteRenderコンポーネントが必要")
	sprite := world.Components.SpriteRender.Get(item).(*gc.SpriteRender)
	assert.Equal(t, "field", sprite.Name, "スプライト名が正しくない")
	assert.Equal(t, 18, sprite.SpriteNumber, "スプライト番号が正しくない")
	assert.Equal(t, gc.DepthNumRug, sprite.Depth, "描画深度が正しくない")

	// ItemLocationOnFieldコンポーネントの確認
	assert.True(t, item.HasComponent(world.Components.ItemLocationOnField), "ItemLocationOnFieldコンポーネントが必要")

	// クリーンアップ
	world.Manager.DeleteEntity(item)
}

func TestSpawnMultipleFieldItems(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 複数のフィールドアイテムを生成
	items := []struct {
		itemName string
		row      gc.Row
		col      gc.Col
	}{
		{"回復薬", gc.Row(1), gc.Col(1)},
		{"手榴弾", gc.Row(2), gc.Col(2)},
		{"ルビー原石", gc.Row(3), gc.Col(3)},
	}

	createdItems := make([]ecs.Entity, 0, len(items))

	for _, itemData := range items {
		item, err := SpawnFieldItem(world, itemData.itemName, itemData.row, itemData.col)
		require.NoError(t, err, "SpawnFieldItem should not return error")
		createdItems = append(createdItems, item)

		// 位置の確認
		gridElement := world.Components.GridElement.Get(item).(*gc.GridElement)
		assert.Equal(t, itemData.row, gridElement.Row, "行位置が正しくない")
		assert.Equal(t, itemData.col, gridElement.Col, "列位置が正しくない")
	}

	// フィールド上のアイテム数を確認
	fieldItemCount := 0
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationOnField,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		fieldItemCount++
	}))

	assert.Equal(t, len(items), fieldItemCount, "フィールド上のアイテム数が正しくない")

	// クリーンアップ
	for _, item := range createdItems {
		world.Manager.DeleteEntity(item)
	}
}
