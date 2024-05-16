package resources

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
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
	// 横グリッド数
	TileWidth gc.Row
	// 縦グリッド数
	TileHeight gc.Col
}

type Tile = gloader.Tile

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYIndex(x int, y int) int {
	return y*int(l.TileWidth) + x
}

// xy座標をタイル座標に変換する
func (l *Level) XYToTileXY(x int, y int) (int, int) {
	tx := x / ec.DungeonTileSize
	ty := y / ec.DungeonTileSize
	return tx, ty
}

func NewLevel(world w.World, newDepth int, width gc.Row, height gc.Col) Level {
	level := Level{
		Depth:      newDepth,
		TileWidth:  width,
		TileHeight: height,
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
