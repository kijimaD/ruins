package loader

import (
	"math/rand"

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
	{
		for {
			x := rand.Intn(int(level.TileWidth))
			y := rand.Intn(int(level.TileHeight))
			tileIdx := level.XYTileIndex(x, y)
			if level.Tiles[tileIdx] == TileFloor {
				SpawnPlayer(world, x*defaultTileSize+defaultTileSize/2, y*defaultTileSize+defaultTileSize/2)
				break
			}
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

// フィールド上に表示されるプレイヤーを生成する
func SpawnPlayer(world w.World, x int, y int) {
	fieldSpriteSheet := (*world.Resources.SpriteSheets)["field"]
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Player:   &gc.Player{},
			SpriteRender: &ec.SpriteRender{
				SpriteSheet:  &fieldSpriteSheet,
				SpriteNumber: 3,
				Depth:        ec.DepthNumTaller,
			},
		})
		loader.AddEntities(world, componentList)
	}
	// カメラ
	{
		componentList := loader.EntityComponentList{}
		componentList.Game = append(componentList.Game, GameComponentList{
			Position: &gc.Position{X: x, Y: y},
			Camera:   &gc.Camera{Scale: 0.1, ScaleTo: 1},
		})
		loader.AddEntities(world, componentList)
	}
}
