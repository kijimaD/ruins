package eui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/styles"
	w "github.com/kijimaD/ruins/lib/world"
)

// NewRowContainer は汎用的なrowコンテナを作成する
func NewRowContainer(opts ...widget.ContainerOpt) *widget.Container {
	return widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(
				widget.NewRowLayout(
					BaseRowLayoutOpts()...,
				),
			),
		}, opts...)...,
	)
}

// NewVerticalContainer は中身が縦並びのコンテナを作成する
func NewVerticalContainer(opts ...widget.ContainerOpt) *widget.Container {
	return widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(
				widget.NewRowLayout(
					append([]widget.RowLayoutOpt{
						widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					},
						BaseRowLayoutOpts()...,
					)...,
				),
			),
		}, opts...)...,
	)
}

// NewItemGridContainer はアイテム系メニューのRootとなる3x3のグリッドコンテナを作成する
func NewItemGridContainer(opts ...widget.ContainerOpt) *widget.Container {
	return widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(
				widget.NewGridLayout(
					// アイテム, スクロール, アイテム性能で3列になっている
					widget.GridLayoutOpts.Columns(3),
					widget.GridLayoutOpts.Spacing(4, 4),
					widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{false, true, false}),
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    4,
						Bottom: 4,
						Left:   4,
						Right:  4,
					}),
				)),
		}, opts...)...,
	)
}

// NewVSplitContainer は縦分割コンテナを作成する
func NewVSplitContainer(top *widget.Container, bottom *widget.Container, opts ...widget.ContainerOpt) *widget.Container {
	split := widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(
				widget.NewGridLayout(
					widget.GridLayoutOpts.Columns(1),
					widget.GridLayoutOpts.Spacing(4, 4),
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    4,
						Bottom: 4,
						Left:   4,
						Right:  4,
					}),
				)),
		}, opts...)...,
	)
	split.AddChild(top)
	split.AddChild(bottom)

	return split
}

// NewWSplitContainer は横分割コンテナを作成する
func NewWSplitContainer(right *widget.Container, left *widget.Container, opts ...widget.ContainerOpt) *widget.Container {
	split := widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(
				widget.NewGridLayout(
					widget.GridLayoutOpts.Columns(2),
					widget.GridLayoutOpts.Spacing(4, 4),
					widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true}),
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    4,
						Bottom: 4,
						Left:   4,
						Right:  4,
					}),
				)),
		}, opts...)...,
	)
	split.AddChild(right)
	split.AddChild(left)

	return split
}

// NewWindowContainer はウィンドウの本体を作成する
func NewWindowContainer(world w.World) *widget.Container {
	res := world.Resources.UIResources

	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
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
				MaxHeight: 500,
			}),
		),
	)
}

// NewWindowHeaderContainer はウィンドウのヘッダーを作成する
func NewWindowHeaderContainer(title string, world w.World) *widget.Container {
	res := world.Resources.UIResources
	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.TitleBar),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	container.AddChild(widget.NewText(
		widget.TextOpts.Text(title, res.Text.TitleFace, styles.TextColor),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
	))

	return container
}

// text ================

// NewMenuText は汎用メニューテキストを作成する（既存との互換性のため維持）
func NewMenuText(title string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	text := widget.NewText(
		widget.TextOpts.Text(title, res.Text.Face, styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)

	return text
}

// NewTitleText はタイトル用テキストを作成する（大きめ、目立つ）
func NewTitleText(text string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.TitleFace, styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewSubtitleText はサブタイトル用テキストを作成する（中サイズ）
func NewSubtitleText(text string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.Face, styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewDescriptionText は説明文用テキストを作成する（小さめ、補助的）
func NewDescriptionText(text string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.SmallFace, styles.ForegroundColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewBodyText は本文用テキストを作成する
func NewBodyText(title string, _ color.RGBA, world w.World) *widget.Text {
	res := world.Resources.UIResources
	text := widget.NewText(
		widget.TextOpts.Text(title, res.Text.Face, styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)

	return text
}

// NewListItemText はリスト項目用テキストを作成する（背景色変更で選択状態を表現）
func NewListItemText(text string, textColor color.RGBA, isSelected bool, world w.World) *widget.Container {
	res := world.Resources.UIResources

	var backgroundColor *image.NineSlice
	if isSelected {
		// 選択中は背景色を付ける
		backgroundColor = image.NewNineSliceColor(styles.ButtonHoverColor)
	} else {
		// 非選択は背景なし（透明）
		backgroundColor = image.NewNineSliceColor(styles.TransparentColor)
	}

	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(backgroundColor),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true, // 横幅を親コンテナに合わせる
			}),
			widget.WidgetOpts.MinSize(120, 0), // 最小横幅を固定
		),
	)

	textWidget := widget.NewText(
		widget.TextOpts.Text(text, res.Text.Face, textColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart, // 左寄せ
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				Padding: widget.Insets{ // 縦パディングを小さく、横パディングは適度に
					Top:    2,
					Bottom: 2,
					Left:   8,
					Right:  8,
				},
			}),
		),
	)

	container.AddChild(textWidget)
	return container
}

// window ================

// NewSmallWindow は小さなウィンドウを作成する
func NewSmallWindow(title *widget.Container, content *widget.Container) *widget.Window {
	return widget.NewWindow(
		widget.WindowOpts.Contents(content),
		widget.WindowOpts.TitleBar(title, 25),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(200, 200),
		widget.WindowOpts.MaxSize(650, 550),
	)
}

// list ================

// NewList はリストウィジェットを作成する
func NewList(entries []any, listOpts []euiext.ListOpt, world w.World) *euiext.List {
	res := world.Resources.UIResources

	return euiext.NewList(
		append([]euiext.ListOpt{
			euiext.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(150, 0),
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionEnd,
					StretchVertical:    true,
					Padding:            widget.NewInsetsSimple(50),
				}),
			)),
			euiext.ListOpts.Entries(entries),
			euiext.ListOpts.ScrollContainerOpts(
				widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
					Idle:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Disabled: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
					Mask:     image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				}),
			),
			euiext.ListOpts.HideHorizontalSlider(),
			euiext.ListOpts.EntryFontFace(*LoadFont(world)),
			euiext.ListOpts.EntryColor(&euiext.ListEntryColor{
				Selected:                   color.NRGBA{R: 0, G: 255, B: 0, A: 255},
				Unselected:                 color.NRGBA{R: 254, G: 255, B: 255, A: 255},
				SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255},
				SelectingBackground:        color.NRGBA{R: 130, G: 130, B: 130, A: 255},
				SelectingFocusedBackground: color.NRGBA{R: 130, G: 140, B: 170, A: 255},
				SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255},
				FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255},
				DisabledUnselected:         color.NRGBA{R: 100, G: 100, B: 100, A: 255},
				DisabledSelected:           color.NRGBA{R: 100, G: 100, B: 100, A: 255},
				DisabledSelectedBackground: color.NRGBA{R: 100, G: 100, B: 100, A: 255},
			}),
			euiext.ListOpts.EntryLabelFunc(func(_ interface{}) string { return "" }),
			euiext.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
			euiext.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
			euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.Image)),
			euiext.ListOpts.SliderOpts(
				widget.SliderOpts.Images(res.List.Track, res.List.Handle),
				widget.SliderOpts.MinHandleSize(res.List.HandleSize),
				widget.SliderOpts.TrackPadding(res.List.TrackPadding),
			),
			euiext.ListOpts.HideHorizontalSlider(),
			euiext.ListOpts.Entries(entries),
			euiext.ListOpts.EntryFontFace(res.List.Face),
			euiext.ListOpts.EntryTextPadding(res.List.EntryPadding),
			euiext.ListOpts.AllowReselect(),
		}, listOpts...)...,
	)
}

// button ================

// NewButton はボタンウィジェットを作成する
func NewButton(text string, world w.World, opts ...widget.ButtonOpt) *widget.Button {
	res := world.Resources.UIResources
	return widget.NewButton(
		append([]widget.ButtonOpt{
			widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			})),
			widget.ButtonOpts.Image(res.Button.Image),
			widget.ButtonOpts.Text(
				text,
				res.Button.Face,
				res.Button.Text,
			),
			widget.ButtonOpts.TextPadding(res.Button.Padding),
		}, opts...)...,
	)
}

// opts ================

// BaseRowLayoutOpts は基本的な行レイアウトオプションを返す
func BaseRowLayoutOpts() []widget.RowLayoutOpt {
	return []widget.RowLayoutOpt{
		widget.RowLayoutOpts.Spacing(4),
		widget.RowLayoutOpts.Padding(widget.Insets{
			Top:    10,
			Bottom: 10,
			Left:   4,
			Right:  4,
		}),
	}
}
