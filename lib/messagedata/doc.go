// Package messagedata はメッセージウィンドウに表示するデータ構造を提供する
//
// # 概要
//
// このパッケージは、ゲーム内のメッセージウィンドウシステムで使用するデータ構造と
// ユーティリティ関数を提供します。メッセージの種類、選択肢、サイズ、コールバック
// などを統一的に管理し、UI層から独立したメッセージデータの操作を可能にします。
//
// # 主要な構造体
//
//   - MessageData: 単一のメッセージを表すデータ構造
//   - Choice: 選択肢のデータ
//   - Size: カスタムサイズ情報
//   - MessageSequence: 連続するメッセージのシーケンス
//   - MessageChain: 選択肢分岐を含む複雑なメッセージフロー
//
// # 基本的な使用方法
//
//	// 会話メッセージ
//	dialog := messagedata.NewDialogMessage("元気ですか？", "キャラクター名")
//
//	// 選択肢付きメッセージ
//	choice := messagedata.NewDialogMessage("どうしますか？", "").
//	    WithChoice("はい", func() { fmt.Println("はいが選ばれました") }).
//	    WithChoice("いいえ", func() { fmt.Println("いいえが選ばれました") })
//
//	// 選択肢から別のメッセージへの分岐
//	battleResult := messagedata.NewEventMessage("戦闘開始！").
//	    Message("激しい戦いが繰り広げられる").
//	    SystemMessage("勝利！")
//
//	escapeResult := messagedata.NewEventMessage("逃走成功").
//	    Message("安全な場所へ到着").
//	    SystemMessage("体力-10")
//
//	choiceWithBranch := messagedata.NewDialogMessage("敵と遭遇した！", "ナレーター").
//	    WithChoiceMessage("戦う", battleResult).
//	    WithChoiceMessage("逃げる", escapeResult)
//
// # メッセージチェーン機能
//
// 複数のメッセージを順次表示したり、選択肢の結果に応じて異なるメッセージ
// シーケンスを表示することができます。
//
//	// 連鎖記法での複数メッセージ定義
//	result := messagedata.NewEventMessage("戦闘開始").
//	    Message("激しい攻防").
//	    DialogMessage("勝利だ！", "主人公").
//	    SystemMessage("経験値+100")
//
//	// 選択肢ごとの結果定義
//	choiceResults := map[string]*messagedata.MessageData{
//	    "戦う": messagedata.NewEventMessage("勇敢に戦った！").
//	        Message("敵が倒れた").
//	        SystemMessage("勝利"),
//	    "逃げる": messagedata.NewEventMessage("逃走成功").
//	        Message("安全な場所に到着").
//	        SystemMessage("体力-10"),
//	}
//
// # MessageSequence
//
// 複雑なメッセージシーケンスを事前に構築する場合:
//
//	battleSeq := messagedata.NewMessageSequence().
//	    EventMessage("戦闘開始").
//	    Message("剣がぶつかり合う").
//	    DialogMessage("この程度か", "敵").
//	    EventMessage("反撃成功").
//	    SystemMessage("勝利")
//
//	msg := messagedata.NewDialogMessage("戦いますか？", "").
//	    WithChoice("はい", func() {}).
//	    Sequence(battleSeq)
//
// # MessageChain
//
// 選択肢分岐を含む複雑なフローの管理:
//
//	chain := messagedata.NewMessageChain()
//
//	question := messagedata.NewDialogMessage("どうしますか？", "NPC").
//	    WithChoice("戦う", func() {}).
//	    WithChoice("逃げる", func() {})
//
//	results := map[string]*messagedata.MessageData{
//	    "戦う": CreateBattleSequence(),
//	    "逃げる": CreateEscapeSequence(),
//	}
//
//	chain.AddChoiceMessage(question, results)
//
// # 責務と使い分け
//
//   - MessageData: 単一メッセージの表現、選択肢分岐対応
//   - MessageSequence: 線形な複数メッセージの連鎖
//   - MessageChain: 複雑な選択肢分岐フローの管理
//
// # 設計原則
//
//   - UI層からの独立: メッセージデータはUI実装に依存しない
//   - 不変性: 作成後のメッセージデータは変更されない
//   - 連鎖性: ビルダーパターンによる直感的なAPI
//   - 拡張性: 新しいメッセージタイプやパターンの追加が容易
//
// # 注意事項
//
//   - WithChoiceMessage により選択肢から新しいメッセージへの分岐が可能
//   - MessageWindowStateが連鎖メッセージとステート遷移を自動処理
//   - 選択肢の分岐とメッセージチェーンを組み合わせた複雑なフローに対応
package messagedata
