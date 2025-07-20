package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	w "github.com/kijimaD/ruins/lib/world"
)

// GameOverState はゲームオーバーのゲームステート
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

// OnPause はステートが一時停止される際に呼ばれる
func (st *GameOverState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *GameOverState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *GameOverState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	bg := (*world.Resources.SpriteSheets)["bg_explosion1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *GameOverState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *GameOverState) Update(_ w.World) es.Transition {
	if st.keyboardInput.IsEnterJustPressedOnce() {
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
	}

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *GameOverState) Draw(_ w.World, screen *ebiten.Image) {
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
