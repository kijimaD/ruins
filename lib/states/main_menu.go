package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type MainMenuState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	mainMenuContainer *widget.Container
}

func (st MainMenuState) String() string {
	return "MainMenu"
}

// State interface ================

func (st *MainMenuState) OnPause(world w.World) {}

func (st *MainMenuState) OnResume(world w.World) {}

func (st *MainMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *MainMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *MainMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}

	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *MainMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *MainMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.mainMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(eui.NewBodyText("Ruins", color.RGBA{255, 255, 255, 255}, world))
	rootContainer.AddChild(st.mainMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *MainMenuState) updateMenuContainer(world w.World) {
	st.mainMenuContainer.RemoveChildren()

	for _, data := range mainMenuTrans {
		data := data
		btn := eui.NewItemButton(
			data.label,
			func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			},
			world,
		)
		st.mainMenuContainer.AddChild(btn)
	}
}

var mainMenuTrans = []struct {
	label string
	trans states.Transition
}{
	{
		label: "イントロ",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&IntroState{}}},
	},
	{
		label: "拠点メニュー",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}},
	},
	{
		label: "raycast field(実装中)",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&RayFieldState{}}},
	},
	{
		label: "終了",
		trans: states.Transition{Type: states.TransQuit},
	},
}
