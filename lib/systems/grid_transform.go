package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"

	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	// スクリーンサイズ x 0.5
	transformOffsetX = 480
	transformOffsetY = -360
	// 1つあたりのタイルサイズ
	tileSize = 32
)

// GridTransformSystem sets transform for grid elements
// TODO: タイルサイズと画面サイズをハードコードしてカメラを実装しているので、画面サイズが変わると壊れる
func GridTransformSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	var playerX = 0
	var playerY = 0
	world.Manager.Join(
		gameComponents.Player,
		gameComponents.GridElement,
		world.Components.Engine.SpriteRender,
		world.Components.Engine.Transform,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
		playerX = gridElement.Line
		playerY = gridElement.Col
	}))

	world.Manager.Join(
		gameComponents.GridElement,
		world.Components.Engine.SpriteRender,
		world.Components.Engine.Transform,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)

		elementSpriteRender := world.Components.Engine.SpriteRender.Get(entity).(*ec.SpriteRender)
		elementTranslation := &world.Components.Engine.Transform.Get(entity).(*ec.Transform).Translation

		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		elementSprite := elementSpriteRender.SpriteSheet.Sprites[elementSpriteRender.SpriteNumber]

		elementTranslation.X = float64(gridElement.Col*elementSprite.Width) + float64(elementSprite.Width)/2 + transformOffsetX - float64(playerY*tileSize)
		elementTranslation.Y = screenHeight - float64(gridElement.Line*elementSprite.Height) - float64(elementSprite.Height)/2 + transformOffsetY + float64(playerX*tileSize)
	}))
}
