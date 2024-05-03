package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raycast"
	"github.com/kijimaD/ruins/lib/spawner"
	gs "github.com/kijimaD/ruins/lib/systems"
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
	// VRTすると画像に数バイトの誤差が出て失敗しているようなのでこの位置
	st.Game.Px = world.Resources.ScreenDimensions.Width / 2
	st.Game.Py = world.Resources.ScreenDimensions.Height / 2
	st.Game.ScreenWidth = world.Resources.ScreenDimensions.Width
	st.Game.ScreenHeight = world.Resources.ScreenDimensions.Height
	st.Game.Prepare()

	spawner.SpawnPlayer(world, 200, 200)
}

func (st *RayFieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *RayFieldState) Update(world w.World) states.Transition {
	st.Game.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&FieldMenuState{}}}
	}

	return states.Transition{}
}

func (st *RayFieldState) Draw(world w.World, screen *ebiten.Image) {
	st.Game.Draw(screen)
	gs.RenderObjectSystem(world, screen)
}
