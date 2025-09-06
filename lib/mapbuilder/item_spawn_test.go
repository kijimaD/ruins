package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpawnFieldItemsIntegration(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// テスト用のマップを生成（小さなマップで簡単にテストできるように）
	chain := NewSmallRoomBuilder(gc.Row(20), gc.Col(20), 12345)
	chain.Build()

	// フィールドアイテムを配置（エラーが発生しないことを確認）
	assert.NoError(t, spawnFieldItems(world, chain))
}
