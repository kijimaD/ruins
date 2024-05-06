package systems

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	m "github.com/kijimaD/ruins/lib/engine/math"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// TODO: ちゃんと保存する
	cameraX = 200
	cameraY = 200
)

func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Position,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		sprite := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

		drawImage(world, screen, sprite, pos)
	}))
}

func drawImage(world w.World, screen *ebiten.Image, spriteRender *ec.SpriteRender, pos *gc.Position) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
	texture := spriteRender.SpriteSheet.Texture
	textureWidth, textureHeight := texture.Image.Size()

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	left := m.Max(0, sprite.X)
	right := m.Min(textureWidth, sprite.X+sprite.Width)
	top := m.Max(0, sprite.Y)
	bottom := m.Min(textureHeight, sprite.Y+sprite.Height)

	op := &spriteRender.Options
	op.GeoM.Reset()
	op.GeoM.Translate(float64(pos.X-sprite.Width/2), float64(pos.Y-sprite.Width/2))
	op.GeoM.Translate(float64(-cameraX), float64(-cameraY))
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(float64(cx), float64(cy))
	screen.DrawImage(texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image), op)
}
