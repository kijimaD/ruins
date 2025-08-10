package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/stretchr/testify/assert"
)

func TestCollisionSystemBattleTransition(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと敵を近い位置に作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0) // 接触する距離

	// 初期状態: イベントなし
	gameResources := world.Resources.Game.(*resources.Dungeon)
	assert.Equal(t, resources.StateEventNone, gameResources.GetStateEvent())

	// CollisionSystemを実行
	CollisionSystem(world)

	// 戦闘開始イベントが設定されることを確認
	assert.Equal(t, resources.StateEventBattleStart, gameResources.GetStateEvent())
}

func TestCollisionSystemNoEventWhenAlreadySet(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと敵を近い位置に作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0)

	// 既に別のイベントが設定されている場合
	gameResources := world.Resources.Game.(*resources.Dungeon)
	gameResources.SetStateEvent(resources.StateEventWarpNext)

	// CollisionSystemを実行
	CollisionSystem(world)

	// イベントが変更されないことを確認
	assert.Equal(t, resources.StateEventWarpNext, gameResources.GetStateEvent())
}

func TestCollisionSystemNoCollision(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと敵を離れた位置に作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 200.0, 200.0) // 接触しない距離

	// CollisionSystemを実行
	CollisionSystem(world)

	// イベントが設定されないことを確認
	gameResources := world.Resources.Game.(*resources.Dungeon)
	assert.Equal(t, resources.StateEventNone, gameResources.GetStateEvent())
}

func TestCollisionSystemMultipleEnemies(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと複数の敵を作成（1つは接触、1つは非接触）
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0) // 接触する敵
	createEnemyEntity(t, world, 200.0, 200.0) // 接触しない敵

	// CollisionSystemを実行
	CollisionSystem(world)

	// 戦闘開始イベントが設定されることを確認
	gameResources := world.Resources.Game.(*resources.Dungeon)
	assert.Equal(t, resources.StateEventBattleStart, gameResources.GetStateEvent())
}

func TestCollisionSystemWithVelocityStop(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// 速度コンポーネント付きでエンティティを作成
	spriteSheet := &gc.SpriteSheet{
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32},
		},
	}

	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 100, Y: 100},
		Operator:    &gc.Operator{},
		FactionType: &gc.FactionAlly,
		Velocity: &gc.Velocity{
			Speed:        2.0,
			ThrottleMode: gc.ThrottleModeFront,
		},
		SpriteRender: &gc.SpriteRender{
			SpriteNumber: 0,
			SpriteSheet:  spriteSheet,
		},
	})
	playerEntities := entities.AddEntities(world, cl)
	playerEntity := playerEntities[0]

	cl2 := entities.ComponentList{}
	cl2.Game = append(cl2.Game, gc.GameComponentList{
		Position:  &gc.Position{X: 110, Y: 110},
		AIMoveFSM: &gc.AIMoveFSM{}, // AI制御された敵として識別
		Velocity: &gc.Velocity{
			Speed:        1.0,
			ThrottleMode: gc.ThrottleModeFront,
		},
		SpriteRender: &gc.SpriteRender{
			SpriteNumber: 0,
			SpriteSheet:  spriteSheet,
		},
	})
	enemyEntities := entities.AddEntities(world, cl2)
	enemyEntity := enemyEntities[0]

	// CollisionSystemを実行
	CollisionSystem(world)

	// 速度が停止されることを確認
	playerVelocity := world.Components.Velocity.Get(playerEntity).(*gc.Velocity)
	enemyVelocity := world.Components.Velocity.Get(enemyEntity).(*gc.Velocity)

	assert.Equal(t, gc.ThrottleModeNope, playerVelocity.ThrottleMode)
	assert.Equal(t, 0.0, playerVelocity.Speed)
	assert.Equal(t, gc.ThrottleModeNope, enemyVelocity.ThrottleMode)
	assert.Equal(t, 0.0, enemyVelocity.Speed)

	// 戦闘開始イベントが設定されることを確認
	gameResources := world.Resources.Game.(*resources.Dungeon)
	assert.Equal(t, resources.StateEventBattleStart, gameResources.GetStateEvent())
}
