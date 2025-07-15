// Package mapbuilder はマップ生成機能を提供する
// 参考: https://bfnightly.bracketproductions.com
package mapbuilder

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// BuilderMap は階層のタイルを作る元になる概念の集合体
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

// IsSpawnableTile は指定タイル座標がスポーン可能かを返す
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
	cx := gc.Pixel(int(tx)*int(utils.TileSize) + int(utils.TileSize)/2)
	cy := gc.Pixel(int(ty)*int(utils.TileSize) + int(utils.TileSize)/2)

	gameComponents := world.Components.Game
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

// UpTile は上にあるタイルを調べる
func (bm BuilderMap) UpTile(idx resources.TileIdx) Tile {
	targetIdx := resources.TileIdx(int(idx) - int(bm.Level.TileWidth))
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// DownTile は下にあるタイルを調べる
func (bm BuilderMap) DownTile(idx resources.TileIdx) Tile {
	targetIdx := int(idx) + int(bm.Level.TileHeight)
	if targetIdx > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// LeftTile は左にあるタイルを調べる
func (bm BuilderMap) LeftTile(idx resources.TileIdx) Tile {
	targetIdx := idx - 1
	if targetIdx < 0 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// RightTile は右にあるタイルを調べる
func (bm BuilderMap) RightTile(idx resources.TileIdx) Tile {
	targetIdx := idx + 1
	if int(targetIdx) > len(bm.Tiles)-1 {
		return TileEmpty
	}

	return bm.Tiles[targetIdx]
}

// AdjacentOrthoAnyFloor は直交する近傍4タイルに床があるか判定する
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

// BuilderChain は階層データBuilderMapに対して適用する生成ロジックを保持する構造体
type BuilderChain struct {
	Starter   *InitialMapBuilder
	Builders  []MetaMapBuilder
	BuildData BuilderMap
}

// NewBuilderChain は新しいビルダーチェーンを作成する
func NewBuilderChain(width gc.Row, height gc.Col) *BuilderChain {
	tileCount := int(width) * int(height)
	tiles := make([]Tile, tileCount)
	for i := range tiles {
		tiles[i] = TileWall
	}

	return &BuilderChain{
		Starter:  nil,
		Builders: []MetaMapBuilder{},
		BuildData: BuilderMap{
			Level: resources.Level{
				TileWidth:  width,
				TileHeight: height,
				TileSize:   utils.TileSize,
				Entities:   make([]ecs.Entity, tileCount),
			},
			Tiles:     tiles,
			Rooms:     []Rect{},
			Corridors: [][]resources.TileIdx{},
		},
	}
}

// StartWith は初期ビルダーを設定する
func (b *BuilderChain) StartWith(initialMapBuilder InitialMapBuilder) {
	b.Starter = &initialMapBuilder
}

// With はメタビルダーを追加する
func (b *BuilderChain) With(metaMapBuilder MetaMapBuilder) {
	b.Builders = append(b.Builders, metaMapBuilder)
}

// Build はビルダーチェーンを実行してマップを生成する
func (b *BuilderChain) Build() {
	if b.Starter == nil {
		log.Fatal("empty starter builder!")
	}
	(*b.Starter).BuildInitial(&b.BuildData)

	for _, meta := range b.Builders {
		meta.BuildMeta(&b.BuildData)
	}
}

// InitialMapBuilder は初期マップをビルドするインターフェース
type InitialMapBuilder interface {
	BuildInitial(*BuilderMap)
}

// MetaMapBuilder はメタ情報をビルドするインターフェース
type MetaMapBuilder interface {
	BuildMeta(*BuilderMap)
}

// SimpleRoomBuilder はシンプルな部屋ビルダーを作成する
func SimpleRoomBuilder(width gc.Row, height gc.Col) *BuilderChain {
	chain := NewBuilderChain(width, height)
	chain.StartWith(RectRoomBuilder{})
	chain.With(RoomDraw{}) // TODO: 暫定的にここで壁を埋めてるので、先に実行する必要がある
	chain.With(LineCorridorBuilder{})

	return chain
}
