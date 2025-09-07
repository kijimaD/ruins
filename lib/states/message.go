package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/colors"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
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

	// UI
	ui *ebitenui.UI
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
		uiConfig.TextColor = colors.TextColor
		uiConfig.PanelImage = res.Panel.Image
		uiConfig.ArrowImage = res.ComboButton.Graphic
		st.uiBuilder = typewriter.NewMessageUIBuilder(st.messageHandler, uiConfig)

		// 初回UIを作成
		st.ui = st.createUI()

		if st.text != "" {
			st.messageHandler.Start(st.text)
		}
	}
}

// createUI はtypewriterのコンテナを組み込んだUIを作成
func (st *MessageState) createUI() *ebitenui.UI {
	// メッセージエリアの高さを設定
	messageHeight := 100

	// GridLayoutを使用して中央配置
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, false, true}), // 上下伸縮、中央固定
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
	)

	// 上部の空きスペース
	topSpacer := widget.NewContainer()
	rootContainer.AddChild(topSpacer)

	// 中央のメッセージエリア
	if st.uiBuilder != nil {
		centerArea := widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(0, messageHeight),
			),
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Padding(widget.Insets{Left: 20}),
			)),
		)

		typewriterContainer := st.uiBuilder.GetContainer()
		centerArea.AddChild(typewriterContainer)
		rootContainer.AddChild(centerArea)
	}

	// 下部の空きスペース
	bottomSpacer := widget.NewContainer()
	rootContainer.AddChild(bottomSpacer)

	return &ebitenui.UI{Container: rootContainer}
}

// OnStop はステートが停止される際に呼ばれる
func (st *MessageState) OnStop(_ w.World) {}

// Update はメッセージステートの更新処理を行う
func (st *MessageState) Update(_ w.World) es.Transition {
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
		// UIBuilderの更新後、UIを再作成
		st.ui = st.createUI()
	}

	// UI更新
	if st.ui != nil {
		st.ui.Update()
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *MessageState) Draw(_ w.World, screen *ebiten.Image) {
	if st.ui != nil {
		st.ui.Draw(screen)
	}
}
