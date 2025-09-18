package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	w "github.com/kijimaD/ruins/lib/world"
)

// GameInfo はHUDの基本ゲーム情報エリア
type GameInfo struct {
	enabled bool
}

// NewGameInfo は新しいHUDGameInfoを作成する
func NewGameInfo() *GameInfo {
	return &GameInfo{
		enabled: true,
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

	// HP情報を左上に描画
	info.drawHealthBar(screen, data.PlayerHP, data.PlayerMaxHP)

	// SP情報をHPの下に描画
	info.drawStaminaBar(screen, data.PlayerSP, data.PlayerMaxSP)

	// 空腹度情報をSPの下に描画
	info.drawHungerBar(screen, data.HungerLevel)

	// フロア情報を描画
	info.drawWhiteText(screen, fmt.Sprintf("floor: B%d", data.FloorNumber), 0, 200)

	// ターン情報を描画（フロア表示の下）
	info.drawWhiteText(screen, fmt.Sprintf("turn: %d", data.TurnNumber), 0, 220)

	// 残りアクションポイント（移動ポイント）を描画
	info.drawWhiteText(screen, fmt.Sprintf("AP: %d", data.PlayerMoves), 0, 240)
}

// drawHealthBar はプレイヤーの体力ゲージを描画する
func (info *GameInfo) drawHealthBar(screen *ebiten.Image, currentHP, maxHP int) {
	// HPゲージの設定
	const (
		baseX    = 10.0  // 左マージン
		y        = 10.0  // 上マージン
		width    = 120.0 // ゲージの幅（短縮）
		height   = 12.0  // ゲージの高さ（細く）
		labelGap = 4.0   // ラベルとゲージの間隔
	)

	// 「HP」ラベルを左に描画
	info.drawWhiteText(screen, "HP", int(baseX), int(y-2))

	// ゲージの開始位置（「HP」ラベルの後）
	gageX := float32(baseX + 20.0) // 「HP」の文字幅分オフセット

	// 背景（黒い枠）を描画
	vector.StrokeRect(screen, gageX-1, float32(y-1), float32(width+2), float32(height+2), 1.0, color.RGBA{0, 0, 0, 255}, false)

	// 背景（暗い赤い領域）を描画
	vector.DrawFilledRect(screen, gageX, float32(y), float32(width), float32(height), color.RGBA{100, 0, 0, 255}, false)

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
		vector.DrawFilledRect(screen, gageX, float32(y), currentWidth, float32(height), barColor, false)
	}

	// 数値をゲージの右に描画
	hpText := fmt.Sprintf("%d/%d", currentHP, maxHP)
	info.drawWhiteText(screen, hpText, int(float32(gageX)+float32(width)+float32(labelGap)), int(y-2))
}

// drawStaminaBar はプレイヤーのスタミナポイントゲージを描画する
func (info *GameInfo) drawStaminaBar(screen *ebiten.Image, currentSP, maxSP int) {
	// SPゲージの設定
	const (
		baseX    = 10.0  // 左マージン
		y        = 28.0  // 上マージン（HPバーの下）
		width    = 120.0 // ゲージの幅
		height   = 12.0  // ゲージの高さ
		labelGap = 4.0   // ラベルとゲージの間隔
	)

	// 「SP」ラベルを左に描画
	info.drawWhiteText(screen, "SP", int(baseX), int(y-2))

	// ゲージの開始位置（「SP」ラベルの後）
	gageX := float32(baseX + 20.0) // 「SP」の文字幅分オフセット

	// 背景（黒い枠）を描画
	vector.StrokeRect(screen, gageX-1, float32(y-1), float32(width+2), float32(height+2), 1.0, color.RGBA{0, 0, 0, 255}, false)

	// 背景（暗い青い領域）を描画
	vector.DrawFilledRect(screen, gageX, float32(y), float32(width), float32(height), color.RGBA{0, 0, 100, 255}, false)

	// SP比率を計算
	if maxSP > 0 {
		spRatio := float32(currentSP) / float32(maxSP)
		if spRatio > 1.0 {
			spRatio = 1.0
		}
		if spRatio < 0.0 {
			spRatio = 0.0
		}

		// 現在のSP（青系のグラデーション）
		var barColor color.RGBA
		if spRatio > 0.5 {
			// 水色から青へ（SP 50%以上）
			intensity := uint8((1.0 - spRatio) * 2.0 * 128)
			barColor = color.RGBA{intensity, 128 + intensity, 255, 255}
		} else {
			// 青から暗い青へ（SP 50%以下）
			intensity := uint8(spRatio * 2.0 * 128)
			barColor = color.RGBA{0, intensity, 128 + intensity, 255}
		}

		// 現在のSPバーを描画
		currentWidth := float32(width) * spRatio
		vector.DrawFilledRect(screen, gageX, float32(y), currentWidth, float32(height), barColor, false)
	}

	// 数値をゲージの右に描画
	spText := fmt.Sprintf("%d/%d", currentSP, maxSP)
	info.drawWhiteText(screen, spText, int(float32(gageX)+float32(width)+float32(labelGap)), int(y-2))
}

// drawHungerBar はプレイヤーの空腹度を描画する
func (info *GameInfo) drawHungerBar(screen *ebiten.Image, hungerLevel string) {
	// 空腹度表示の設定
	const (
		baseX = 10.0 // 左マージン
		y     = 46.0 // 上マージン（SPバーの下）
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

	info.drawWhiteText(screen, hungerText, int(baseX), int(y-2))
}

// drawWhiteText は通常の文字でテキストを描画するヘルパー関数
func (info *GameInfo) drawWhiteText(screen *ebiten.Image, text string, x, y int) {
	// 通常の文字を描画
	ebitenutil.DebugPrintAt(screen, text, x, y)
}
