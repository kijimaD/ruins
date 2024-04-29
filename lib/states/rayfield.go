package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raycast"
)

type RayFieldState struct {
	Game raycast.Game
}

func (st RayFieldState) String() string {
	return "RayField"
}

// State interface ================

func (st *RayFieldState) OnPause(world w.World) {}

func (st *RayFieldState) OnResume(world w.World) {}

func (st *RayFieldState) OnStart(world w.World) {
	// 画像に数バイトの誤差が出て、VRTで失敗しているようなので
	// world.Resources.ScreenDimensions.Width / 4
	// world.Resources.ScreenDimensions.Height / 4
	st.Game.Px = 20
	st.Game.Py = 20
	st.Game.ScreenWidth = world.Resources.ScreenDimensions.Width
	st.Game.ScreenHeight = world.Resources.ScreenDimensions.Height
	st.Game.Prepare()
}

func (st *RayFieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *RayFieldState) Update(world w.World) states.Transition {
	st.Game.Update()

	return states.Transition{}
}

func (st *RayFieldState) Draw(world w.World, screen *ebiten.Image) {
	st.Game.Draw(screen)
}
