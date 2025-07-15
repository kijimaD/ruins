package styles

import "image/color"

// RGB はRGB値からcolor.RGBAを作成する
func RGB(rgb uint64) color.RGBA {
	return color.RGBA{
		R: uint8((rgb & (0xFF << (8 * 2))) >> (8 * 2)),
		G: uint8((rgb & (0xFF << (8 * 1))) >> (8 * 1)),
		B: uint8((rgb & (0xFF << (8 * 0))) >> (8 * 0)),
		A: 0xFF,
	}
}

var (
	// TransparentColor は透明色を表す
	TransparentColor = color.RGBA{}
	// PrimaryColor は主要色を表す
	PrimaryColor = RGB(0x9dd793)
	// SecondaryColor はサブ色を表す
	SecondaryColor = RGB(0x9dd793)
	// TextColor は地のテキスト色を表す
	TextColor = RGB(0xf5f5f5)
	// ForegroundColor は前景色を表す
	ForegroundColor = RGB(0xa9a9a9)
	// BackgroundColor は背景色を表す
	BackgroundColor = RGB(0x000000)
	// DebugColor はデバッグ色を表す
	DebugColor = RGB(0x0000FF)
	// TransBlackColor は透過黒背景色を表す
	TransBlackColor = color.RGBA{0, 0, 0, 140}

	// WindowBodyColor はウィンドウ本体色を表す
	WindowBodyColor   = RGB(0x808080)
	// WindowHeaderColor はウィンドウヘッダー色を表す
	WindowHeaderColor = RGB(0x939393)

	// ButtonIdleColor はボタン通常色を表す
	ButtonIdleColor     = RGB(0xaaaaaa)
	// ButtonHoverColor はボタンホバー色を表す
	ButtonHoverColor    = RGB(0x828296)
	// ButtonPressedColor はボタン押下色を表す
	ButtonPressedColor  = RGB(0x646478)
	// ButtonDisabledColor はボタン無効色を表す
	ButtonDisabledColor = RGB(0x555555)

	// SuccessColor は成功色を表す
	SuccessColor = RGB(0x198754)
	// DangerColor は危険色を表す
	DangerColor  = RGB(0xdc3545)

	// FireColor は炎色（赤）を表す
	FireColor    = RGB(0xc44303)
	// ThunderColor は雷色（暗青）を表す
	ThunderColor = RGB(0x4169e1)
	// ChillColor は冷気色（明青）を表す
	ChillColor   = RGB(0x00ffff)
	// PhotonColor は光子色（黄）を表す
	PhotonColor  = RGB(0xffff00)
)
