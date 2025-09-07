package gamelog

import "image/color"

// ゲームで使用する色定義
var (
	// 基本色
	ColorWhite   = color.RGBA{255, 255, 255, 255}
	ColorBlack   = color.RGBA{0, 0, 0, 255}
	ColorRed     = color.RGBA{255, 0, 0, 255}
	ColorGreen   = color.RGBA{0, 255, 0, 255}
	ColorBlue    = color.RGBA{0, 0, 255, 255}
	ColorYellow  = color.RGBA{255, 255, 0, 255}
	ColorCyan    = color.RGBA{0, 255, 255, 255}
	ColorMagenta = color.RGBA{255, 0, 255, 255}

	// よく使われる色のバリエーション
	ColorLightGray = color.RGBA{192, 192, 192, 255}
	ColorDarkGray  = color.RGBA{128, 128, 128, 255}
	ColorOrange    = color.RGBA{255, 165, 0, 255}
	ColorPurple    = color.RGBA{128, 0, 128, 255}
)

// NamedColor はRGBから色を作成する関数
func NamedColor(r, g, b uint8) color.RGBA {
	return color.RGBA{r, g, b, 255}
}
