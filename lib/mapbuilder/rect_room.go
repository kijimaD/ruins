package mapbuilder

import (
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
)

// RectRoomBuilder は長方形の部屋を作成する
type RectRoomBuilder struct{}

// BuildInitial は初期ビルドを行う
func (b RectRoomBuilder) BuildInitial(buildData *BuilderMap) {
	b.BuildRooms(buildData)
}

// BuildRooms は部屋をビルドする
func (b RectRoomBuilder) BuildRooms(buildData *BuilderMap) {
	maxRooms := 4 + rand.Intn(10)
	rooms := []Rect{}
	for i := 0; i < maxRooms; i++ {
		x := rand.Intn(int(buildData.Level.TileWidth))
		y := rand.Intn(int(buildData.Level.TileHeight))
		w := 2 + rand.Intn(8)
		h := 2 + rand.Intn(8)
		newRoom := Rect{
			X1: gc.Row(x),
			X2: gc.Row(mathutil.Min(x+w, int(buildData.Level.TileWidth))),
			Y1: gc.Col(y),
			Y2: gc.Col(mathutil.Min(y+h, int(buildData.Level.TileHeight))),
		}
		rooms = append(rooms, newRoom)
	}

	buildData.Rooms = rooms
}
