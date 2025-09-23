// Package messagewindow はストーリーメッセージ表示用のモーダルウィンドウを提供する
//
// # Overview
//
// messagewindowパッケージは、ストーリー進行やイベント時に表示される
// Elona風のメッセージウィンドウシステムを提供します。
// モーダルウィンドウとして表示され、プレイヤーの操作を待ちます。
//
// # Features
//
//   - ストーリーメッセージの表示（テキスト、画像対応）
//   - スキップ可能なページング機能
//   - キーボード・マウスでの操作
//   - 選択肢システム対応
//   - カスタマイズ可能な外観とレイアウト
//
// # Basic Usage
//
//	// シンプルなメッセージウィンドウ
//	window := messagewindow.NewBuilder().
//		Message("冒険者よ、ようこそ！").
//		Build()
//
//	// ゲームループで更新・描画
//	window.Update(input)
//	window.Draw(screen)
//
//	// ウィンドウが閉じられたかチェック
//	if window.IsClosed() {
//		// 次の処理へ
//	}
//
// # Future Extensions
//
//	// 選択肢システムの使用例
//	window := messagewindow.NewBuilder().
//		Message("この中から選べ。").
//		Choice("戦う", func() { fight() }).
//		Choice("逃げる", func() { escape() }).
//		Choice("交渉する", func() { negotiate() }).
//		Build()
//
// # Design Principles
//
//   - Builder Pattern: 柔軟で読みやすい設定
//   - Event-Driven: コールバックによる拡張性
//   - Single Responsibility: メッセージ表示に特化
//   - Separation of Concerns: UIとロジックの分離
//   - Future-Proof: 選択肢システムへの拡張を考慮
package messagewindow
