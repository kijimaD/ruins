package resources

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/spawner"
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

// フィールド上でのイベント
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
	// 階数
	Depth int
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数(正方形)
	TileSize int
}

type Tile = gloader.Tile

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYIndex(x int, y int) int {
	return y*int(l.TileWidth) + x
}

// ステージ幅。横の全体ピクセル数
func (l *Level) Width() int {
	return int(l.TileWidth) * l.TileSize
}

// ステージ縦。縦の全体ピクセル数
func (l *Level) Height() int {
	return int(l.TileHeight) * l.TileSize
}

const defaultTileSize = 32

func NewLevel(world w.World, newDepth int, width gc.Row, height gc.Col) Level {
	level := Level{
		Depth:      newDepth,
		TileWidth:  width,
		TileHeight: height,
		TileSize:   defaultTileSize,
	}
	for j := 0; j < int(width); j++ {
		for k := 0; k < int(height); k++ {
			spawner.SpawnFloor(world, gc.Row(k), gc.Col(j))
		}
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
