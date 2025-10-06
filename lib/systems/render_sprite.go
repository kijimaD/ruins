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

// RenderSpriteSystem は (下) タイル -> 暗闇 -> 光源グロー -> 影 -> スプライト (上) の順に表示する
func RenderSpriteSystem(world w.World, screen *ebiten.Image) {
	// 現在の視界データを取得（全描画で使用）
	visibilityData := GetCurrentVisibilityData()

	// シャドウ画像を初期化
	initializeShadowImages()

	// 各種描画処理
	renderGridTiles(world, screen, visibilityData)             // 床タイルを描画
	renderDistanceBasedDarkness(world, screen, visibilityData) // 床タイルに暗闇オーバーレイ
	renderLightSourceGlow(world, screen, visibilityData)       // 床タイルに光源グロー
	renderShadows(world, screen, visibilityData)               // 影を描画
	renderSprites(world, screen, visibilityData)               // 物体を描画
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

// renderGridTiles は床タイルを描画する
func renderGridTiles(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	iSprite := 0
	entities := make([]ecs.Entity, world.Manager.Join(world.Components.SpriteRender, world.Components.GridElement).Size())
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.Prop.Not(),
		world.Components.TurnBased.Not(),
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

// renderLightSourceGlow は光源タイルに明るいオーバーレイを描画する
func renderLightSourceGlow(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	// カメラ位置とスケールを取得
	var cameraPos gc.Position
	cameraScale := 1.0

	world.Manager.Join(
		world.Components.Camera,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		cameraGridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		cameraPos = gc.Position{
			X: gc.Pixel(int(cameraGridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2),
			Y: gc.Pixel(int(cameraGridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2),
		}
		camera := world.Components.Camera.Get(entity).(*gc.Camera)
		cameraScale = camera.Scale
	}))

	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height
	tileSize := int(consts.TileSize)

	// 光源エンティティに明るいオーバーレイを描画
	world.Manager.Join(
		world.Components.LightSource,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		lightSource := world.Components.LightSource.Get(entity).(*gc.LightSource)
		if !lightSource.Enabled {
			return
		}

		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 視界チェック - 視界内のもののみ描画
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", gridElement.X, gridElement.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				return
			}
		}

		// タイルの画面座標を計算
		worldX := float64(int(gridElement.X) * tileSize)
		worldY := float64(int(gridElement.Y) * tileSize)
		screenX := (worldX-float64(cameraPos.X))*cameraScale + float64(screenWidth)/2
		screenY := (worldY-float64(cameraPos.Y))*cameraScale + float64(screenHeight)/2

		// 光源色で明るいオーバーレイを作成
		cacheKey := fmt.Sprintf("glow_%d,%d,%d", lightSource.Color.R, lightSource.Color.G, lightSource.Color.B)
		glowImg, exists := spriteImageCache[cacheKey]
		if !exists {
			glowImg = ebiten.NewImage(tileSize, tileSize)
			glowColor := color.RGBA{
				R: uint8(float64(lightSource.Color.R) * 0.6),
				G: uint8(float64(lightSource.Color.G) * 0.5),
				B: uint8(float64(lightSource.Color.B) * 0.3),
				A: 80,
			}
			glowImg.Fill(glowColor)
			if len(spriteImageCache) < 1000 {
				spriteImageCache[cacheKey] = glowImg
			}
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(cameraScale, cameraScale)
		op.GeoM.Translate(screenX, screenY)
		screen.DrawImage(glowImg, op)
	}))
}

// renderSprites はスプライトを描画する
func renderSprites(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	var entities []ecs.Entity

	// Props を収集
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.Prop,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entities = append(entities, entity)
	}))

	// Movers を収集
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.TurnBased,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entities = append(entities, entity)
	}))

	sort.Slice(entities, func(i, j int) bool {
		spriteRender1 := world.Components.SpriteRender.Get(entities[i]).(*gc.SpriteRender)
		spriteRender2 := world.Components.SpriteRender.Get(entities[j]).(*gc.SpriteRender)
		return spriteRender1.Depth < spriteRender2.Depth
	})

	// 描画
	for _, entity := range entities {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 視界チェック - 視界内のもののみ描画
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", gridElement.X, gridElement.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				continue
			}
		}

		// 光源チェック - 光がある場所のみ描画（完全に暗い場所は描画しない）
		lightInfo := getCachedLightInfo(world, int(gridElement.X), int(gridElement.Y))
		if lightInfo.Darkness >= 1.0 {
			continue // 完全に暗い場所は描画しない
		}

		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		pos := &gc.Position{
			X: gc.Pixel(int(gridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2),
			Y: gc.Pixel(int(gridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2),
		}
		drawImage(world, screen, spriteRender, pos, 0)
	}
}

// renderShadows は物体と壁の影を描画する
func renderShadows(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	// 物体の影
	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// TurnBased または Prop を持つエンティティのみ
		if !entity.HasComponent(world.Components.TurnBased) && !entity.HasComponent(world.Components.Prop) {
			return
		}

		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 視界チェック
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

	// 壁の影（下タイルが床の場合のみ）
	tileMap := make(map[gc.GridElement]ecs.Entity)
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(e ecs.Entity) {
		ge := world.Components.GridElement.Get(e).(*gc.GridElement)
		tileMap[*ge] = e
	}))

	world.Manager.Join(
		world.Components.SpriteRender,
		world.Components.GridElement,
		world.Components.BlockView,
		world.Components.BlockPass,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// 視界チェック
		if visibilityData != nil {
			tileKey := fmt.Sprintf("%d,%d", grid.X, grid.Y)
			if tileData, exists := visibilityData[tileKey]; !exists || !tileData.Visible {
				return
			}
		}

		spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

		// 高さのあるものだけが影を落とす
		if spriteRender.Depth > gc.DepthNumTaller {
			return
		}

		// 下のタイルを検索
		belowPos := gc.GridElement{X: grid.X, Y: grid.Y + 1}
		belowTileEntity, foundBelow := tileMap[belowPos]

		if !foundBelow {
			return
		}

		belowSpriteRender, ok := world.Components.SpriteRender.Get(belowTileEntity).(*gc.SpriteRender)
		if !ok || belowSpriteRender.Depth != gc.DepthNumFloor {
			return // 下が床でなければ影を描画しない
		}

		// 下のタイルが壁でないことも確認（壁->床->壁の場合は影を描画しない）
		if belowTileEntity.HasComponent(world.Components.BlockView) && belowTileEntity.HasComponent(world.Components.BlockPass) {
			return
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(int(grid.X)*int(consts.TileSize)), float64(int(grid.Y)*int(consts.TileSize)+int(consts.TileSize)))
		SetTranslate(world, op)
		if wallShadowImage != nil {
			screen.DrawImage(wallShadowImage, op)
		}
	}))
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
	if cfg.ShowMonitor && spriteRender.SpriteSheetName == "tile" && strings.HasPrefix(spriteRender.SpriteKey, "dirt_") {
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
