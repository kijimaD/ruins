package styled

import (
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/ruins/lib/colors"
	w "github.com/kijimaD/ruins/lib/world"
)

// LoadButtonImage はボタンイメージを読み込む
// TODO: いい感じにしたい
func LoadButtonImage() *widget.ButtonImage {
	idle := e_image.NewNineSliceColor(colors.ButtonIdleColor)
	hover := e_image.NewNineSliceColor(colors.ButtonHoverColor)
	pressed := e_image.NewNineSliceColor(colors.ButtonPressedColor)
	pressedHover := e_image.NewNineSliceColor(colors.ButtonPressedColor)
	disabled := e_image.NewNineSliceColor(colors.ButtonPressedColor)

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
	face := (*world.Resources.Faces)["kappa"]

	return &face
}
