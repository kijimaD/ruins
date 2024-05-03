package eui

import (
	"image/color"
	"math"

	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/styles"
)

func NewEmptyContainer() *widget.Container {
	return widget.NewContainer()
}

// 汎用的なrowコンテナ
func NewRowContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(8),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   4,
				Right:  4,
			}),
		)))
}

// 中身が縦並びのコンテナ
func NewVerticalContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Bottom: 10,
					Left:   10,
					Right:  10,
				}),
			)),
	)
}

// ポーズ用の半透明なコンテナ。コピペがひどいのでどうにかする
func NewVerticalTransContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Bottom: 10,
					Left:   10,
					Right:  10,
				}),
			)),
	)
}

// アイテム系メニューのRootとなる3x3のグリッドコンテナ
func NewItemGridContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				// アイテム, スクロール, アイテム性能で3列になっている
				widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Spacing(2, 0),
				widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    20,
					Bottom: 20,
					Left:   20,
					Right:  20,
				}),
			)),
	)
}

// 縦分割コンテナ
func NewVSplitContainer(top *widget.Container, bottom *widget.Container) *widget.Container {
	split := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(2, 0),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    2,
					Bottom: 2,
					Left:   2,
					Right:  2,
				}),
			)),
	)
	split.AddChild(top)
	split.AddChild(bottom)

	return split
}

// 横分割コンテナ
func NewWSplitContainer(right *widget.Container, left *widget.Container) *widget.Container {
	split := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(2, 0),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    2,
					Bottom: 2,
					Left:   2,
					Right:  2,
				}),
			)),
	)
	split.AddChild(right)
	split.AddChild(left)

	return split
}

// スクロールコンテナとスクロールバー
func NewScrollContainer(content widget.HasWidget) (*widget.ScrollContainer, *widget.Slider) {
	scrollContainer := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.StretchContentWidth(),
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: e_image.NewNineSliceColor(styles.ForegroundColor),
			Mask: e_image.NewNineSliceColor(styles.ForegroundColor),
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

// スクロールコンテナの中身になるコンテナ
func NewScrollContentContainer() *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(2),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    4,
					Bottom: 4,
					Left:   4,
					Right:  4,
				}),
			)))
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

// ウィンドウのヘッダー
func NewWindowHeaderContainer(title string, world w.World) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.WindowHeaderColor)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	container.AddChild(widget.NewText(
		widget.TextOpts.Text(title, *LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	return container
}

// text ================

func NewMenuText(title string, world w.World) *widget.Text {
	text := widget.NewText(
		widget.TextOpts.Text(title, *LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)

	return text
}

func NewBodyText(title string, color color.RGBA, world w.World) *widget.Text {
	text := widget.NewText(
		widget.TextOpts.Text(title, *LoadFont(world), color),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)

	return text
}

// window ================

// ウィンドウ
func NewSmallWindow(title *widget.Container, content *widget.Container) *widget.Window {
	return widget.NewWindow(
		widget.WindowOpts.Contents(content),
		widget.WindowOpts.TitleBar(title, 25),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(200, 200),
		widget.WindowOpts.MaxSize(300, 400),
	)
}

// button ================

func NewItemButton(text string, f func(args *widget.ButtonClickedEventArgs), world w.World) *widget.Button {
	return widget.NewButton(
		widget.ButtonOpts.Image(LoadButtonImage()),
		widget.ButtonOpts.Text(text,
			*LoadFont(world),
			&widget.ButtonTextColor{
				Idle: styles.TextColor,
			},
		),
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.ClickedHandler(f),
	)
}
