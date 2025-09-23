package messagedata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDialogMessage(t *testing.T) {
	t.Parallel()

	t.Run("基本的な会話メッセージの作成", func(t *testing.T) {
		t.Parallel()

		text := "こんにちは"
		speaker := "テストキャラクター"
		msg := NewDialogMessage(text, speaker)

		assert.Equal(t, text, msg.Text)
		assert.Equal(t, speaker, msg.Speaker)
		assert.Empty(t, msg.Choices)
		assert.Nil(t, msg.Size)
		assert.Nil(t, msg.OnComplete)
		assert.Empty(t, msg.NextMessages)
	})

	t.Run("空の話者名でも作成可能", func(t *testing.T) {
		t.Parallel()

		msg := NewDialogMessage("メッセージ", "")

		assert.Equal(t, "メッセージ", msg.Text)
		assert.Equal(t, "", msg.Speaker)
	})

	t.Run("空のテキストでも作成可能", func(t *testing.T) {
		t.Parallel()

		msg := NewDialogMessage("", "キャラクター")

		assert.Equal(t, "", msg.Text)
		assert.Equal(t, "キャラクター", msg.Speaker)
	})
}

func TestNewSystemMessage(t *testing.T) {
	t.Parallel()

	t.Run("システムメッセージの作成", func(t *testing.T) {
		t.Parallel()

		text := "ゲームが保存されました"
		msg := NewSystemMessage(text)

		assert.Equal(t, text, msg.Text)
		assert.Equal(t, "システム", msg.Speaker)
		assert.Empty(t, msg.Choices)
		assert.Nil(t, msg.Size)
		assert.Nil(t, msg.OnComplete)
		assert.Empty(t, msg.NextMessages)
	})

	t.Run("空のテキストでも作成可能", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("")

		assert.Equal(t, "", msg.Text)
		assert.Equal(t, "システム", msg.Speaker)
	})
}

func TestMessageDataBuilderMethods(t *testing.T) {
	t.Parallel()

	t.Run("WithSpeakerメソッド", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("テスト").WithSpeaker("システム")

		assert.Equal(t, "システム", msg.Speaker)
		assert.Equal(t, "テスト", msg.Text)
	})

	t.Run("WithSizeメソッド", func(t *testing.T) {
		t.Parallel()

		msg := NewDialogMessage("テスト", "").WithSize(800, 600)

		require.NotNil(t, msg.Size)
		assert.Equal(t, 800, msg.Size.Width)
		assert.Equal(t, 600, msg.Size.Height)
	})

	t.Run("WithOnCompleteメソッド", func(t *testing.T) {
		t.Parallel()

		callbackExecuted := false
		callback := func() {
			callbackExecuted = true
		}

		msg := NewSystemMessage("テスト").WithOnComplete(callback)

		require.NotNil(t, msg.OnComplete)
		msg.OnComplete()
		assert.True(t, callbackExecuted)
	})

	t.Run("メソッドチェーンの組み合わせ", func(t *testing.T) {
		t.Parallel()

		callbackExecuted := false
		msg := NewDialogMessage("テストメッセージ", "初期話者").
			WithSpeaker("最終話者").
			WithSize(1024, 768).
			WithOnComplete(func() { callbackExecuted = true })

		assert.Equal(t, "テストメッセージ", msg.Text)
		assert.Equal(t, "最終話者", msg.Speaker)
		require.NotNil(t, msg.Size)
		assert.Equal(t, 1024, msg.Size.Width)
		assert.Equal(t, 768, msg.Size.Height)
		require.NotNil(t, msg.OnComplete)
		msg.OnComplete()
		assert.True(t, callbackExecuted)
	})
}

func TestChoiceMethods(t *testing.T) {
	t.Parallel()

	t.Run("WithChoiceメソッド", func(t *testing.T) {
		t.Parallel()

		actionExecuted := false
		action := func() {
			actionExecuted = true
		}

		msg := NewDialogMessage("選択してください", "").
			WithChoice("選択肢1", action)

		require.Len(t, msg.Choices, 1)
		choice := msg.Choices[0]
		assert.Equal(t, "選択肢1", choice.Text)
		assert.Equal(t, "", choice.Description)
		assert.Nil(t, choice.MessageData)
		assert.False(t, choice.Disabled)

		require.NotNil(t, choice.Action)
		choice.Action()
		assert.True(t, actionExecuted)
	})

	t.Run("WithChoiceDescriptionメソッド", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("選択してください").
			WithChoiceDescription("攻撃", "敵を攻撃する", func() {})

		require.Len(t, msg.Choices, 1)
		choice := msg.Choices[0]
		assert.Equal(t, "攻撃", choice.Text)
		assert.Equal(t, "敵を攻撃する", choice.Description)
		assert.NotNil(t, choice.Action)
	})

	t.Run("WithChoiceMessageメソッド", func(t *testing.T) {
		t.Parallel()

		resultMessage := NewSystemMessage("結果メッセージ")
		msg := NewDialogMessage("選択してください", "").
			WithChoiceMessage("選択肢", resultMessage)

		require.Len(t, msg.Choices, 1)
		choice := msg.Choices[0]
		assert.Equal(t, "選択肢", choice.Text)
		assert.Equal(t, resultMessage, choice.MessageData)
	})

	t.Run("WithChoiceMessageDescriptionメソッド", func(t *testing.T) {
		t.Parallel()

		resultMessage := NewSystemMessage("処理完了")
		msg := NewDialogMessage("何をしますか？", "").
			WithChoiceMessageDescription("実行", "処理を実行します", resultMessage)

		require.Len(t, msg.Choices, 1)
		choice := msg.Choices[0]
		assert.Equal(t, "実行", choice.Text)
		assert.Equal(t, "処理を実行します", choice.Description)
		assert.Equal(t, resultMessage, choice.MessageData)
	})

	t.Run("複数の選択肢の追加", func(t *testing.T) {
		t.Parallel()

		msg := NewDialogMessage("どうしますか？", "NPC").
			WithChoice("はい", func() {}).
			WithChoice("いいえ", func() {}).
			WithChoiceDescription("詳細", "詳しく見る", func() {})

		assert.Len(t, msg.Choices, 3)
		assert.Equal(t, "はい", msg.Choices[0].Text)
		assert.Equal(t, "いいえ", msg.Choices[1].Text)
		assert.Equal(t, "詳細", msg.Choices[2].Text)
		assert.Equal(t, "詳しく見る", msg.Choices[2].Description)
	})

	t.Run("nilアクションでも追加可能", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("テスト").WithChoice("選択肢", nil)

		require.Len(t, msg.Choices, 1)
		assert.Nil(t, msg.Choices[0].Action)
	})
}

func TestMessageChaining(t *testing.T) {
	t.Parallel()

	t.Run("DialogMessageメソッドによる連鎖", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("システムメッセージ").
			DialogMessage("会話メッセージ", "キャラクター")

		require.Len(t, msg.NextMessages, 1)
		nextMsg := msg.NextMessages[0]
		assert.Equal(t, "会話メッセージ", nextMsg.Text)
		assert.Equal(t, "キャラクター", nextMsg.Speaker)
	})

	t.Run("SystemMessageメソッドによる連鎖", func(t *testing.T) {
		t.Parallel()

		msg := NewDialogMessage("会話", "キャラクター").
			SystemMessage("システム通知")

		require.Len(t, msg.NextMessages, 1)
		nextMsg := msg.NextMessages[0]
		assert.Equal(t, "システム通知", nextMsg.Text)
	})

	t.Run("SystemMessageメソッドによる連鎖（追加）", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("開始").
			SystemMessage("イベント発生")

		require.Len(t, msg.NextMessages, 1)
		nextMsg := msg.NextMessages[0]
		assert.Equal(t, "イベント発生", nextMsg.Text)
		assert.Equal(t, "システム", nextMsg.Speaker)
	})

	t.Run("複数メッセージの連鎖", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("戦闘開始").
			SystemMessage("攻撃").
			DialogMessage("やったか？", "主人公").
			SystemMessage("勝利！")

		assert.Len(t, msg.NextMessages, 3)
		assert.Equal(t, "攻撃", msg.NextMessages[0].Text)
		assert.Equal(t, "やったか？", msg.NextMessages[1].Text)
		assert.Equal(t, "主人公", msg.NextMessages[1].Speaker)
		assert.Equal(t, "勝利！", msg.NextMessages[2].Text)
	})

	t.Run("HasNextMessagesメソッド", func(t *testing.T) {
		t.Parallel()

		// 次のメッセージがない場合
		msg1 := NewSystemMessage("単体メッセージ")
		assert.False(t, msg1.HasNextMessages())

		// 次のメッセージがある場合
		msg2 := NewSystemMessage("最初").SystemMessage("次")
		assert.True(t, msg2.HasNextMessages())
	})

	t.Run("GetNextMessagesメソッド", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("最初").
			SystemMessage("2番目").
			SystemMessage("3番目")

		nextMessages := msg.GetNextMessages()
		assert.Len(t, nextMessages, 2)
		assert.Equal(t, "2番目", nextMessages[0].Text)
		assert.Equal(t, "3番目", nextMessages[1].Text)
	})
}

func TestSize(t *testing.T) {
	t.Parallel()

	t.Run("Sizeの基本構造", func(t *testing.T) {
		t.Parallel()

		size := &Size{Width: 800, Height: 600}
		assert.Equal(t, 800, size.Width)
		assert.Equal(t, 600, size.Height)
	})

	t.Run("WithSizeで負の値も設定可能", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("テスト").WithSize(-100, -50)
		require.NotNil(t, msg.Size)
		assert.Equal(t, -100, msg.Size.Width)
		assert.Equal(t, -50, msg.Size.Height)
	})

	t.Run("WithSizeでゼロ値も設定可能", func(t *testing.T) {
		t.Parallel()

		msg := NewSystemMessage("テスト").WithSize(0, 0)
		require.NotNil(t, msg.Size)
		assert.Equal(t, 0, msg.Size.Width)
		assert.Equal(t, 0, msg.Size.Height)
	})
}

func TestChoice(t *testing.T) {
	t.Parallel()

	t.Run("Choiceの基本構造", func(t *testing.T) {
		t.Parallel()

		resultMsg := NewSystemMessage("結果")
		choice := Choice{
			Text:        "選択肢",
			Description: "説明",
			Action:      func() {},
			MessageData: resultMsg,
			Disabled:    true,
		}

		assert.Equal(t, "選択肢", choice.Text)
		assert.Equal(t, "説明", choice.Description)
		assert.NotNil(t, choice.Action)
		assert.Equal(t, resultMsg, choice.MessageData)
		assert.True(t, choice.Disabled)
	})

	t.Run("Choiceの初期値", func(t *testing.T) {
		t.Parallel()

		choice := Choice{}

		assert.Equal(t, "", choice.Text)
		assert.Equal(t, "", choice.Description)
		assert.Nil(t, choice.Action)
		assert.Nil(t, choice.MessageData)
		assert.False(t, choice.Disabled)
	})
}

func TestComplexScenarios(t *testing.T) {
	t.Parallel()

	t.Run("選択肢分岐を含む複雑なメッセージフロー", func(t *testing.T) {
		t.Parallel()

		// 戦闘結果のメッセージシーケンス
		battleResult := NewSystemMessage("戦闘開始").
			SystemMessage("激しい攻防").
			DialogMessage("勝利だ！", "主人公").
			SystemMessage("経験値+100")

		// 逃走結果のメッセージ
		escapeResult := NewSystemMessage("逃走成功").
			SystemMessage("体力-10")

		// 選択肢付きメッセージ
		encounterMsg := NewDialogMessage("敵に遭遇した！", "ナレーター").
			WithChoiceMessage("戦う", battleResult).
			WithChoiceMessage("逃げる", escapeResult)

		// 検証
		assert.Equal(t, "敵に遭遇した！", encounterMsg.Text)
		assert.Equal(t, "ナレーター", encounterMsg.Speaker)
		assert.Len(t, encounterMsg.Choices, 2)

		// 戦うの選択肢
		fightChoice := encounterMsg.Choices[0]
		assert.Equal(t, "戦う", fightChoice.Text)
		require.NotNil(t, fightChoice.MessageData)
		assert.Equal(t, "戦闘開始", fightChoice.MessageData.Text)
		assert.Len(t, fightChoice.MessageData.NextMessages, 3)

		// 逃げるの選択肢
		escapeChoice := encounterMsg.Choices[1]
		assert.Equal(t, "逃げる", escapeChoice.Text)
		require.NotNil(t, escapeChoice.MessageData)
		assert.Equal(t, "逃走成功", escapeChoice.MessageData.Text)
		assert.Len(t, escapeChoice.MessageData.NextMessages, 1)
	})

	t.Run("全機能を組み合わせたメッセージ", func(t *testing.T) {
		t.Parallel()

		completeCalled := false
		actionCalled := false

		msg := NewDialogMessage("複雑なテスト", "テストキャラクター").
			WithSpeaker("最終キャラクター").
			WithSize(1024, 768).
			WithChoice("アクション", func() { actionCalled = true }).
			WithChoiceDescription("説明付き", "詳細な説明", func() {}).
			WithOnComplete(func() { completeCalled = true }).
			SystemMessage("次のメッセージ").
			SystemMessage("システム通知").
			DialogMessage("最後の会話", "別キャラクター")

		// 基本設定の確認
		assert.Equal(t, "複雑なテスト", msg.Text)
		assert.Equal(t, "最終キャラクター", msg.Speaker)
		require.NotNil(t, msg.Size)
		assert.Equal(t, 1024, msg.Size.Width)
		assert.Equal(t, 768, msg.Size.Height)

		// 選択肢の確認
		assert.Len(t, msg.Choices, 2)
		assert.Equal(t, "アクション", msg.Choices[0].Text)
		assert.Equal(t, "説明付き", msg.Choices[1].Text)
		assert.Equal(t, "詳細な説明", msg.Choices[1].Description)

		// 連鎖メッセージの確認
		assert.Len(t, msg.NextMessages, 3)
		assert.Equal(t, "次のメッセージ", msg.NextMessages[0].Text)
		assert.Equal(t, "システム通知", msg.NextMessages[1].Text)
		assert.Equal(t, "最後の会話", msg.NextMessages[2].Text)
		assert.Equal(t, "別キャラクター", msg.NextMessages[2].Speaker)

		// コールバックの実行
		require.NotNil(t, msg.OnComplete)
		msg.OnComplete()
		assert.True(t, completeCalled)

		// アクションの実行
		require.NotNil(t, msg.Choices[0].Action)
		msg.Choices[0].Action()
		assert.True(t, actionCalled)
	})
}
