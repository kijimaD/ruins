package resources

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/utils/vutil"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

type Game struct {
	StateEvent StateEvent
	Level      Level
}

// ステート上でのイベント
type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
)

// Tileはsystemなどでも使う。systemから直接gloaderを扱わせたくないので、ここでエクスポートする
const (
	TilePlayer     = gloader.TilePlayer
	TileWall       = gloader.TileWall
	TileWarpNext   = gloader.TileWarpNext
	TileWarpEscape = gloader.TileWarpEscape
	TileEmpty      = gloader.TileEmpty
)

// 現在の階層
type Level struct {
	// タイル群
	Tiles vutil.Vec2d[Tile]
	// 階数
	Depth int
}

type Tile = gloader.Tile

// UpdateGameLayoutはゲームウィンドウサイズを更新する
func UpdateGameLayout(world w.World) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
