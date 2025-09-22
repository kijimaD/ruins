// Package states はゲームの導入テキストを表示するステート
package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/typewriter"
	w "github.com/kijimaD/ruins/lib/world"
)

// IntroState はイントロのゲームステート
type IntroState struct {
	es.BaseState[w.World]
	ui            *ebitenui.UI
	currentText   string
	currentIndex  int
	texts         []string
	bg            *ebiten.Image
	keyboardInput input.KeyboardInput

	// typewriter関連フィールド
	messageHandler   *typewriter.MessageHandler
	uiBuilder        *typewriter.MessageUIBuilder
	messageContainer *widget.Container
}

func (st IntroState) String() string {
	return "Intro"
}

// 客観的に状況、世界観、使命を語らせる
var introTexts = []string{
	"「『虚脱症』の患者がまた一人運ばれてきました。」",
	"「どうして急に増えているんでしょうね、この病気。」",
	"「原因も治療法もさっぱりだ。」",
	"「あの少年、毎日来ているそうですね。」",
	"「母親の病気を治すために\n『遺跡』に挑もうとしているらしい。」",
	"「『珠』の伝説を信じているんでしょうか。」",
	"「...何もかも戦争で失われた後だ。\n皆、すがるものを求めているのでしょう。」",
}

var introBgImages = []string{
	"bg_urban1",
	"bg_urban1",
	"bg_urban1",
	"bg_crystal1",
	"bg_crystal1",
	"bg_jungle1",
	"bg_jungle1",
}

// State interface ================

var _ es.State[w.World] = &IntroState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *IntroState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *IntroState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *IntroState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// 初期化
	st.texts = introTexts
	st.currentIndex = 0
	st.currentText = ""

	// 最初の背景を設定
	if len(introBgImages) > 0 {
		spriteSheet := (*world.Resources.SpriteSheets)[introBgImages[0]]
		st.bg = spriteSheet.Texture.Image
	}

	// MessageHandlerを初期化
	st.messageHandler = typewriter.NewMessageHandler(typewriter.DialogConfig(), st.keyboardInput)

	// UIBuilderを初期化
	res := world.Resources.UIResources
	uiConfig := typewriter.DefaultUIConfig()
	uiConfig.TextFace = (*world.Resources.Faces)["kappa"]
	uiConfig.TextColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	if res != nil {
		uiConfig.ArrowImage = res.ComboButton.Graphic
	}
	st.uiBuilder = typewriter.NewMessageUIBuilder(st.messageHandler, uiConfig)

	// コールバックを設定
	st.messageHandler.SetOnComplete(func() bool {
		// 次のテキストに進む
		st.currentIndex++
		if st.currentIndex < len(st.texts) {
			// 背景を更新
			if st.currentIndex < len(introBgImages) {
				spriteSheet := (*world.Resources.SpriteSheets)[introBgImages[st.currentIndex]]
				st.bg = spriteSheet.Texture.Image
			}
			// 次のメッセージを開始
			st.messageHandler.Start(st.texts[st.currentIndex])
			return false // まだ完了していない
		}
		return true // 全て完了
	})

	// UIを初期化
	st.ui = st.initUI(world)

	// 最初のメッセージを開始
	if len(st.texts) > 0 {
		st.messageHandler.Start(st.texts[0])
	}
}

// OnStop はステートが停止される際に呼ばれる
func (st *IntroState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *IntroState) Update(world w.World) es.Transition[w.World] {
	// Escapeキーでスキップ
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}}
	}

	// typewriter更新（入力処理も含む）
	shouldComplete := st.messageHandler.Update()
	if shouldComplete {
		// 全てのテキストが完了
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewMainMenuState}}
	}

	// UIBuilderが存在する場合はUI更新
	if st.uiBuilder != nil {
		st.uiBuilder.Update()
		// UIBuilderの更新後、UIを再作成
		st.ui = st.initUI(world)
	}

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *IntroState) Draw(_ w.World, screen *ebiten.Image) {
	// ebitenui で背景をいい感じにするにはどうすればよいのだろう
	opts := &ebiten.DrawImageOptions{}
	if st.bg != nil {
		screen.DrawImage(st.bg, opts)
	}
	st.ui.Draw(screen)
}

// ================

func (st *IntroState) initUI(world w.World) *ebitenui.UI {
	// 画面幅を取得
	screenWidth := world.Resources.ScreenDimensions.Width

	// AnchorLayoutで縦中央配置を実現
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// GridLayoutコンテナを縦方向少し上に配置
	gridContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{true}),
			widget.GridLayoutOpts.Spacing(0, 0),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
			}),
		),
	)

	// 左スペーサー（固定幅）
	leftSpacer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(screenWidth/8, 0), // 画面の1/8
		),
	)

	// 中央コンテナ（伸縮）
	st.messageContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)

	// 右スペーサー（固定幅）
	rightSpacer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(screenWidth/8, 0), // 画面の1/8
		),
	)

	gridContainer.AddChild(leftSpacer)
	gridContainer.AddChild(st.messageContainer)
	gridContainer.AddChild(rightSpacer)

	rootContainer.AddChild(gridContainer)
	st.updateMessageContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *IntroState) updateMessageContainer(world w.World) {
	st.messageContainer.RemoveChildren()

	// UIBuilderが存在する場合はそのコンテナを使用（message stateと同様）
	if st.uiBuilder != nil {
		typewriterContainer := st.uiBuilder.GetContainer()
		if typewriterContainer != nil {
			st.messageContainer.AddChild(typewriterContainer)
			return
		}
	}

	// フォールバック: 従来のテキスト表示
	textWidget := widget.NewText(
		widget.TextOpts.Text(st.currentText, (*world.Resources.Faces)["kappa"], color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	)

	st.messageContainer.AddChild(textWidget)
}
