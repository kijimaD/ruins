package systems

import (
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func AIMoveSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Position,
		gameComponents.AIMoveFSM,
		gameComponents.AIRoaming,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		fsm := gameComponents.AIMoveFSM.Get(entity).(*gc.AIMoveFSM)
		diff := time.Now().Sub(fsm.LastStateChange)
		if diff.Seconds() > 2 {
			fsm.LastStateChange = time.Now()

			pos := gameComponents.Position.Get(entity).(*gc.Position)
			pos.Angle += 1
			pos.Speed += 1
		}
	}))
}
