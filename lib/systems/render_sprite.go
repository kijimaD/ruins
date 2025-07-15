package systems

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"

	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	wallShadowImage  *ebiten.Image // 壁が落とす影
	moverShadowImage *ebiten.Image // 動く物体が落とす影
)

// RenderSpriteSystem は (下) タイル -> 影 -> スプライト (上) の順に表示する
func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	gameComponents := world.Components.Game.(*gc.Components)
	gameResources := world.Resources.Game.(*resources.Game)

	// 初回のみ生成
	if wallShadowImage == nil {
		wallShadowImage = ebiten.NewImage(int(utils.TileSize), int(utils.TileSize/2))
		wallShadowImage.Fill(color.RGBA{0, 0, 0, 80})
	}
	if moverShadowImage == nil {
		moverShadowImage = ebiten.NewImage(int(utils.TileSize-6-2), int(utils.TileSize/2))
		moverShadowImage.Fill(color.RGBA{0, 0, 0, 120})
	}
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
			spriteRender1 := gameComponents.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
			spriteRender2 := gameComponents.SpriteRender.Get(entities[j]).(*gc.SpriteRender)
			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			// タイル描画
			gridElement := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*gc.SpriteRender)
			tileSize := gameResources.Level.TileSize
			pos := &gc.Position{
				X: gc.Pixel(int(gridElement.Row)*int(tileSize) + int(tileSize/2)),
				Y: gc.Pixel(int(gridElement.Col)*int(tileSize) + int(tileSize/2)),
			}
			drawImage(world, screen, spriteRender, pos, 0)
		}
	}
	{
		// 移動体の影。影をキャストする用のコンポーネントを追加したほうがよさそう
		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			pos := gameComponents.Position.Get(entity).(*gc.Position)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(int(pos.X)-12), float64(pos.Y))
			utils.SetTranslate(world, op)
			if moverShadowImage != nil {
				screen.DrawImage(moverShadowImage, op)
			}
		}))
	}
	{
		// 壁の影。影をキャストする用のコンポーネントを追加したほうがよさそう
		// 下のタイルがフロアであれば追加する
		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.GridElement,
			gameComponents.BlockView,
			gameComponents.BlockPass,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			grid := gameComponents.GridElement.Get(entity).(*gc.GridElement)
			gameResources := world.Resources.Game.(*resources.Game)
			belowTileIdx := gameResources.Level.XYTileIndex(grid.Row, grid.Col+1)
			if (belowTileIdx < 0) || (int(belowTileIdx) > len(gameResources.Level.Entities)-1) {
				return
			}
			belowTileEntity := gameResources.Level.Entities[int(belowTileIdx)]
			belowSpriteRender, ok := gameComponents.SpriteRender.Get(belowTileEntity).(*gc.SpriteRender)
			if ok {
				if belowSpriteRender.Depth == gc.DepthNumFloor {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(int(grid.Row)*int(utils.TileSize)), float64(int(grid.Col)*int(utils.TileSize)+int(utils.TileSize)))
					utils.SetTranslate(world, op)
					if wallShadowImage != nil {
						screen.DrawImage(wallShadowImage, op)
					}
				}
			}
		}))
	}
	{
		// 移動体
		iSprite := 0
		entities := make([]ecs.Entity, world.Manager.Join(gameComponents.SpriteRender, gameComponents.Position, gameComponents.Velocity).Size())
		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.Velocity,
			gameComponents.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			entities[iSprite] = entity
			iSprite++
		}))
		sort.Slice(entities, func(i, j int) bool {
			spriteRender1 := gameComponents.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
			spriteRender2 := gameComponents.SpriteRender.Get(entities[j]).(*gc.SpriteRender)

			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			// 座標描画
			velocity := gameComponents.Velocity.Get(entity).(*gc.Velocity)
			pos := gameComponents.Position.Get(entity).(*gc.Position)
			spriteRender := gameComponents.SpriteRender.Get(entity).(*gc.SpriteRender)
			drawImage(world, screen, spriteRender, pos, velocity.Angle)
		}
	}
}

func getImage(spriteRender *gc.SpriteRender) *ebiten.Image {
	var result *ebiten.Image
	key := fmt.Sprintf("%s/%d", spriteRender.SpriteSheet.Name, spriteRender.SpriteNumber)
	if v, ok := spriteImageCache[key]; ok {
		result = v
	} else {
		// テクスチャから欲しいスプライトを切り出す
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		texture := spriteRender.SpriteSheet.Texture
		textureWidth := texture.Image.Bounds().Dx()
		textureHeight := texture.Image.Bounds().Dy()

		left := max(0, sprite.X)
		right := min(textureWidth, sprite.X+sprite.Width)
		top := max(0, sprite.Y)
		bottom := min(textureHeight, sprite.Y+sprite.Height)

		result = texture.Image.SubImage(image.Rect(left, top, right, bottom)).(*ebiten.Image)
		spriteImageCache[key] = result
	}

	return result
}

func drawImage(world w.World, screen *ebiten.Image, spriteRender *gc.SpriteRender, pos *gc.Position, angle float64) {
	sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]

	op := &spriteRender.Options
	op.GeoM.Reset()                                                       // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(-sprite.Width/2), float64(-sprite.Width/2)) // 回転軸を画像の中心にする
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	utils.SetTranslate(world, op)
	screen.DrawImage(getImage(spriteRender), op)
}
