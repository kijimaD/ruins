// Package eui はEbitenUIコンポーネントに対するスタイル付きのゲーム固有ラッパー関数を提供する。
//
// # Overview
//
// euiパッケージは、EbitenUIライブラリの基本コンポーネントに対して
// プロジェクト固有のスタイルとデザインを適用したヘルパー関数群を提供します。
// エンティティと関わらない、UI表示に特化した基本的なパーツを提供します。
//
// # Package Hierarchy
//
// このプロジェクトのUIアーキテクチャは3層構造になっています：
//
//	widgets/     ← 業務ロジック付きの高レベルコンポーネント
//	   ↓ 使用
//	eui/         ← プロジェクト固有スタイルの中レベルコンポーネント（このパッケージ）
//	   ↓ 使用
//	ebitenui/    ← 外部ライブラリの低レベルコンポーネント
//
// # Responsibilities
//
// euiパッケージの責務：
//   - EbitenUIコンポーネントのプロジェクト固有スタイル適用
//   - ゲームリソース（フォント、色、画像）との統合
//   - 基本的なレイアウトコンテナの提供
//   - 静的なUI要素のファクトリ関数群
//   - 一貫したデザインシステムの維持
//
// # Usage vs Other Packages
//
// ## euiパッケージを使う場合
//   - 基本的なレイアウトコンテナが欲しい（NewRowContainer, NewVerticalContainer）
//   - プロジェクト統一スタイルのボタンやテキストが欲しい（NewButton, NewMenuText）
//   - 静的な表示のみで状態管理は不要
//   - ゲームリソース（World）を使った表示をしたい
//   - 簡単なヘルパー関数で十分
//
// ## widgetsパッケージを使う場合
//   - メニュー、ダイアログ、フォームなど複雑な操作が必要
//   - キーボードナビゲーションが必要
//   - 状態管理が必要（選択状態、入力データなど）
//   - ビジネスロジックとの連携が必要
//   - 単体テストを書きたい
//
// # Example
//
// euiパッケージの典型的な使用例：
//
//	// 基本的なレイアウト構築
//	container := eui.NewVerticalContainer()
//
//	// プロジェクトスタイルのボタン作成
//	button := eui.NewButton("クリック", world)
//	container.AddChild(button)
//
//	// ゲーム固有スタイルのテキスト
//	title := eui.NewMenuText("タイトル", world)
//	container.AddChild(title)
//
//	// 分割レイアウト
//	leftPanel := eui.NewVerticalContainer()
//	rightPanel := eui.NewVerticalContainer()
//	splitContainer := eui.NewWSplitContainer(leftPanel, rightPanel)
//
// # Design Principles
//
//   - Styling Consistency: プロジェクト全体で一貫したデザイン
//   - Resource Integration: ゲームリソースとの統合
//   - Simplicity: シンプルなファクトリ関数
//   - Stateless: 状態を持たない純粋な関数群
//   - EbitenUI Compatibility: EbitenUIとの高い互換性
package eui
