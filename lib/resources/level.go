package resources

import (
	"log"
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

type Game struct {
	// フィールド上で発生したイベント。各stateで補足されて処理される
	StateEvent StateEvent
	// 現在階のフィールド情報
	Level loader.Level
	// 階層数
	Depth int
}

const defaultTileSize = 32

func NewLevel(world w.World, newDepth int, width gc.Row, height gc.Col) loader.Level {
	tileCount := int(width) * int(height)
	level := loader.Level{
		TileWidth:  width,
		TileHeight: height,
		TileSize:   defaultTileSize,
		Entities:   make([]ecs.Entity, tileCount),
	}
	chain := mapbuilder.SimpleRoomBuilder()
	chain.Build()

	// ワープホールを生成する
	// FIXME: たまに届かない位置に生成される
	failCountWarpNext := 0
	for {
		if failCountWarpNext > 1000 {
			log.Fatal("ワープホールの生成に失敗した")
		}
		x := rand.Intn(int(level.TileWidth))
		y := rand.Intn(int(level.TileHeight))
		tileIdx := level.XYTileIndex(x, y)
		if chain.BuildData.Tiles[tileIdx] == mapbuilder.TileFloor {
			chain.BuildData.Tiles[tileIdx] = mapbuilder.TileWarpNext

			break
		}
		failCountWarpNext++
	}
	// プレイヤーを配置する
	failCountPlayer := 0
	for {
		if failCountPlayer > 1000 {
			log.Fatal("プレイヤーの生成に失敗した")
		}
		x := rand.Intn(int(level.TileWidth))
		y := rand.Intn(int(level.TileHeight))
		tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
		if chain.BuildData.Tiles[tileIdx] == mapbuilder.TileFloor {
			SpawnPlayer(world, x*defaultTileSize+defaultTileSize/2, y*defaultTileSize+defaultTileSize/2)
			break
		}
		failCountPlayer++
	}

	// tilesを元にエンティティを生成する
	for i, t := range chain.BuildData.Tiles {
		x, y := chain.BuildData.Level.XYTileCoord(i)
		switch t {
		case mapbuilder.TileFloor:
			chain.BuildData.Level.Entities[i] = SpawnFloor(world, gc.Row(x), gc.Col(y))
		case mapbuilder.TileWall:
			chain.BuildData.Level.Entities[i] = SpawnFieldWall(world, gc.Row(x), gc.Col(y))
		case mapbuilder.TileWarpNext:
			chain.BuildData.Level.Entities[i] = SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}

// フィールド上でのイベント
type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpNext   = StateEvent("WARP_NEXT")
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
)

// UpdateGameLayoutはゲームウィンドウサイズを更新する
func UpdateGameLayout(world w.World) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
