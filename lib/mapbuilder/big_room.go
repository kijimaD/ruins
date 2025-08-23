package mapbuilder

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// BigRoomBuilder は大部屋を生成するビルダー
// マップ全体の大部分を1つの部屋にする
type BigRoomBuilder struct{}

// BuildInitial は初期マップをビルドする
func (b BigRoomBuilder) BuildInitial(buildData *BuilderMap) {
	// マップの境界を考慮して大きな部屋を1つ作成
	// 外周に1タイル分の壁を残す
	margin := 2
	room := Rect{
		X1: gc.Row(margin),
		Y1: gc.Col(margin),
		X2: gc.Row(int(buildData.Level.TileWidth) - margin - 1),
		Y2: gc.Col(int(buildData.Level.TileHeight) - margin - 1),
	}

	// 部屋をリストに追加
	buildData.Rooms = append(buildData.Rooms, room)
}

// BigRoomDraw は大部屋を描画するビルダー
// RoomDrawと同じロジックだが、大部屋専用の描画処理を提供
type BigRoomDraw struct{}

// BuildMeta は大部屋をタイルに描画する
func (b BigRoomDraw) BuildMeta(buildData *BuilderMap) {
	for _, room := range buildData.Rooms {
		// 部屋の内部を床タイルで埋める
		for x := room.X1; x <= room.X2; x++ {
			for y := room.Y1; y <= room.Y2; y++ {
				idx := buildData.Level.XYTileIndex(x, y)
				buildData.Tiles[idx] = TileFloor
			}
		}

		// 部屋の境界を壁で囲む
		// 上辺と下辺
		for x := room.X1 - 1; x <= room.X2+1; x++ {
			// 上辺
			if y := room.Y1 - 1; y >= 0 {
				idx := buildData.Level.XYTileIndex(x, y)
				if buildData.Tiles[idx] != TileFloor {
					buildData.Tiles[idx] = TileWall
				}
			}
			// 下辺
			if y := room.Y2 + 1; int(y) < int(buildData.Level.TileHeight) {
				idx := buildData.Level.XYTileIndex(x, y)
				if buildData.Tiles[idx] != TileFloor {
					buildData.Tiles[idx] = TileWall
				}
			}
		}

		// 左辺と右辺
		for y := room.Y1; y <= room.Y2; y++ {
			// 左辺
			if x := room.X1 - 1; x >= 0 {
				idx := buildData.Level.XYTileIndex(x, y)
				if buildData.Tiles[idx] != TileFloor {
					buildData.Tiles[idx] = TileWall
				}
			}
			// 右辺
			if x := room.X2 + 1; int(x) < int(buildData.Level.TileWidth) {
				idx := buildData.Level.XYTileIndex(x, y)
				if buildData.Tiles[idx] != TileFloor {
					buildData.Tiles[idx] = TileWall
				}
			}
		}
	}
}

// BigRoomWithPillarsBuilder は柱付き大部屋を生成するビルダー
// 大部屋の中に規則的に柱を配置する
type BigRoomWithPillarsBuilder struct {
	// 柱の間隔
	PillarSpacing int
}

// BuildInitial は初期マップをビルドする
func (b BigRoomWithPillarsBuilder) BuildInitial(buildData *BuilderMap) {
	// 通常の大部屋を作成
	bigRoom := BigRoomBuilder{}
	bigRoom.BuildInitial(buildData)
}

// BigRoomWithPillarsDraw は柱付き大部屋を描画するビルダー
type BigRoomWithPillarsDraw struct {
	// 柱の間隔（デフォルト: 4）
	PillarSpacing int
}

// BuildMeta は柱付き大部屋をタイルに描画する
func (b BigRoomWithPillarsDraw) BuildMeta(buildData *BuilderMap) {
	// まず大部屋を描画
	bigRoomDraw := BigRoomDraw{}
	bigRoomDraw.BuildMeta(buildData)

	// 柱の間隔を設定（デフォルト値）
	spacing := b.PillarSpacing
	if spacing <= 0 {
		spacing = 4
	}

	// 部屋内に柱を配置
	for _, room := range buildData.Rooms {
		// 柱の開始位置を計算（部屋の中心から対称に配置）
		startX := int(room.X1) + spacing
		startY := int(room.Y1) + spacing

		// 規則的に柱を配置
		for x := startX; x < int(room.X2); x += spacing + 1 {
			for y := startY; y < int(room.Y2); y += spacing + 1 {
				idx := buildData.Level.XYTileIndex(gc.Row(x), gc.Col(y))
				buildData.Tiles[idx] = TileWall
			}
		}
	}
}
