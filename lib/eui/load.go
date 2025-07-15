package eui

import (
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/styles"
)

// LoadButtonImage はボタンイメージを読み込む
// TODO: いい感じにしたい
func LoadButtonImage() *widget.ButtonImage {
	idle := e_image.NewNineSliceColor(styles.ButtonIdleColor)
	hover := e_image.NewNineSliceColor(styles.ButtonHoverColor)
	pressed := e_image.NewNineSliceColor(styles.ButtonPressedColor)
	pressedHover := e_image.NewNineSliceColor(styles.ButtonPressedColor)
	disabled := e_image.NewNineSliceColor(styles.ButtonPressedColor)

	return &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressedHover,
		Disabled:     disabled,
	}
}

// LoadFont はフォントを読み込む
// TODO: いい感じにしたい
func LoadFont(world w.World) *text.Face {
	face := (*world.Resources.DefaultFaces)["kappa"]

	return &face
}
