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

	return container
}

// createMenuButton はメニューボタンを作成する
func (b *UIBuilder) createMenuButton(menu *Menu, index int, item Item) *widget.Button {
	// ボタンの初期フォーカス状態を設定
	isFocused := index == menu.GetFocusedIndex()

	// 半透明のボタン画像を作成
	buttonImage := b.createTransparentButtonImage()

	btn := eui.NewButton(
		item.Label,
		b.world,
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			menu.SetFocusedIndex(index)
			menu.selectCurrent()
		}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 32), // ボタンの最小サイズを設定（横幅200、縦幅30）
		),
		widget.ButtonOpts.Image(buttonImage), // 半透明背景を設定
	)

	// 無効化されたアイテムの処理
	if item.Disabled {
		btn.GetWidget().Disabled = true
	}

	// 初期フォーカス設定（無効化されていない場合のみ）
	if isFocused && !item.Disabled {
		btn.Focus(true)
	}

	return btn
}

// UpdateFocus はメニューのフォーカス表示を更新する
func (b *UIBuilder) UpdateFocus(menu *Menu) {
	if len(menu.itemWidgets) == 0 {
		return
	}

	// 全てのボタンのフォーカスを解除
	for i, w := range menu.itemWidgets {
		// widget.ButtonはinterFaceでありポインタ型ではないため、型アサーションを修正
		if btn, ok := w.(interface{ Focus(bool) }); ok {
			btn.Focus(i == menu.GetFocusedIndex())
		}
	}
}

// createTransparentButtonImage は半透明のボタン画像を作成する
func (b *UIBuilder) createTransparentButtonImage() *widget.ButtonImage {
	// アイドル状態: 半透明の黒（アルファ値128）
	idle := image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 0, A: 128})

	// ホバー状態: 少し明るい半透明の灰色
	hover := image.NewNineSliceColor(color.NRGBA{R: 60, G: 60, B: 60, A: 160})

	// プレス状態: さらに明るい半透明の灰色
	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 180})

	// 無効状態: 暗い半透明
	disabled := image.NewNineSliceColor(color.NRGBA{R: 30, G: 30, B: 30, A: 80})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
}
