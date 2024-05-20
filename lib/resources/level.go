package resources

import (
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
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
	Level Level
	// 階層数
	Depth int
}

// 現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数。タイルは正方形のため、縦横で同じピクセル数になる
	TileSize int
	// タイル群。地図生成の処理を保持するのに使う
	Tiles []Tile
	// タイルエンティティ群
	Entities []ecs.Entity
}

const defaultTileSize = 32

func NewLevel(world w.World, newDepth int, width gc.Row, height gc.Col) Level {
	tileCount := int(width) * int(height)
	level := Level{
		TileWidth:  width,
		TileHeight: height,
		TileSize:   defaultTileSize,
		Tiles:      make([]Tile, tileCount),
		Entities:   make([]ecs.Entity, tileCount),
	}

	for i, _ := range level.Tiles {
		level.Tiles[i] = TileFloor
	}
	// 壁を生成する
	{
		n := rand.Intn(14)
		n += 6 // 6 ~ 20
		for i := 0; i < n; i++ {
			x := rand.Intn(int(level.TileWidth))
			y := rand.Intn(int(level.TileHeight))
			tileIdx := level.XYTileIndex(x, y)
			level.Tiles[tileIdx] = TileWall
		}
	}
	// ワープホールを生成する
	{
		x := rand.Intn(int(level.TileWidth))
		y := rand.Intn(int(level.TileHeight))
		tileIdx := level.XYTileIndex(x, y)
		level.Tiles[tileIdx] = TileWarpNext
	}

	// tilesを元にエンティティを生成する
	for i, t := range level.Tiles {
		x, y := level.XYTileCoord(i)
		switch t {
		case TileFloor:
			level.Entities[i] = SpawnFloor(world, gc.Row(x), gc.Col(y))
		case TileWall:
			level.Entities[i] = SpawnFieldWall(world, gc.Row(x), gc.Col(y))
		case TileWarpNext:
			level.Entities[i] = SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		}
	}

	// プレイヤー配置
	for {
		x := rand.Intn(int(level.TileWidth))
		y := rand.Intn(int(level.TileHeight))
		tileIdx := level.XYTileIndex(x, y)
		if level.Tiles[tileIdx] == TileFloor {
			SpawnPlayer(world, x*defaultTileSize+defaultTileSize/2, y*defaultTileSize+defaultTileSize/2)
			break
		}
	}

	return level
}

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx int, ty int) int {
	return ty*int(l.TileWidth) + tx
}

// タイルスライスのインデックスからタイル座標を求める
func (l *Level) XYTileCoord(idx int) (int, int) {
	x := idx % int(l.TileWidth)
	y := idx / int(l.TileWidth)
	return x, y
}

// xy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x int, y int) ecs.Entity {
	tx := x / l.TileSize
	ty := y / l.TileSize
	idx := l.XYTileIndex(tx, ty)

	return l.Entities[idx]
}

// ステージ幅。横の全体ピクセル数
func (l *Level) Width() int {
	return int(l.TileWidth) * l.TileSize
}

// ステージ縦。縦の全体ピクセル数
func (l *Level) Height() int {
	return int(l.TileHeight) * l.TileSize
}

// フィールドのタイル
type Tile uint8

const (
	TileEmpty Tile = 0
	TileFloor Tile = 1 << iota
	TileWall
	TileWarpNext
)

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
