// Package messagedata はメッセージウィンドウに表示するデータ構造を提供する
package messagedata

import (
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
	w "github.com/kijimaD/ruins/lib/world"
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
	MessageData *MessageData           // 選択肢を選んだ時に表示するメッセージ
	Transition  es.Transition[w.World] // 選択肢を選んだ時のステート遷移
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

// MessageChain は連続して表示するメッセージのチェーン（選択肢分岐対応）
type MessageChain struct {
	messages     []*MessageData
	currentIndex int
	choiceMap    map[string][]*MessageData // 選択肢テキスト -> 次のメッセージ群
}

// NewMessageChain は新しいメッセージチェーンを作成する
func NewMessageChain() *MessageChain {
	return &MessageChain{
		messages:  make([]*MessageData, 0),
		choiceMap: make(map[string][]*MessageData),
	}
}

// AddMessage はチェーンにメッセージを追加する
func (mc *MessageChain) AddMessage(message *MessageData) *MessageChain {
	mc.messages = append(mc.messages, message)
	return mc
}

// AddChoiceMessage は選択肢付きメッセージと各選択肢の結果メッセージを追加する
func (mc *MessageChain) AddChoiceMessage(questionMessage *MessageData, choiceResults map[string]*MessageData) *MessageChain {
	// 質問メッセージを追加
	mc.messages = append(mc.messages, questionMessage)

	// 各選択肢の結果メッセージを追加
	for choiceText, resultMessage := range choiceResults {
		// 結果メッセージとその連鎖メッセージを収集
		messageChain := mc.collectMessageChain(resultMessage)
		mc.choiceMap[choiceText] = messageChain
	}

	return mc
}

// AddChoiceMessageExtended は選択肢付きメッセージと各選択肢の結果メッセージ群を追加
func (mc *MessageChain) AddChoiceMessageExtended(questionMessage *MessageData, choiceSequences map[string][]*MessageData) *MessageChain {
	// 質問メッセージを追加
	mc.messages = append(mc.messages, questionMessage)

	// 各選択肢のメッセージ群をマップに保存
	for choiceText, messageSequence := range choiceSequences {
		mc.choiceMap[choiceText] = messageSequence
	}

	return mc
}

// collectMessageChain はメッセージとその連鎖メッセージを全て収集する
func (mc *MessageChain) collectMessageChain(message *MessageData) []*MessageData {
	result := []*MessageData{message}
	result = append(result, message.GetNextMessages()...)
	return result
}

// GetFirstMessage は最初のメッセージを取得する
func (mc *MessageChain) GetFirstMessage() *MessageData {
	if len(mc.messages) == 0 {
		return nil
	}
	mc.currentIndex = 0
	return mc.messages[0]
}

// GetNextMessage は次のメッセージを取得する
func (mc *MessageChain) GetNextMessage() *MessageData {
	if mc.currentIndex+1 >= len(mc.messages) {
		return nil
	}
	mc.currentIndex++
	return mc.messages[mc.currentIndex]
}

// GetMessagesForChoice は選択肢に対応するメッセージ群を取得
func (mc *MessageChain) GetMessagesForChoice(choiceText string) []*MessageData {
	if messages, exists := mc.choiceMap[choiceText]; exists {
		return messages
	}
	return nil
}

// HasMoreMessages はまだ表示すべきメッセージがあるかを確認する
func (mc *MessageChain) HasMoreMessages() bool {
	return mc.currentIndex+1 < len(mc.messages)
}

// GetCurrentIndex は現在のメッセージインデックスを取得する
func (mc *MessageChain) GetCurrentIndex() int {
	return mc.currentIndex
}
