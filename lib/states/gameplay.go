package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

type GamePlayState struct{}

// State interface ================

func (st *GamePlayState) OnPause(world w.World) {}

func (st *GamePlayState) OnResume(world w.World) {}

func (st *GamePlayState) OnStart(world w.World) {}

func (st *GamePlayState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *GamePlayState) Update(world w.World) states.Transition {
	return states.Transition{}
}
