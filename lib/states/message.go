package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/styles"
)

type MessageState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	text     string
	textFunc *func() string
}

func (st MessageState) String() string {
	return "Message"
}

// State interface ================

func (st *MessageState) OnPause(world w.World) {}

func (st *MessageState) OnResume(world w.World) {}

func (st *MessageState) OnStart(world w.World) {}

func (st *MessageState) OnStop(world w.World) {}

func (st *MessageState) Update(world w.World) states.Transition {
	st.ui = st.reloadUI(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return states.Transition{Type: states.TransPop}
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
