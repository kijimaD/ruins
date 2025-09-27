package mapplanner

// Tile はマップの基盤構造を表すタイル
//
// ## タイルの概念
// マップは規定数のタイル（TileWidth × TileHeight）で正方形に構成されます。
// 各位置には必ず1つのタイルが存在し、タイル自体は不変です。
// タイルの上にエンティティ（壁、床、アイテム、NPCなど）を配置していきます。
//
// ## マップ生成での役割
// マップ生成時に使用されるタイルタイプです。
// 生成後はmapspawnerでエンティティに変換されます。
//
// ## 通行可否の対応表
//   - TileFloor, TileWarpNext, TileWarpEscape: 通行可能
//   - TileWall: 通行不可（BlockPassコンポーネント付きエンティティに変換）
//   - TileEmpty: 空のタイル（デフォルト、通常は使用されない）
type Tile uint8

const (
	// TileEmpty は空のタイル（デフォルト状態、通常は床タイルに設定される）
	TileEmpty Tile = iota
	// TileFloor は床タイル（通行可能、基本的な歩行可能エリア）
	TileFloor
	// TileWall は壁タイル（通行不可、障害物エンティティに変換される）
	TileWall
	// TileWarpNext は次の階に向かうワープタイル（通行可能、進行ポータルエンティティに変換される）
	TileWarpNext
	// TileWarpEscape は脱出用ワープタイル（通行可能、脱出ポータルエンティティに変換される）
	TileWarpEscape
)
