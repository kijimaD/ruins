// Package mapspawner はMapPlanに基づいてECSエンティティを生成する機能を提供する
//
// このパッケージの責務:
// - MapPlanからLevelオブジェクトの生成
// - EntitySpecからECSエンティティの生成
// - スポーン処理のエラーハンドリング
//
// 使い分け:
// - mapplaner: マップ構造の計画・設計
// - mapspawner: 実際のECSエンティティ生成
//
// 依存関係:
// - worldhelper: 各種エンティティ生成関数
// - mapplaner: MapPlan, EntitySpec等の型定義
package mapspawner
