package systems

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"

	"github.com/kijimaD/ruins/lib/consts"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	wallShadowImage  *ebiten.Image // 壁が落とす影
	moverShadowImage *ebiten.Image // 動く物体が落とす影
)

// SetTranslate はカメラを考慮した画像配置オプションをセットする
// TODO: ズーム率を追加する
func SetTranslate(world w.World, op *ebiten.DrawImageOptions) {
	var camera *gc.Camera
	var cPos *gc.Position
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		camera = world.Components.Camera.Get(entity).(*gc.Camera)
		cPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置
	op.GeoM.Translate(float64(-cPos.X), float64(-cPos.Y))
	op.GeoM.Scale(camera.Scale, camera.Scale)
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
}

// RenderSpriteSystem は (下) タイル -> 影 -> スプライト (上) の順に表示する
func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	// 現在の視界データを取得（全描画で使用）
	visibilityData := GetCurrentVisibilityData()

	// 初回のみ生成
	if wallShadowImage == nil {
		wallShadowImage = ebiten.NewImage(int(consts.TileSize), int(consts.TileSize/2))
		wallShadowImage.Fill(color.RGBA{0, 0, 0, 80})
	}
	if moverShadowImage == nil {
		moverShadowImage = ebiten.NewImage(int(consts.TileSize-6-2), int(consts.TileSize/2))
		moverShadowImage.Fill(color.RGBA{0, 0, 0, 120})
	}
	{
		// グリッド
		iSprite := 0
		entities := make([]ecs.Entity, world.Manager.Join(world.Components.SpriteRender, world.Components.GridElement).Size())
		world.Manager.Join(
			world.Components.SpriteRender,
			world.Components.GridElement,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			entities[iSprite] = entity
			iSprite++
		}))
		sort.Slice(entities, func(i, j int) bool {
			spriteRender1 := world.Components.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
			spriteRender2 := world.Components.SpriteRender.Get(entities[j]).(*gc.SpriteRender)
			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			// タイル描画
			gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

			// 視界チェック - 視界内または探索済みのタイルのみ描画
			if visibilityData != nil {
				tileKey := fmt.Sprintf("%d,%d", gridElement.Row, gridElement.Col)
				if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
					// 探索済みかどうかもチェック
					if !gameResources.ExploredTiles[tileKey] {
						continue // 未探索かつ視界外のタイルは描画しない
					}
				}
			}

			spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
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
			world.Components.SpriteRender,
			world.Components.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			pos := world.Components.Position.Get(entity).(*gc.Position)

			// 移動体の影も視界チェック
			if visibilityData != nil {
				tileX := int(pos.X) / 32
				tileY := int(pos.Y) / 32
				tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

				// 視界内でない場合は影も描画しない
				if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
					return
				}
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(int(pos.X)-12), float64(pos.Y))
			SetTranslate(world, op)
			if moverShadowImage != nil {
				screen.DrawImage(moverShadowImage, op)
			}
		}))
	}
	{
		// 壁の影。影をキャストする用のコンポーネントを追加したほうがよさそう
		// 下のタイルがフロアであれば追加する
		world.Manager.Join(
			world.Components.SpriteRender,
			world.Components.GridElement,
			world.Components.BlockView,
			world.Components.BlockPass,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

			// 壁の影も視界チェック - 視界内または探索済みのタイルのみ描画
			if visibilityData != nil {
				tileKey := fmt.Sprintf("%d,%d", grid.Row, grid.Col)
				if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
					// 探索済みかどうかもチェック
					if !gameResources.ExploredTiles[tileKey] {
						return // 未探索かつ視界外の影は描画しない
					}
				}
			}

			gameResources := world.Resources.Dungeon.(*resources.Dungeon)
			belowTileIdx := gameResources.Level.XYTileIndex(grid.Row, grid.Col+1)
			if (belowTileIdx < 0) || (int(belowTileIdx) > len(gameResources.Level.Entities)-1) {
				return
			}
			belowTileEntity := gameResources.Level.Entities[int(belowTileIdx)]
			belowSpriteRender, ok := world.Components.SpriteRender.Get(belowTileEntity).(*gc.SpriteRender)
			if ok {
				if belowSpriteRender.Depth == gc.DepthNumFloor {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(int(grid.Row)*int(consts.TileSize)), float64(int(grid.Col)*int(consts.TileSize)+int(consts.TileSize)))
					SetTranslate(world, op)
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
		entities := make([]ecs.Entity, world.Manager.Join(world.Components.SpriteRender, world.Components.Position, world.Components.Velocity).Size())
		world.Manager.Join(
			world.Components.SpriteRender,
			world.Components.Velocity,
			world.Components.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			entities[iSprite] = entity
			iSprite++
		}))
		sort.Slice(entities, func(i, j int) bool {
			spriteRender1 := world.Components.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
			spriteRender2 := world.Components.SpriteRender.Get(entities[j]).(*gc.SpriteRender)

			return spriteRender1.Depth < spriteRender2.Depth
		})
		for _, entity := range entities {
			pos := world.Components.Position.Get(entity).(*gc.Position)

			// 移動体の視界チェック - 視界内のもののみ描画
			if visibilityData != nil {
				// エンティティの位置からタイル座標を計算
				tileX := int(pos.X) / 32
				tileY := int(pos.Y) / 32
				tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

				// 視界内でない場合は描画しない
				if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
					continue
				}
			}

			// 座標描画
			velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
			spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
			drawImage(world, screen, spriteRender, pos, velocity.Angle)
		}
	}
}

func getImage(world w.World, spriteRender *gc.SpriteRender) *ebiten.Image {
	var result *ebiten.Image
	key := fmt.Sprintf("%s/%d", spriteRender.Name, spriteRender.SpriteNumber)
	if v, ok := spriteImageCache[key]; ok {
		result = v
	} else {
		// Resourcesからスプライトシートを取得
		if world.Resources.SpriteSheets == nil {
			return nil
		}
		spriteSheet, exists := (*world.Resources.SpriteSheets)[spriteRender.Name]
		if !exists {
			return nil
		}

		// テクスチャから欲しいスプライトを切り出す
		if spriteRender.SpriteNumber >= len(spriteSheet.Sprites) {
			return nil
		}
		sprite := spriteSheet.Sprites[spriteRender.SpriteNumber]
		texture := spriteSheet.Texture
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
	// Resourcesからスプライトシートを取得
	if world.Resources.SpriteSheets == nil {
		return
	}
	spriteSheet, exists := (*world.Resources.SpriteSheets)[spriteRender.Name]
	if !exists {
		return
	}

	if spriteRender.SpriteNumber >= len(spriteSheet.Sprites) {
		return
	}
	sprite := spriteSheet.Sprites[spriteRender.SpriteNumber]

	op := &spriteRender.Options
	op.GeoM.Reset()                                                       // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(-sprite.Width/2), float64(-sprite.Width/2)) // 回転軸を画像の中心にする
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	SetTranslate(world, op)
	screen.DrawImage(getImage(world, spriteRender), op)
}
