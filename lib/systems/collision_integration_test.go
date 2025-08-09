package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestCollisionSystemBattleTransition(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと敵を近い位置に作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0) // 接触する距離

	// 初期状態: イベントなし
	gameResources := world.Resources.Game.(*resources.Game)
	assert.Equal(t, resources.StateEventNone, gameResources.StateEvent)

	// CollisionSystemを実行
	CollisionSystem(world)

	// 戦闘開始イベントが設定されることを確認
	assert.Equal(t, resources.StateEventBattleStart, gameResources.StateEvent)
}

func TestCollisionSystemNoEventWhenAlreadySet(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// プレイヤーと敵を近い位置に作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0)

	// 既に別のイベントが設定されている場合
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext

	// CollisionSystemを実行
	CollisionSystem(world)

	// イベントが変更されないことを確認
	assert.Equal(t, resources.StateEventWarpNext, gameResources.StateEvent)
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
	gameResources := world.Resources.Game.(*resources.Game)
	assert.Equal(t, resources.StateEventNone, gameResources.StateEvent)
}

func TestBattleEncounterApply(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// テスト用エンティティを作成
	createPlayerEntity(t, world, 100.0, 100.0)
	createEnemyEntity(t, world, 110.0, 110.0)

	// エンティティを取得
	var playerEntity, enemyEntity ecs.Entity
	world.Manager.Join(world.Components.Position, world.Components.Operator).Visit(ecs.Visit(func(e ecs.Entity) {
		playerEntity = e
	}))
	world.Manager.Join(world.Components.Position, world.Components.FactionEnemy).Visit(ecs.Visit(func(e ecs.Entity) {
		enemyEntity = e
	}))

	// BattleEncounterエフェクトを直接テスト
	battleEffect := &effects.BattleEncounter{
		FieldEnemyEntity: enemyEntity, // フィールド上の敵シンボル
	}

	// Scopeにプレイヤーエンティティを設定してApplyを実行
	scope := &effects.Scope{Creator: &playerEntity}
	err := battleEffect.Apply(world, scope)
	require.NoError(t, err)

	// 戦闘開始イベントが設定されることを確認
	gameResources := world.Resources.Game.(*resources.Game)
	assert.Equal(t, resources.StateEventBattleStart, gameResources.StateEvent)
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
	gameResources := world.Resources.Game.(*resources.Game)
	assert.Equal(t, resources.StateEventBattleStart, gameResources.StateEvent)
}

func TestCollisionSystemWithVelocityStop(t *testing.T) {
	t.Parallel()
	world := createTestWorldWithResources(t)

	// 速度コンポーネント付きでエンティティを作成
	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 100, Y: 100},
		Operator:    &gc.Operator{},
		FactionType: &gc.FactionAlly,
		Velocity: &gc.Velocity{
			Speed:        2.0,
			ThrottleMode: gc.ThrottleModeFront,
		},
	})
	playerEntities := entities.AddEntities(world, cl)
	playerEntity := playerEntities[0]

	cl2 := entities.ComponentList{}
	cl2.Game = append(cl2.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 110, Y: 110},
		FactionType: &gc.FactionEnemy,
		Velocity: &gc.Velocity{
			Speed:        1.0,
			ThrottleMode: gc.ThrottleModeFront,
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
	gameResources := world.Resources.Game.(*resources.Game)
	assert.Equal(t, resources.StateEventBattleStart, gameResources.StateEvent)
}
