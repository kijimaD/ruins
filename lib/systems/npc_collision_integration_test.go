package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/worldhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollisionWithSpawnedNPC(t *testing.T) {
	t.Parallel()
	// 実際にSpawnNPCで生成された敵とプレイヤーの衝突をテスト
	world, err := game.InitWorld(960, 720)
	require.NoError(t, err)

	// Gameリソースを初期化
	gameResource := &resources.Dungeon{}
	gameResource.SetStateEvent(resources.StateEventNone)
	world.Resources.Dungeon = gameResource

	// SpriteSheetsを初期化
	spriteSheets := make(map[string]gc.SpriteSheet)
	spriteSheets["field"] = gc.SpriteSheet{
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32}, // インデックス0
			{Width: 32, Height: 32}, // インデックス1
			{Width: 32, Height: 32}, // インデックス2
			{Width: 32, Height: 32}, // インデックス3
			{Width: 32, Height: 32}, // インデックス4
			{Width: 32, Height: 32}, // インデックス5
			{Width: 32, Height: 32}, // インデックス6 (NPCのスプライト)
		},
	}
	world.Resources.SpriteSheets = &spriteSheets

	// プレイヤーを生成
	CreatePlayerEntity(t, world, 100.0, 100.0)

	// NPCを近い位置に生成（敵として動作）
	require.NoError(t, worldhelper.SpawnNPC(world, 110, 110))

	// 初期状態: イベントなし
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	assert.Equal(t, resources.StateEventNone, gameResources.GetStateEvent())

	// CollisionSystemを実行
	CollisionSystem(world)

	// 戦闘開始イベントが設定されることを確認
	assert.Equal(t, resources.StateEventBattleStart, gameResources.GetStateEvent())
}
