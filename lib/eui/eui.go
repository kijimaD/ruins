package eui

import (
	"image/color"
	"math"

	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
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
