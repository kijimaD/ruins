// Package mapspawner はEntityPlanに基づいてECSエンティティを生成する機能を提供する
//
// このパッケージの責務:
// - EntityPlanからLevelオブジェクトの生成
// - EntitySpecからECSエンティティの生成
// - スプライト番号の補完などスポーン時の最終調整
// - スポーン処理のエラーハンドリング
//
// 使い分け:
// - mapplaner: マップ構造の計画・設計、プランナーチェーン制御
// - mapspawner: 実際のECSエンティティ生成、スポーン時の最終調整
//
// 責務境界:
// - mapplanerが全てのplanning（計画）を担当
// - mapspawnerは純粋にspawning（生成）のみを担当
// - PlanAndSpawn関数はmapplanerのBuildPlanWithIntegrationを呼び出してEntityPlanを取得
//
// 依存関係:
// - worldhelper: 各種エンティティ生成関数
// - mapplaner: EntityPlan, EntitySpec等の型定義
package mapspawner
