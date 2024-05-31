// 参考: https://bfnightly.bracketproductions.com
package mapbuilder

import (
	"log"

	"github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/loader"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 地図上のタイルを作る元になる概念の集合体
type BuilderMap struct {
	Level     loader.Level
	Tiles     []Tile
	Rooms     []Rect
	Corridors [][]int
}

// 上にあるタイルを調べる
func (bm BuilderMap) UpTile(idx int) Tile {
	targetIdx := idx - int(bm.Level.TileWidth)
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 下にあるタイルを調べる
func (bm BuilderMap) DownTile(idx int) Tile {
	targetIdx := idx + int(bm.Level.TileHeight)
	if targetIdx > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 右にあるタイルを調べる
func (bm BuilderMap) LeftTile(idx int) Tile {
	targetIdx := idx - 1
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 左にあるタイルを調べる
func (bm BuilderMap) RightTile(idx int) Tile {
	targetIdx := idx + 1
	if targetIdx > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 直交する近傍4タイルに床があるか判定する
func (bm BuilderMap) AdjacentOrthoAnyFloor(idx int) bool {
	return bm.UpTile(idx) == TileFloor ||
		bm.DownTile(idx) == TileFloor ||
		bm.RightTile(idx) == TileFloor ||
		bm.LeftTile(idx) == TileFloor ||
		bm.UpTile(idx) == TileWarpNext ||
		bm.DownTile(idx) == TileWarpNext ||
		bm.RightTile(idx) == TileWarpNext ||
		bm.LeftTile(idx) == TileWarpNext
}

type BuilderChain struct {
	Starter   *InitialMapBuilder
	Builders  []MetaMapBuilder
	BuildData BuilderMap
}

func NewBuilderChain(width components.Row, height components.Col) *BuilderChain {
	tileCount := int(width) * int(height)
	tiles := make([]Tile, tileCount)
	for i, _ := range tiles {
		tiles[i] = TileWall
	}

	return &BuilderChain{
		Starter:  nil,
		Builders: []MetaMapBuilder{},
		BuildData: BuilderMap{
			Level: loader.Level{
				TileWidth:  components.Row(width),
				TileHeight: components.Col(height),
				TileSize:   32,
				Entities:   make([]ecs.Entity, tileCount),
			},
			Tiles:     tiles,
			Rooms:     []Rect{},
			Corridors: [][]int{},
		},
	}
}

func (b *BuilderChain) StartWith(initialMapBuilder InitialMapBuilder) {
	b.Starter = &initialMapBuilder
}

func (b *BuilderChain) With(metaMapBuilder MetaMapBuilder) {
	b.Builders = append(b.Builders, metaMapBuilder)
}

func (b *BuilderChain) Build() {
	if b.Starter == nil {
		log.Fatal("empty starter builder!")
	}
	(*b.Starter).BuildInitial(&b.BuildData)

	for _, meta := range b.Builders {
		meta.BuildMeta(&b.BuildData)
	}
}

type InitialMapBuilder interface {
	BuildInitial(*BuilderMap)
}

type MetaMapBuilder interface {
	BuildMeta(*BuilderMap)
}

func SimpleRoomBuilder(width components.Row, height components.Col) *BuilderChain {
	chain := NewBuilderChain(width, height)
	chain.StartWith(RectRoomBuilder{})
	chain.With(RoomDraw{}) // TODO: 暫定的にここで壁を埋めてるので、先に実行する必要がある
	chain.With(LineCorridorBuilder{})

	return chain
}
