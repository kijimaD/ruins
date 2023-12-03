package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/states"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
	"github.com/kijimaD/sokotwo/lib/resources"
	gs "github.com/kijimaD/sokotwo/lib/systems"
)

type FieldState struct{}

// State interface ================

func (st *FieldState) OnPause(world w.World) {}

func (st *FieldState) OnResume(world w.World) {}

func (st *FieldState) OnStart(world w.World) {
	packageData := utils.Try(gloader.LoadPackage("forest"))

	world.Resources.Game = &resources.Game{Package: packageData}
	resources.InitLevel(world, 1)
}

func (st *FieldState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *FieldState) Update(world w.World) states.Transition {
	gs.GridTransformSystem(world)
	gs.GridUpdateSystem(world)
	gs.MoveSystem(world)

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventWarpEscape:
		gameResources.StateEvent = resources.StateEventNone
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}
	}

	return states.Transition{}
}
