package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func RenderObjectSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		switch {
		case entity.HasComponent(gameComponents.Player):

			vector.DrawFilledRect(screen, float32(pos.X)-16, float32(pos.Y)-16, 32, 32, color.White, true)
		default:
			vector.DrawFilledRect(screen, float32(pos.X)-8, float32(pos.Y)-8, 16, 16, color.RGBA{0, 255, 0, 255}, true)
		}
	}))
}
