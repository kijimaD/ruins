package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	w "github.com/kijimaD/ruins/lib/world"
)

// GameInfo はHUDの基本ゲーム情報エリア
type GameInfo struct {
	bodyFace    text.Face
	headingFace text.Face // 階層表示用の大きなフォント
	enabled     bool
}

// NewGameInfo は新しいHUDGameInfoを作成する
func NewGameInfo(bodyFace text.Face, headingFace text.Face) *GameInfo {
	return &GameInfo{
		bodyFace:    bodyFace,
		headingFace: headingFace,
		enabled:     true,
	}
}

// Update はゲーム情報エリアを更新する
func (info *GameInfo) Update(_ w.World) {
	// 現在は更新処理なし
}

// Draw はゲーム情報エリアを描画する
func (info *GameInfo) Draw(screen *ebiten.Image, data GameInfoData) {
	if !info.enabled {
		return
	}

	// HP情報
	info.drawHealthBar(screen, data.PlayerHP, data.PlayerMaxHP)

	// SP情報
	info.drawStaminaBar(screen, data.PlayerSP, data.PlayerMaxSP)

	// EP情報
	info.drawElectricityBar(screen, data.PlayerEP, data.PlayerMaxEP)

	// 空腹度
	info.drawHungerBar(screen, data.HungerLevel)

	// ターン情報
	drawOutlinedText(screen, fmt.Sprintf("turn: %d", data.TurnNumber), info.bodyFace, 0, 150, color.White)

	// 残りアクションポイント
	drawOutlinedText(screen, fmt.Sprintf("AP: %d", data.PlayerMoves), info.bodyFace, 0, 170, color.White)

	// フロア情報（最後に描画して最前面に表示）
	info.drawFloorNumber(screen, data)
}

// drawFloorNumber は階層番号を大きく描画する
// 右上（ミニマップの上）に配置
func (info *GameInfo) drawFloorNumber(screen *ebiten.Image, data GameInfoData) {
	const (
		marginRight = 10.0 // 右マージン
		marginTop   = 5.0  // 上マージン
	)

	floorText := fmt.Sprintf("%3dF", data.FloorNumber)

	// テキストの幅を測定
	textWidth, _ := text.Measure(floorText, info.headingFace, 0)

	// 右上に配置
	x := float64(data.ScreenDimensions.Width) - textWidth - marginRight
	y := marginTop

	drawOutlinedText(screen, floorText, info.headingFace, x, y, color.White)
}

// drawHealthBar はプレイヤーの体力ゲージを描画する
func (info *GameInfo) drawHealthBar(screen *ebiten.Image, currentHP, maxHP int) {
	// HPゲージの設定
	const (
		baseX    = 10.0  // 左マージン
		y        = 50.0  // 上マージン
		width    = 120.0 // ゲージの幅
		height   = 20.0  // ゲージの高さ
		labelGap = 4.0   // ラベルとゲージの間隔
	)

	// ゲージの開始位置
	gageX := float32(baseX)

	// 背景（黒い枠）を描画
	vector.StrokeRect(screen, gageX-1, float32(y-1), float32(width+2), float32(height+2), 1.0, color.RGBA{0, 0, 0, 255}, false)

	// 背景（暗い赤い領域）を描画
	vector.FillRect(screen, gageX, float32(y), float32(width), float32(height), color.RGBA{100, 0, 0, 255}, false)

	// HP比率を計算
	if maxHP > 0 {
		hpRatio := float32(currentHP) / float32(maxHP)
		if hpRatio > 1.0 {
			hpRatio = 1.0
		}
		if hpRatio < 0.0 {
			hpRatio = 0.0
		}

		// 現在のHP（緑〜赤のグラデーション）
		var barColor color.RGBA
		if hpRatio > 0.5 {
			// 緑から黄色へ（HP 50%以上）
			intensity := uint8((1.0 - hpRatio) * 2.0 * 255)
			barColor = color.RGBA{intensity, 255, 0, 255}
		} else {
			// 黄色から赤へ（HP 50%以下）
			intensity := uint8(hpRatio * 2.0 * 255)
			barColor = color.RGBA{255, intensity, 0, 255}
		}

		// 現在のHPバーを描画
		currentWidth := float32(width) * hpRatio
		vector.FillRect(screen, gageX, float32(y), currentWidth, float32(height), barColor, false)
	}

	// 数値をゲージの中央に描画
	hpText := fmt.Sprintf("%d/%d", currentHP, maxHP)
	textWidth, _ := text.Measure(hpText, info.bodyFace, 0)
	textX := float64(gageX) + float64(width)/2 - textWidth/2
	textY := y + float64(height)/2 - 6.0 // フォントサイズ16の場合の調整値
	drawOutlinedText(screen, hpText, info.bodyFace, textX, textY, color.White)
}

// drawStaminaBar はプレイヤーのスタミナポイントゲージを描画する
func (info *GameInfo) drawStaminaBar(screen *ebiten.Image, currentSP, maxSP int) {
	// SPゲージの設定
	const (
		baseX    = 10.0  // 左マージン
		y        = 76.0  // 上マージン
		width    = 120.0 // ゲージの幅
		height   = 20.0  // ゲージの高さ
		labelGap = 4.0   // ラベルとゲージの間隔
	)

	// ゲージの開始位置
	gageX := float32(baseX)

	// 背景（黒い枠）を描画
	vector.StrokeRect(screen, gageX-1, float32(y-1), float32(width+2), float32(height+2), 1.0, color.RGBA{0, 0, 0, 255}, false)

	// 背景（暗いグレー領域）を描画
	vector.FillRect(screen, gageX, float32(y), float32(width), float32(height), color.RGBA{100, 100, 100, 255}, false)

	// SP比率を計算
	if maxSP > 0 {
		spRatio := float32(currentSP) / float32(maxSP)
		if spRatio > 1.0 {
			spRatio = 1.0
		}
		if spRatio < 0.0 {
			spRatio = 0.0
		}

		var barColor color.RGBA
		if spRatio > 0.5 {
			// 明るい黄色・オレンジ（SP 50%以上）
			barColor = color.RGBA{255, 200, 0, 255}
		} else {
			// やや暗い黄色・オレンジ（SP 50%以下）
			intensity := uint8(spRatio * 2.0 * 200)
			barColor = color.RGBA{255, intensity, 0, 255}
		}

		// 現在のSPバーを描画
		currentWidth := float32(width) * spRatio
		vector.FillRect(screen, gageX, float32(y), currentWidth, float32(height), barColor, false)
	}

	// 数値をゲージの中央に描画（垂直方向にも中央配置）
	spText := fmt.Sprintf("%d/%d", currentSP, maxSP)
	textWidth, _ := text.Measure(spText, info.bodyFace, 0)
	textX := float64(gageX) + float64(width)/2 - textWidth/2
	textY := y + float64(height)/2 - 6.0 // フォントサイズ16の場合の調整値
	drawOutlinedText(screen, spText, info.bodyFace, textX, textY, color.White)
}

// drawElectricityBar はプレイヤーの電力ポイントゲージを描画する
func (info *GameInfo) drawElectricityBar(screen *ebiten.Image, currentEP, maxEP int) {
	// EPゲージの設定
	const (
		baseX    = 10.0  // 左マージン
		y        = 102.0 // 上マージン
		width    = 120.0 // ゲージの幅
		height   = 20.0  // ゲージの高さ
		labelGap = 4.0   // ラベルとゲージの間隔
	)

	// ゲージの開始位置
	gageX := float32(baseX)

	// 背景（黒い枠）を描画
	vector.StrokeRect(screen, gageX-1, float32(y-1), float32(width+2), float32(height+2), 1.0, color.RGBA{0, 0, 0, 255}, false)

	// 背景（暗い青い領域）を描画
	vector.FillRect(screen, gageX, float32(y), float32(width), float32(height), color.RGBA{0, 0, 80, 255}, false)

	// EP比率を計算
	if maxEP > 0 {
		epRatio := float32(currentEP) / float32(maxEP)
		if epRatio > 1.0 {
			epRatio = 1.0
		}
		if epRatio < 0.0 {
			epRatio = 0.0
		}

		// 現在のEP（青系のグラデーション、電力らしい色）
		var barColor color.RGBA
		if epRatio > 0.5 {
			// シアンから青へ（EP 50%以上）
			intensity := uint8((1.0 - epRatio) * 2.0 * 100)
			barColor = color.RGBA{intensity, 200, 255, 255}
		} else {
			// 青から暗い青へ（EP 50%以下）
			intensity := uint8(epRatio * 2.0 * 200)
			barColor = color.RGBA{0, intensity, 100 + uint8(epRatio*155), 255}
		}

		// 現在のEPバーを描画
		currentWidth := float32(width) * epRatio
		vector.FillRect(screen, gageX, float32(y), currentWidth, float32(height), barColor, false)
	}

	// 数値をゲージの中央に描画（垂直方向にも中央配置）
	epText := fmt.Sprintf("%d/%d", currentEP, maxEP)
	textWidth, _ := text.Measure(epText, info.bodyFace, 0)
	textX := float64(gageX) + float64(width)/2 - textWidth/2
	textY := y + float64(height)/2 - 6.0 // フォントサイズ16の場合の調整値
	drawOutlinedText(screen, epText, info.bodyFace, textX, textY, color.White)
}

// drawHungerBar はプレイヤーの空腹度を描画する
func (info *GameInfo) drawHungerBar(screen *ebiten.Image, hungerLevel string) {
	// 空腹度表示の設定
	const (
		baseX = 10.0  // 左マージン
		y     = 128.0 // 上マージン（EPバーの下、高さ変更分を調整）
	)

	// 空腹度レベルのテキストを描画
	hungerText := fmt.Sprintf("Hunger %s", hungerLevel)

	// TODO: 空腹度レベルに応じて色を変える
	switch hungerLevel {
	case "Full":
		// 通常の白色
	case "Normal":
		// 通常の白色
	case "Hungry":
		// やや警告の色
	case "Starving":
		// 危険な色
	}

	drawOutlinedText(screen, hungerText, info.bodyFace, float64(baseX), y-2, color.White)
}
