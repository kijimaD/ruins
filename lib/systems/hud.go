package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// HUDSystem はゲームの HUD 情報を描画する
func HUDSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Game.(*resources.Game)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", gameResources.Depth), 0, 200)

	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed: %.2f", velocity.Speed), 0, 220)
	}))
}
