package messagewindow

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
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
	window      *widget.Window

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

	// 初期化
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
	mainContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	w.ui = &ebitenui.UI{
		Container: mainContainer,
	}
	w.createAndAddWindow()
}

// createAndAddWindow はウィンドウを作成してUIに追加する
func (w *Window) createAndAddWindow() {
	// ウィンドウコンテナ
	windowContainer := w.createWindowContainer()

	// ウィンドウサイズを選択肢に応じて調整
	windowSize := w.calculateWindowSize()

	// タイトルバー付きウィンドウを作成
	titleContainer := w.createTitleContainer()
	w.window = styled.NewSmallWindow(
		titleContainer,
		windowContainer,
		widget.WindowOpts.CloseMode(widget.NONE), // クリックで閉じない
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.MinSize(windowSize.Width, windowSize.Height),
		widget.WindowOpts.MaxSize(windowSize.Width, windowSize.Height), // 固定サイズ
	)

	// 明示的に画面中央に配置
	screenWidth := w.world.Resources.ScreenDimensions.Width
	screenHeight := w.world.Resources.ScreenDimensions.Height
	x := (screenWidth - windowSize.Width) / 2
	y := (screenHeight - windowSize.Height) / 2
	w.window.SetLocation(image.Rect(x, y, x+windowSize.Width, y+windowSize.Height))

	// UIにウィンドウを追加
	w.ui.AddWindow(w.window)
}

// createTitleContainer はタイトルコンテナを作成する
func (w *Window) createTitleContainer() *widget.Container {
	title := ""
	if w.content.SpeakerName != "" {
		title = w.content.SpeakerName
	}
	return styled.NewWindowHeaderContainer(title, w.world.Resources.UIResources)
}

// calculateWindowSize は選択肢に応じてウィンドウサイズを計算する
func (w *Window) calculateWindowSize() WindowSize {
	baseHeight := w.config.Size.Height

	// 選択肢がある場合は高さを追加
	if w.hasChoices && len(w.content.Choices) > 0 {
		// 選択肢1つあたり約30px、最低400px確保
		choiceHeight := len(w.content.Choices) * 30
		minHeightWithChoices := 400
		if baseHeight+choiceHeight > minHeightWithChoices {
			baseHeight = baseHeight + choiceHeight
		} else {
			baseHeight = minHeightWithChoices
		}
	}

	return WindowSize{
		Width:  w.config.Size.Width,
		Height: baseHeight,
	}
}

// createWindowContainer はウィンドウコンテナを作成する
func (w *Window) createWindowContainer() *widget.Container {
	windowContainer := styled.NewWindowContainer(w.world.Resources.UIResources)

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
	} else {
		// 選択肢がない場合は Enter プロンプトを表示する
		enterPrompt := w.createEnterPrompt()
		windowContainer.AddChild(enterPrompt)
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

// createEnterPrompt はEnter待ちプロンプトを作成する
func (w *Window) createEnterPrompt() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(0),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 15, Right: 10}),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter, // 中央寄せ
			}),
		),
	)

	// プロンプトテキスト
	promptText := "Enter"

	prompt := styled.NewListItemText(
		promptText,
		color.RGBA{255, 255, 255, 255}, // 白色テキスト
		true,                           // 選択状態（背景色付き）
		w.world.Resources.UIResources,
	)

	container.AddChild(prompt)
	return container
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
			ID:       choice.Text,
			Label:    choice.Text,
			Disabled: choice.Disabled,
			UserData: i, // インデックスを保存
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
