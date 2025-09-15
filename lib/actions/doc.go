// Package actions はアクションとアクティビティの実装を提供する。
//
// # 概要
//
// このパッケージは、ゲーム内のあらゆるアクション（移動、攻撃、休息など）の
// 具体的な実装を提供する。CDDAスタイルの中断可能なアクティビティシステムを採用し、
// 即座実行と継続実行の両方のアクションを統一的に管理する。
//
// # 責務
//
// - アクションの具体的な実装とロジック
// - アクティビティライフサイクル管理（開始、実行、中断、再開、完了）
// - アクションコストの定義と管理
// - ターン管理システムとの連携
//
// # 使い分け
//
// ## ActionAPI
// - 全アクションの統一エントリーポイント
// - アクション実行の中心ハブ
// - ターン管理システムとの連携
//
// ## ActivityManager
// - アクティビティの状態管理
// - 中断・再開機能の提供
// - 複数アクティビティの同時管理
//
// ## 個別アクション実装
// - **MoveActivity**: 移動アクション（即座実行）
// - **AttackActivity**: 攻撃アクション（即座実行）
// - **RestActivity**: 休息アクション（継続実行、中断可能）
// - **WaitActivity**: 待機アクション（継続実行、中断可能）
//
// # 他パッケージとの関係
//
// ```
// systems → actions.ActionAPI → アクション実行
//
//	↓
//
// actions → turns.ConsumePlayerMoves → コスト消費
// ```
//
// ## 責務の境界
//
// - **actions**: どのような行動をするか（What Action & How）
// - **turns**: いつ実行するかの制御（When）
// - **systems**: 何の入力から実行するか（What Input）
//
// # アーキテクチャ
//
// ## 2層構造
//
// ### 1. 即座実行アクション（1ターンで完了）
// - 移動、攻撃、アイテム拾得など
// - `TurnsTotal = 1`で即座に完了
// - シンプルな実行ロジック
//
// ### 2. 継続実行アクション（複数ターンにわたる）
// - 休息、読書、クラフトなど
// - `TurnsTotal > 1`で段階的に実行
// - 中断・再開機能あり
//
// ## 統一インターフェース
//
// 全てのアクションは`Activity`構造体を通じて統一的に管理される：
//
//	type Activity struct {
//		Type       ActivityType  // アクション種別
//		State      ActivityState // 実行状態
//		TurnsTotal int          // 必要ターン数
//		TurnsLeft  int          // 残りターン数
//		// ...
//	}
//
// # 設計原則
//
// 1. **統一性**: 全アクションを同じインターフェースで管理
// 2. **拡張性**: 新しいアクションの追加が容易
// 3. **中断可能性**: 必要に応じてアクションを中断・再開
// 4. **責務分離**: アクションロジックとターン管理を分離
// 5. **検証**: 実行前・再開前の条件チェック
//
// # 使用例
//
//	// ActionAPIを通じた統一的なアクション実行
//	actionAPI := actions.NewActionAPI()
//
//	// 即座実行アクション（移動）
//	result, err := actionAPI.QuickMove(player, destination, world)
//
//	// 継続実行アクション（休息）
//	result, err := actionAPI.StartRest(player, 10, world)
//
//	// アクティビティの管理
//	actionAPI.InterruptActivity(player, "戦闘開始")
//	actionAPI.ResumeActivity(player, world)
//
//	// ターン毎の処理
//	actionAPI.ProcessTurn(world)
//
// # CDDAとの対応関係
//
// このパッケージの設計は Cataclysm: Dark Days Ahead の activity_actor システムを参考にしている：
//
// - CDDAのactivity_actor → Activity構造体
// - CDDAのdo_turn() → DoTurn()メソッド
// - CDDAのfinish() → Complete()メソッド
// - CDDAのcanceled() → Interrupt()メソッド
// - CDDAのmove_cost → アクションコスト概念
//
// # 拡張方法
//
// 新しいアクションを追加する場合：
//
// 1. ActivityTypeに新しい定数を追加
// 2. activityInfosに情報を追加
// 3. 具体的な実装ファイルを作成（例：new_action.go）
// 4. ActionAPI.createActivityに分岐を追加
// 5. 必要に応じてActivity.DoTurnに処理を追加
package actions
