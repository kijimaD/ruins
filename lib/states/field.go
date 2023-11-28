package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

type FieldState struct{}

// State interface ================

func (st *FieldState) OnPause(world w.World) {}

func (st *FieldState) OnResume(world w.World) {}

func (st *FieldState) OnStart(world w.World) {}

func (st *FieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *FieldState) Update(world w.World) states.Transition {
	return states.Transition{}
}
