package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestMoveSystemEnemyCollisionFix(t *testing.T) {
	t.Parallel()
	// 移動システムで敵との衝突が阻害されないことを確認するテスト
	world := createTestWorldForCollision(t)

	// テスト用のスプライトシートをResourcesに追加
	if world.Resources.SpriteSheets == nil {
		sheets := make(map[string]gc.SpriteSheet)
		world.Resources.SpriteSheets = &sheets
	}
	(*world.Resources.SpriteSheets)["test"] = gc.SpriteSheet{
		Name: "test",
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32},
		},
	}

	// Levelを初期化してAtEntity呼び出しを成功させる
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.Level = resources.Level{
		TileWidth:  100,
		TileHeight: 100,
		TileSize:   32,
		Entities:   make([]ecs.Entity, 10000), // 100x100のタイル配列
	}

	// プレイヤーを作成
	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 100, Y: 100},
		Operator:    &gc.Operator{},
		FactionType: &gc.FactionAlly,
		Velocity: &gc.Velocity{
			Speed:        2.0,
			Angle:        90, // 右方向に移動（90度）
			ThrottleMode: gc.ThrottleModeFront,
			MaxSpeed:     2.0, // 最高速度を設定
		},
		SpriteRender: &gc.SpriteRender{
			Name:         "test",
			SpriteNumber: 0,
		},
	})
	playerEntities := entities.AddEntities(world, cl)
	playerEntity := playerEntities[0]

	// 敵を作成（プレイヤーの移動先に配置、BlockPassを持つ）
	cl2 := entities.ComponentList{}
	cl2.Game = append(cl2.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 130, Y: 100}, // プレイヤーの移動先
		FactionType: &gc.FactionEnemy,
		BlockPass:   &gc.BlockPass{}, // 通行阻止コンポーネント
		SpriteRender: &gc.SpriteRender{
			Name:         "test",
			SpriteNumber: 0,
		},
	})
	enemyEntities := entities.AddEntities(world, cl2)
	enemyEntity := enemyEntities[0]

	// 初期位置を記録
	initialPlayerPos := world.Components.Position.Get(playerEntity).(*gc.Position)
	initialX := initialPlayerPos.X
	t.Logf("移動前のプレイヤー位置: (%v, %v)", initialX, initialPlayerPos.Y)

	// MoveSystemを実行（プレイヤーが右方向に移動）
	MoveSystem(world)

	// プレイヤーが移動できていることを確認（敵との衝突で阻害されていない）
	finalPlayerPos := world.Components.Position.Get(playerEntity).(*gc.Position)
	t.Logf("移動後のプレイヤー位置: (%v, %v)", finalPlayerPos.X, finalPlayerPos.Y)
	assert.Greater(t, finalPlayerPos.X, initialX, "プレイヤーは敵がBlockPassを持っていても移動できるべき")

	// 敵は元の位置に留まっていることを確認
	enemyPos := world.Components.Position.Get(enemyEntity).(*gc.Position)
	assert.Equal(t, gc.Pixel(130), enemyPos.X)
}

func TestMoveSystemWallCollisionStillWorks(t *testing.T) {
	t.Parallel()
	// 壁との衝突阻害が正常に動作することを確認するテスト
	world := createTestWorldForCollision(t)

	// テスト用のスプライトシートをResourcesに追加
	if world.Resources.SpriteSheets == nil {
		sheets := make(map[string]gc.SpriteSheet)
		world.Resources.SpriteSheets = &sheets
	}
	(*world.Resources.SpriteSheets)["test"] = gc.SpriteSheet{
		Name: "test",
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32},
		},
	}

	// Levelを初期化してAtEntity呼び出しを成功させる
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.Level = resources.Level{
		TileWidth:  100,
		TileHeight: 100,
		TileSize:   32,
		Entities:   make([]ecs.Entity, 10000), // 100x100のタイル配列
	}

	// プレイヤーを作成
	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: 100, Y: 100},
		Operator:    &gc.Operator{},
		FactionType: &gc.FactionAlly,
		Velocity: &gc.Velocity{
			Speed:        2.0,
			Angle:        90, // 右方向に移動（90度）
			ThrottleMode: gc.ThrottleModeFront,
			MaxSpeed:     2.0, // 最高速度を設定
		},
		SpriteRender: &gc.SpriteRender{
			Name:         "test",
			SpriteNumber: 0,
		},
	})
	playerEntities := entities.AddEntities(world, cl)
	playerEntity := playerEntities[0]

	// 壁を作成（プレイヤーの移動先に配置、BlockPassを持ち、敵ではない）
	cl2 := entities.ComponentList{}
	cl2.Game = append(cl2.Game, gc.GameComponentList{
		Position:  &gc.Position{X: 102, Y: 100}, // プレイヤーのすぐ近く
		BlockPass: &gc.BlockPass{},              // 通行阻止コンポーネント
		// FactionEnemyコンポーネントなし = 壁
		SpriteRender: &gc.SpriteRender{
			Name:         "test",
			SpriteNumber: 0,
		},
	})
	wallEntities := entities.AddEntities(world, cl2)
	wallEntity := wallEntities[0]

	// 初期位置を記録
	initialPlayerPos := world.Components.Position.Get(playerEntity).(*gc.Position)
	initialX := initialPlayerPos.X

	// 壁エンティティのコンポーネント確認のためのデバッグ
	t.Logf("壁エンティティにBlockPassがあるか: %v", wallEntity.HasComponent(world.Components.BlockPass))
	t.Logf("壁エンティティにFactionEnemyがあるか: %v", wallEntity.HasComponent(world.Components.FactionEnemy))

	// MoveSystemを実行（プレイヤーが右方向に移動を試みる）
	MoveSystem(world)

	// プレイヤーが移動できていないことを確認（壁との衝突で阻害される）
	finalPlayerPos := world.Components.Position.Get(playerEntity).(*gc.Position)
	t.Logf("移動後のプレイヤー位置: (%v, %v)", finalPlayerPos.X, finalPlayerPos.Y)
	assert.Equal(t, initialX, finalPlayerPos.X, "プレイヤーは壁との衝突で移動を阻害されるべき")
}
