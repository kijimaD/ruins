// Package menu は状態管理を持つ再利用可能なキーボードナビゲーション対応メニューウィジェットを提供する。
//
// # Overview
//
// menuパッケージは、キーボードナビゲーション対応の高機能メニューウィジェットを提供します。
// 設定駆動でメニュー項目を定義し、コールバック機能でアプリケーションロジックと連携します。
// 状態管理、イベント処理、UI描画を分離した設計により、テスタブルで再利用可能な実装を実現しています。
//
// # Features
//
//   - キーボードナビゲーション（矢印キー、Tab、Enter、Escape）
//   - 循環ナビゲーション（オプション）
//   - 無効項目のスキップ
//   - グリッド表示対応（2D ナビゲーション）
//   - 水平/垂直レイアウト
//   - 設定駆動の宣言的定義
//   - イベント駆動のコールバック
//   - 完全にテスト可能な設計
//
// # Basic Usage
//
//	// メニュー項目を定義
//	items := []menu.Item{
//		{ID: "start", Label: "ゲーム開始"},
//		{ID: "load", Label: "ロード"},
//		{ID: "exit", Label: "終了"},
//	}
//
//	// メニューを作成
//	m := menu.NewMenu(menu.MenuConfig{
//		Items:          items,
//		WrapNavigation: true,
//		Orientation:    menu.Vertical,
//	}, menu.MenuCallbacks{
//		OnSelect: func(index int, item menu.Item) {
//			switch item.ID {
//			case "start":
//				startNewGame()
//			case "exit":
//				exitApplication()
//			}
//		},
//		OnCancel: func() {
//			goBack()
//		},
//	})
//
//	// ゲームループ内で更新
//	func Update() {
//		m.Update(keyboardInput)
//	}
//
// # Grid Layout
//
// メニューはグリッド表示にも対応しており、2Dナビゲーションが可能です：
//
//	config := menu.MenuConfig{
//		Items:   items,
//		Columns: 3, // 3列のグリッド
//	}
//
//	// 左右の矢印キーで水平移動、上下で垂直移動
//
// # Keyboard Controls
//
//   - ↑/↓: 垂直ナビゲーション
//   - ←/→: グリッド表示時の水平ナビゲーション
//   - Tab: 次の項目へ移動
//   - Shift+Tab: 前の項目へ移動
//   - Enter: 項目選択
//   - Escape: キャンセル
//
// # Testing
//
// menuパッケージは完全にテスト可能な設計になっています：
//
//	func TestMyMenu(t *testing.T) {
//		mockInput := input.NewMockKeyboardInput()
//		menu := menu.NewMenu(config, callbacks)
//
//		// キー入力をシミュレート
//		mockInput.SetKeyJustPressed(ebiten.KeyArrowDown, true)
//		menu.Update(mockInput)
//
//		// 結果を検証
//		assert.Equal(t, 1, menu.GetFocusedIndex())
//	}
//
// # UI Integration
//
// メニューのUI要素はMenuUIBuilderを通じて構築されます：
//
//	builder := menu.NewMenuUIBuilder(world)
//	container := builder.BuildUI(menu)
//
//	// EbitenUIのコンテナとして使用可能
//	parentContainer.AddChild(container)
//
// # Configuration Options
//
// MenuConfig で様々な設定が可能です：
//
//   - Items: メニュー項目のリスト
//   - InitialIndex: 初期フォーカス位置
//   - WrapNavigation: 端での循環ナビゲーション
//   - Orientation: 垂直/水平レイアウト
//   - Columns: グリッド表示時の列数
//
// # Callback Events
//
// MenuCallbacks で以下のイベントを処理できます：
//
//   - OnSelect: 項目選択時
//   - OnCancel: キャンセル時（Escape キー）
//   - OnFocusChange: フォーカス変更時
//   - OnHover: マウスホバー時（将来実装予定）
//
// # Design Principles
//
//   - Configuration over Code: 設定による宣言的定義
//   - Separation of Concerns: ロジック・状態・UIの分離
//   - Testability: モック可能な依存関係
//   - Accessibility: キーボードファーストなナビゲーション
//   - Flexibility: 様々なレイアウトとユースケースに対応
package menu
