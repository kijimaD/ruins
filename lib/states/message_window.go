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
	messageData   *messagedata.MessageData
	messageWindow *messagewindow.Window
}

func (st MessageWindowState) String() string {
	return "MessageWindow"
}

var _ es.State[w.World] = &MessageWindowState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageWindowState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageWindowState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageWindowState) OnStart(world w.World) {
	// メッセージデータからキュー対応メッセージウィンドウを構築
	st.messageWindow = messagewindow.NewBuilder(world).Build(st.messageData)
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageWindowState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *MessageWindowState) Update(_ w.World) es.Transition[w.World] {
	if st.messageWindow != nil {
		st.messageWindow.Update()

		// メッセージウィンドウが閉じられた場合はステートをポップ
		if st.messageWindow.IsClosed() {
			return es.Transition[w.World]{Type: es.TransPop}
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
