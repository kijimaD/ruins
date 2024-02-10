package systems

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
)

func MoveSystem(world w.World) {
	moveUpAction := world.Resources.InputHandler.Actions[resources.MoveUpAction]
	moveDownAction := world.Resources.InputHandler.Actions[resources.MoveDownAction]
	moveLeftAction := world.Resources.InputHandler.Actions[resources.MoveLeftAction]
	moveRightAction := world.Resources.InputHandler.Actions[resources.MoveRightAction]

	switch {
	case moveUpAction:
		resources.Move(world, resources.MovementUp)
	case moveDownAction:
		resources.Move(world, resources.MovementDown)
	case moveLeftAction:
		resources.Move(world, resources.MovementLeft)
	case moveRightAction:
		resources.Move(world, resources.MovementRight)
	}
}
