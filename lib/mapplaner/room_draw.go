package mapplanner

import gc "github.com/kijimaD/ruins/lib/components"

// RoomDraw は部屋を描画するビルダー
type RoomDraw struct{}

// BuildMeta はメタデータをビルドする
func (b RoomDraw) BuildMeta(buildData *MetaPlan) {
	b.build(buildData)
}

func (b RoomDraw) build(buildData *MetaPlan) {
	for _, room := range buildData.Rooms {
		b.rectangle(buildData, room)
	}
}

func (b RoomDraw) rectangle(buildData *MetaPlan, room gc.Rect) {
	for x := room.X1; x <= room.X2; x++ {
		for y := room.Y1; y <= room.Y2; y++ {
			idx := buildData.Level.XYTileIndex(x, y)
			if 0 < int(idx) && int(idx) < int(buildData.Level.TileWidth)*int(buildData.Level.TileHeight)-1 {
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}
