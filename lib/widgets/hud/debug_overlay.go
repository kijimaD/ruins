package hud

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	w "github.com/kijimaD/ruins/lib/world"
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

// Update はデバッグオーバーレイを更新する
func (overlay *DebugOverlay) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はデバッグオーバーレイを描画する
func (overlay *DebugOverlay) Draw(screen *ebiten.Image, data DebugOverlayData) {
	if !overlay.enabled || !data.Enabled {
		return
	}

	// AI状態を描画
	for _, aiState := range data.AIStates {
		textOffsetY := 30.0
		ebitenutil.DebugPrintAt(screen, aiState.StateText, int(aiState.ScreenX)-20, int(aiState.ScreenY-textOffsetY))
	}

	// 視界範囲を描画
	for _, visionRange := range data.VisionRanges {
		overlay.drawVisionCircle(screen, float32(visionRange.ScreenX), float32(visionRange.ScreenY), visionRange.ScaledRadius)
	}
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
