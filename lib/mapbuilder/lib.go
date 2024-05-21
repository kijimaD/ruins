// 参考: https://bfnightly.bracketproductions.com
package mapbuilder

import (
	"log"

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

type BuilderChain struct {
	Starter   *InitialMapBuilder
	Builders  []MetaMapBuilder
	BuildData BuilderMap
}

func NewBuilderChain() *BuilderChain {
	tileCount := int(20) * int(20)
	tiles := make([]Tile, tileCount)
	for i, _ := range tiles {
		tiles[i] = TileWall
	}

	return &BuilderChain{
		Starter:  nil,
		Builders: []MetaMapBuilder{},
		BuildData: BuilderMap{
			Level: loader.Level{
				// 仮の値
				TileWidth:  20,
				TileHeight: 20,
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

func SimpleRoomBuilder() *BuilderChain {
	chain := NewBuilderChain()
	chain.StartWith(SimpleMapBuilder{})
	chain.With(RoomDraw{})

	return chain
}
