package mapbuilder

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// RectRoomBuilder は長方形の部屋を作成する
type RectRoomBuilder struct{}

// BuildInitial は初期ビルドを行う
func (b RectRoomBuilder) BuildInitial(buildData *BuilderMap) {
	b.BuildRooms(buildData)
}

// BuildRooms は部屋をビルドする
func (b RectRoomBuilder) BuildRooms(buildData *BuilderMap) {
	maxRooms := 4 + buildData.RandomSource.Intn(10)
	rooms := []gc.Rect{}
	for i := 0; i < maxRooms; i++ {
		x := buildData.RandomSource.Intn(int(buildData.Level.TileWidth))
		y := buildData.RandomSource.Intn(int(buildData.Level.TileHeight))
		w := 2 + buildData.RandomSource.Intn(8)
		h := 2 + buildData.RandomSource.Intn(8)
		newRoom := gc.Rect{
			X1: gc.Tile(x),
			X2: gc.Tile(min(x+w, int(buildData.Level.TileWidth))),
			Y1: gc.Tile(y),
			Y2: gc.Tile(min(y+h, int(buildData.Level.TileHeight))),
		}
		rooms = append(rooms, newRoom)
	}

	buildData.Rooms = rooms
}
