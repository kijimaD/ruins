package mapbuilder

type RoomDraw struct{}

func (b RoomDraw) BuildMeta(buildData *BuilderMap) {
	b.build(buildData)
}

func (b RoomDraw) build(buildData *BuilderMap) {
	// 全体を埋める
	// TODO: 移動する
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	for _, room := range buildData.Rooms {
		b.rectangle(buildData, room)
	}
}

func (b RoomDraw) rectangle(buildData *BuilderMap, room Rect) {
	for x := room.X1; x <= room.X2; x++ {
		for y := room.Y1; y <= room.Y2; y++ {
			idx := buildData.Level.XYTileIndex(x, y)
			if 0 < int(idx) && int(idx) < int(buildData.Level.TileWidth)*int(buildData.Level.TileHeight)-1 {
				buildData.Tiles[idx] = TileFloor
			}
		}
	}
}
