package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func HUDSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Game.(*resources.Game)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", gameResources.Depth), 0, 200)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Operator,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed: %.2f", pos.Speed), 0, 220)
	}))
}
