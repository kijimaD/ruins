package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
)

type DungeonMenuState struct {
	ui                   *ebitenui.UI
	trans                *states.Transition
	dungeonMenuContainer *widget.Container
}

func (st DungeonMenuState) String() string {
	return "DungeonMenu"
}

// State interface ================

var _ es.State = &DungeonMenuState{}

func (st *DungeonMenuState) OnPause(world w.World) {}

func (st *DungeonMenuState) OnResume(world w.World) {}

func (st *DungeonMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *DungeonMenuState) OnStop(world w.World) {}

func (st *DungeonMenuState) Update(world w.World) states.Transition {
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

func (st *DungeonMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *DungeonMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)
	st.dungeonMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(st.dungeonMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *DungeonMenuState) updateMenuContainer(world w.World) {
	st.dungeonMenuContainer.RemoveChildren()

	for _, data := range dungeonMenuTrans {
		data := data
		btn := eui.NewButton(
			data.label,
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			}),
		)
		st.dungeonMenuContainer.AddChild(btn)
	}
}

var dungeonMenuTrans = []struct {
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
