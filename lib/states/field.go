package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
)

type FieldState struct{}

// State interface ================

func (st *FieldState) OnPause(world w.World) {}

func (st *FieldState) OnResume(world w.World) {}

func (st *FieldState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Field)
}

func (st *FieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *FieldState) Update(world w.World) states.Transition {
	return states.Transition{}
}
