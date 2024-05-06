package systems

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	m "github.com/kijimaD/ruins/lib/engine/math"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils/camera"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// TODO: ちゃんと保存する
	CameraX = 200
	CameraY = 200
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

	left := m.Max(0, sprite.X)
	right := m.Min(textureWidth, sprite.X+sprite.Width)
	top := m.Max(0, sprite.Y)
	bottom := m.Min(textureHeight, sprite.Y+sprite.Height)

	op := &spriteRender.Options
	op.GeoM.Reset() // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(pos.X-sprite.Width/2), float64(pos.Y-sprite.Width/2))
	camera.SetTranslate(world, op, -CameraX, -CameraY)
	screen.DrawImage(texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image), op)
}
