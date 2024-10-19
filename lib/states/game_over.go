package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type GameOverState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	// 背景
	bg *ebiten.Image
}

func (st GameOverState) String() string {
	return "GameOverState"
}

// State interface ================

func (st *GameOverState) OnPause(world w.World) {}

func (st *GameOverState) OnResume(world w.World) {}

func (st *GameOverState) OnStart(world w.World) {
	bg := (*world.Resources.SpriteSheets)["bg_explosion1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

func (st *GameOverState) OnStop(world w.World) {}

func (st *GameOverState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	st.ui.Update()

	return states.Transition{Type: states.TransNone}
}

func (st *GameOverState) Draw(world w.World, screen *ebiten.Image) {
	if st.bg != nil {
		screen.DrawImage(st.bg, &ebiten.DrawImageOptions{})
	}

	st.ui.Draw(screen)
}

// ================

func (st *GameOverState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()

	return &ebitenui.UI{Container: rootContainer}
}
