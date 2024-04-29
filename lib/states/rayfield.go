package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

type RayFieldState struct{}

func (st RayFieldState) String() string {
	return "RayField"
}

// State interface ================

func (st *RayFieldState) OnPause(world w.World) {}

func (st *RayFieldState) OnResume(world w.World) {}

func (st *RayFieldState) OnStart(world w.World) {
}

func (st *RayFieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *RayFieldState) Update(world w.World) states.Transition {
	return states.Transition{}
}

func (st *RayFieldState) Draw(world w.World, screen *ebiten.Image) {}
