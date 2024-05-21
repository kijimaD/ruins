package mapbuilder

import (
	"math/rand"
)

type SimpleMapBuilder struct{}

func (b SimpleMapBuilder) BuildInitial(buildData *BuilderMap) {
	b.BuildRooms(buildData)
}

func (b SimpleMapBuilder) BuildRooms(buildData *BuilderMap) {
	const maxRooms = 30
	rooms := []Rect{}
	for i := 0; i < maxRooms; i++ {
		x := 1 + rand.Intn(18)
		y := 1 + rand.Intn(18)
		w := 1 + rand.Intn(10)
		h := 1 + rand.Intn(10)
		newRoom := Rect{x, y, w, h}
		rooms = append(rooms, newRoom)
	}

	buildData.Rooms = rooms
}
