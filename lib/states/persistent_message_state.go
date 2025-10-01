package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/widgets/messagewindow"
	w "github.com/kijimaD/ruins/lib/world"
)

// PersistentMessageState は選択肢実行後もウィンドウを開いたままにするMessageStateラッパー
type PersistentMessageState struct {
	MessageState
}

func (st PersistentMessageState) String() string {
	return "PersistentMessage"
}

var _ es.State[w.World] = &PersistentMessageState{}

// Update はゲームステートの更新処理を行う
func (st *PersistentMessageState) Update(_ w.World) es.Transition[w.World] {
	// Escapeキーで明示的にPopする
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition[w.World]{Type: es.TransPop}
	}

	if st.messageWindow != nil {
		st.messageWindow.Update()

		if st.messageWindow.IsClosed() {
			// BaseStateで設定された遷移を優先確認
			if transition := st.ConsumeTransition(); transition.Type != es.TransNone {
				return transition
			}
			// PersistentMessageStateは自動Popしない
			return es.Transition[w.World]{Type: es.TransNone}
		}
		// MessageWindowがアクティブな間は何もしない
		return es.Transition[w.World]{Type: es.TransNone}
	}

	return st.ConsumeTransition()
}

// OnResume はステートが再開される際に呼ばれる
func (st *PersistentMessageState) OnResume(world w.World) {
	// メッセージウィンドウを強制的に再構築
	if st.messageData != nil {
		st.messageWindow = messagewindow.NewBuilder(world).Build(st.messageData)
	}
}

// NewPersistentMessageState は永続的なメッセージステートを作成する
func NewPersistentMessageState(messageData *messagedata.MessageData) *PersistentMessageState {
	return &PersistentMessageState{
		MessageState: MessageState{
			messageData: messageData,
		},
	}
}
