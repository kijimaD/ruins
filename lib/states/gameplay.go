package states

import (
	"github.com/x-hgg-x/goecsengine/states"
	w "github.com/x-hgg-x/goecsengine/world"
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
