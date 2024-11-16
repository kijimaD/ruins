// ゲームの導入テキストを表示するステート
package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/utils/msg"
)

type IntroState struct {
	ui    *ebitenui.UI
	trans *states.Transition
	queue msg.Queue
	cycle int
	bg    *ebiten.Image

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

func (st *IntroState) OnPause(world w.World) {}

func (st *IntroState) OnResume(world w.World) {}

func (st *IntroState) OnStart(world w.World) {
	st.queue = msg.NewQueueFromText(introText)
	st.ui = st.initUI(world)
}

func (st *IntroState) OnStop(world w.World) {}

func (st *IntroState) Update(world w.World) states.Transition {
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
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		queueResult = st.queue.Pop()
	case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
		queueResult = st.queue.Pop()
	case inpututil.IsKeyJustPressed(ebiten.KeyEscape):
		// debug
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}

	switch queueResult {
	case msg.QueueStateFinish:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}

	st.updateMessageContainer(world)
	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *IntroState) Draw(world w.World, screen *ebiten.Image) {
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
