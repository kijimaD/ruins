package resources

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
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
	TileEmpty = gloader.TileEmpty
	TileWall  = gloader.TileWall
)

// 現在の階層
type Level struct {
	// タイル群
	Tiles []Tile
	// 階数
	Depth int
	// 横グリッド数
	Width int
	// 縦グリッド数
	Height int
}

type Tile = gloader.Tile

func (l *Level) XYIndex(x int, y int) int {
	return y*l.Width + x
}

func NewLevel(newDepth int, width int, height int) Level {
	tileCount := width * height
	level := Level{
		Tiles:  make([]Tile, tileCount),
		Depth:  newDepth,
		Width:  width,
		Height: height,
	}
	for i, _ := range level.Tiles {
		level.Tiles[i] = TileEmpty
	}

	return level
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
