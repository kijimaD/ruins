// Package mapbuilder はマップ生成機能を提供する
// 参考: https://bfnightly.bracketproductions.com
package mapbuilder

import (
	"log"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
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
	// RandomSource はシード値による再現可能なランダム生成を提供する
	RandomSource *RandomSource
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
	cx := gc.Pixel(int(tx)*int(consts.TileSize) + int(consts.TileSize)/2)
	cy := gc.Pixel(int(ty)*int(consts.TileSize) + int(consts.TileSize)/2)

	world.Manager.Join(
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := world.Components.Position.Get(entity).(*gc.Position)
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
	targetIdx := int(idx) + int(bm.Level.TileWidth)
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

// AdjacentAnyFloor は直交・斜めを含む近傍8タイルに床があるか判定する
func (bm BuilderMap) AdjacentAnyFloor(idx resources.TileIdx) bool {
	x, y := bm.Level.XYTileCoord(idx)
	width := int(bm.Level.TileWidth)
	height := int(bm.Level.TileHeight)

	// 8方向の隣接タイル座標をチェック
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1}, // 上段
		{0, -1}, {0, 1}, // 中段（中心を除く）
		{1, -1}, {1, 0}, {1, 1}, // 下段
	}

	for _, dir := range directions {
		nx, ny := int(x)+dir[0], int(y)+dir[1]

		// 境界チェック
		if nx < 0 || nx >= width || ny < 0 || ny >= height {
			continue
		}

		neighborIdx := bm.Level.XYTileIndex(gc.Row(nx), gc.Col(ny))
		tile := bm.Tiles[neighborIdx]

		if tile == TileFloor || tile == TileWarpNext || tile == TileWarpEscape {
			return true
		}
	}

	return false
}

// GetWallType は近傍パターンから適切な壁タイプを判定する
func (bm BuilderMap) GetWallType(idx resources.TileIdx) WallType {
	// 4方向の隣接タイルの床状況をチェック
	upFloor := bm.isFloorOrWarp(bm.UpTile(idx))
	downFloor := bm.isFloorOrWarp(bm.DownTile(idx))
	leftFloor := bm.isFloorOrWarp(bm.LeftTile(idx))
	rightFloor := bm.isFloorOrWarp(bm.RightTile(idx))

	// 単純なケース：一方向のみに床がある場合
	if singleWallType := bm.checkSingleDirectionWalls(upFloor, downFloor, leftFloor, rightFloor); singleWallType != WallTypeGeneric {
		return singleWallType
	}

	// 角のケース：2方向に床がある場合
	if cornerWallType := bm.checkCornerWalls(upFloor, downFloor, leftFloor, rightFloor); cornerWallType != WallTypeGeneric {
		return cornerWallType
	}

	// 複雑なパターンまたは判定不可の場合
	return WallTypeGeneric
}

// checkSingleDirectionWalls は単一方向に床がある場合の壁タイプを返す
func (bm BuilderMap) checkSingleDirectionWalls(upFloor, downFloor, leftFloor, rightFloor bool) WallType {
	if downFloor && !upFloor && !leftFloor && !rightFloor {
		return WallTypeTop // 下に床がある → 上壁
	}
	if upFloor && !downFloor && !leftFloor && !rightFloor {
		return WallTypeBottom // 上に床がある → 下壁
	}
	if rightFloor && !upFloor && !downFloor && !leftFloor {
		return WallTypeLeft // 右に床がある → 左壁
	}
	if leftFloor && !upFloor && !downFloor && !rightFloor {
		return WallTypeRight // 左に床がある → 右壁
	}
	return WallTypeGeneric
}

// checkCornerWalls は2方向に床がある場合の壁タイプを返す
func (bm BuilderMap) checkCornerWalls(upFloor, downFloor, leftFloor, rightFloor bool) WallType {
	if downFloor && rightFloor && !upFloor && !leftFloor {
		return WallTypeTopLeft // 下右に床 → 左上角
	}
	if downFloor && leftFloor && !upFloor && !rightFloor {
		return WallTypeTopRight // 下左に床 → 右上角
	}
	if upFloor && rightFloor && !downFloor && !leftFloor {
		return WallTypeBottomLeft // 上右に床 → 左下角
	}
	if upFloor && leftFloor && !downFloor && !rightFloor {
		return WallTypeBottomRight // 上左に床 → 右下角
	}
	return WallTypeGeneric
}

// isFloorOrWarp は床またはワープタイルかを判定する
func (bm BuilderMap) isFloorOrWarp(tile Tile) bool {
	return tile == TileFloor || tile == TileWarpNext || tile == TileWarpEscape
}

// BuilderChain は階層データBuilderMapに対して適用する生成ロジックを保持する構造体
type BuilderChain struct {
	Starter   *InitialMapBuilder
	Builders  []MetaMapBuilder
	BuildData BuilderMap
}

// NewBuilderChain はシード値を指定してビルダーチェーンを作成する
// シードが0の場合はランダムなシードを生成する
func NewBuilderChain(width gc.Row, height gc.Col, seed uint64) *BuilderChain {
	tileCount := int(width) * int(height)
	tiles := make([]Tile, tileCount)

	// シードが0の場合はランダムなシードを生成
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}

	return &BuilderChain{
		Starter:  nil,
		Builders: []MetaMapBuilder{},
		BuildData: BuilderMap{
			Level: resources.Level{
				TileWidth:  width,
				TileHeight: height,
				TileSize:   consts.TileSize,
				Entities:   make([]ecs.Entity, tileCount),
			},
			Tiles:        tiles,
			Rooms:        []Rect{},
			Corridors:    [][]resources.TileIdx{},
			RandomSource: NewRandomSource(seed),
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
// タイルへの描画は行わず、構造体フィールドの値を初期化するだけ
type InitialMapBuilder interface {
	BuildInitial(*BuilderMap)
}

// MetaMapBuilder はメタ情報をビルドするインターフェース
type MetaMapBuilder interface {
	BuildMeta(*BuilderMap)
}

// NewSmallRoomBuilder はシンプルな小部屋ビルダーを作成する
func NewSmallRoomBuilder(width gc.Row, height gc.Col, seed uint64) *BuilderChain {
	chain := NewBuilderChain(width, height, seed)
	chain.StartWith(RectRoomBuilder{})
	chain.With(NewFillAll(TileWall))      // 全体を壁で埋める
	chain.With(RoomDraw{})                // 部屋を描画
	chain.With(LineCorridorBuilder{})     // 廊下を作成
	chain.With(NewBoundaryWall(TileWall)) // 最外周を壁で囲む

	return chain
}
