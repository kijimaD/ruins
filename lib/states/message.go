package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
)

// FIXME: 最後のpopが行われたときに、遷移先でもenterが押された扱いになる...
// 最後のenterを押す → 元のstateに戻る → 遷移先でenterが押される
type MessageState struct {
	ui            *ebitenui.UI
	trans         *states.Transition
	keyboardInput input.KeyboardInput

	text     string
	textFunc *func() string
}

func (st MessageState) String() string {
	return "Message"
}

// State interface ================

var _ es.State = &MessageState{}

func (st *MessageState) OnPause(world w.World) {}

func (st *MessageState) OnResume(world w.World) {}

func (st *MessageState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
}

func (st *MessageState) OnStop(world w.World) {}

func (st *MessageState) Update(world w.World) states.Transition {
	st.ui = st.reloadUI(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	if st.keyboardInput.IsEnterJustPressedOnce() {
		return states.Transition{Type: states.TransPop}
	}

	if st.textFunc != nil {
		f := *st.textFunc
		st.text = f()
		st.textFunc = nil
	}

	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *MessageState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *MessageState) reloadUI(world w.World) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5)))),
	)
	res := world.Resources.UIResources
	text := widget.NewText(

		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.TextOpts.Text(st.text, res.Text.Face, styles.TextColor),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	)
	rootContainer.AddChild(text)

	return &ebitenui.UI{Container: rootContainer}
}
