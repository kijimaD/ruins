package mapbuilder

// Tile はフィールドのタイル
type Tile uint8

const (
	// TileEmpty は空のタイル
	TileEmpty Tile = iota
	// TileFloor は床タイル
	TileFloor
	// TileWall は壁タイル
	TileWall
	// TileWarpNext は次の階に向かうワープタイル
	TileWarpNext
	// TileWarpEscape は脱出用ワープタイル
	TileWarpEscape
)
