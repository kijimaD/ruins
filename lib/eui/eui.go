package eui

import (
	"image/color"
	"math"

	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/styles"
)

func EmptyContainer() *widget.Container {
	return widget.NewContainer()
}

func NewScrollContainer(content widget.HasWidget) (*widget.ScrollContainer, *widget.Slider) {
	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			Mask: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
		}),
	)
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ContentRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}
	trackPadding := widget.Insets{4, 20, 20, 4}
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			&widget.SliderTrackImage{
				Idle:  e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 0}),
				Hover: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 0}),
			},
			&widget.ButtonImage{
				Idle:    e_image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   e_image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: e_image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
		widget.SliderOpts.TrackPadding(trackPadding),
	)
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	return scrollContainer, vSlider
}

// 前面に開くウィンドウ用のコンテナ。色が違ったりする
func NewWindowContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.WindowBodyColor)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
				Left:   10,
				Right:  10,
			}),
			widget.RowLayoutOpts.Spacing(2),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxHeight: 160,
			}),
		),
	)
}

// 左上のメニュータイトル
// 「インベントリ」や「仲間」や「装備」とか
func NewMenuTitle(title string, world w.World) *widget.Text {
	text := widget.NewText(
		widget.TextOpts.Text(title, LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)

	return text
}
