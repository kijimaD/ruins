package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
)

type GameOverState struct {
	es.BaseState
	ui            *ebitenui.UI
	keyboardInput input.KeyboardInput

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
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	bg := (*world.Resources.SpriteSheets)["bg_explosion1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

func (st *GameOverState) OnStop(world w.World) {}

func (st *GameOverState) Update(world w.World) es.Transition {
	if st.keyboardInput.IsEnterJustPressedOnce() {
		return es.Transition{Type: es.TransSwitch, NewStates: []es.State{&MainMenuState{}}}
	}

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
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
