package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// SingleRoomPlanner は1つの大きな部屋を作成する（テスト用）
type SingleRoomPlanner struct{}

// PlanInitial は初期プランを行う
func (b SingleRoomPlanner) PlanInitial(planData *MetaPlan) error {
	// マップの中央に大きな部屋を1つ作成
	width := planData.Level.TileWidth
	height := planData.Level.TileHeight

	// 境界から2タイル内側に部屋を作成
	room := gc.Rect{
		X1: gc.Tile(2),
		Y1: gc.Tile(2),
		X2: width - 2,
		Y2: height - 2,
	}

	planData.Rooms = []gc.Rect{room}
	return nil
}

// SingleRoomDraw は1部屋を描画する
type SingleRoomDraw struct{}

// PlanMeta は1部屋を描画する
func (d SingleRoomDraw) PlanMeta(planData *MetaPlan) {
	for _, room := range planData.Rooms {
		// 部屋の内部を床タイルで埋める
		for y := room.Y1 + 1; y < room.Y2; y++ {
			for x := room.X1 + 1; x < room.X2; x++ {
				idx := planData.Level.XYTileIndex(x, y)
				planData.Tiles[idx] = planData.GenerateTile("Floor")
			}
		}
	}
}
