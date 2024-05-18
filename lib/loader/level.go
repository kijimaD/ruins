package loader

import (
	"math/rand"
	"regexp"

	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

const (
	exteriorSpriteNumber   = 0
	wallSpriteNumber       = 1
	floorSpriteNumber      = 2
	playerSpriteNumber     = 3
	warpNextSpriteNumber   = 4
	warpEscapeSpriteNumber = 5
)

// const (
// 	// フロア
// 	charFloor = ' '
// 	// 壁
// 	charWall = '#'
// 	// 操作するプレイヤー
// 	charPlayer = '@'
// 	// 壁より外側の埋める部分
// 	charExterior = '_'
// 	// 次の階層へ
// 	charWarpNext = 'O'
// 	// 脱出
// 	charWarpEscape = 'X'
// )

var regexpValidChars = regexp.MustCompile(`^[ #@+_OX]+$`)

// 現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数。タイルは正方形のため、縦横で同じピクセル数になる
	TileSize int
	// タイル群
	Grid []ecs.Entity
}

func NewLevel(world w.World, newDepth int, width gc.Row, height gc.Col) Level {
	tiles := make([]ecs.Entity, 0, int(width)*int(height))

	tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(0)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(0)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(0)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(0)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(4), gc.Col(0)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(0)))

	tiles = append(tiles, SpawnFieldWall(world, gc.Row(0), gc.Col(1)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(1)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(1)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(1)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(4), gc.Col(1)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(1)))

	tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(2)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(2)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(2)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(2)))
	tiles = append(tiles, SpawnFieldWall(world, gc.Row(4), gc.Col(2)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(2)))

	tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(3)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(3)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(3)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(3)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(4), gc.Col(3)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(3)))

	tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(4)))
	tiles = append(tiles, SpawnFieldWall(world, gc.Row(1), gc.Col(4)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(4)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(4)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(4), gc.Col(4)))
	tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(4)))

	r := rand.Intn(2)
	if r == 0 {
		tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(5)))
		tiles = append(tiles, SpawnFieldWarpNext(world, gc.Row(3), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(4), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(5)))
	} else {
		tiles = append(tiles, SpawnFloor(world, gc.Row(0), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(1), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(2), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(3), gc.Col(5)))
		tiles = append(tiles, SpawnFieldWarpNext(world, gc.Row(4), gc.Col(5)))
		tiles = append(tiles, SpawnFloor(world, gc.Row(5), gc.Col(5)))
	}

	level := Level{
		TileWidth:  width,
		TileHeight: height,
		TileSize:   defaultTileSize,
		Grid:       tiles,
	}

	return level
}

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx int, ty int) int {
	return ty*int(l.TileWidth) + tx
}

// xy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x int, y int) ecs.Entity {
	tx := x / l.TileSize
	ty := y / l.TileSize
	idx := l.XYTileIndex(tx, ty)

	return l.Grid[idx]
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

// フィールドのタイル
type Tile uint8

const (
	TileEmpty Tile = 0
	TileFloor Tile = 1 << iota
	TileWall
	TileWarpNext
)

// ================

// フィールド上に表示される床を生成する
func SpawnFloor(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 2,
			Depth:        ec.DepthNumFloor,
		},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される壁を生成する
func SpawnFieldWall(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 1,
			Depth:        ec.DepthNumTaller,
		},
		BlockView: &gc.BlockView{},
		BlockPass: &gc.BlockPass{},
	})

	return loader.AddEntities(world, componentList)[0]
}

// フィールド上に表示される階段を生成する
func SpawnFieldWarpNext(world w.World, x gc.Row, y gc.Col) ecs.Entity {
	SpawnFloor(world, x, y) // 下敷き描画

	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	componentList := loader.EntityComponentList{}
	componentList.Game = append(componentList.Game, GameComponentList{
		GridElement: &gc.GridElement{Row: x, Col: y},
		SpriteRender: &ec.SpriteRender{
			SpriteSheet:  &fieldSpriteSheet,
			SpriteNumber: 4,
			Depth:        ec.DepthNumRug,
		},
		Warp: &gc.Warp{Mode: gc.WarpModeNext},
	})

	return loader.AddEntities(world, componentList)[0]
}
