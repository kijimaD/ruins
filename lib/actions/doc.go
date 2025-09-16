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
// # ActionとActivityの使い分け
//
// ## Action（アクション）
// **概念レベル**：ユーザーの意図や指示を表現する抽象的な概念
// - **ActionAPI**: すべてのアクションの統一エントリーポイント
// - **ActionParams**: アクション実行時のパラメータ
// - **ActionResult**: アクション実行結果
// - **ActionPoints**: アクションポイント（ターンコスト）
//
// 責務：プレイヤーやAIからの「何をしたいか」の指示を受け取り、パラメータを
// 受け渡して結果を返却する。ターンコストの管理も行う。
//
// ## Activity（アクティビティ）
// **実装レベル**：具体的な処理の実行単位
// - **ActivityType**: 実行可能なアクティビティの種別（移動、攻撃、休息など）
// - **Activity struct**: 実行中のアクティビティの状態データ
// - **ActivityInterface**: アクティビティの実行ロジック
// - **各種Activity実装**: MoveActivity, AttackActivity, RestActivityなど
//
// 責務：アクションの具体的な実行処理、状態管理（実行中、一時停止、完了、
// キャンセル）、ライフサイクル管理（Validate, Start, DoTurn, Finish, Canceled）、
// 継続実行とターン管理を行う。
//
// ## 関係性
// ```
// ユーザー入力 → Action（何をしたいか） → Activity（どう実行するか）
// ```
// Actionは外部インターフェース（API層）として機能し、Activityは内部実装
// （実行エンジン層）として機能する。この分離により、外部からは単純な
// Actionで指示でき、内部では複雑な継続実行・中断・再開処理をActivityで
// 管理できる。
//
// ## 主要コンポーネントと使い分け
//
// ### ActionAPI（外部インターフェース層）
// **用途**：外部システムからアクションを実行する際の統一エントリーポイント
// - systems/tile_input_system.go から使用（プレイヤー入力処理）
// - aiinput/processor.go から使用（AI行動処理）
//
// 責務：
// - アクティビティの作成（パラメータからActivityを生成）
// - コスト計算（APコストとターン数の計算）
// - 便利メソッドの提供（QuickMove, QuickAttack, StartRestなど）
// - ターン管理システムとの連携
//
// ### ActivityManager（内部実行エンジン）
// **用途**：アクティビティの実行管理（内部実装用）
// - ActionAPIの内部でのみ使用
// - 直接使用する場合はテスト目的のみ推奨
//
// 責務：
// - 実行中アクティビティの状態管理
// - ライフサイクル制御（Start, DoTurn, Finish, Canceled）
// - 実行可能性の検証
// - ActivityActorとの連携（具体的な実装の呼び出し）
//
// **重要**：外部パッケージからはActionAPIを使用してください。
// ActivityManagerは内部実装の詳細であり、直接使用は推奨しません。
//
// ### 個別Activity実装
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
// ### 1. 即座実行アクション（短期間で完了）
// - 移動、攻撃、アイテム拾得など
// - `TurnsTotal = 1`で残りAP1でも1ターンで完了
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
