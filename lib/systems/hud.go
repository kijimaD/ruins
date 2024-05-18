package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
)

func HUDSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Game.(*resources.Game)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("B%d", gameResources.Depth), 0, 200)
}
