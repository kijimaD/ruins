package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// 段階的な暗闇用の画像キャッシュ（透明度レベルを増加）
	darknessCacheImages []*ebiten.Image

	// プレイヤー位置キャッシュ（4px移動ごとに更新）
	playerPositionCache struct {
		lastPlayerX    gc.Pixel
		lastPlayerY    gc.Pixel
		visibilityData map[string]TileVisibility
		isInitialized  bool
	}

	// レイキャスト結果のキャッシュ
	raycastCache = make(map[string]bool) // "x1,y1,x2,y2" -> visible
)

// abs は絶対値を返す
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ClearVisionCaches は全ての視界関連キャッシュをクリアする（階移動時などに使用）
func ClearVisionCaches() {
	// プレイヤー位置キャッシュをクリア
	playerPositionCache.isInitialized = false
	playerPositionCache.visibilityData = nil

	// レイキャストキャッシュをクリア
	raycastCache = make(map[string]bool)
}

// VisionSystem はタイルごとの視界を管理し、暗闇を描画する
func VisionSystem(world w.World, screen *ebiten.Image) {
	// プレイヤー位置を取得
	var playerGridElement *gc.GridElement
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerGridElement = world.Components.GridElement.Get(entity).(*gc.GridElement)
	}))

	if playerGridElement == nil {
		return
	}

	// タイル座標をピクセル座標に変換
	playerPos := &gc.Position{
		X: gc.Pixel(int(playerGridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2),
		Y: gc.Pixel(int(playerGridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2),
	}

	// 移動ごとの視界更新判定
	const updateThreshold = 4
	needsUpdate := !playerPositionCache.isInitialized ||
		abs(int(playerPos.X-playerPositionCache.lastPlayerX)) >= updateThreshold ||
		abs(int(playerPos.Y-playerPositionCache.lastPlayerY)) >= updateThreshold

	var visibilityData map[string]TileVisibility

	if needsUpdate {
		// Dungeonリソースから探索済みマップを取得

		// タイルの可視性マップを更新
		visionRadius := gc.Pixel(320)
		visibilityData = calculateTileVisibilityWithDistance(world, playerPos.X, playerPos.Y, visionRadius)

		// 視界内のタイルを探索済みとしてマーク
		for _, tileData := range visibilityData {
			if tileData.Visible {
				// GridElementを直接キーとして使用
				gridElement := gc.GridElement{X: gc.Tile(tileData.Col), Y: gc.Tile(tileData.Row)}
				world.Resources.Dungeon.ExploredTiles[gridElement] = true
			}
		}

		// キャッシュ更新
		playerPositionCache.lastPlayerX = playerPos.X
		playerPositionCache.lastPlayerY = playerPos.Y
		playerPositionCache.visibilityData = visibilityData
		playerPositionCache.isInitialized = true
	} else {
		// キャッシュされたデータを使用
		visibilityData = playerPositionCache.visibilityData
	}

	// 距離に応じた段階的暗闇を描画
	drawDistanceBasedDarkness(world, screen, visibilityData)
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

	// 視界範囲を分割して段階的処理（視界範囲最適化）
	maxTileDistance := int(radius)/tileSize + 2

	// タイルベース視界判定（Dark Days Ahead風）

	for dx := -maxTileDistance; dx <= maxTileDistance; dx++ {
		for dy := -maxTileDistance; dy <= maxTileDistance; dy++ {
			tileX := playerTileX + dx
			tileY := playerTileY + dy

			// 早期距離チェック（枝払い）
			if abs(dx) > maxTileDistance || abs(dy) > maxTileDistance {
				continue
			}

			// タイルの中心座標を計算
			tileCenterX := float64(tileX*tileSize + tileSize/2)
			tileCenterY := float64(tileY*tileSize + tileSize/2)

			// プレイヤーからタイル中心への距離をチェック（平方根計算の最適化）
			dxF := tileCenterX - float64(playerX)
			dyF := tileCenterY - float64(playerY)
			distanceSquared := dxF*dxF + dyF*dyF
			radiusSquared := float64(radius) * float64(radius)

			tileKey := fmt.Sprintf("%d,%d", tileX, tileY)

			// 視界範囲内のタイルのみ処理
			if distanceSquared <= radiusSquared {
				distanceToTile := math.Sqrt(distanceSquared)

				// Dark Days Ahead風の統一されたタイルベース視界判定
				visible := isTileVisibleByRaycast(world, float64(playerX), float64(playerY), tileCenterX, tileCenterY)

				// 距離に応じた暗闇の計算
				darkness := calculateDarknessByDistance(distanceToTile, float64(radius))

				visibilityMap[tileKey] = TileVisibility{
					Row:      tileY,
					Col:      tileX,
					Visible:  visible,
					Distance: distanceToTile,
					Darkness: darkness,
				}
			}
			// 視界外のタイルは処理しない（最適化）
		}
	}

	return visibilityMap
}

// isTileVisibleByRaycast はタイルベース視界判定
func isTileVisibleByRaycast(world w.World, playerX, playerY, targetX, targetY float64) bool {
	// キャッシュキーを生成
	px := int(playerX/4) * 4
	py := int(playerY/4) * 4
	tx := int(targetX/4) * 4
	ty := int(targetY/4) * 4
	cacheKey := fmt.Sprintf("%d,%d,%d,%d", px, py, tx, ty)

	// キャッシュから結果をチェック
	if result, exists := raycastCache[cacheKey]; exists {
		return result
	}

	// タイル座標に変換
	const tileSize = 32.0
	playerTileX := int(playerX / tileSize)
	playerTileY := int(playerY / tileSize)
	targetTileX := int(targetX / tileSize)
	targetTileY := int(targetY / tileSize)

	// 同じタイルまたは隣接タイルは常に見える
	if abs(targetTileX-playerTileX) <= 1 && abs(targetTileY-playerTileY) <= 1 {
		raycastCache[cacheKey] = true
		return true
	}

	// ブレゼンハムのライン描画アルゴリズムでタイルベースの視線判定
	result := bresenhamLineOfSight(world, playerTileX, playerTileY, targetTileX, targetTileY)

	// 結果をキャッシュ
	if len(raycastCache) < 15000 {
		raycastCache[cacheKey] = result
	}

	return result
}

// bresenhamLineOfSight はブレゼンハムアルゴリズムを使用したタイルベース視線判定
func bresenhamLineOfSight(world w.World, x0, y0, x1, y1 int) bool {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy
	x, y := x0, y0

	for {
		// ターゲットに到達したら見える
		if x == x1 && y == y1 {
			return true
		}

		// 現在のタイルが壁かチェック（ターゲット以外）
		if x != x1 || y != y1 {
			tileCenterX := float64(x*int(consts.TileSize) + int(consts.TileSize)/2)
			tileCenterY := float64(y*int(consts.TileSize) + int(consts.TileSize)/2)
			if isBlockedByWall(world, gc.Pixel(tileCenterX), gc.Pixel(tileCenterY)) {
				return false // 壁に遮られている
			}
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// calculateDarknessByDistance は距離に応じた暗闇レベルを計算する
func calculateDarknessByDistance(distance, maxRadius float64) float64 {
	if distance <= 0 {
		return 0.0 // プレイヤーの位置は完全に明るい
	}

	// 距離の正規化 (0.0-1.0)
	normalizedDistance := distance / maxRadius
	if normalizedDistance >= 1.0 {
		return 0.95 // 最遠距離でも完全に真っ暗にはしない
	}

	// 滑らかな二次カーブによる減衰（中心が明るく、外側に向かって滑らかに暗くなる）
	// 0.2までは完全に明るい（コア照明領域）
	if normalizedDistance <= 0.2 {
		return 0.0
	}

	// 0.2から1.0にかけて滑らかに暗くなる
	// 二次関数: y = ((x-0.2) / 0.8)^1.5 * 0.95
	adjustedDistance := (normalizedDistance - 0.2) / 0.8
	darkness := math.Pow(adjustedDistance, 1.5) * 0.95

	return darkness
}

// drawDistanceBasedDarkness は距離に応じた段階的暗闇を描画する
func drawDistanceBasedDarkness(world w.World, screen *ebiten.Image, visibilityData map[string]TileVisibility) {
	tileSize := int(consts.TileSize)

	// カメラ位置とスケールを取得
	var cameraPos gc.Position
	cameraScale := 1.0 // デフォルトスケール

	// カメラのGridElementから位置を取得
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

	// 段階的暗闃用の画像を初期化（キャッシュ）
	if len(darknessCacheImages) == 0 {
		initializeDarknessCache(tileSize)
	}

	// 画面上に表示されるタイル範囲を計算
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// スケールを考慮した実際の表示範囲を計算
	actualScreenWidth := int(float64(screenWidth) / cameraScale)
	actualScreenHeight := int(float64(screenHeight) / cameraScale)

	// カメラオフセットを考慮した画面範囲
	leftEdge := int(cameraPos.X) - actualScreenWidth/2
	rightEdge := int(cameraPos.X) + actualScreenWidth/2
	topEdge := int(cameraPos.Y) - actualScreenHeight/2
	bottomEdge := int(cameraPos.Y) + actualScreenHeight/2

	// タイル範囲に変換
	startTileX := leftEdge/tileSize - 1
	endTileX := rightEdge/tileSize + 1
	startTileY := topEdge/tileSize - 1
	endTileY := bottomEdge/tileSize + 1

	// 距離に応じた段階的暗闇を描画
	for tileX := startTileX; tileX <= endTileX; tileX++ {
		for tileY := startTileY; tileY <= endTileY; tileY++ {
			tileKey := fmt.Sprintf("%d,%d", tileX, tileY)
			gridElement := gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}

			var darkness float64

			// 視界データをチェック
			if tileData, exists := visibilityData[tileKey]; exists {
				if tileData.Visible {
					// 視界内: 距離に応じた暗闇レベル
					darkness = tileData.Darkness
				} else {
					// 視界範囲内だが見えないタイル: 完全に暗い
					darkness = 1.0
				}
			} else {
				// 視界範囲外: 探索済みかチェック
				if explored := world.Resources.Dungeon.ExploredTiles[gridElement]; explored {
					// 探索済み視界外タイル: 真っ黒
					darkness = 1.0
				} else {
					// 未探索タイル: 描画しない（完全に隠れる）
					continue
				}
			}

			// 暗闇レベルが0より大きい場合のみ描画
			if darkness > 0.0 {
				// タイルの画面座標を計算（スケールを考慮）
				worldX := float64(tileX * tileSize)
				worldY := float64(tileY * tileSize)
				screenX := (worldX-float64(cameraPos.X))*cameraScale + float64(screenWidth)/2
				screenY := (worldY-float64(cameraPos.Y))*cameraScale + float64(screenHeight)/2

				// 暗闇レベルに応じた画像を描画
				drawDarknessAtLevelWithScale(screen, screenX, screenY, darkness, cameraScale)
			}
		}
	}
}

// initializeDarknessCache は段階的暗闃用の画像キャッシュを初期化する
func initializeDarknessCache(tileSize int) {
	// tileSizeが0以下の場合は初期化しない
	if tileSize <= 0 {
		return
	}

	// より多くの暗闇レベルの画像を作成（10段階）
	darknessCacheImages = make([]*ebiten.Image, 11)

	// 各暗闇レベルの画像を作成
	darknessCacheImages[0] = nil // 0.0: 暗闇なし

	// 0.1から1.0まで0.1刻みで10段階作成
	for i := 1; i <= 10; i++ {
		darkness := float64(i) * 0.1
		alpha := uint8(darkness * 255) // 透明度を0-255に変換

		darknessCacheImages[i] = ebiten.NewImage(tileSize, tileSize)

		// 最外側（80%以上の暗さ）は真っ黒、それ以外は暖色系
		if darkness >= 0.8 {
			// 最外側は純黒
			darknessCacheImages[i].Fill(color.RGBA{0, 0, 0, alpha})
		} else {
			// 内側は暖色系の暗闇（微かな茶色がかった暗闇）
			r := uint8(float64(alpha) * 0.15) // 透明度の15%の赤成分
			g := uint8(float64(alpha) * 0.10) // 透明度の10%の緑成分
			b := uint8(float64(alpha) * 0.05) // 透明度の5%の青成分
			darknessCacheImages[i].Fill(color.RGBA{r, g, b, alpha})
		}
	}
}

// drawDarknessAtLevelWithScale は指定した暗闇レベルで暗闇を描画する（スケール対応）
func drawDarknessAtLevelWithScale(screen *ebiten.Image, x, y, darkness, scale float64) {
	if darkness <= 0.0 {
		return // 暗闇なし
	}

	// 暗闇レベルを配列インデックスに変換（0.1刻みで10段階）
	index := int(darkness*10 + 0.5) // 四捨五入
	if index > 10 {
		index = 10
	}
	if index < 1 {
		index = 1
	}

	darknessImage := darknessCacheImages[index]
	if darknessImage != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(x, y)
		screen.DrawImage(darknessImage, op)
	}
}

// GetCurrentVisibilityData は現在の視界データを返す（レンダリング用）
func GetCurrentVisibilityData() map[string]TileVisibility {
	if playerPositionCache.isInitialized {
		return playerPositionCache.visibilityData
	}
	return nil
}

// isBlockedByWall は直接的な壁チェック
func isBlockedByWall(world w.World, x, y gc.Pixel) bool {
	fx, fy := float64(x), float64(y)

	// GridElement + BlockView のチェック（32x32タイル）
	tileX := int(fx / 32)
	tileY := int(fy / 32)

	// より正確なチェック: GridElementから直接確認
	blocked := false
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.BlockView,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(grid.X) == tileX && int(grid.Y) == tileY {
			blocked = true
		}
	}))

	if blocked {
		return true
	}

	// GridElement + BlockView のチェック
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.BlockView,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// タイル単位での視界ブロックチェック
		tilePixelX := float64(int(gridElement.X)*int(consts.TileSize) + int(consts.TileSize)/2)
		tilePixelY := float64(int(gridElement.Y)*int(consts.TileSize) + int(consts.TileSize)/2)

		// タイル境界での判定
		left := tilePixelX - float64(consts.TileSize)/2
		right := tilePixelX + float64(consts.TileSize)/2
		top := tilePixelY - float64(consts.TileSize)/2
		bottom := tilePixelY + float64(consts.TileSize)/2

		if fx >= left && fx <= right && fy >= top && fy <= bottom {
			blocked = true
		}
	}))

	return blocked
}
