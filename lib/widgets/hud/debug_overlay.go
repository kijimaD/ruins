package hud

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// DebugOverlay はAI情報のデバッグ表示エリア
type DebugOverlay struct {
	enabled bool
}

// NewDebugOverlay は新しいHUDDebugOverlayを作成する
func NewDebugOverlay() *DebugOverlay {
	return &DebugOverlay{
		enabled: true,
	}
}

// SetEnabled は有効/無効を設定する
func (overlay *DebugOverlay) SetEnabled(enabled bool) {
	overlay.enabled = enabled
}

// Update はデバッグオーバーレイを更新する
func (overlay *DebugOverlay) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はデバッグオーバーレイを描画する（AI情報を表示）
func (overlay *DebugOverlay) Draw(world w.World, screen *ebiten.Image) {
	if !overlay.enabled {
		return
	}

	// AI情報を描画
	overlay.drawAIVisionRanges(world, screen)
	overlay.drawAIStates(world, screen)
	overlay.drawAIMovementDirections(world, screen)
}

// drawAIStates はAIエンティティのステートをスプライトの近くに表示する
func (overlay *DebugOverlay) drawAIStates(world w.World, screen *ebiten.Image) {
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
func (overlay *DebugOverlay) drawAIVisionRanges(world w.World, screen *ebiten.Image) {
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
		overlay.drawVisionCircle(screen, float32(screenX), float32(screenY), scaledRadius)
	}))
}

// drawVisionCircle は指定した位置と半径で視界円を描画する
func (overlay *DebugOverlay) drawVisionCircle(screen *ebiten.Image, centerX, centerY, radius float32) {
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
func (overlay *DebugOverlay) drawAIMovementDirections(world w.World, screen *ebiten.Image) {
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
			overlay.drawDirectionArrow(screen, screenX, screenY, velocity.Angle, velocity.Speed, cameraScale)
		}
	}))
}

// drawDirectionArrow は指定した位置に進行方向の矢印を描画する
func (overlay *DebugOverlay) drawDirectionArrow(screen *ebiten.Image, x, y, angle, speed, cameraScale float64) {
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

	vector.StrokeLine(screen, float32(endX), float32(endY), float32(leftX), float32(leftY), strokeWidth, color.RGBA{0, 255, 0, 255}, false)
	vector.StrokeLine(screen, float32(endX), float32(endY), float32(rightX), float32(rightY), strokeWidth, color.RGBA{0, 255, 0, 255}, false)
}
