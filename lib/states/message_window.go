package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
	w "github.com/kijimaD/ruins/lib/world"
)

// MessageWindowState はメッセージウィンドウを表示する専用ステート
type MessageWindowState struct {
	es.BaseState[w.World]
	messageData       *messagedata.MessageData
	messageWindow     *messagewindow.Window
	messageQueue      []*messagedata.MessageData // 連鎖メッセージキュー
	pendingTransition *es.Transition[w.World]    // 選択肢から設定される遷移
}

func (st MessageWindowState) String() string {
	return "MessageWindow"
}

var _ es.State[w.World] = &MessageWindowState{}

// NewMessageWindowState はメッセージデータを受け取って新しいMessageWindowStateを作成する
func NewMessageWindowState(messageData *messagedata.MessageData) *MessageWindowState {
	return &MessageWindowState{
		messageData: messageData,
	}
}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageWindowState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageWindowState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageWindowState) OnStart(world w.World) {
	// 連鎖メッセージをキューに追加
	if st.messageData.HasNextMessages() {
		st.messageQueue = append(st.messageQueue, st.messageData.GetNextMessages()...)
	}

	// メッセージデータからメッセージウィンドウを構築
	st.buildMessageWindow(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageWindowState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *MessageWindowState) Update(world w.World) es.Transition[w.World] {
	if st.messageWindow != nil {
		st.messageWindow.Update()

		// メッセージウィンドウが閉じられた場合
		if st.messageWindow.IsClosed() {
			// 選択肢からの遷移が保留されている場合は、それを優先
			if st.pendingTransition != nil {
				transition := *st.pendingTransition
				st.pendingTransition = nil
				return transition
			}
			return st.handleMessageClosed(world)
		}
	}

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageWindowState) Draw(_ w.World, screen *ebiten.Image) {
	if st.messageWindow != nil {
		st.messageWindow.Draw(screen)
	}
}

// buildMessageWindow はメッセージデータからメッセージウィンドウを構築する
func (st *MessageWindowState) buildMessageWindow(world w.World) {
	builder := messagewindow.NewBuilder(world).
		Message(st.messageData.Text)

	// 話者が設定されている場合
	if st.messageData.Speaker != "" {
		builder = builder.Speaker(st.messageData.Speaker)
	}

	// 選択肢が設定されている場合
	for _, choice := range st.messageData.Choices {
		// 選択肢のコピーを作成（クロージャのキャプチャ問題を回避）
		choiceCopy := choice

		// アクション関数を構築
		var actionFunc func()
		if choiceCopy.MessageData != nil {
			// メッセージデータが設定されている場合は新しいMessageWindowStateをpush
			actionFunc = func() {
				if choiceCopy.Action != nil {
					choiceCopy.Action()
				}
				// MessageWindowStateを新しくpush
				transition := es.Transition[w.World]{
					Type: es.TransPush,
					NewStateFuncs: []es.StateFactory[w.World]{
						func() es.State[w.World] { return NewMessageWindowState(choiceCopy.MessageData) },
					},
				}
				st.pendingTransition = &transition
			}
		} else {
			// 通常のActionのみの場合
			actionFunc = choiceCopy.Action
		}

		builder = builder.Choice(choiceCopy.Text, actionFunc)
	}

	st.messageWindow = builder.Build()
}

// handleMessageClosed はメッセージが閉じられた時の処理
func (st *MessageWindowState) handleMessageClosed(world w.World) es.Transition[w.World] {
	// 完了コールバックがあれば実行
	if st.messageData.OnComplete != nil {
		st.messageData.OnComplete()
	}

	// キューに次のメッセージがある場合は表示
	if len(st.messageQueue) > 0 {
		nextMessage := st.messageQueue[0]
		st.messageQueue = st.messageQueue[1:]

		// 次のメッセージの連鎖メッセージもキューに追加
		if nextMessage.HasNextMessages() {
			st.messageQueue = append(nextMessage.GetNextMessages(), st.messageQueue...)
		}

		// 次のメッセージを表示
		st.messageData = nextMessage
		st.buildMessageWindow(world)
		return es.Transition[w.World]{Type: es.TransNone}
	}

	// 全てのメッセージが完了した場合、ステートをポップ
	return es.Transition[w.World]{Type: es.TransPop}
}
