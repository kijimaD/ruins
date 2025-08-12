package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// 段階的な暗闇用の画像キャッシュ
	darknessCacheImages []*ebiten.Image
)

// VisionSystem はタイルごとの視界を管理し、暗闇を描画する
func VisionSystem(world w.World, screen *ebiten.Image) {
	// プレイヤー位置を取得
	var playerPos *gc.Position
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	if playerPos == nil {
		return
	}

	var exploredMap *gc.ExploredMap
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
		world.Components.ExploredMap,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		exploredMap = world.Components.ExploredMap.Get(entity).(*gc.ExploredMap)
	}))

	if exploredMap == nil {
		// ExploredMapが見つからない場合はエラー
		return
	}

	// タイルの可視性マップを更新
	visionRadius := gc.Pixel(320)
	visibilityData := calculateTileVisibilityWithDistance(world, playerPos.X, playerPos.Y, visionRadius)

	// 視界内のタイルを探索済みとしてマーク
	visibleCount := 0
	for _, tileData := range visibilityData {
		if tileData.Visible {
			tileKey := fmt.Sprintf("%d,%d", tileData.Row, tileData.Col)
			exploredMap.ExploredTiles[tileKey] = true
			visibleCount++
		}
	}

	// 距離に応じた段階的暗闇を描画
	drawGradualDarknessOverlay(world, screen, visibilityData)
}

func visionVertices(num int, x gc.Pixel, y gc.Pixel, r gc.Pixel) []ebiten.Vertex {
	vs := []ebiten.Vertex{}
	for i := 0; i < num; i++ {
		rate := float64(i) / float64(num)
		cr := 0.0
		cg := 0.0
		cb := 0.0
		vs = append(vs, ebiten.Vertex{
			DstX:   float32(float64(r)*math.Cos(2*math.Pi*rate)) + float32(x),
			DstY:   float32(float64(r)*math.Sin(2*math.Pi*rate)) + float32(y),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(cr),
			ColorG: float32(cg),
			ColorB: float32(cb),
			ColorA: 1,
		})
	}

	vs = append(vs, ebiten.Vertex{
		DstX:   float32(x),
		DstY:   float32(y),
		SrcX:   0,
		SrcY:   0,
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0,
	})

	return vs
}

// TileVisibility はタイルの可視性を表す
type TileVisibility struct {
	Row      int
	Col      int
	Visible  bool
	Distance float64
	Darkness float64 // 0.0（明るい）から 1.0（真っ暗）
}

// calculateTileVisibilityWithDistance はレイキャストでタイルごとの可視性と距離を計算する
func calculateTileVisibilityWithDistance(world w.World, playerX, playerY, radius gc.Pixel) map[string]TileVisibility {
	tileSize := 32 // タイルサイズ（固定値、実際はgameResourcesから取得すべき）

	visibilityMap := make(map[string]TileVisibility)

	// プレイヤーの位置からタイル座標を計算
	playerTileX := int(playerX) / tileSize
	playerTileY := int(playerY) / tileSize

	// 視界範囲内のタイルをチェック
	maxTileDistance := int(radius)/tileSize + 2

	for dx := -maxTileDistance; dx <= maxTileDistance; dx++ {
		for dy := -maxTileDistance; dy <= maxTileDistance; dy++ {
			tileX := playerTileX + dx
			tileY := playerTileY + dy

			// タイルの中心座標を計算
			tileCenterX := float64(tileX*tileSize + tileSize/2)
			tileCenterY := float64(tileY*tileSize + tileSize/2)

			// プレイヤーからタイル中心への距離をチェック
			distanceToTile := math.Sqrt(
				math.Pow(tileCenterX-float64(playerX), 2) +
					math.Pow(tileCenterY-float64(playerY), 2))

			tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

			if distanceToTile <= float64(radius) {
				// レイキャストでタイルが見えるかチェック
				visible := isTileVisibleByRaycast(world, float64(playerX), float64(playerY), tileCenterX, tileCenterY)

				// 距離に応じた暗闇の計算
				darkness := calculateDarknessByDistance(distanceToTile, float64(radius))

				visibilityMap[tileKey] = TileVisibility{
					Row:      tileX,
					Col:      tileY,
					Visible:  visible,
					Distance: distanceToTile,
					Darkness: darkness,
				}
			} else {
				visibilityMap[tileKey] = TileVisibility{
					Row:      tileX,
					Col:      tileY,
					Visible:  false,
					Distance: distanceToTile,
					Darkness: 1.0, // 視界外は完全に暗い
				}
			}
		}
	}

	return visibilityMap
}

// isTileVisibleByRaycast はレイキャストでタイルが見えるかチェック
func isTileVisibleByRaycast(world w.World, playerX, playerY, targetX, targetY float64) bool {
	// プレイヤーからターゲットへのベクトル
	dx := targetX - playerX
	dy := targetY - playerY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		return true // プレイヤーの位置は常に見える
	}

	// ターゲットタイル自体がBlockViewを持つかチェック
	targetIsWall := isBlockedByWall(world, gc.Pixel(targetX), gc.Pixel(targetY))

	// プレイヤーからターゲットまでの間にあるBlockViewタイルの数を計算
	blockViewCount := countBlockViewTilesBetween(world, playerX, playerY, targetX, targetY)

	if targetIsWall {
		// ターゲット自体が壁の場合、途中のBlockViewタイルが1つ以下なら見える
		return blockViewCount <= 1
	}

	// 通常のタイル（床など）の場合は、途中にBlockViewタイルがないときのみ見える
	return blockViewCount == 0
}

// countBlockViewTilesBetween はプレイヤーとターゲット間のBlockViewタイル数を数える
func countBlockViewTilesBetween(world w.World, playerX, playerY, targetX, targetY float64) int {
	dx := targetX - playerX
	dy := targetY - playerY
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		return 0
	}

	// 正規化
	stepX := dx / distance
	stepY := dy / distance

	blockViewCount := 0
	visitedTiles := make(map[string]bool) // 同じタイルを重複カウントしないため

	// レイキャストでBlockViewタイルを数える
	const stepSize = 2.0
	for step := stepSize; step < distance-stepSize; step += stepSize {
		rayX := playerX + stepX*step
		rayY := playerY + stepY*step

		// タイル座標に変換
		tileSize := 32.0
		tileX := int(rayX / tileSize)
		tileY := int(rayY / tileSize)
		tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

		// 既にチェック済みのタイルはスキップ
		if visitedTiles[tileKey] {
			continue
		}
		visitedTiles[tileKey] = true

		// タイルの中心座標
		tileCenterX := float64(tileX)*tileSize + tileSize/2
		tileCenterY := float64(tileY)*tileSize + tileSize/2

		// このタイルがBlockViewを持つかチェック
		if isBlockedByWall(world, gc.Pixel(tileCenterX), gc.Pixel(tileCenterY)) {
			blockViewCount++
		}
	}

	return blockViewCount
}

// calculateDarknessByDistance は距離に応じた暗闇レベルを計算する
func calculateDarknessByDistance(distance, maxRadius float64) float64 {
	if distance <= 0 {
		return 0.0 // プレイヤーの位置は完全に明るい
	}

	// 距離の正規化 (0.0-1.0)
	normalizedDistance := distance / maxRadius

	// 段階的な暗闇レベルを計算
	if normalizedDistance <= 0.3 {
		return 0.0 // 近い範囲は完全に明るい
	}
	if normalizedDistance <= 0.6 {
		return 0.2 // 中範囲は少し暗い
	}
	if normalizedDistance <= 0.8 {
		return 0.5 // 遠範囲は中程度暗い
	}
	return 0.8 // 端の方はかなり暗い
}

// drawGradualDarknessOverlay は距離に応じた段階的暗闇を描画する
func drawGradualDarknessOverlay(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	tileSize := 32 // タイルサイズ

	// カメラ位置を取得
	var cameraPos gc.Position
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		cameraPos = *world.Components.Position.Get(entity).(*gc.Position)
	}))

	// 段階的暗闃用の画像を初期化（キャッシュ）
	if len(darknessCacheImages) == 0 {
		initializeDarknessCache(tileSize)
	}

	// 画面上に表示されるタイル範囲を計算
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// カメラオフセットを考慮した画面範囲
	leftEdge := int(cameraPos.X) - screenWidth/2
	rightEdge := int(cameraPos.X) + screenWidth/2
	topEdge := int(cameraPos.Y) - screenHeight/2
	bottomEdge := int(cameraPos.Y) + screenHeight/2

	// タイル範囲に変換
	startTileX := leftEdge/tileSize - 1
	endTileX := rightEdge/tileSize + 1
	startTileY := topEdge/tileSize - 1
	endTileY := bottomEdge/tileSize + 1

	// 各タイルをチェックして暗闇を描画
	for tileX := startTileX; tileX <= endTileX; tileX++ {
		for tileY := startTileY; tileY <= endTileY; tileY++ {
			tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

			// タイルの暗闇レベルを取得
			darkness := 1.0 // デフォルトは完全に暗い
			if tileData, exists := visibilityData[tileKey]; exists {
				if tileData.Visible {
					darkness = tileData.Darkness
				} else {
					darkness = 1.0 // 見えないタイルは完全に暗い
				}
			}

			// 暗闇レベルが0より大きい場合のみ描画
			if darkness > 0.0 {
				// タイルの画面座標を計算
				screenX := float64(tileX*tileSize) - float64(cameraPos.X) + float64(screenWidth)/2
				screenY := float64(tileY*tileSize) - float64(cameraPos.Y) + float64(screenHeight)/2

				// 暗闇レベルに応じた画像を選択して描画
				drawDarknessAtLevel(screen, screenX, screenY, darkness)
			}
		}
	}
}

// initializeDarknessCache は段階的暗闃用の画像キャッシュを初期化する
func initializeDarknessCache(tileSize int) {
	// 異なる暗闇レベルの画像を作成
	darknessCacheImages = make([]*ebiten.Image, 5)

	// 各暗闇レベルの画像を作成
	darknessCacheImages[0] = nil // 0.0: 暗闇なし

	// 0.2: 少し暗い
	darknessCacheImages[1] = ebiten.NewImage(tileSize, tileSize)
	darknessCacheImages[1].Fill(color.RGBA{0, 0, 0, 51}) // 約20%の透明度

	// 0.5: 中程度暗い
	darknessCacheImages[2] = ebiten.NewImage(tileSize, tileSize)
	darknessCacheImages[2].Fill(color.RGBA{0, 0, 0, 128}) // 約50%の透明度

	// 0.8: かなり暗い
	darknessCacheImages[3] = ebiten.NewImage(tileSize, tileSize)
	darknessCacheImages[3].Fill(color.RGBA{0, 0, 0, 204}) // 約80%の透明度

	// 1.0: 完全に暗い
	darknessCacheImages[4] = ebiten.NewImage(tileSize, tileSize)
	darknessCacheImages[4].Fill(color.RGBA{0, 0, 0, 255}) // 完全不透明
}

// drawDarknessAtLevel は指定した暗闇レベルで暗闇を描画する
func drawDarknessAtLevel(screen *ebiten.Image, x, y, darkness float64) {
	var darknessImage *ebiten.Image

	// 暗闇レベルに応じた画像を選択
	if darkness <= 0.0 {
		return // 暗闇なし
	} else if darkness <= 0.2 {
		darknessImage = darknessCacheImages[1]
	} else if darkness <= 0.5 {
		darknessImage = darknessCacheImages[2]
	} else if darkness <= 0.8 {
		darknessImage = darknessCacheImages[3]
	} else {
		darknessImage = darknessCacheImages[4]
	}

	if darknessImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(darknessImage, op)
	}
}

// isBlockedByWall は指定した位置に視界を遮る壁があるかチェックする
func isBlockedByWall(world w.World, x, y gc.Pixel) bool {
	var blocked bool

	// Position を持つエンティティをチェック
	world.Manager.Join(
		world.Components.Position,
		world.Components.BlockView,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if blocked {
			return
		}

		pos := world.Components.Position.Get(entity).(*gc.Position)
		sprite := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		spriteInfo := sprite.SpriteSheet.Sprites[sprite.SpriteNumber]

		// スプライトの境界ボックス
		left := float64(pos.X) - float64(spriteInfo.Width)/2
		right := float64(pos.X) + float64(spriteInfo.Width)/2
		top := float64(pos.Y) - float64(spriteInfo.Height)/2
		bottom := float64(pos.Y) + float64(spriteInfo.Height)/2

		if float64(x) >= left && float64(x) <= right &&
			float64(y) >= top && float64(y) <= bottom {
			blocked = true
		}
	}))

	// GridElement を持つエンティティもチェック
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.BlockView,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if blocked {
			return
		}

		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)
		sprite := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
		spriteInfo := sprite.SpriteSheet.Sprites[sprite.SpriteNumber]

		// グリッド位置をピクセル座標に変換
		gridX := int(grid.Row) * spriteInfo.Width
		gridY := int(grid.Col) * spriteInfo.Height

		// グリッドの境界ボックス
		left := float64(gridX)
		right := float64(gridX + spriteInfo.Width)
		top := float64(gridY)
		bottom := float64(gridY + spriteInfo.Height)

		if float64(x) >= left && float64(x) <= right &&
			float64(y) >= top && float64(y) <= bottom {
			blocked = true
		}
	}))

	return blocked
}
