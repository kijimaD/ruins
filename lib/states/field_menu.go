package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type FieldMenuState struct {
	selection int

	ui                 *ebitenui.UI
	trans              *states.Transition
	fieldMenuContainer *widget.Container
}

func (st FieldMenuState) String() string {
	return "FieldMenu"
}

// State interface ================

func (st *FieldMenuState) OnPause(world w.World) {}

func (st *FieldMenuState) OnResume(world w.World) {}

func (st *FieldMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *FieldMenuState) OnStop(world w.World) {}

func (st *FieldMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}

	st.ui.Update()
	if st.trans != nil {
		next := *st.trans
		st.trans = nil

		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *FieldMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *FieldMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalTransContainer()
	st.fieldMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(st.fieldMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *FieldMenuState) updateMenuContainer(world w.World) {
	st.fieldMenuContainer.RemoveChildren()

	for _, data := range fieldMenuTrans {
		data := data
		btn := eui.NewItemButton(
			data.label,
			func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			},
			world,
		)
		st.fieldMenuContainer.AddChild(btn)
	}
}

var fieldMenuTrans = []struct {
	label string
	trans states.Transition
}{
	{
		label: "閉じる",
		trans: states.Transition{Type: states.TransPop},
	},
	{
		label: "終了",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}},
	},
}
