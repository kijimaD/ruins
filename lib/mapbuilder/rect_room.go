package mapbuilder

import (
	"math/rand"

	"github.com/kijimaD/ruins/lib/utils/mathutil"
)

// 長方形の部屋を作成する
type RectRoomBuilder struct{}

func (b RectRoomBuilder) BuildInitial(buildData *BuilderMap) {
	b.BuildRooms(buildData)
}

func (b RectRoomBuilder) BuildRooms(buildData *BuilderMap) {
	maxRooms := 4 + rand.Intn(10)
	rooms := []Rect{}
	for i := 0; i < maxRooms; i++ {
		x := rand.Intn(int(buildData.Level.TileWidth))
		y := rand.Intn(int(buildData.Level.TileHeight))
		w := 2 + rand.Intn(8)
		h := 2 + rand.Intn(8)
		newRoom := Rect{
			X1: x,
			X2: mathutil.Min(x+w, int(buildData.Level.TileWidth)),
			Y1: y,
			Y2: mathutil.Min(y+h, int(buildData.Level.TileHeight)),
		}
		rooms = append(rooms, newRoom)
	}

	buildData.Rooms = rooms
}
