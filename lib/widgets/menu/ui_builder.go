package menu

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/eui"
	w "github.com/kijimaD/ruins/lib/world"
)

// UIBuilder はメニューのUI要素を構築する
type UIBuilder struct {
	world w.World
}

// NewUIBuilder はUIビルダーを作成する
func NewUIBuilder(world w.World) *UIBuilder {
	return &UIBuilder{
		world: world,
	}
}

// BuildUI はメニューのUI要素を構築する
func (b *UIBuilder) BuildUI(menu *Menu) *widget.Container {
	var container *widget.Container

	if menu.config.Columns > 0 {
		// グリッド表示
		container = b.buildGridUI(menu)
	} else if menu.config.Orientation == Horizontal {
		// 水平リスト表示
		container = b.buildHorizontalUI(menu)
	} else {
		// 垂直リスト表示
		container = b.buildVerticalUI(menu)
	}

	menu.SetContainer(container)
	return container
}

// buildVerticalUI は垂直リスト表示のUIを構築する
func (b *UIBuilder) buildVerticalUI(menu *Menu) *widget.Container {
	container := eui.NewVerticalContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	for i, item := range menu.config.Items {
		btn := b.createMenuButton(menu, i, item)
		container.AddChild(btn)
		menu.itemWidgets = append(menu.itemWidgets, btn)
	}

	b.UpdateFocus(menu)

	return container
}

// buildHorizontalUI は水平リスト表示のUIを構築する
func (b *UIBuilder) buildHorizontalUI(menu *Menu) *widget.Container {
	container := eui.NewRowContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	for i, item := range menu.config.Items {
		btn := b.createMenuButton(menu, i, item)
		container.AddChild(btn)
		menu.itemWidgets = append(menu.itemWidgets, btn)
	}

	b.UpdateFocus(menu)

	return container
}

// buildGridUI はグリッド表示のUIを構築する
func (b *UIBuilder) buildGridUI(menu *Menu) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(menu.config.Columns),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
				widget.GridLayoutOpts.Spacing(2, 2),
			),
		),
	)

	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	for i, item := range menu.config.Items {
		btn := b.createMenuButton(menu, i, item)
		container.AddChild(btn)
		menu.itemWidgets = append(menu.itemWidgets, btn)
	}

	b.UpdateFocus(menu)

	return container
}

// createMenuButton はメニューボタンを作成する
func (b *UIBuilder) createMenuButton(menu *Menu, index int, item Item) *widget.Button {
	res := b.world.Resources.UIResources

	btn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(100, 28),
		),
		widget.ButtonOpts.Image(res.Button.Image),
		widget.ButtonOpts.Text(
			item.Label,
			res.Button.Face,
			res.Button.Text,
		),
		widget.ButtonOpts.TextPadding(res.Button.Padding),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter), // 左寄せ
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			menu.SetFocusedIndex(index)
			menu.selectCurrent()
		}),
	)

	// 無効化されたアイテムの処理
	if item.Disabled {
		btn.GetWidget().Disabled = true
	}

	return btn
}

// UpdateFocus はメニューのフォーカス表示を更新する
// カーソルで選択中の要素だけボタンを変える。マウスのhoverでは色が変わらないようにしている
// カーソル移動は独自実装なので、UIを対応させるために必要
func (b *UIBuilder) UpdateFocus(menu *Menu) {
	if len(menu.itemWidgets) == 0 {
		return
	}

	// 全てのボタンのフォーカスを更新
	for i, w := range menu.itemWidgets {
		if btn, ok := w.(*widget.Button); ok {
			isFocused := i == menu.GetFocusedIndex()

			// フォーカス状態に応じてボタンの画像を更新
			if isFocused {
				// フォーカス時: より明るい背景色
				focusedImage := b.createFocusedButtonImage()
				btn.Image = focusedImage
			} else {
				// 非フォーカス時: 通常の半透明背景
				normalImage := b.createTransparentButtonImage()
				btn.Image = normalImage
			}
		}
	}
}

// createTransparentButtonImage は半透明のボタン画像を作成する
// マウスホバーでは色が変わらず、キーボード操作（フォーカス）でのみ反応する
func (b *UIBuilder) createTransparentButtonImage() *widget.ButtonImage {
	// アイドル状態: 透明
	idle := image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 0, A: 0})

	// プレス状態: さらに明るい半透明の灰色
	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 180})

	// 無効状態: 暗い半透明
	disabled := image.NewNineSliceColor(color.NRGBA{R: 30, G: 30, B: 30, A: 16})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    idle, // フォーカス時と同じ色でマウスホバー効果を無効化
		Pressed:  pressed,
		Disabled: disabled,
	}
}

// createFocusedButtonImage はフォーカス時の明るいボタン画像を作成する
func (b *UIBuilder) createFocusedButtonImage() *widget.ButtonImage {
	// フォーカス時: より明るい半透明の灰色
	focused := image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 0, A: 120})

	// プレス状態: さらに明るい色
	pressed := image.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 120, A: 200})

	// 無効状態: 暗い半透明
	disabled := image.NewNineSliceColor(color.NRGBA{R: 30, G: 30, B: 30, A: 80})

	return &widget.ButtonImage{
		Idle:     focused,
		Hover:    focused, // フォーカス時と同じ色でマウスホバー効果を無効化
		Pressed:  pressed,
		Disabled: disabled,
	}
}
