package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/states"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
	"github.com/kijimaD/sokotwo/lib/resources"
)

type FieldState struct{}

// State interface ================

func (st *FieldState) OnPause(world w.World) {}

func (st *FieldState) OnResume(world w.World) {}

func (st *FieldState) OnStart(world w.World) {
	packageData := utils.Try(gloader.LoadPackage("forest"))

	// prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	// loader.AddEntities(world, prefabs.Field)

	world.Resources.Game = &resources.Game{Package: packageData}
}

func (st *FieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *FieldState) Update(world w.World) states.Transition {
	return states.Transition{}
}
