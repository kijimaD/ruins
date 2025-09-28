package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// RectRoomPlanner は長方形の部屋を作成する
type RectRoomPlanner struct{}

// PlanInitial は初期プランを行う
func (b RectRoomPlanner) PlanInitial(planData *MetaPlan) error {
	b.PlanRooms(planData)
	return nil
}

// PlanRooms は部屋をプランする
func (b RectRoomPlanner) PlanRooms(planData *MetaPlan) {
	maxRooms := 4 + planData.RandomSource.Intn(10)
	rooms := []gc.Rect{}
	for i := 0; i < maxRooms; i++ {
		x := planData.RandomSource.Intn(int(planData.Level.TileWidth))
		y := planData.RandomSource.Intn(int(planData.Level.TileHeight))
		w := 2 + planData.RandomSource.Intn(8)
		h := 2 + planData.RandomSource.Intn(8)
		newRoom := gc.Rect{
			X1: gc.Tile(x),
			X2: gc.Tile(min(x+w, int(planData.Level.TileWidth))),
			Y1: gc.Tile(y),
			Y2: gc.Tile(min(y+h, int(planData.Level.TileHeight))),
		}
		rooms = append(rooms, newRoom)
	}

	planData.Rooms = rooms
}
