package mapbuilder

// フィールドのタイル
type Tile uint8

const (
	TileEmpty Tile = iota
	TileFloor
	TileWall
	TileWarpNext
	TileWarpEscape
)
