package systems

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	m "github.com/kijimaD/ruins/lib/engine/math"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func RenderObjectSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(
		gameComponents.Position,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		sprite := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

		spriteWidth := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float32(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		// オブジェクトの足元の影。とりあえず矩形。スプライトの白黒画像を下に表示するのが望ましい
		vector.DrawFilledRect(screen, float32(pos.X)-16, float32(pos.Y)-16, spriteWidth, spriteHeight, color.RGBA{0, 0, 0, 100}, true)

		drawImage(screen, sprite, pos)
	}))
}

func drawImage(screen *ebiten.Image, spriteRender *ec.SpriteRender, pos *gc.Position) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
	texture := spriteRender.SpriteSheet.Texture
	textureWidth, textureHeight := texture.Image.Size()

	left := m.Max(0, sprite.X)
	right := m.Min(textureWidth, sprite.X+sprite.Width)
	top := m.Max(0, sprite.Y)
	bottom := m.Min(textureHeight, sprite.Y+sprite.Height)

	op := spriteRender.Options
	op.GeoM.Translate(float64(pos.X-sprite.Width/2), float64(pos.Y-sprite.Width/2))
	screen.DrawImage(texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image), &op)
}
