package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/typewriter"
	w "github.com/kijimaD/ruins/lib/world"
)

// MessageState はメッセージ表示用のステート
type MessageState struct {
	es.BaseState
	keyboardInput input.KeyboardInput

	text     string
	textFunc *func() string

	// タイプライター機能
	messageHandler *typewriter.MessageHandler
	uiBuilder      *typewriter.MessageUIBuilder
}

func (st MessageState) String() string {
	return "Message"
}

// State interface ================

var _ es.State = &MessageState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MessageState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MessageState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *MessageState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// タイプライター初期化
	if st.messageHandler == nil {
		// MessageHandlerを初期化
		st.messageHandler = typewriter.NewMessageHandler(typewriter.BattleConfig(), st.keyboardInput)

		// UIBuilderを初期化
		res := world.Resources.UIResources
		uiConfig := typewriter.DefaultUIConfig()
		uiConfig.TextFace = res.Text.Face
		uiConfig.TextColor = styles.TextColor
		uiConfig.PanelImage = res.Panel.Image
		uiConfig.ArrowImage = res.ComboButton.Graphic
		st.uiBuilder = typewriter.NewMessageUIBuilder(st.messageHandler, uiConfig)

		if st.text != "" {
			st.messageHandler.Start(st.text)
		}
	}
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageState) OnStop(_ w.World) {}

// Update はメッセージステートの更新処理を行う
func (st *MessageState) Update(world w.World) es.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransQuit}
	}

	// タイプライター処理
	if st.messageHandler != nil {
		// MessageHandlerに処理を委譲し、完了状態をチェック
		shouldComplete := st.messageHandler.Update()
		if shouldComplete {
			return es.Transition{Type: es.TransPop}
		}
	}

	// textFunc による動的テキスト更新
	if st.textFunc != nil {
		f := *st.textFunc
		newText := f()
		st.textFunc = nil

		// 新しいテキストの場合
		if newText != st.text {
			st.text = newText

			if st.messageHandler != nil {
				st.messageHandler.Start(st.text)
			}
		}
	}

	// UIBuilderが存在する場合はUI更新
	if st.uiBuilder != nil {
		st.uiBuilder.Update()
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageState) Draw(_ w.World, screen *ebiten.Image) {
	if st.uiBuilder != nil && st.uiBuilder.GetUI() != nil {
		st.uiBuilder.GetUI().Draw(screen)
	}
}
