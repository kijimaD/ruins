package systems

import (
	"math/rand/v2"
	"time"

	"github.com/kijimaD/ruins/lib/components"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func AIInputSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Velocity,
		gameComponents.Position,
		gameComponents.AIMoveFSM,
		gameComponents.AIRoaming,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		roaming := gameComponents.AIRoaming.Get(entity).(*gc.AIRoaming)
		velocity := gameComponents.Velocity.Get(entity).(*gc.Velocity)
		if time.Now().Sub(roaming.StartSubState).Seconds() > roaming.DurationSubState.Seconds() {
			roaming.StartSubState = time.Now()
			roaming.DurationSubState = time.Second * time.Duration(rand.IntN(3))

			var subState components.AIRoamingSubState
			switch rand.IntN(2) {
			case 0:
				subState = components.AIRoamingWaiting
			case 1:
				subState = components.AIRoamingDriving
			}

			switch subState {
			case components.AIRoamingWaiting:
				// TODO: スロットルみたいな移動用関数を作ってゆるやかに変化させるべきである
				velocity.Speed = 0
				velocity.Angle += float64(rand.IntN(91))
			case components.AIRoamingDriving:
				velocity.Speed = 1
			}
		}
	}))
}
