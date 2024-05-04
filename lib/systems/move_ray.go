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

	padding := 20
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// // +1/-1 is to stop player before it reaches the border
	if pos.X >= screenWidth-padding {
		pos.X = screenWidth - padding - 1
	}

	if pos.X <= padding {
		pos.X = padding + 1
	}

	if pos.Y >= screenHeight-padding {
		pos.Y = screenHeight - padding - 1
	}

	if pos.Y <= padding {
		pos.Y = padding + 1
	}
}
