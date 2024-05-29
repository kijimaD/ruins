package mapbuilder

import (
	"math/rand"
)

// 長方形の部屋を作成する
type RectRoomBuilder struct{}

func (b RectRoomBuilder) BuildInitial(buildData *BuilderMap) {
	b.BuildRooms(buildData)
}

func (b RectRoomBuilder) BuildRooms(buildData *BuilderMap) {
	const maxRooms = 8
	rooms := []Rect{}
	for i := 0; i < maxRooms; i++ {
		x := 0 + rand.Intn(16)
		y := 0 + rand.Intn(16)
		w := 2 + rand.Intn(2)
		h := 2 + rand.Intn(2)
		newRoom := Rect{x, y, x + w, y + h}
		rooms = append(rooms, newRoom)
	}

	buildData.Rooms = rooms
}
