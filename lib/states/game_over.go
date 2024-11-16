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
	"github.com/kijimaD/ruins/lib/styles"
)

type GameOverState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	// 背景
	bg *ebiten.Image
}

func (st GameOverState) String() string {
	return "GameOver"
}

// State interface ================

var _ es.State = &GameOverState{}

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

	res := world.Resources.UIResources
	rootContainer.AddChild(widget.NewText(widget.TextOpts.Text("GAME OVER...", res.Text.BigTitleFace, styles.TextColor)))

	return &ebitenui.UI{Container: rootContainer}
}
