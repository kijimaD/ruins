package messagewindow

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

// Window はメッセージウィンドウ
type Window struct {
	config   Config
	content  MessageContent
	world    w.World
	onClose  func()
	onChoice func(choice Choice)

	// 状態管理
	isOpen      bool
	ui          *ebitenui.UI
	initialized bool

	// 入力管理
	keyboardInput input.KeyboardInput

	// 選択肢システム用
	choiceMenu    *menu.Menu
	choiceBuilder *menu.UIBuilder
	hasChoices    bool
}

// Update はウィンドウを更新する
func (w *Window) Update() {
	if !w.isOpen {
		return
	}

	// 初回のみ初期化
	if !w.initialized {
		// キーボード入力インスタンスを初期化
		w.keyboardInput = input.GetSharedKeyboardInput()
		// 選択肢状態を設定
		w.hasChoices = len(w.content.Choices) > 0
		if w.hasChoices {
			w.initChoiceMenu()
		}
		// UI初期化（選択肢設定後に実行）
		w.initUI()
		w.initialized = true
	}

	// 選択肢メニューがある場合は優先的に処理
	if w.hasChoices && w.choiceMenu != nil {
		w.choiceMenu.Update(w.keyboardInput)
	} else {
		// 通常のキーボード入力処理
		w.handleKeyboardInput()
	}

	// UI更新
	if w.ui != nil {
		w.ui.Update()
	}
}

// Draw はウィンドウを描画する
func (w *Window) Draw(screen *ebiten.Image) {
	if !w.isOpen || w.ui == nil {
		return
	}

	// 背景オーバーレイを描画
	if w.config.ShowBackground {
		w.drawBackground(screen)
	}

	// ウィンドウを描画
	w.ui.Draw(screen)
}

// IsOpen はウィンドウが開いているかを返す
func (w *Window) IsOpen() bool {
	return w.isOpen
}

// IsClosed はウィンドウが閉じているかを返す
func (w *Window) IsClosed() bool {
	return !w.isOpen
}

// Close はウィンドウを閉じる
func (w *Window) Close() {
	if !w.isOpen {
		return
	}

	w.isOpen = false

	// コールバック実行
	if w.onClose != nil {
		w.onClose()
	}
}

// initUI はUIを初期化する
func (w *Window) initUI() {
	// メインコンテナを作成
	mainContainer := w.createMainContainer()

	// UIを初期化
	w.ui = &ebitenui.UI{
		Container: mainContainer,
	}
}

// createMainContainer はメインコンテナを作成する
func (w *Window) createMainContainer() *widget.Container {
	// メインコンテナ（ウィンドウ全体）
	mainContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
	)

	// ウィンドウコンテナ
	windowContainer := w.createWindowContainer()

	// 位置を設定したラッパーコンテナ
	positionWrapper := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
			widget.WidgetOpts.MinSize(w.config.Size.Width, w.config.Size.Height),
		),
	)

	positionWrapper.AddChild(windowContainer)
	mainContainer.AddChild(positionWrapper)

	return mainContainer
}

// createWindowContainer はウィンドウコンテナを作成する
func (w *Window) createWindowContainer() *widget.Container {
	// ウィンドウの背景とボーダーを持つコンテナ
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(0),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    w.config.WindowStyle.Padding.Top,
					Bottom: w.config.WindowStyle.Padding.Bottom,
					Left:   w.config.WindowStyle.Padding.Left,
					Right:  w.config.WindowStyle.Padding.Right,
				}),
			),
		),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(w.config.WindowStyle.BackgroundColor),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(w.config.Size.Width, w.config.Size.Height),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// 話者名を表示（会話メッセージの場合）
	if w.content.SpeakerName != "" {
		speakerText := styled.NewListItemText(
			w.content.SpeakerName,
			w.config.TextStyle.Color,
			false,
			w.world.Resources.UIResources,
		)
		windowContainer.AddChild(speakerText)
	}

	// メッセージテキストを追加
	messageText := styled.NewListItemText(
		w.content.Text,
		w.config.TextStyle.Color,
		false,
		w.world.Resources.UIResources,
	)
	windowContainer.AddChild(messageText)

	// 選択肢表示エリア
	if w.hasChoices && w.choiceMenu != nil {
		choicesContainer := w.createChoicesContainer()
		windowContainer.AddChild(choicesContainer)
	}

	// アクションエリア（閉じるボタンなど）
	if w.config.ActionStyle.ShowCloseButton {
		actionContainer := w.createActionContainer()
		windowContainer.AddChild(actionContainer)
	}

	return windowContainer
}

// createChoicesContainer は選択肢コンテナを作成する
func (w *Window) createChoicesContainer() *widget.Container {
	// Menuコンポーネント用のUIを構築
	if w.choiceBuilder != nil && w.choiceMenu != nil {
		return w.choiceBuilder.BuildUI(w.choiceMenu)
	}

	// フォールバック: 選択肢を直接表示
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(5),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 10}),
			),
		),
	)

	// デバッグ用: 選択肢を直接テキストとして表示
	for i, choice := range w.content.Choices {
		choiceText := styled.NewListItemText(
			choice.Text,
			w.config.TextStyle.Color,
			i == 0, // 最初の選択肢を選択状態として表示
			w.world.Resources.UIResources,
		)
		container.AddChild(choiceText)
	}

	return container
}

// createActionContainer はアクションコンテナを作成する
func (w *Window) createActionContainer() *widget.Container {
	actionContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 10}),
			),
		),
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(w.config.ActionStyle.ActionAreaColor),
		),
	)

	// 閉じるボタン
	closeText := styled.NewListItemText(
		w.config.ActionStyle.CloseButtonText,
		w.config.ActionStyle.ActionTextColor,
		false,
		w.world.Resources.UIResources,
	)
	actionContainer.AddChild(closeText)

	return actionContainer
}

// handleKeyboardInput はキーボード入力を処理する（選択肢がない場合）
func (w *Window) handleKeyboardInput() {
	// Enterキーの重複押下を防ぐためにグローバル管理システムを使用
	if w.isEnterKeyPressed() {
		w.Close()
		return
	}

	// その他のスキップ可能キーをチェック（Enterを除く）
	for _, key := range w.config.SkippableKeys {
		if key != ebiten.KeyEnter && w.keyboardInput.IsKeyJustPressed(key) {
			w.Close()
			return
		}
	}
}

// isEnterKeyPressed はEnterキーが重複押下されないように管理された形で押されたかを返す
func (w *Window) isEnterKeyPressed() bool {
	// 設定でEnterがスキップ可能キーに含まれている場合のみチェック
	for _, key := range w.config.SkippableKeys {
		if key == ebiten.KeyEnter {
			return w.keyboardInput.IsEnterJustPressedOnce()
		}
	}
	return false
}

// initChoiceMenu は選択肢メニューを初期化する
func (w *Window) initChoiceMenu() {
	if len(w.content.Choices) == 0 {
		return
	}

	// Menu用のアイテムを作成
	items := make([]menu.Item, len(w.content.Choices))
	for i, choice := range w.content.Choices {
		items[i] = menu.Item{
			ID:          choice.Text,
			Label:       choice.Text,
			Description: choice.Description,
			Disabled:    choice.Disabled,
			UserData:    i, // インデックスを保存
		}
	}

	// Menu設定
	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      10,
		ShowPageIndicator: false,
	}

	// Menuコールバック
	callbacks := menu.Callbacks{
		OnSelect: func(index int, _ menu.Item) {
			w.selectChoice(index)
		},
		OnCancel: func() {
			w.Close()
		},
		OnFocusChange: func(_, _ int) {
			// フォーカス変更時にUIを更新
			if w.choiceBuilder != nil {
				w.choiceBuilder.UpdateFocus(w.choiceMenu)
			}
		},
	}

	// Menuを作成
	w.choiceMenu = menu.NewMenu(config, callbacks)
	w.choiceBuilder = menu.NewUIBuilder(w.world)
}

// selectChoice は選択肢を選択する
func (w *Window) selectChoice(index int) {
	if index < 0 || index >= len(w.content.Choices) {
		return
	}

	choice := w.content.Choices[index]

	// コールバック実行
	if w.onChoice != nil {
		w.onChoice(choice)
	}

	// アクション実行
	if choice.Action != nil {
		choice.Action()
	}

	// ウィンドウを閉じる
	w.Close()
}

// drawBackground は背景オーバーレイを描画する
func (w *Window) drawBackground(screen *ebiten.Image) {
	// 半透明の黒い背景を描画
	bounds := screen.Bounds()
	overlay := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	overlay.Fill(color.RGBA{0, 0, 0, 120})
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
}
