// Package actions はゲーム内のアクション（行動）システムを提供する。プレイヤーの意図を表現する
//
// ## 主要機能
//
// ### アクション定義 (types.go)
// - ActionID enumによるアクション種別管理
// - 基本的なアクションメタデータ（名前、移動コスト、ターン消費等）
// - CDDAのaction.hを参考にしたシンプルな設計
//
// ### 実行可能性チェック (validation.go)
// - CanExecuteAction()による事前チェック
// - CDDAのcan_interact_at()システムを参考
// - 位置・状態・条件による実行可能性判定
//
// ### アクション実行 (executor.go)
// - Execute()による実際のアクション実行
// - Effectシステムとの統合
// - エラーハンドリングとフィードバック
//
// ### アクティビティシステム (activity.go)
// - 時間消費を伴うアクション（休息、移動等）
// - CDDAのplayer_activityを参考にした継続処理
// - 中断機能とターンベース処理
//
// ## 設計思想
//
// - **シンプル性**: CDDAの複雑性を排除し、RPGに必要な機能に集約
// - **拡張性**: 新アクション追加が容易なenum + interface設計
// - **統合性**: 既存のeffectsシステムとの密接な連携
// - **段階的実装**: 基本→時間消費→UI改善の順で構築
//
// ## 使用例
//
//	// アクション実行可能性チェック
//	if actions.CanExecuteAction(actions.ActionAttack, playerPos, world) {
//		// アクション実行
//		result := actions.Execute(actions.ActionAttack, player, target, world)
//		if result.Success {
//			// 成功処理
//		}
//	}
//
//	// 時間消費アクション
//	activity := &actions.Activity{
//		Type: actions.ActivityRest,
//		TurnsTotal: 10,
//		TurnsLeft: 10,
//	}
//	player.SetCurrentActivity(activity)
package actions
