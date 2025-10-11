package menu

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/widgets/styled"
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
	mainContainer := styled.NewVerticalContainer()

	// ページインジケーターを追加
	pageText := menu.GetPageIndicatorText()
	if pageText != "" {
		pageIndicator := b.CreatePageIndicator(menu)
		mainContainer.AddChild(pageIndicator)
	}

	// メニューアイテムのコンテナ
	menuContainer := styled.NewVerticalContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	// 表示する項目のみを追加（スクロール対応）
	visibleItems, indices := menu.GetVisibleItems()
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
	mainContainer := styled.NewVerticalContainer()

	// ページインジケーターを追加
	pageText := menu.GetPageIndicatorText()
	if pageText != "" {
		pageIndicator := b.CreatePageIndicator(menu)
		mainContainer.AddChild(pageIndicator)
	}

	// メニューアイテムのコンテナ
	menuContainer := styled.NewRowContainer()
	menu.itemWidgets = make([]widget.PreferredSizeLocateableWidget, 0)

	// 表示する項目のみを追加（ペジネーション）
	visibleItems, indices := menu.GetVisibleItems()
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
// 追加ラベルがある場合は、Button + ラベル群をコンテナでまとめて返す
func (b *UIBuilder) CreateMenuButton(menu *Menu, index int, item Item) widget.PreferredSizeLocateableWidget {
	res := b.world.Resources.UIResources

	// フォーカス状態をチェック
	isFocused := index == menu.GetFocusedIndex()

	// ボタン画像を作成
	var buttonImage *widget.ButtonImage
	if isFocused {
		buttonImage = b.createFocusedButtonImage()
	} else {
		buttonImage = b.createTransparentButtonImage()
	}

	// ボタンを作成
	btn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: len(item.AdditionalLabels) == 0, // 追加ラベルがなければStretch
			}),
			widget.WidgetOpts.MinSize(100, 28),
		),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(item.Label, &res.Button.Face, res.Button.Text),
		widget.ButtonOpts.TextPadding(&res.Button.Padding),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ButtonOpts.ClickedHandler(func(_ *widget.ButtonClickedEventArgs) {
			menu.SetFocusedIndex(index)
			menu.selectCurrent()
		}),
	)

	// 無効化されたアイテムの処理
	if item.Disabled {
		btn.GetWidget().Disabled = true
	}

	// 追加ラベルがない場合はボタンをそのまま返す
	if len(item.AdditionalLabels) == 0 {
		return btn
	}

	// 追加ラベルがある場合は、コンテナでまとめる
	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	// ボタンを追加
	container.AddChild(btn)

	// 右側: 追加ラベル群（右寄せ）
	for _, label := range item.AdditionalLabels {
		additionalText := widget.NewText(
			widget.TextOpts.Text(label, &res.Button.Face, res.Button.Text.Idle),
			widget.TextOpts.Position(widget.TextPositionEnd, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionEnd,
				}),
			),
		)
		container.AddChild(additionalText)
	}

	return container
}

// UpdateFocus はメニューのフォーカス表示を更新する
// カーソルで選択中の要素だけボタンを変える。マウスのhoverでは色が変わらないようにしている
// カーソル移動は独自実装なので、UIを対応させるために必要
func (b *UIBuilder) UpdateFocus(menu *Menu) {
	if len(menu.itemWidgets) == 0 {
		return
	}

	// 表示中の項目とそのインデックスを取得
	_, indices := menu.GetVisibleItems()

	// 全てのボタンのフォーカスを更新
	for i, w := range menu.itemWidgets {
		if i >= len(indices) {
			continue
		}

		originalIndex := indices[i]
		isFocused := originalIndex == menu.GetFocusedIndex()

		// ボタンの場合
		if btn, ok := w.(*widget.Button); ok {
			// フォーカス状態に応じてボタンの画像を更新
			if isFocused {
				// フォーカス時: より明るい背景色
				focusedImage := b.createFocusedButtonImage()
				btn.SetImage(focusedImage)
			} else {
				// 非フォーカス時: 通常の半透明背景
				normalImage := b.createTransparentButtonImage()
				btn.SetImage(normalImage)
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
	// フォーカス時: 半透明の灰色
	focused := image.NewNineSliceColor(consts.ButtonHoverColor)

	// プレス状態: さらに明るい色
	pressed := image.NewNineSliceColor(consts.ButtonPressedColor)

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

	pageText := menu.GetPageIndicatorText()

	return widget.NewText(
		widget.TextOpts.Text(pageText, &res.Text.SmallFace, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(300, 20),
		),
	)
}
