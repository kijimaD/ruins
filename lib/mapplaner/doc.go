// Package mapplaner provides map planning functionality.
//
// このパッケージは階層マップの生成機能を提供します：
//   - タイルベースのマップ生成
//   - 各種マップアルゴリズム（部屋、洞窟、森林、遺跡など）
//   - エンティティ配置（NPC、アイテム、ワープホール）
//
// ## タイル定義
//
// マップ生成で使用されるタイルタイプ：
//   - TileEmpty: 空のタイル
//   - TileFloor: 床タイル（通行可能）
//   - TileWall: 壁タイル（通行不可）
//   - TileWarpNext: 次階層へのワープタイル（通行可能）
//   - TileWarpEscape: 脱出用ワープタイル（通行可能）
//
// ## 通行可否判定
//
// マップ生成時には PathFinder.IsWalkable() で通行可否を判定します：
//   - 通行可能: TileFloor, TileWarpNext, TileWarpEscape
//   - 通行不可: TileWall
//
// 生成されたタイルはlevel.goでエンティティに変換され、実行時はエンティティベースの
// 通行可否判定（movement.CanMoveTo）が使用されます。
//
// ## マップ生成の流れ
//
// 1. PlannerChainによる段階的生成
// 2. タイル配列の構築
// 3. エンティティへの変換
// 4. NPC・アイテム配置
package mapplaner
