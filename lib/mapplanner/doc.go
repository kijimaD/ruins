// Package mapplanner provides map planning functionality.
//
// このパッケージは階層マップの生成機能を提供します：
//   - タイルベースのマップ生成
//   - 各種マップアルゴリズム（部屋、洞窟、森林、遺跡など）
//   - タイルとエンティティの配置計画作成
//
// ## マップ構造の概念
//
// マップは規定数のタイル（TileWidth × TileHeight）で正方形に構成されます。
// 各位置には必ず1つのタイルが存在し、タイル自体は不変です。
// タイルの上にエンティティ（壁、床、アイテム、NPCなど）を配置していきます。
//
// ## 主要データ構造の違い
//
// - **MetaPlan**: マップ生成プロセス中の中間データ
//   - タイル配列（[]raw.TileRaw）、部屋情報、廊下情報、乱数生成器を含む
//   - PlannerChain内で段階的に構築される
//   - 生成アルゴリズムで使用される作業用データ
//
// - **EntityPlan**: エンティティ生成用の最終配置計画
//   - EntitySpecのリストとして詳細な配置計画を管理
//   - MetaPlanから BuildPlan() で生成される
//   - mapspawnerで実際のECSエンティティ生成に使用される
//
// ## タイル定義
//
// ### 基本タイルタイプ
// マップ生成で使用される標準タイルタイプ：
//   - planData.GenerateTile("Empty"): 空のタイル（デフォルト状態）
//   - planData.GenerateTile("Floor"): 床タイル（通行可能）
//   - planData.GenerateTile("Wall"): 壁タイル（通行不可）
//   - TileWater: 水タイル（通行可能だが特殊）
//   - TileDoor: 扉タイル（開閉可能な通路）
//   - TilePit: 落とし穴タイル（歩くと落下）
//
// ### TOMLベースタイル定義システム
// 新しいタイル定義システムでは、TOMLファイルでタイルの種類と属性を定義できます：
//
//	[[tile]]
//	Name = "Floor"
//	Description = "床タイル - 移動可能な基本的なタイル"
//	Type = "FLOOR"
//	Walkable = true
//	BlocksLOS = false
//
// TileMasterクラスを使用してタイル定義を読み込み・管理：
//   - LoadTileFromFile(): TOMLファイルからタイル定義を読み込み
//   - GenerateTile(): 名前指定でタイルオブジェクトを生成
//
// ## エンティティ
//
// エンティティとして実装されます：
//   - 床タイル + ワープポータルエンティティ = ワープ機能のある場所
//
// ## 通行可否判定
//
// マップ生成時にはタイルの Walkable フィールドで通行可否を判定します：
//   - 通行可能: planData.GenerateTile("Floor")（Walkable=true）
//   - 通行不可: planData.GenerateTile("Wall")（Walkable=false）, planData.GenerateTile("Empty")（Walkable=false）
//
// ## 配置計画の作成
//
// mapplannerは以下の2つのアプローチで配置計画を作成します：
//
// ### 1. タイルベース生成（標準）
//   - PlannerChainでタイル配列を生成
//   - mapspawnerでタイルからエンティティ配置計画（EntityPlan）を自動生成
//   - タイルタイプに応じて対応するエンティティ（床、壁、ワープホールなど）を配置
//
// ### 2. 文字列ベース生成（高度）
//   - 文字列マップから直接タイルとエンティティの配置計画を作成
//   - NPCやアイテムなどの詳細なエンティティ配置も可能
//   - EntityPlanで明示的にエンティティ配置計画（EntitySpec）を管理
//
// ## マップ生成の流れ
//
// ### タイルベース生成の場合
// 1. タイル配列の初期化（全てplanData.GenerateTile("Empty")）
// 2. PlannerChainによる段階的タイル配置（MetaPlan）
// 3. MetaPlan.BuildPlan()でEntityPlan生成
// 4. mapspawner.SpawnLevelで実際のECSエンティティ生成
//
// ### 文字列ベース生成の場合
// 1. 文字列マップ定義からタイル・エンティティ配置計画を作成
// 2. EntityPlanにEntitySpecとして詳細なエンティティ配置計画を記録
// 3. mapspawner.SpawnLevelでEntityPlanに基づくECSエンティティ生成
//
// いずれの場合も実行時はエンティティベースの通行可否判定（movement.CanMoveTo）を使用
package mapplanner
