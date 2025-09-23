// Package messagedata はメッセージウィンドウに表示するデータ構造を提供する
package messagedata

import (
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
)

// MessageData はメッセージウィンドウに表示するデータ
type MessageData struct {
	Text         string
	Speaker      string
	Type         messagewindow.MessageType
	Choices      []Choice
	Size         *Size
	OnComplete   func()         // メッセージ完了時のコールバック
	NextMessages []*MessageData // 次に表示するメッセージ群
}

// Choice は選択肢のデータ
type Choice struct {
	Text        string
	Description string
	Action      func()
	MessageData *MessageData // 選択肢を選んだ時に表示するメッセージ
	Disabled    bool
}

// Size はカスタムサイズ
type Size struct {
	Width  int
	Height int
}

// NewDialogMessage は会話メッセージを作成する
func NewDialogMessage(text, speaker string) *MessageData {
	return &MessageData{
		Text:    text,
		Speaker: speaker,
		Type:    messagewindow.TypeDialog,
	}
}

// NewSystemMessage はシステムメッセージを作成する
func NewSystemMessage(text string) *MessageData {
	return &MessageData{
		Text: text,
		Type: messagewindow.TypeSystem,
	}
}

// NewEventMessage はイベントメッセージを作成する
func NewEventMessage(text string) *MessageData {
	return &MessageData{
		Text: text,
		Type: messagewindow.TypeEvent,
	}
}

// WithSpeaker は話者を設定する
func (m *MessageData) WithSpeaker(speaker string) *MessageData {
	m.Speaker = speaker
	return m
}

// WithSize はカスタムサイズを設定する
func (m *MessageData) WithSize(width, height int) *MessageData {
	m.Size = &Size{Width: width, Height: height}
	return m
}

// WithChoice は選択肢を追加する
func (m *MessageData) WithChoice(text string, action func()) *MessageData {
	m.Choices = append(m.Choices, Choice{
		Text:   text,
		Action: action,
	})
	return m
}

// WithChoiceDescription は説明付き選択肢を追加する
func (m *MessageData) WithChoiceDescription(text, description string, action func()) *MessageData {
	m.Choices = append(m.Choices, Choice{
		Text:        text,
		Description: description,
		Action:      action,
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

// WithChoiceMessageDescription は説明付き選択肢にメッセージを関連付けて追加する
func (m *MessageData) WithChoiceMessageDescription(text, description string, messageData *MessageData) *MessageData {
	m.Choices = append(m.Choices, Choice{
		Text:        text,
		Description: description,
		MessageData: messageData,
	})
	return m
}

// WithOnComplete は完了時のコールバックを設定する
func (m *MessageData) WithOnComplete(callback func()) *MessageData {
	m.OnComplete = callback
	return m
}

// MessageSequence は連続するメッセージのシーケンス
type MessageSequence struct {
	messages []*MessageData
}

// NewMessageSequence は新しいメッセージシーケンスを作成
func NewMessageSequence() *MessageSequence {
	return &MessageSequence{
		messages: make([]*MessageData, 0),
	}
}

// DialogMessage は会話メッセージを追加
func (ms *MessageSequence) DialogMessage(text, speaker string) *MessageSequence {
	ms.messages = append(ms.messages, NewDialogMessage(text, speaker))
	return ms
}

// SystemMessage はシステムメッセージを追加
func (ms *MessageSequence) SystemMessage(text string) *MessageSequence {
	ms.messages = append(ms.messages, NewSystemMessage(text))
	return ms
}

// EventMessage はイベントメッセージを追加
func (ms *MessageSequence) EventMessage(text string) *MessageSequence {
	ms.messages = append(ms.messages, NewEventMessage(text))
	return ms
}

// GetMessages は全メッセージを取得
func (ms *MessageSequence) GetMessages() []*MessageData {
	return ms.messages
}

// MessageDataのチェーンメソッド

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

// EventMessage はイベントメッセージを連鎖
func (m *MessageData) EventMessage(text string) *MessageData {
	m.NextMessages = append(m.NextMessages, NewEventMessage(text))
	return m
}

// Sequence はメッセージシーケンスを連鎖
func (m *MessageData) Sequence(sequence *MessageSequence) *MessageData {
	m.NextMessages = append(m.NextMessages, sequence.GetMessages()...)
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
