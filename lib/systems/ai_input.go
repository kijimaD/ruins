package systems

import (
	"math/rand/v2"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AIInputSystem は AI制御されたエンティティの入力処理を行う
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
		if time.Since(roaming.StartSubState).Seconds() > roaming.DurationSubState.Seconds() {
			roaming.StartSubState = time.Now()
			roaming.DurationSubState = time.Second * time.Duration(rand.IntN(3))

			var subState gc.AIRoamingSubState
			switch rand.IntN(2) {
			case 0:
				subState = gc.AIRoamingWaiting
			case 1:
				subState = gc.AIRoamingDriving
			}

			switch subState {
			case gc.AIRoamingWaiting:
				velocity.ThrottleMode = gc.ThrottleModeNope
				velocity.Angle += float64(rand.IntN(91))
			case gc.AIRoamingDriving:
				velocity.ThrottleMode = gc.ThrottleModeFront
			}
		}
	}))
}
