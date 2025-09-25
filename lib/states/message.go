package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
	w "github.com/kijimaD/ruins/lib/world"
)

// MessageState はメッセージを表示する専用ステート
type MessageState struct {
	es.BaseState[w.World]
	messageData   *messagedata.MessageData
	messageWindow *messagewindow.Window
}

func (st MessageState) String() string {
	return "Message"
}

var _ es.State[w.World] = &MessageState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageState) OnStart(world w.World) {
	// メッセージデータからキュー対応メッセージウィンドウを構築
	st.messageWindow = messagewindow.NewBuilder(world).Build(st.messageData)
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *MessageState) Update(_ w.World) es.Transition[w.World] {
	if st.messageWindow != nil {
		st.messageWindow.Update()

		if st.messageWindow.IsClosed() {
			// TransitionFactoryが設定されている場合はそれを使用
			if st.messageData != nil && st.messageData.TransitionFactory != nil {
				return st.messageData.TransitionFactory()
			}
			// デフォルトはステートをポップ
			return es.Transition[w.World]{Type: es.TransPop}
		}
	}

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageState) Draw(_ w.World, screen *ebiten.Image) {
	if st.messageWindow != nil {
		st.messageWindow.Draw(screen)
	}
}
