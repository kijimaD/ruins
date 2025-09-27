package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// SingleRoomPlanner は1つの大きな部屋を作成する（テスト用）
type SingleRoomPlanner struct{}

// BuildInitial は初期ビルドを行う
func (b SingleRoomPlanner) BuildInitial(buildData *MetaPlan) {
	// マップの中央に大きな部屋を1つ作成
	width := buildData.Level.TileWidth
	height := buildData.Level.TileHeight

	// 境界から2タイル内側に部屋を作成
	room := gc.Rect{
		X1: gc.Tile(2),
		Y1: gc.Tile(2),
		X2: width - 2,
		Y2: height - 2,
	}

	buildData.Rooms = []gc.Rect{room}
}

// SingleRoomDraw は1部屋を描画する
type SingleRoomDraw struct{}

// BuildMeta は1部屋を描画する
func (d SingleRoomDraw) BuildMeta(buildData *MetaPlan) {
	for _, room := range buildData.Rooms {
		// 部屋の内部を床タイルで埋める
		for y := room.Y1 + 1; y < room.Y2; y++ {
			for x := room.X1 + 1; x < room.X2; x++ {
				idx := buildData.Level.XYTileIndex(x, y)
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}
