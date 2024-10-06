package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type BattleState struct {
	ui    *ebitenui.UI
	trans *states.Transition
}

func (st BattleState) String() string {
	return "Battle"
}

// State interface ================

func (st *BattleState) OnPause(world w.World) {}

func (st *BattleState) OnResume(world w.World) {}

func (st *BattleState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *BattleState) OnStop(world w.World) {}

func (st *BattleState) Update(world w.World) states.Transition {
	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	st.ui.Update()

	return states.Transition{Type: states.TransNone}
}

func (st *BattleState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *BattleState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalTransContainer()

	return &ebitenui.UI{Container: rootContainer}
}
