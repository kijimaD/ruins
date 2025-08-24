package menu

import (
	"fmt"
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

	if menu.config.Orientation == Horizontal {
		// 水平リスト表示
		container = b.buildHorizontalUI(menu)
	} else {
		// 垂直リスト表示
		container = b.buildVerticalUI(menu)
	}

	menu.SetContainer(container)
	menu.SetUIBuilder(b) // UIビルダーを設定
	return container
}

// buildVerticalUI は垂直リスト表示のUIを構築する
func (b *UIBuilder) buildVerticalUI(menu *Menu) *widget.Container {
	mainContainer := eui.NewVerticalContainer()

	// ページインジケーターを追加
	if menu.config.ShowPageIndicator && menu.config.ItemsPerPage > 0 && menu.GetTotalPages() > 1 {
		pageIndicator := b.CreatePageIndicator(menu)
		mainContainer.AddChild(pageIndicator)
	}

	// メニューアイテムのコンテナ
	menuContainer := eui.NewVerticalContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	// 表示する項目のみを追加（スクロール対応）
	visibleItems, indices := menu.GetVisibleItemsWithIndices()
	for i, item := range visibleItems {
		originalIndex := indices[i]
		btn := b.CreateMenuButton(menu, originalIndex, item)
		menuContainer.AddChild(btn)
		menu.itemWidgets = append(menu.itemWidgets, btn)
	}

	mainContainer.AddChild(menuContainer)
	b.UpdateFocus(menu)

	return mainContainer
}

// buildHorizontalUI は水平リスト表示のUIを構築する
func (b *UIBuilder) buildHorizontalUI(menu *Menu) *widget.Container {
	mainContainer := eui.NewVerticalContainer()

	// ページインジケーターを追加
	if menu.config.ShowPageIndicator && menu.config.ItemsPerPage > 0 && menu.GetTotalPages() > 1 {
		pageIndicator := b.CreatePageIndicator(menu)
		mainContainer.AddChild(pageIndicator)
	}

	// メニューアイテムのコンテナ
	menuContainer := eui.NewRowContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	// 表示する項目のみを追加（ペジネーション）
	visibleItems, indices := menu.GetVisibleItemsWithIndices()
	for i, item := range visibleItems {
		originalIndex := indices[i]
		btn := b.CreateMenuButton(menu, originalIndex, item)
		menuContainer.AddChild(btn)
		menu.itemWidgets = append(menu.itemWidgets, btn)
	}

	mainContainer.AddChild(menuContainer)
	b.UpdateFocus(menu)

	return mainContainer
}

// CreateMenuButton はメニューボタンを作成する
func (b *UIBuilder) CreateMenuButton(menu *Menu, index int, item Item) *widget.Button {
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

	// 表示中の項目とそのインデックスを取得
	_, indices := menu.GetVisibleItemsWithIndices()

	// 全てのボタンのフォーカスを更新
	for i, w := range menu.itemWidgets {
		if btn, ok := w.(*widget.Button); ok && i < len(indices) {
			originalIndex := indices[i]
			isFocused := originalIndex == menu.GetFocusedIndex()

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

// CreatePageIndicator はページインジケーターを作成する
func (b *UIBuilder) CreatePageIndicator(menu *Menu) *widget.Text {
	res := b.world.Resources.UIResources

	pageText := fmt.Sprintf("%d/%d",
		menu.GetCurrentPage(), menu.GetTotalPages())

	return widget.NewText(
		widget.TextOpts.Text(pageText, res.Text.SmallFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(300, 20),
		),
	)
}
