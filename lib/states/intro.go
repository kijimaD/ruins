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
	"github.com/kijimaD/ruins/lib/msg"
	w "github.com/kijimaD/ruins/lib/world"
)

// IntroState はイントロのゲームステート
type IntroState struct {
	es.BaseState
	ui            *ebitenui.UI
	queue         msg.Queue
	cycle         int
	bg            *ebiten.Image
	keyboardInput input.KeyboardInput

	messageContainer *widget.Container
}

func (st IntroState) String() string {
	return "Intro"
}

var introText = `
[image source="bg_urban1"]

遺跡。[p]
粗末な装備で怪物と財宝に満ちた遺跡に挑み、[p]

[image source="bg_crystal1"]
[wait time="500"]

得られたささいな品から生活を発展させた。[p]
国家の成立と工業の進展とともに、[l]
遺跡から得られる技術が利用できることがわかってくると、[p]
しばしば支配権をめぐって戦争が行われるようになった。[p]`

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
	st.queue = msg.NewQueueFromText(introText)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *IntroState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *IntroState) Update(world w.World) es.Transition {
	var queueResult msg.QueueState

	if v, ok := st.queue.Head().(*msg.ChangeBg); ok {
		spriteSheet := (*world.Resources.SpriteSheets)[v.Source]
		st.bg = spriteSheet.Texture.Image
	}

	if st.cycle%2 == 0 {
		queueResult = st.queue.RunHead()
		st.cycle = 0
	}
	st.cycle++

	switch {
	case st.keyboardInput.IsEnterJustPressedOnce():
		queueResult = st.queue.Pop()
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		queueResult = st.queue.Pop()
	case st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape):
		// debug
		return es.Transition{Type: es.TransSwitch, NewStates: []es.State{&MainMenuState{}}}
	}

	switch queueResult {
	case msg.QueueStateFinish:
		return es.Transition{Type: es.TransSwitch, NewStates: []es.State{&MainMenuState{}}}
	}

	st.updateMessageContainer(world)
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
	st.messageContainer.AddChild(eui.NewMenuText(st.queue.Display(), world))
}
