package eui

import (
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/styles"
	"golang.org/x/image/font"
)

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

// TODO: いい感じにしたい
func LoadFont(world w.World) font.Face {
	opts := truetype.Options{
		Size: 24,
		DPI:  72,
	}

	return truetype.NewFace((*world.Resources.Fonts)["kappa"].Font, &opts)
}
