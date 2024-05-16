package systems

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	m "github.com/kijimaD/ruins/lib/engine/math"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"

	"github.com/kijimaD/ruins/lib/utils/camera"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	// TODO: ↓的な方法でソートしたほうがパフォーマンスがよさそう
	// Sort by increasing values of depth
	// sort.Slice(spritesDepths, func(i, j int) bool {
	// 	return spritesDepths[i].depth < spritesDepths[j].depth
	// })
	for _, v := range ec.DepthNums {
		// タイル描画
		world.Manager.Join(
			gameComponents.GridElement,
			gameComponents.SpriteRender,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
			tileSize := gameResources.Level.TileSize
			pos := &gc.Position{
				X: int(gridElement.Row)*tileSize + tileSize/2,
				Y: int(gridElement.Col)*tileSize + tileSize/2,
			}
			if spriteRender.Depth != v {
				return
			}
			drawImage(world, screen, spriteRender, pos)
		}))

		// 座標描画
		world.Manager.Join(
			gameComponents.Position,
			gameComponents.SpriteRender,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			pos := gameComponents.Position.Get(entity).(*gc.Position)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
			if spriteRender.Depth != v {
				return
			}
			drawImage(world, screen, spriteRender, pos)
		}))
	}
}

func drawImage(world w.World, screen *ebiten.Image, spriteRender *ec.SpriteRender, pos *gc.Position) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
	texture := spriteRender.SpriteSheet.Texture
	textureWidth, textureHeight := texture.Image.Size()

	// テクスチャから欲しいスプライトを切り出す
	left := m.Max(0, sprite.X)
	right := m.Min(textureWidth, sprite.X+sprite.Width)
	top := m.Max(0, sprite.Y)
	bottom := m.Min(textureHeight, sprite.Y+sprite.Height)

	op := &spriteRender.Options
	op.GeoM.Reset()                                                       // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(-sprite.Width/2), float64(-sprite.Width/2)) // 回転軸を画像の中心にする
	op.GeoM.Rotate(pos.Angle)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	camera.SetTranslate(world, op)
	screen.DrawImage(texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image), op)
}
