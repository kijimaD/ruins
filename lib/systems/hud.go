package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kijimaD/ruins/lib/camera"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	// AI視界可視化用の円形画像
	aiVisionCircleImage *ebiten.Image
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
	}

	// ミニマップを描画
	drawMinimap(world, screen)
}

// drawAIStates はAIエンティティのステートをスプライトの近くに表示する
func drawAIStates(world w.World, screen *ebiten.Image) {
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

		// スプライトの上にテキストを表示
		// カメラのオフセットを考慮して座標を計算
		var cameraPos gc.Position
		world.Manager.Join(
			world.Components.Camera,
			world.Components.Position,
		).Visit(ecs.Visit(func(camEntity ecs.Entity) {
			cameraPos = *world.Components.Position.Get(camEntity).(*gc.Position)
		}))

		// 画面座標に変換
		screenX := float64(position.X-cameraPos.X) + float64(world.Resources.ScreenDimensions.Width)/2
		screenY := float64(position.Y-cameraPos.Y) + float64(world.Resources.ScreenDimensions.Height)/2

		// スプライトの上20ピクセルの位置にテキスト表示
		ebitenutil.DebugPrintAt(screen, stateText, int(screenX)-20, int(screenY)-30)
	}))
}

// drawAIVisionRanges はデバッグ時にAIの視界範囲を描画する
func drawAIVisionRanges(world w.World, screen *ebiten.Image) {
	// AI視界可視化用の円形画像を初期化（1回だけ）
	if aiVisionCircleImage == nil {
		// 円の直径を最大視界距離の2倍（300ピクセル）に設定
		size := 300
		aiVisionCircleImage = ebiten.NewImage(size, size)

		// 半透明の緑色で円を描画
		radius := float64(size / 2)
		center := float64(size / 2)

		// 円形の頂点を作成
		vertices := []ebiten.Vertex{}
		indices := []uint16{}

		// 中心点
		vertices = append(vertices, ebiten.Vertex{
			DstX:   float32(center),
			DstY:   float32(center),
			SrcX:   0,
			SrcY:   0,
			ColorR: 0.0,
			ColorG: 1.0,
			ColorB: 0.0,
			ColorA: 0.3, // 半透明
		})

		// 円周上の点
		circlePoints := 32
		for i := 0; i < circlePoints; i++ {
			angle := 2 * math.Pi * float64(i) / float64(circlePoints)
			x := center + radius*math.Cos(angle)
			y := center + radius*math.Sin(angle)

			vertices = append(vertices, ebiten.Vertex{
				DstX:   float32(x),
				DstY:   float32(y),
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
		aiVisionCircleImage.DrawTriangles(vertices, indices, whiteImg, opt)
	}

	// 各NPCの視界を描画
	world.Manager.Join(
		world.Components.Position,
		world.Components.AIVision,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		position := world.Components.Position.Get(entity).(*gc.Position)
		vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

		// 視界範囲に応じてスケールを調整
		scale := vision.ViewDistance / 300.0 // 300は基準の視界距離

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(
			float64(position.X)-(300.0*scale/2.0), // 中心に配置
			float64(position.Y)-(300.0*scale/2.0),
		)
		camera.SetTranslate(world, op)

		screen.DrawImage(aiVisionCircleImage, op)
	}))
}

// drawMinimap はミニマップを画面右上に描画する
func drawMinimap(world w.World, screen *ebiten.Image) {
	// プレイヤー位置とExploredMapを取得
	var playerPos *gc.Position
	var exploredMap *gc.ExploredMap

	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
		world.Components.ExploredMap,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
		exploredMap = world.Components.ExploredMap.Get(entity).(*gc.ExploredMap)
	}))

	if playerPos == nil || exploredMap == nil {
		return // 必要なコンポーネントがない場合は描画しない
	}

	// 探索済みタイルがない場合でも、ミニマップの枠だけは表示する
	if len(exploredMap.ExploredTiles) == 0 {
		// 空のミニマップを描画
		drawEmptyMinimap(world, screen)
		return
	}

	// ミニマップの設定
	minimapWidth := 150
	minimapHeight := 150
	minimapScale := 3 // 1タイルを3ピクセルで表現
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
	for tileKey := range exploredMap.ExploredTiles {
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
	minimapWidth := 150
	minimapHeight := 150
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
