package resources

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
)

type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

// Tileはsystemなどでも使う。systemから直接gloaderを扱わせたくないので、ここでエクスポートする
const (
	TilePlayer     = gloader.TilePlayer
	TileWall       = gloader.TileWall
	TileWarpNext   = gloader.TileWarpNext
	TileWarpEscape = gloader.TileWarpEscape
	TileEmpty      = gloader.TileEmpty
)

type Level struct {
	CurrentNum int
}

type Tile = gloader.Tile

type Game struct {
	StateEvent StateEvent
	Level      Level
}

// UpdateGameLayoutはゲームウィンドウサイズを更新する
func UpdateGameLayout(world w.World) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
