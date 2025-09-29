package systems

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	w "github.com/kijimaD/ruins/lib/world"

	"github.com/kijimaD/ruins/lib/consts"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	wallShadowImage  *ebiten.Image // 壁が落とす影
	moverShadowImage *ebiten.Image // 動く物体が落とす影
)

var spriteImageCache map[string]*ebiten.Image = map[string]*ebiten.Image{}

// SetTranslate はカメラを考慮した画像配置オプションをセットする
// TODO: ズーム率を追加する
func SetTranslate(world w.World, op *ebiten.DrawImageOptions) {
	var camera *gc.Camera
	var cPos *gc.Position
	var cGridElement *gc.GridElement

	// ピクセル座標のカメラを先に確認
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		camera = world.Components.Camera.Get(entity).(*gc.Camera)
		cPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	// グリッド座標のカメラを確認
	world.Manager.Join(
		world.Components.Camera,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		camera = world.Components.Camera.Get(entity).(*gc.Camera)
		cGridElement = world.Components.GridElement.Get(entity).(*gc.GridElement)
	}))

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置の設定
	if cGridElement != nil {
		// グリッド座標をピクセル座標に変換
		tilePixelX := float64(int(cGridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2)
		tilePixelY := float64(int(cGridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2)
		op.GeoM.Translate(-tilePixelX, -tilePixelY)
	} else if cPos != nil {
		op.GeoM.Translate(float64(-cPos.X), float64(-cPos.Y))
	}

	if camera != nil {
		op.GeoM.Scale(camera.Scale, camera.Scale)
	}
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
}

// RenderSpriteSystem は (下) タイル -> 影 -> スプライト (上) の順に表示する
func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	// 現在の視界データを取得（全描画で使用）
	visibilityData := GetCurrentVisibilityData()

	// シャドウ画像を初期化
	initializeShadowImages()

	// 各種描画処理
	renderGridTiles(world, screen, visibilityData)
	renderMoverShadows(world, screen, visibilityData)
	renderWallShadows(world, screen, visibilityData)
	renderMovers(world, screen, visibilityData)
}

// initializeShadowImages は影画像を初期化する
func initializeShadowImages() {
	if wallShadowImage == nil {
		wallWidth := int(consts.TileSize)
		wallHeight := int(consts.TileSize / 2)
		if wallWidth > 0 && wallHeight > 0 {
			wallShadowImage = ebiten.NewImage(wallWidth, wallHeight)
			wallShadowImage.Fill(color.RGBA{0, 0, 0, 80})
		}
	}
	if moverShadowImage == nil {
		moverWidth := int(consts.TileSize - 6 - 2)
		moverHeight := int(consts.TileSize / 2)
		if moverWidth > 0 && moverHeight > 0 {
			moverShadowImage = ebiten.NewImage(moverWidth, moverHeight)
			moverShadowImage.Fill(color.RGBA{0, 0, 0, 120})
		}
	}
}

// renderGridTiles はグリッドタイルを描画する
func renderGridTiles(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {

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
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 視界チェック - 視界内または探索済みのタイルのみ描画
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", gridElement.X, gridElement.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				// 探索済みかどうかもチェック
				if !world.Resources.Dungeon.ExploredTiles[*gridElement] {
					continue // 未探索かつ視界外のタイルは描画しない
				}
			}
		}

		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		pos := &gc.Position{
			X: gc.Pixel(int(gridElement.X)*int(consts.TileSize) + int(consts.TileSize/2)),
			Y: gc.Pixel(int(gridElement.Y)*int(consts.TileSize) + int(consts.TileSize/2)),
		}
		drawImage(world, screen, spriteRender, pos, 0)
	}
}

// renderMoverShadows は移動体の影を描画する
func renderMoverShadows(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.TurnBased,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 移動体の影も視界チェック
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", gridElement.X, gridElement.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				return
			}
		}

		// グリッド座標をピクセル座標に変換
		pixelX := float64(int(gridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2 - 12)
		pixelY := float64(int(gridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pixelX, pixelY)
		SetTranslate(world, op)
		if moverShadowImage != nil {
			screen.DrawImage(moverShadowImage, op)
		}
	}))
}

// renderWallShadows は壁の影を描画する
func renderWallShadows(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {

	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.BlockView,
		world.Components.BlockPass,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 壁の影も視界チェック - 視界内または探索済みのタイルのみ描画
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", grid.X, grid.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				// 探索済みかどうかもチェック
				if !world.Resources.Dungeon.ExploredTiles[*grid] {
					return // 未探索かつ視界外の影は描画しない
				}
			}
		}

		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

		// 高さのあるもの（壁など）だけが影を落とす
		if spriteRender.Depth != gc.DepthNumTaller {
			return
		}

		belowTileIdx := world.Resources.Dungeon.Level.XYTileIndex(grid.X, grid.Y+1)
		if (belowTileIdx < 0) || (int(belowTileIdx) > len(world.Resources.Dungeon.Level.Entities)-1) {
			return
		}
		belowTileEntity := world.Resources.Dungeon.Level.Entities[int(belowTileIdx)]
		belowSpriteRender, ok := world.Components.SpriteRender.Get(belowTileEntity).(*gc.SpriteRender)
		if !ok || belowSpriteRender.Depth != gc.DepthNumFloor {
			return // 下が床でなければ影を描画しない
		}

		// 下のタイルが壁でないことも確認（壁->床->壁の場合は影を描画しない）
		if belowTileEntity.HasComponent(world.Components.BlockView) && belowTileEntity.HasComponent(world.Components.BlockPass) {
			return // 下のタイルも壁なら影を描画しない
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(int(grid.X)*int(consts.TileSize)), float64(int(grid.Y)*int(consts.TileSize)+int(consts.TileSize)))
		SetTranslate(world, op)
		if wallShadowImage != nil {
			screen.DrawImage(wallShadowImage, op)
		}
	}))
}

// renderMovers は移動体を描画する
func renderMovers(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	iSprite := 0
	entities := make([]ecs.Entity, world.Manager.Join(world.Components.SpriteRender, world.Components.GridElement, world.Components.TurnBased).Size())
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.TurnBased,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entities[iSprite] = entity
		iSprite++
	}))
	sort.Slice(entities[:iSprite], func(i, j int) bool {
		spriteRender1 := world.Components.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
		spriteRender2 := world.Components.SpriteRender.Get(entities[j]).(*gc.SpriteRender)
		return spriteRender1.Depth < spriteRender2.Depth
	})
	for i := 0; i < iSprite; i++ {
		entity := entities[i]
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// グリッド座標からピクセル座標に変換
		pixelPos := &gc.Position{
			X: gc.Pixel(int(gridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2),
			Y: gc.Pixel(int(gridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2),
		}

		// 移動体の視界チェック - 視界内のもののみ描画
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", gridElement.X, gridElement.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				continue
			}
		}

		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		drawImage(world, screen, spriteRender, pixelPos, 0) // 移動体は回転なし
	}
}

func getImage(world w.World, spriteRender *gc.SpriteRender) *ebiten.Image {
	var result *ebiten.Image
	key := fmt.Sprintf("%s/%s", spriteRender.SpriteSheetName, spriteRender.SpriteKey)
	if v, ok := spriteImageCache[key]; ok {
		result = v
	} else {
		// Resourcesからスプライトシートを取得
		if world.Resources.SpriteSheets == nil {
			return nil
		}
		spriteSheet, exists := (*world.Resources.SpriteSheets)[spriteRender.SpriteSheetName]
		if !exists {
			return nil
		}

		// スプライトキーからスプライトを取得
		sprite, exists := spriteSheet.Sprites[spriteRender.SpriteKey]
		if !exists {
			return nil
		}

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
	spriteSheet, exists := (*world.Resources.SpriteSheets)[spriteRender.SpriteSheetName]
	if !exists {
		return
	}

	sprite, exists := spriteSheet.Sprites[spriteRender.SpriteKey]
	if !exists {
		return
	}

	op := &spriteRender.Options
	op.GeoM.Reset()                                                       // FIXME: Resetがないと非表示になる。なぜ?
	op.GeoM.Translate(float64(-sprite.Width/2), float64(-sprite.Width/2)) // 回転軸を画像の中心にする
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(float64(pos.X), float64(pos.Y))
	SetTranslate(world, op)
	screen.DrawImage(getImage(world, spriteRender), op)

	// デバッグ用：スプライト番号表示(土だけ)
	cfg := config.Get()
	if cfg.ShowMonitor && spriteRender.SpriteSheetName == "field" && strings.HasPrefix(spriteRender.SpriteKey, "dirt_") {
		// dirt_X から番号を抽出
		number := strings.TrimPrefix(spriteRender.SpriteKey, "dirt_")

		// カメラ変換を考慮したテキスト位置を計算
		textOp := &ebiten.DrawImageOptions{}
		textOp.GeoM.Translate(float64(pos.X-8), float64(pos.Y-8)) // タイルの左上付近に表示
		SetTranslate(world, textOp)

		// テキスト表示位置を逆変換で求める
		screenX, screenY := textOp.GeoM.Apply(0, 0)
		ebitenutil.DebugPrintAt(screen, number, int(screenX), int(screenY))
	}
}
