package systems

import (
	"fmt"
	"image"
	"sort"

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

	{
		// グリッド
		iSprite := 0
		entities := make([]ecs.Entity, world.Manager.Join(gameComponents.SpriteRender, gameComponents.GridElement).Size())
		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.GridElement,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			entities[iSprite] = entity
			iSprite++
		}))
		sort.Slice(entities, func(i, j int) bool {
			spriteRender1 := gameComponents.SpriteRender.Get(entities[i]).(*ec.SpriteRender)
			spriteRender2 := gameComponents.SpriteRender.Get(entities[j]).(*ec.SpriteRender)
			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			// タイル描画
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
			tileSize := gameResources.Level.TileSize
			pos := &gc.Position{
				X: gc.Pixel(int(gridElement.Row)*int(tileSize) + int(tileSize/2)),
				Y: gc.Pixel(int(gridElement.Col)*int(tileSize) + int(tileSize/2)),
			}
			drawImage(world, screen, spriteRender, pos)
		}
	}
	{
		// 移動体
		iSprite := 0
		entities := make([]ecs.Entity, world.Manager.Join(gameComponents.SpriteRender, gameComponents.Position).Size())
		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			entities[iSprite] = entity
			iSprite++
		}))
		sort.Slice(entities, func(i, j int) bool {
			spriteRender1 := gameComponents.SpriteRender.Get(entities[i]).(*ec.SpriteRender)
			spriteRender2 := gameComponents.SpriteRender.Get(entities[j]).(*ec.SpriteRender)

			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			// 座標描画
			pos := gameComponents.Position.Get(entity).(*gc.Position)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
			drawImage(world, screen, spriteRender, pos)
		}
	}
}

func getImage(spriteRender *ec.SpriteRender) *ebiten.Image {
	var result *ebiten.Image
	key := fmt.Sprintf("%s/%d", spriteRender.SpriteSheet.Name, spriteRender.SpriteNumber)
	if v, ok := spriteImageCache[key]; ok {
		result = v
	} else {
		// テクスチャから欲しいスプライトを切り出す
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		texture := spriteRender.SpriteSheet.Texture
		textureWidth, textureHeight := texture.Image.Size()

		left := m.Max(0, sprite.X)
		right := m.Min(textureWidth, sprite.X+sprite.Width)
		top := m.Max(0, sprite.Y)
		bottom := m.Min(textureHeight, sprite.Y+sprite.Height)

		result = texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image)
		spriteImageCache[key] = result
	}

	return result
}

func drawImage(world w.World, screen *ebiten.Image, spriteRender *ec.SpriteRender, pos *gc.Position) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]

	op := &spriteRender.Options
	op.GeoM.Reset()                                                       // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(-sprite.Width/2), float64(-sprite.Width/2)) // 回転軸を画像の中心にする
	op.GeoM.Rotate(pos.Angle)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	camera.SetTranslate(world, op)
	screen.DrawImage(getImage(spriteRender), op)
}
