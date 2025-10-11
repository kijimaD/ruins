package hud

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	w "github.com/kijimaD/ruins/lib/world"
)

// CurrencyDisplay は遺光表示ウィジェット
type CurrencyDisplay struct {
	world   w.World
	enabled bool
}

// NewCurrencyDisplay は新しいCurrencyDisplayを作成する
func NewCurrencyDisplay(world w.World) *CurrencyDisplay {
	return &CurrencyDisplay{
		world:   world,
		enabled: true,
	}
}

// SetEnabled は表示の有効/無効を設定する
func (c *CurrencyDisplay) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// Update は更新処理（現在は何もしない）
func (c *CurrencyDisplay) Update(_ w.World) {
	// 必要に応じて更新処理を追加
}

// Draw は遺光を描画する
func (c *CurrencyDisplay) Draw(screen *ebiten.Image, data CurrencyData) {
	if !c.enabled {
		return
	}

	// 画面サイズを取得
	screenWidth := data.ScreenDimensions.Width
	screenHeight := data.ScreenDimensions.Height

	// 遺光テキスト
	currencyText := fmt.Sprintf("◆ %d", data.Currency)

	// UIリソースからフォントを取得
	face := c.world.Resources.UIResources.Text.Face

	// テキストの幅を計算
	textWidth, _ := text.Measure(currencyText, face, 0)

	// メッセージウィンドウの位置を計算
	fixedHeight := data.Config.LogAreaMargin*2 + data.Config.MaxLogLines*data.Config.LineHeight + data.Config.YPadding*2
	logAreaY := screenHeight - fixedHeight

	// メッセージウィンドウの上端の上に配置（右寄せ）
	currencyX := float64(screenWidth-data.Config.LogAreaMargin) - textWidth
	currencyY := float64(logAreaY - 25) // メッセージウィンドウの上に十分なスペースを取って表示

	// テキストを描画
	op := &text.DrawOptions{}
	op.GeoM.Translate(currencyX, currencyY)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, currencyText, face, op)
}
