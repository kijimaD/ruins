package mapplaner

// Tile はフィールドのタイル
//
// マップ生成時に使用されるタイルタイプです。
// 生成後はlevel.goでエンティティに変換されます。
//
// 通行可否の対応表：
//   - TileFloor, TileWarpNext, TileWarpEscape: 通行可能
//   - TileWall: 通行不可（BlockPassコンポーネント付きエンティティに変換）
type Tile uint8

const (
	// TileEmpty は空のタイル
	TileEmpty Tile = iota
	// TileFloor は床タイル（通行可能）
	TileFloor
	// TileWall は壁タイル（通行不可）
	TileWall
	// TileWarpNext は次の階に向かうワープタイル（通行可能）
	TileWarpNext
	// TileWarpEscape は脱出用ワープタイル（通行可能）
	TileWarpEscape
)
