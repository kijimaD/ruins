// Package messagedata はメッセージウィンドウに表示するデータ構造を提供する
package messagedata

import (
	"image/color"

	w "github.com/kijimaD/ruins/lib/world"
)

// MessageData はメッセージウィンドウに表示するデータ
type MessageData struct {
	Speaker          string
	Choices          []Choice
	OnComplete       func()          // メッセージ完了時のコールバック
	NextMessages     []*MessageData  // 次に表示するメッセージ群
	TextSegmentLines [][]TextSegment // 行ごとの色付きテキストセグメント
}

// TextSegment は色付きテキストのセグメント
type TextSegment struct {
	Text            string
	Color           *color.RGBA // nilの場合はデフォルト色
	BackgroundColor *color.RGBA // nilの場合は背景なし
}

// Choice は選択肢のデータ
type Choice struct {
	Text        string
	Action      func(w.World) error // 選択時に実行する
	MessageData *MessageData        // 選択肢を選んだ時に表示するメッセージ
	Disabled    bool
}

// NewDialogMessage は会話メッセージを作成する
func NewDialogMessage(text, speaker string) *MessageData {
	msg := &MessageData{
		Speaker: speaker,
	}
	msg.AddText(text)
	return msg
}

// NewSystemMessage はシステムメッセージを作成する
func NewSystemMessage(text string) *MessageData {
	msg := &MessageData{
		Speaker: "システム",
	}
	msg.AddText(text)
	return msg
}

// WithChoice は選択肢を追加する
func (m *MessageData) WithChoice(text string, action func(w.World) error) *MessageData {
	m.Choices = append(m.Choices, Choice{
		Text:   text,
		Action: action,
	})
	return m
}

// WithChoiceMessage は選択肢にメッセージを関連付けて追加する
func (m *MessageData) WithChoiceMessage(text string, messageData *MessageData) *MessageData {
	m.Choices = append(m.Choices, Choice{
		Text:        text,
		MessageData: messageData,
	})
	return m
}

// WithOnComplete は完了時のコールバックを設定する
func (m *MessageData) WithOnComplete(callback func()) *MessageData {
	m.OnComplete = callback
	return m
}

// DialogMessage は会話メッセージを連鎖
func (m *MessageData) DialogMessage(text, speaker string) *MessageData {
	m.NextMessages = append(m.NextMessages, NewDialogMessage(text, speaker))
	return m
}

// SystemMessage はシステムメッセージを連鎖
func (m *MessageData) SystemMessage(text string) *MessageData {
	m.NextMessages = append(m.NextMessages, NewSystemMessage(text))
	return m
}

// HasNextMessages は次のメッセージがあるかを確認
func (m *MessageData) HasNextMessages() bool {
	return len(m.NextMessages) > 0
}

// GetNextMessages は次のメッセージ群を取得
func (m *MessageData) GetNextMessages() []*MessageData {
	return m.NextMessages
}

// ensureCurrentLine は現在の行が存在することを保証する
func (m *MessageData) ensureCurrentLine() {
	if len(m.TextSegmentLines) == 0 {
		m.TextSegmentLines = append(m.TextSegmentLines, []TextSegment{})
	}
}

// AddText は通常テキストを追加する
func (m *MessageData) AddText(text string) *MessageData {
	m.ensureCurrentLine()
	currentLineIdx := len(m.TextSegmentLines) - 1
	m.TextSegmentLines[currentLineIdx] = append(m.TextSegmentLines[currentLineIdx], TextSegment{Text: text})
	return m
}

// AddNewLine は改行を追加する（新しい行を作成）
func (m *MessageData) AddNewLine() *MessageData {
	m.TextSegmentLines = append(m.TextSegmentLines, []TextSegment{})
	return m
}

// AddKeyword はキーワード（赤色背景）テキストを追加する
func (m *MessageData) AddKeyword(text string) *MessageData {
	m.ensureCurrentLine()
	importantColor := color.RGBA{255, 100, 100, 255}
	importantBgColor := color.RGBA{80, 20, 20, 180}
	currentLineIdx := len(m.TextSegmentLines) - 1
	m.TextSegmentLines[currentLineIdx] = append(m.TextSegmentLines[currentLineIdx], TextSegment{
		Text:            text,
		Color:           &importantColor,
		BackgroundColor: &importantBgColor,
	})
	return m
}
