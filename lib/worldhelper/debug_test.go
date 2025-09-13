package worldhelper

import (
	"testing"

	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestInitDebugData(t *testing.T) {
	t.Parallel()
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// 初期状態ではパーティメンバーは0人
	memberCount := 0
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.Player,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 0, memberCount, "初期状態ではパーティメンバーは0人であるべき")

	// デバッグデータ初期化実行
	InitDebugData(world)

	// 初期化後はプレイヤーが1人いるはず
	memberCount = 0
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.Player,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 1, memberCount, "デバッグ初期化後はプレイヤーが1人いるべき")

	// 2回目の実行では何も追加されないことを確認
	InitDebugData(world)
	memberCount = 0
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.Player,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 1, memberCount, "2回目の実行ではプレイヤー数は変わらないべき")

	// アイテムが生成されていることを確認
	ironAmount := GetAmount("鉄", world)
	assert.Greater(t, ironAmount, 0, "鉄のアイテムが生成されているべき")
}
