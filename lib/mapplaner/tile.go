package mapplanner

// TileType はタイルの種類を表すenum
type TileType uint8

const (
	// TileTypeEmpty は空のタイル（デフォルト状態）
	TileTypeEmpty TileType = iota
	// TileTypeFloor は床タイル
	TileTypeFloor
	// TileTypeWall は壁タイル
	TileTypeWall //  将来的にはエンティティに移行予定
	// TileTypeWarpNext は次の階に向かうワープタイル（将来的にはエンティティに移行予定）
	TileTypeWarpNext
	// TileTypeWarpEscape は脱出用ワープタイル（将来的にはエンティティに移行予定）
	TileTypeWarpEscape
)

// Tile はマップの基盤構造を表すタイル
//
// ## タイルの概念
// マップは規定数のタイル（TileWidth × TileHeight）で正方形に構成される。
// 各位置には必ず1つのタイルが存在し、タイル自体は不変である。
// タイルの上にエンティティ（壁、アイテム、NPCなど）を配置する。
//
// ## マップ生成での役割
// マップ生成時に使用されるタイルタイプ。
// 生成後はmapspawnerでエンティティに変換される。
type Tile struct {
	Type      TileType // タイルの種類（識別用）
	Walkable  bool     // 移動可能かどうか
	BlocksLOS bool     // 視線を遮るかどうか（将来の拡張用）
}

// タイル定数の定義
var (
	// TileEmpty は空のタイル
	TileEmpty = Tile{Type: TileTypeEmpty, Walkable: false, BlocksLOS: false}
	// TileFloor は床タイル
	TileFloor = Tile{Type: TileTypeFloor, Walkable: true, BlocksLOS: false}
	// TileWall は壁タイル
	TileWall = Tile{Type: TileTypeWall, Walkable: false, BlocksLOS: true}
	// TileWarpNext は次の階に向かうワープタイル（将来的にはエンティティに移行予定）
	TileWarpNext = Tile{Type: TileTypeWarpNext, Walkable: true, BlocksLOS: false}
	// TileWarpEscape は脱出用ワープタイル（将来的にはエンティティに移行予定）
	TileWarpEscape = Tile{Type: TileTypeWarpEscape, Walkable: true, BlocksLOS: false}
)
