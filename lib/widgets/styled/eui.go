package styled

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/colors"
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
		widget.TextOpts.Text(title, res.Text.TitleFace, colors.TextColor),
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
		widget.TextOpts.Text(title, res.Text.Face, colors.TextColor),
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
		widget.TextOpts.Text(text, res.Text.TitleFace, colors.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewSubtitleText はサブタイトル用テキストを作成する（中サイズ）
func NewSubtitleText(text string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.Face, colors.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewDescriptionText は説明文用テキストを作成する（小さめ、補助的）
func NewDescriptionText(text string, world w.World) *widget.Text {
	res := world.Resources.UIResources
	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.SmallFace, colors.ForegroundColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// NewPageIndicator は右寄せのページインジケーターを作成する
func NewPageIndicator(text string, world w.World) *widget.Container {
	res := world.Resources.UIResources

	// 透明な背景のコンテナを作成（NewListItemTextと同じパターン）
	backgroundColor := image.NewNineSliceColor(colors.TransparentColor)

	container := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(backgroundColor),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true, // 横幅を親コンテナに合わせる
			}),
			widget.WidgetOpts.MinSize(120, 0), // NewListItemTextと同じ最小横幅
		),
	)

	// 右寄せのテキスト
	textWidget := widget.NewText(
		widget.TextOpts.Text(text, res.Text.SmallFace, colors.ForegroundColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd, // 右寄せ
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				Padding: widget.Insets{
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

// NewBodyText は本文用テキストを作成する
func NewBodyText(title string, _ color.RGBA, world w.World) *widget.Text {
	res := world.Resources.UIResources
	text := widget.NewText(
		widget.TextOpts.Text(title, res.Text.Face, colors.TextColor),
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
		backgroundColor = image.NewNineSliceColor(colors.ButtonHoverColor)
	} else {
		// 非選択は背景なし（透明）
		backgroundColor = image.NewNineSliceColor(colors.TransparentColor)
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

// NewFragmentText は色付きログフラグメント専用のテキストを作成する（文字数分だけの幅）
func NewFragmentText(text string, textColor color.RGBA, world w.World) *widget.Text {
	res := world.Resources.UIResources

	return widget.NewText(
		widget.TextOpts.Text(text, res.Text.Face, textColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: false, // 横幅を伸ばさない
			}),
			// MinSizeは指定しない（テキストの自然な幅を使用）
		),
	)
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

// NewMessageList はメッセージ表示用のリストウィジェットを作成する（戦闘ログなど用）
func NewMessageList(entries []any, world w.World, opts ...widget.ListOpt) *widget.List {
	res := world.Resources.UIResources

	// メッセージ表示用のデフォルトオプション
	defaultOpts := []widget.ListOpt{
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(world.Resources.ScreenDimensions.Width-100, 280),
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				HorizontalPosition: widget.GridLayoutPositionCenter,
				VerticalPosition:   widget.GridLayoutPositionEnd,
				MaxHeight:          280,
			}),
		)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.Images(res.List.Track, res.List.Handle),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(4)),
		),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			if str, ok := e.(string); ok {
				return str
			}
			return ""
		}),
		widget.ListOpts.EntrySelectedHandler(func(_ *widget.ListEntrySelectedEventArgs) {}),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                  colors.TextColor,
			Unselected:                colors.TextColor,
			SelectedBackground:        colors.ButtonHoverColor,
			SelectedFocusedBackground: colors.ButtonHoverColor,
		}),
		widget.ListOpts.EntryFontFace(res.Text.Face),
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ListOpts.EntryTextPadding(widget.Insets{
			Top:    4,
			Bottom: 4,
			Left:   16,
			Right:  16,
		}),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(res.List.ImageTrans),
		),
		widget.ListOpts.HideHorizontalSlider(),
	}

	// カスタムオプションを追加
	allOpts := append(defaultOpts, opts...)

	return widget.NewList(allOpts...)
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
