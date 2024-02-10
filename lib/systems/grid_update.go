package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"

	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// GridUpdateSystem updates grid elements
func GridUpdateSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	playerIndex := -1
	for iTile, tile := range gameResources.Level.Grid.Data {
		switch {
		case tile.Contains(resources.TilePlayer):
			playerIndex = iTile
		}
	}

	levelWidth := gameResources.Level.Grid.NCols
	levelHeight := gameResources.Level.Grid.NRows

	paddingRow := (gameResources.GridLayout.Height - levelHeight) / 2
	paddingCol := (gameResources.GridLayout.Width - levelWidth) / 2

	world.Manager.Join(gameComponents.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		switch {
		case entity.HasComponent(gameComponents.Player):
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			gridElement.Line = paddingRow + playerIndex/levelWidth
			gridElement.Col = paddingCol + playerIndex%levelWidth
		}
	}))
}
