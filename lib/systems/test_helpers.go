package systems

import (
	"fmt"
	"testing"

	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/require"

	gc "github.com/kijimaD/ruins/lib/components"
)

// spriteSize はテスト用のスプライトサイズ構造体
type spriteSize struct {
	width  int
	height int
}

// CreateTestWorldWithResources はテスト用のワールドを作成し、必要なリソースを初期化する
func CreateTestWorldWithResources(t *testing.T) w.World {
	t.Helper()

	components := &gc.Components{}
	world, err := w.InitWorld(components)
	require.NoError(t, err)

	// Gameリソースを初期化
	gameResource := &resources.Dungeon{}
	gameResource.SetStateEvent(resources.StateEventNone)
	world.Resources.Dungeon = gameResource

	// TurnManagerを初期化
	turnManager := turns.NewTurnManager()
	world.Resources.TurnManager = turnManager

	return world
}

// CreateTestWorldForCollision は衝突判定テスト用のワールドを作成する
func CreateTestWorldForCollision(t *testing.T) w.World {
	t.Helper()

	components := &gc.Components{}
	world, err := w.InitWorld(components)
	require.NoError(t, err)

	// Gameリソースを初期化してpanic回避
	gameResource := &resources.Dungeon{}
	gameResource.SetStateEvent(resources.StateEventNone)
	world.Resources.Dungeon = gameResource

	// TurnManagerを初期化
	turnManager := turns.NewTurnManager()
	world.Resources.TurnManager = turnManager

	return world
}

// CreatePlayerEntity は指定された位置にプレイヤーエンティティを作成する
func CreatePlayerEntity(t *testing.T, world w.World, x, y float64) {
	t.Helper()

	// テスト用のスプライトシートを作成してResourcesに追加
	if world.Resources.SpriteSheets == nil {
		sheets := make(map[string]gc.SpriteSheet)
		world.Resources.SpriteSheets = &sheets
	}
	(*world.Resources.SpriteSheets)["test"] = gc.SpriteSheet{
		Name: "test",
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32}, // 標準サイズ
		},
	}

	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:    &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
		FactionType: &gc.FactionAlly,
		SpriteRender: &gc.SpriteRender{
			SpriteNumber: 0,
			Name:         "test",
		},
	})
	entities.AddEntities(world, cl)
}

// CreateEnemyEntity は指定された位置に敵エンティティを作成する
func CreateEnemyEntity(t *testing.T, world w.World, x, y float64) {
	t.Helper()

	// テスト用のスプライトシートを作成してResourcesに追加
	if world.Resources.SpriteSheets == nil {
		sheets := make(map[string]gc.SpriteSheet)
		world.Resources.SpriteSheets = &sheets
	}
	(*world.Resources.SpriteSheets)["test"] = gc.SpriteSheet{
		Name: "test",
		Sprites: []gc.Sprite{
			{Width: 32, Height: 32}, // 標準サイズ
		},
	}

	cl := entities.ComponentList{}
	cl.Game = append(cl.Game, gc.GameComponentList{
		Position:  &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
		AIMoveFSM: &gc.AIMoveFSM{}, // AI制御された敵として識別
		SpriteRender: &gc.SpriteRender{
			SpriteNumber: 0,
			Name:         "test",
		},
	})
	entities.AddEntities(world, cl)
}

// CreateEntityWithSprite は指定されたスプライトサイズでエンティティを作成する
func CreateEntityWithSprite(t *testing.T, world w.World, x, y float64, width, height int, isPlayer bool) {
	t.Helper()

	// 一意なスプライトシート名を生成
	var sheetName string
	if isPlayer {
		sheetName = fmt.Sprintf("player_%dx%d", width, height)
	} else {
		sheetName = fmt.Sprintf("enemy_%dx%d", width, height)
	}

	// テスト用のスプライトシートを作成してResourcesに追加
	if world.Resources.SpriteSheets == nil {
		sheets := make(map[string]gc.SpriteSheet)
		world.Resources.SpriteSheets = &sheets
	}
	(*world.Resources.SpriteSheets)[sheetName] = gc.SpriteSheet{
		Name: sheetName,
		Sprites: []gc.Sprite{
			{Width: width, Height: height}, // インデックス0のスプライト
		},
	}

	// テスト用のスプライト情報を作成
	spriteRender := &gc.SpriteRender{
		SpriteNumber: 0,
		Name:         sheetName,
	}

	cl := entities.ComponentList{}
	if isPlayer {
		cl.Game = append(cl.Game, gc.GameComponentList{
			Position:     &gc.Position{X: gc.Pixel(x), Y: gc.Pixel(y)},
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

// CreateEntityWithSpriteSize はspriteSize構造体を使ってエンティティを作成する
func CreateEntityWithSpriteSize(t *testing.T, world w.World, x, y float64, size spriteSize, isPlayer bool) {
	t.Helper()
	CreateEntityWithSprite(t, world, x, y, size.width, size.height, isPlayer)
}

// CreateEntityWithGridPosition はグリッド座標でエンティティを作成する
func CreateEntityWithGridPosition(t *testing.T, world w.World, gridX, gridY int, isPlayer bool) {
	t.Helper()

	var cl entities.ComponentList

	if isPlayer {
		cl.Game = append(cl.Game, gc.GameComponentList{
			GridElement: &gc.GridElement{X: gc.Tile(gridX), Y: gc.Tile(gridY)},
			FactionType: &gc.FactionAlly,
		})
	} else {
		cl.Game = append(cl.Game, gc.GameComponentList{
			GridElement: &gc.GridElement{X: gc.Tile(gridX), Y: gc.Tile(gridY)},
			AIMoveFSM:   &gc.AIMoveFSM{},
			FactionType: &gc.FactionEnemy,
		})
	}

	entities.AddEntities(world, cl)
}
