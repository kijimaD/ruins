// Package levelgen はゲームレベル生成の統合機能を提供する
//
// このパッケージの責務:
// - mapplanerとmapspawnerを組み合わせたレベル生成
// - NPCやアイテム配置のオーケストレーション
// - プレイヤー配置とマップ接続性の検証
//
// 使い分け:
// - mapplaner: マップ構造の計画・設計
// - mapspawner: エンティティの生成
// - levelgen: 全体的なゲームレベルの生成統合
//
// 依存関係:
// - mapplaner: マップ構造生成
// - mapspawner: エンティティ生成
// - worldhelper: プレイヤー配置、NPC/アイテム生成
package levelgen
