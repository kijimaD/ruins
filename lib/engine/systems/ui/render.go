package uisystem

import (
	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/utils"
	w "github.com/kijimaD/ruins/lib/engine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// RenderUISystem draws text entities
func RenderUISystem(world w.World, screen *ebiten.Image) {
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		textData := world.Components.Engine.Text.Get(entity).(*c.Text)
		uiTransform := world.Components.Engine.UITransform.Get(entity).(*c.UITransform)

		// Compute dot offset
		x, y := utils.Try2(c.ComputeDotOffset(textData.Text, textData.FontFace, uiTransform.Pivot))

		// Draw text
		screenWidth := world.Resources.ScreenDimensions.Width
		screenHeight := world.Resources.ScreenDimensions.Height

		offsetX, offsetY := uiTransform.ComputeOriginOffset(screenWidth, screenHeight)
		text.Draw(screen, textData.Text, textData.FontFace, uiTransform.Translation.X+offsetX-x, screenHeight-uiTransform.Translation.Y-offsetY-y, textData.Color)
	}))
}
