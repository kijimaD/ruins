package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"

	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	transformOffsetX = 0
	transformOffsetY = -80
)

// GridTransformSystem sets transform for grid elements
func GridTransformSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.GridElement, world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		elementSpriteRender := world.Components.Engine.SpriteRender.Get(entity).(*ec.SpriteRender)
		elementTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		elementSprite := elementSpriteRender.SpriteSheet.Sprites[elementSpriteRender.SpriteNumber]

		elementTranslation.X = float64(gridElement.Col*elementSprite.Width) + float64(elementSprite.Width)/2 + transformOffsetX
		elementTranslation.Y = screenHeight - float64(gridElement.Line*elementSprite.Height) - float64(elementSprite.Height)/2 + transformOffsetY
	}))
}
