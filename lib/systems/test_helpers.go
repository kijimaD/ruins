package systems

import (
	"testing"

	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/require"

	gc "github.com/kijimaD/ruins/lib/components"
)

// createTestWorldWithResources はテスト用のワールドを作成し、必要なリソースを初期化する
func createTestWorldWithResources(t *testing.T) w.World {
	t.Helper()

	components := &gc.Components{}
	world, err := w.InitWorld(components)
	require.NoError(t, err)

	// Gameリソースを初期化
	world.Resources.Game = &resources.Game{
		StateEvent: resources.StateEventNone,
	}

	return world
}

// createTestWorldForCollision は衝突判定テスト用のワールドを作成する
func createTestWorldForCollision(t *testing.T) w.World {
	t.Helper()

	components := &gc.Components{}
	world, err := w.InitWorld(components)
	require.NoError(t, err)

	// Gameリソースを初期化してpanic回避
	world.Resources.Game = &resources.Game{
		StateEvent: resources.StateEventNone,
	}

	return world
}

// createPlayerEntity は指定された位置にプレイヤーエンティティを作成する
func createPlayerEntity(t *testing.T, world w.World, x, y float64) {
	t.Helper()

	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
		Operator:    &gc.Operator{},
		FactionType: &gc.FactionAlly,
	})
	entities.AddEntities(world, cl)
}

// createEnemyEntity は指定された位置に敵エンティティを作成する
func createEnemyEntity(t *testing.T, world w.World, x, y float64) {
	t.Helper()

	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
		FactionType: &gc.FactionEnemy,
	})
	entities.AddEntities(world, cl)
}

// createEntityWithSprite は指定されたスプライトサイズでエンティティを作成する
func createEntityWithSprite(t *testing.T, world w.World, x, y float64, width, height int, isPlayer bool) {
	t.Helper()

	// テスト用のスプライト情報を作成
	spriteRender := &gc.SpriteRender{
		SpriteNumber: 0,
		SpriteSheet: &gc.SpriteSheet{
			Sprites: []gc.Sprite{
				{Width: width, Height: height}, // インデックス0のスプライト
			},
		},
	}

	cl := entities.ComponentList{}
	if isPlayer {
		cl.Game = append(cl.Game, gc.GameComponentList{
			Position:     &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
			Operator:     &gc.Operator{},
			FactionType:  &gc.FactionAlly,
			SpriteRender: spriteRender,
		})
	} else {
		cl.Game = append(cl.Game, gc.GameComponentList{
			Position:     &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
			FactionType:  &gc.FactionEnemy,
			SpriteRender: spriteRender,
		})
	}

	entities.AddEntities(world, cl)
}

// createEntityWithSpriteSize はspriteSize構造体を使ってエンティティを作成する
func createEntityWithSpriteSize(t *testing.T, world w.World, x, y float64, size spriteSize, isPlayer bool) {
	t.Helper()
	createEntityWithSprite(t, world, x, y, size.width, size.height, isPlayer)
}
