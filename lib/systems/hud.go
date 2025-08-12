package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// HUDSystem はゲームの HUD 情報を描画する
func HUDSystem(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("floor: B%d", gameResources.Depth), 0, 200)

	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("speed: %.2f", velocity.Speed), 0, 220)
	}))

	// AI デバッグ表示フラグが有効な時のみAI情報表示
	cfg := config.Get()
	if cfg.ShowAIDebug {
		drawAIVisionRanges(world, screen)
		drawAIStates(world, screen)
		drawAIMovementDirections(world, screen)
	}

	// ミニマップを描画
	drawMinimap(world, screen)
}

// drawAIStates はAIエンティティのステートをスプライトの近くに表示する
func drawAIStates(world w.World, screen *ebiten.Image) {
	// カメラ位置とスケールを取得
	var cameraPos gc.Position
	var cameraScale float64
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(camEntity ecs.Entity) {
		cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		camera := world.Components.Camera.Get(camEntity).(*gc.Camera)
		cameraScale = camera.Scale
	}))

	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	world.Manager.Join(
		world.Components.Position,
		world.Components.AIMoveFSM,
		world.Components.AIRoaming,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		roaming := world.Components.AIRoaming.Get(entity).(*gc.AIRoaming)

		// AIの現在の状態を判定
		var stateText string
		if entity.HasComponent(world.Components.AIChasing) {
			stateText = "CHASING"
		} else {
			switch roaming.SubState {
			case gc.AIRoamingWaiting:
				stateText = "WAITING"
			case gc.AIRoamingDriving:
				stateText = "ROAMING"
			case gc.AIRoamingChasing:
				stateText = "CHASING"
			default:
				stateText = "UNKNOWN"
			}
		}

		// カメラスケールを考慮した画面座標に変換
		screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenWidth)/2
		screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenHeight)/2

		// スプライトの上にテキスト表示（テキスト位置もスケールを考慮）
		textOffsetY := 30.0 * cameraScale
		ebitenutil.DebugPrintAt(screen, stateText, int(screenX)-20, int(screenY-textOffsetY))
	}))
}

// drawAIVisionRanges はデバッグ時にAIの視界範囲を描画する
func drawAIVisionRanges(world w.World, screen *ebiten.Image) {
	// カメラ位置とスケールを取得
	var cameraPos gc.Position
	var cameraScale float64
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(camEntity ecs.Entity) {
		cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		camera := world.Components.Camera.Get(camEntity).(*gc.Camera)
		cameraScale = camera.Scale
	}))

	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// 各NPCの視界を描画
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIVision,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

		// カメラスケールを考慮した画面座標に変換
		screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenWidth)/2
		screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenHeight)/2

		// カメラスケールを考慮した視界円の半径
		scaledRadius := float32(float64(vision.ViewDistance) * cameraScale)

		// AIVision.ViewDistanceを直接使用して視界円を描画
		drawVisionCircle(screen, float32(screenX), float32(screenY), scaledRadius)
	}))
}

// drawVisionCircle は指定した位置と半径で視界円を描画する
func drawVisionCircle(screen *ebiten.Image, centerX, centerY, radius float32) {
	// 円周上の点数
	circlePoints := 32
	vertices := []ebiten.Vertex{}
	indices := []uint16{}

	// 中心点
	vertices = append(vertices, ebiten.Vertex{
		DstX:   centerX,
		DstY:   centerY,
		SrcX:   0,
		SrcY:   0,
		ColorR: 0.0,
		ColorG: 1.0,
		ColorB: 0.0,
		ColorA: 0.3, // 半透明
	})

	// 円周上の点
	for i := 0; i < circlePoints; i++ {
		angle := 2 * math.Pi * float64(i) / float64(circlePoints)
		x := centerX + radius*float32(math.Cos(angle))
		y := centerY + radius*float32(math.Sin(angle))

		vertices = append(vertices, ebiten.Vertex{
			DstX:   x,
			DstY:   y,
			SrcX:   0,
			SrcY:   0,
			ColorR: 0.0,
			ColorG: 1.0,
			ColorB: 0.0,
			ColorA: 0.3,
		})

		// 三角形のインデックス
		if i < circlePoints {
			indices = append(indices, 0, uint16(i+1), uint16((i+1)%circlePoints+1))
		}
	}

	// 円を描画
	opt := &ebiten.DrawTrianglesOptions{}
	// 1x1ピクセルの白い画像を作成
	whiteImg := ebiten.NewImage(1, 1)
	whiteImg.Fill(color.White)
	screen.DrawTriangles(vertices, indices, whiteImg, opt)
}

// drawAIMovementDirections はAIの進行方向を矢印で表示する
func drawAIMovementDirections(world w.World, screen *ebiten.Image) {
	// カメラ位置とスケールを取得
	var cameraPos gc.Position
	var cameraScale float64
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(camEntity ecs.Entity) {
		cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		camera := world.Components.Camera.Get(camEntity).(*gc.Camera)
		cameraScale = camera.Scale
	}))

	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// AIエンティティの移動方向を描画
	world.Manager.Join(
		world.Components.Position,
		world.Components.Velocity,
		world.Components.AIMoveFSM,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)

		// カメラスケールを考慮した画面座標に変換
		screenX := (float64(position.X)-float64(cameraPos.X))*cameraScale + float64(screenWidth)/2
		screenY := (float64(position.Y)-float64(cameraPos.Y))*cameraScale + float64(screenHeight)/2

		// 移動方向がある場合のみ描画
		if velocity.Speed > 0 && velocity.ThrottleMode == gc.ThrottleModeFront {
			drawDirectionArrow(screen, screenX, screenY, velocity.Angle, velocity.Speed, cameraScale)
		}
	}))
}

// drawDirectionArrow は指定した位置に進行方向の矢印を描画する
func drawDirectionArrow(screen *ebiten.Image, x, y, angle, speed, cameraScale float64) {
	// 矢印の長さを速度に応じて調整（最小20、最大60ピクセル）、カメラスケールも考慮
	baseLength := 20.0 + speed*20.0
	if baseLength > 60 {
		baseLength = 60
	}
	length := baseLength * cameraScale

	// 角度をラジアンに変換
	radians := angle * math.Pi / 180

	// 矢印の先端位置
	endX := x + length*math.Cos(radians)
	endY := y + length*math.Sin(radians)

	// 線の太さもカメラスケールに応じて調整
	strokeWidth := float32(2.0 * cameraScale)
	if strokeWidth < 1.0 {
		strokeWidth = 1.0
	}

	// メインラインを描画（緑色）
	vector.StrokeLine(screen, float32(x), float32(y), float32(endX), float32(endY), strokeWidth, color.RGBA{0, 255, 0, 255}, false)

	// 矢印の頭部を描画
	arrowHeadLength := 10.0 * cameraScale
	leftAngle := radians + 2.5  // 約145度
	rightAngle := radians - 2.5 // 約-145度

	leftX := endX + arrowHeadLength*math.Cos(leftAngle)
	leftY := endY + arrowHeadLength*math.Sin(leftAngle)
	rightX := endX + arrowHeadLength*math.Cos(rightAngle)
	rightY := endY + arrowHeadLength*math.Sin(rightAngle)

	// 矢印の頭部ラインを描画
	vector.StrokeLine(screen, float32(endX), float32(endY), float32(leftX), float32(leftY), strokeWidth, color.RGBA{0, 255, 0, 255}, false)
	vector.StrokeLine(screen, float32(endX), float32(endY), float32(rightX), float32(rightY), strokeWidth, color.RGBA{0, 255, 0, 255}, false)
}

// drawMinimap はミニマップを画面右上に描画する
func drawMinimap(world w.World, screen *ebiten.Image) {
	// プレイヤー位置を取得
	var playerPos *gc.Position

	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	if playerPos == nil {
		return // プレイヤーが見つからない場合は描画しない
	}

	// Dungeonリソースから探索済みマップとミニマップ設定を取得
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	// 探索済みタイルがない場合でも、ミニマップの枠だけは表示する
	if len(gameResources.ExploredTiles) == 0 {
		// 空のミニマップを描画
		drawEmptyMinimap(world, screen)
		return
	}

	// ミニマップの設定
	minimapWidth := gameResources.Minimap.Width
	minimapHeight := gameResources.Minimap.Height
	minimapScale := gameResources.Minimap.Scale // 1タイルをscaleピクセルで表現
	screenWidth := world.Resources.ScreenDimensions.Width
	minimapX := screenWidth - minimapWidth - 10 // 画面右端から10ピクセル内側
	minimapY := 10                              // 画面上端から10ピクセル下

	// ミニマップの背景を描画（半透明の黒い四角）
	minimapBg := ebiten.NewImage(minimapWidth, minimapHeight)
	minimapBg.Fill(color.RGBA{0, 0, 0, 128}) // 半透明の黒
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(minimapX), float64(minimapY))
	screen.DrawImage(minimapBg, op)

	// プレイヤーの現在位置をタイル座標に変換
	tileSize := 32 // タイルサイズ
	playerTileX := int(playerPos.X) / tileSize
	playerTileY := int(playerPos.Y) / tileSize

	// ミニマップの中心をプレイヤー位置に合わせる
	centerX := minimapX + minimapWidth/2
	centerY := minimapY + minimapHeight/2

	// 探索済みタイルを描画
	for tileKey := range gameResources.ExploredTiles {
		var tileX, tileY int
		if _, err := fmt.Sscanf(tileKey, "%d,%d", &tileX, &tileY); err != nil {
			continue
		}

		// プレイヤー位置からの相対位置を計算
		relativeX := tileX - playerTileX
		relativeY := tileY - playerTileY

		// ミニマップ上の座標を計算
		mapX := float32(centerX + relativeX*minimapScale)
		mapY := float32(centerY + relativeY*minimapScale)

		// ミニマップの範囲内かチェック
		if mapX >= float32(minimapX) && mapX <= float32(minimapX+minimapWidth-minimapScale) &&
			mapY >= float32(minimapY) && mapY <= float32(minimapY+minimapHeight-minimapScale) {

			// タイルのタイプに応じて色を決定
			tileColor := getTileColorForMinimap(world, tileX, tileY)

			// 小さな四角形でタイルを表現
			vector.DrawFilledRect(screen, mapX, mapY, float32(minimapScale), float32(minimapScale), tileColor, false)
		}
	}

	// プレイヤーの位置を赤い点で表示
	playerMapX := float32(centerX)
	playerMapY := float32(centerY)
	vector.DrawFilledCircle(screen, playerMapX, playerMapY, 2, color.RGBA{255, 0, 0, 255}, false)

	// ミニマップの枠を描画
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY-1), 1, float32(minimapHeight+2), color.RGBA{255, 255, 255, 255}, false)            // 左
	vector.DrawFilledRect(screen, float32(minimapX+minimapWidth), float32(minimapY-1), 1, float32(minimapHeight+2), color.RGBA{255, 255, 255, 255}, false) // 右
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY-1), float32(minimapWidth+2), 1, color.RGBA{255, 255, 255, 255}, false)             // 上
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY+minimapHeight), float32(minimapWidth+2), 1, color.RGBA{255, 255, 255, 255}, false) // 下
}

// getTileColorForMinimap はタイルの種類に応じてミニマップ上の色を返す
func getTileColorForMinimap(world w.World, tileX, tileY int) color.RGBA {
	// そのタイル位置に実際にエンティティが存在するかチェック
	hasWall := false
	hasFloor := false

	// GridElement を持つエンティティをチェック
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		grid := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// グリッドの座標がタイル座標と一致するかチェック
		if int(grid.Row) == tileX && int(grid.Col) == tileY {
			// このタイルにエンティティが存在する
			if entity.HasComponent(world.Components.BlockView) {
				hasWall = true
			} else {
				hasFloor = true
			}
		}
	}))

	// 実際にエンティティが存在する場合のみ描画
	if hasWall {
		return color.RGBA{100, 100, 100, 255} // 壁は灰色
	} else if hasFloor {
		return color.RGBA{200, 200, 200, 128} // 床は薄い灰色
	}

	// 何もない場所は描画しない（透明）
	return color.RGBA{0, 0, 0, 0} // 透明
}

// drawEmptyMinimap は空のミニマップ（枠のみ）を描画する
func drawEmptyMinimap(world w.World, screen *ebiten.Image) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	minimapWidth := gameResources.Minimap.Width
	minimapHeight := gameResources.Minimap.Height
	screenWidth := world.Resources.ScreenDimensions.Width
	minimapX := screenWidth - minimapWidth - 10
	minimapY := 10

	// ミニマップの背景を描画（半透明の黒い四角）
	minimapBg := ebiten.NewImage(minimapWidth, minimapHeight)
	minimapBg.Fill(color.RGBA{0, 0, 0, 128})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(minimapX), float64(minimapY))
	screen.DrawImage(minimapBg, op)

	// ミニマップの枠を描画
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY-1), 1, float32(minimapHeight+2), color.RGBA{255, 255, 255, 255}, false)            // 左
	vector.DrawFilledRect(screen, float32(minimapX+minimapWidth), float32(minimapY-1), 1, float32(minimapHeight+2), color.RGBA{255, 255, 255, 255}, false) // 右
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY-1), float32(minimapWidth+2), 1, color.RGBA{255, 255, 255, 255}, false)             // 上
	vector.DrawFilledRect(screen, float32(minimapX-1), float32(minimapY+minimapHeight), float32(minimapWidth+2), 1, color.RGBA{255, 255, 255, 255}, false) // 下

	// 中央に"No Data"テキストを表示
	ebitenutil.DebugPrintAt(screen, "No Data", minimapX+50, minimapY+70)
}
