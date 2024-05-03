package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// raycast用move
func MoveRaySystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	var pos *gc.Position
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	// if inpututil.IsKeyJustPressed(ebiten.KeyR) {
	// 	g.showRays = !g.showRays
	// }

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		pos.X += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		pos.Y += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		pos.X -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		pos.Y -= 2
	}

	// // +1/-1 is to stop player before it reaches the border
	// if g.Px >= g.ScreenWidth-padding {
	// 	g.Px = g.ScreenWidth - padding - 1
	// }

	// if g.Px <= padding {
	// 	g.Px = padding + 1
	// }

	// if g.Py >= g.ScreenHeight-padding {
	// 	g.Py = g.ScreenHeight - padding - 1
	// }

	// if g.Py <= padding {
	// 	g.Py = padding + 1
	// }
}
