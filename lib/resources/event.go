package resources

import (
	ecs "github.com/x-hgg-x/goecs/v2"
)

// StateEventType はイベントの種類を表す
type StateEventType string

const (
	// StateEventTypeNone はイベントなしを表す
	StateEventTypeNone = StateEventType("NONE")
	// StateEventTypeWarpNext は次の階層への移動を表す
	StateEventTypeWarpNext = StateEventType("WARP_NEXT")
	// StateEventTypeWarpEscape は脱出を表す
	StateEventTypeWarpEscape = StateEventType("WARP_ESCAPE")
	// StateEventTypeShowDialog は会話メッセージの表示を表す
	StateEventTypeShowDialog = StateEventType("SHOW_DIALOG")
)

// StateEvent はフィールド上でのイベント。ステート遷移が発生する
type StateEvent interface {
	Type() StateEventType
}

// NoneEvent はイベントなしを表す
type NoneEvent struct{}

// Type はイベントタイプを返す
func (e NoneEvent) Type() StateEventType {
	return StateEventTypeNone
}

// WarpNextEvent は次の階層への移動を表す
type WarpNextEvent struct{}

// Type はイベントタイプを返す
func (e WarpNextEvent) Type() StateEventType {
	return StateEventTypeWarpNext
}

// WarpEscapeEvent は脱出を表す
type WarpEscapeEvent struct{}

// Type はイベントタイプを返す
func (e WarpEscapeEvent) Type() StateEventType {
	return StateEventTypeWarpEscape
}

// ShowDialogEvent は会話メッセージの表示を表す
type ShowDialogEvent struct {
	MessageKey    string
	SpeakerEntity ecs.Entity
}

// Type はイベントタイプを返す
func (e ShowDialogEvent) Type() StateEventType {
	return StateEventTypeShowDialog
}
