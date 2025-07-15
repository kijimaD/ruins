package worldhelper

import (
	"testing"

	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestInitDebugData(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game

	// 初期状態では味方メンバーは0人
	memberCount := 0
	world.Manager.Join(
		gameComponents.FactionAlly,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 0, memberCount, "初期状態では味方メンバーは0人であるべき")

	// デバッグデータ初期化実行
	InitDebugData(world)

	// 初期化後は味方メンバーが3人いるはず
	memberCount = 0
	world.Manager.Join(
		gameComponents.FactionAlly,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 3, memberCount, "デバッグ初期化後は味方メンバーが3人いるべき")

	// 2回目の実行では何も追加されないことを確認
	InitDebugData(world)
	memberCount = 0
	world.Manager.Join(
		gameComponents.FactionAlly,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))
	assert.Equal(t, 3, memberCount, "2回目の実行では味方メンバー数は変わらないべき")

	// アイテムが生成されていることを確認
	ironAmount := GetAmount("鉄", world)
	assert.Greater(t, ironAmount, 0, "鉄のアイテムが生成されているべき")
}
