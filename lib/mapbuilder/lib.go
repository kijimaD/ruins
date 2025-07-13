// 参考: https://bfnightly.bracketproductions.com
package mapbuilder

import (
	"log"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 階層のタイルを作る元になる概念の集合体
type BuilderMap struct {
	// 階層情報
	Level resources.Level
	// 階層を構成するタイル群。長さはステージの大きさで決まる
	Tiles []Tile
	// 部屋群。部屋は長方形の移動可能な空間のことをいう。
	// 部屋はタイルの集合体である
	Rooms []Rect
	// 廊下群。廊下は部屋と部屋をつなぐ移動可能な空間のことをいう。
	// 廊下はタイルの集合体である
	Corridors [][]resources.TileIdx
}

// 指定タイル座標がスポーン可能かを返す
// スポーンチェックは地図生成時にしか使わないだろう
func (bm BuilderMap) IsSpawnableTile(world w.World, tx gc.Row, ty gc.Col) bool {
	idx := bm.Level.XYTileIndex(tx, ty)
	tile := bm.Tiles[idx]
	if tile != TileFloor {
		return false
	}

	if bm.existEntityOnTile(world, tx, ty) {
		return false
	}

	return true
}

// 指定タイル座標にエンティティがすでにあるかを返す
// MEMO: 階層生成時スポーンさせるときは、タイルの座標中心にスポーンさせている。Positionを持つエンティティの数ぶんで検証できる
func (bm BuilderMap) existEntityOnTile(world w.World, tx gc.Row, ty gc.Col) bool {
	isExist := false
	cx := components.Pixel(int(tx)*int(utils.TileSize) + int(utils.TileSize)/2)
	cy := components.Pixel(int(ty)*int(utils.TileSize) + int(utils.TileSize)/2)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		if pos.X == cx && pos.Y == cy {
			isExist = true

			return
		}
	}))

	return isExist
}

// 上にあるタイルを調べる
func (bm BuilderMap) UpTile(idx resources.TileIdx) Tile {
	targetIdx := resources.TileIdx(int(idx) - int(bm.Level.TileWidth))
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 下にあるタイルを調べる
func (bm BuilderMap) DownTile(idx resources.TileIdx) Tile {
	targetIdx := int(idx) + int(bm.Level.TileHeight)
	if targetIdx > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 右にあるタイルを調べる
func (bm BuilderMap) LeftTile(idx resources.TileIdx) Tile {
	targetIdx := idx - 1
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 左にあるタイルを調べる
func (bm BuilderMap) RightTile(idx resources.TileIdx) Tile {
	targetIdx := idx + 1
	if int(targetIdx) > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// 直交する近傍4タイルに床があるか判定する
func (bm BuilderMap) AdjacentOrthoAnyFloor(idx resources.TileIdx) bool {
	return bm.UpTile(idx) == TileFloor ||
		bm.DownTile(idx) == TileFloor ||
		bm.RightTile(idx) == TileFloor ||
		bm.LeftTile(idx) == TileFloor ||
		bm.UpTile(idx) == TileWarpNext ||
		bm.DownTile(idx) == TileWarpNext ||
		bm.RightTile(idx) == TileWarpNext ||
		bm.LeftTile(idx) == TileWarpNext
}

// 階層データBuilderMapに対して適用する生成ロジックを保持する構造体
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
			Level: resources.Level{
				TileWidth:  components.Row(width),
				TileHeight: components.Col(height),
				TileSize:   utils.TileSize,
				Entities:   make([]ecs.Entity, tileCount),
			},
			Tiles:     tiles,
			Rooms:     []Rect{},
			Corridors: [][]resources.TileIdx{},
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
