package styles

import "image/color"

func RGB(rgb uint64) color.RGBA {
	return color.RGBA{
		R: uint8((rgb & (0xFF << (8 * 2))) >> (8 * 2)),
		G: uint8((rgb & (0xFF << (8 * 1))) >> (8 * 1)),
		B: uint8((rgb & (0xFF << (8 * 0))) >> (8 * 0)),
		A: 0xFF,
	}
}

var (
	TransparentColor = color.RGBA{}
	// 主要
	PrimaryColor = RGB(0x9dd793)
	// サブ
	SecondaryColor = RGB(0x9dd793)
	// 地のテキスト
	TextColor = RGB(0xf5f5f5)
	// 前
	ForegroundColor = RGB(0xa9a9a9)
	// 背景
	BackgroundColor = RGB(0x000000)
	// デバッグ
	DebugColor = RGB(0x0000FF)

	// ウィンドウ
	WindowBodyColor   = RGB(0x808080)
	WindowHeaderColor = RGB(0x939393)

	// ボタン
	ButtonIdleColor    = RGB(0xaaaaaa)
	ButtonHoverColor   = RGB(0x828296)
	ButtonPressedColor = RGB(0x646478)

	SuccessColor = RGB(0x198754)
	DangerColor  = RGB(0xdc3545)
)
