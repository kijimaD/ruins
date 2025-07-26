// Package states はゲームの導入テキストを表示するステート
package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/typewriter"
	w "github.com/kijimaD/ruins/lib/world"
)

// IntroState はイントロのゲームステート
type IntroState struct {
	es.BaseState
	ui            *ebitenui.UI
	currentText   string
	currentIndex  int
	texts         []string
	bg            *ebiten.Image
	keyboardInput input.KeyboardInput

	// typewriter関連フィールド
	messageHandler   *typewriter.MessageHandler
	messageContainer *widget.Container
}

func (st IntroState) String() string {
	return "Intro"
}

var introTexts = []string{
	"戦争が終わり、街には復興の槌音が響く。",
	"古い言い伝えによると、地下深くに眠る遺跡の最下層には珠があり、どんな願いも叶えるとされる。",
	"多くの人は迷信だと笑うが、遺跡の不思議な技術を見れば、完全に否定することもできない。",
	"母が倒れてから、もう三ヶ月になる。",
	"虚脱症。原因不明の病気で、現代医学では治療法が確立されていない。",
	"一部では『珠の力で治る』という話もあるが、医学界では相手にされていない。",
	"それでも、俺には他に方法がない。",
	"探索者登録番号二八四七、十七歳男性。",
	"目的：遺跡探索および珠の回収。",
	"母さん、必ず帰る。",
}

var introBgImages = []string{
	"bg_urban1",
	"bg_urban1",
	"bg_urban1",
	"bg_crystal1",
	"bg_crystal1",
	"bg_crystal1",
	"bg_crystal1",
	"bg_jungle1",
	"bg_jungle1",
	"bg_jungle1",
}

// State interface ================

var _ es.State = &IntroState{}

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

	// コールバックを設定
	st.messageHandler.SetOnUpdateUI(func(text string) {
		st.currentText = text
		st.updateMessageContainer(world)
	})

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
func (st *IntroState) Update(world w.World) es.Transition {
	// Escapeキーでスキップ
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
	}

	// マウスクリックでスキップ
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// 現在のメッセージを完了させて次に進む
		st.currentIndex++
		if st.currentIndex < len(st.texts) {
			// 背景を更新
			if st.currentIndex < len(introBgImages) {
				spriteSheet := (*world.Resources.SpriteSheets)[introBgImages[st.currentIndex]]
				st.bg = spriteSheet.Texture.Image
			}
			// 次のメッセージを開始
			st.messageHandler.Start(st.texts[st.currentIndex])
		} else {
			// 全て完了
			return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
		}
	}

	// typewriter更新
	shouldComplete := st.messageHandler.Update()
	if shouldComplete {
		// 全てのテキストが完了
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
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
	rootContainer := eui.NewRowContainer()
	st.messageContainer = eui.NewRowContainer()
	rootContainer.AddChild(st.messageContainer)

	st.updateMessageContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *IntroState) updateMessageContainer(world w.World) {
	st.messageContainer.RemoveChildren()
	st.messageContainer.AddChild(eui.NewMenuText(st.currentText, world))
}
