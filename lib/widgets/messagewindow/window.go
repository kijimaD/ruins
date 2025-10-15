package messagewindow

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/inputmapper"
	"github.com/kijimaD/ruins/lib/messagedata"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

// Window はメッセージウィンドウを表す
type Window struct {
	config   Config
	content  MessageContent
	world    w.World
	onClose  func()
	onChoice func(choice Choice)

	isOpen      bool
	ui          *ebitenui.UI
	initialized bool
	window      *widget.Window

	// 選択肢がある場合、メニューシステムでページング可能な選択肢一覧を表示
	choiceMenu      *menu.Menu
	choiceBuilder   *menu.UIBuilder
	hasChoices      bool
	currentMenuPage int
	needsUIRebuild  bool // ページ変更時のUI再構築フラグ

	// 複数メッセージを順番に表示
	queueManager   *QueueManager
	currentMessage *messagedata.MessageData
}

// Update はウィンドウを更新する
func (w *Window) Update() {
	if !w.isOpen {
		return
	}

	if !w.initialized {
		w.hasChoices = len(w.content.Choices) > 0
		if w.hasChoices {
			w.initChoiceMenu()
		}
		w.initUI()
		w.initialized = true
	}

	if w.hasChoices && w.choiceMenu != nil {
		w.choiceMenu.Update()

		if w.needsUIRebuild {
			w.rebuildUI()
			w.needsUIRebuild = false
		}
	} else {
		if action, ok := w.HandleInput(); ok {
			w.DoAction(action)
		}
	}

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
// キューに次のメッセージがある場合は閉じずに次を表示する
func (w *Window) Close() {
	if !w.isOpen {
		return
	}

	if w.currentMessage != nil && w.currentMessage.OnComplete != nil {
		w.currentMessage.OnComplete()
	}

	if w.queueManager != nil && w.queueManager.HasNext() {
		w.showNextMessage()
		return
	}

	w.isOpen = false

	if w.onClose != nil {
		w.onClose()
	}
}

// showNextMessage はキューから次のメッセージを取り出して表示する
func (w *Window) showNextMessage() {
	if w.queueManager == nil {
		return
	}

	nextMessage := w.queueManager.Dequeue()
	if nextMessage == nil {
		w.isOpen = false
		if w.onClose != nil {
			w.onClose()
		}
		return
	}

	w.currentMessage = nextMessage
	w.updateContentFromMessage(nextMessage)

	// UI再初期化
	w.initialized = false
	w.ui = nil
	w.window = nil
	w.choiceMenu = nil
	w.choiceBuilder = nil
}

// updateContentFromMessage はMessageDataから表示コンテンツを更新する
func (w *Window) updateContentFromMessage(msg *messagedata.MessageData) {
	w.content.Text = msg.Text
	w.content.SpeakerName = msg.Speaker

	w.content.Choices = make([]Choice, len(msg.Choices))
	for i, choice := range msg.Choices {
		w.content.Choices[i] = Choice{
			Text: choice.Text,
			Action: func() {
				if choice.Action != nil {
					if err := choice.Action(w.world); err != nil {
						// TODO: エラーハンドリング改善
						panic(err)
					}
				}
				// 選択肢に関連メッセージがある場合はキュー先頭に追加して即座に表示
				if choice.MessageData != nil {
					if w.queueManager == nil {
						w.queueManager = NewQueueManager()
					}
					w.queueManager.EnqueueFront(choice.MessageData)
				}
			},
		}
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

	// ウィンドウ位置を計算
	x, y := w.calculateWindowPosition(windowSize)
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

// calculateWindowPosition はウィンドウの表示位置を計算する
func (w *Window) calculateWindowPosition(windowSize WindowSize) (x, y int) {
	screenWidth := w.world.Resources.ScreenDimensions.Width
	screenHeight := w.world.Resources.ScreenDimensions.Height

	x = (screenWidth - windowSize.Width) / 2

	numChoices := len(w.content.Choices)
	if w.hasChoices && numChoices > 0 {
		if numChoices <= 3 {
			y = (screenHeight - windowSize.Height) / 2
		} else if numChoices <= 8 {
			y = screenHeight / 3
		} else {
			y = screenHeight / 5
		}
	} else {
		y = (screenHeight - windowSize.Height) / 2
	}

	margin := 30
	if y+windowSize.Height > screenHeight-margin {
		y = screenHeight - windowSize.Height - margin
	}
	if y < margin {
		y = margin
	}

	return x, y
}

// calculateWindowSize は選択肢に応じてウィンドウサイズを計算する
func (w *Window) calculateWindowSize() WindowSize {
	baseHeight := w.config.Size.Height

	// 選択肢がある場合は高さを追加
	if w.hasChoices && len(w.content.Choices) > 0 {
		numChoices := len(w.content.Choices)

		// 選択肢1つあたりの高さを計算（フォントサイズ + 余白）
		choiceItemHeight := 35
		choiceHeight := numChoices * choiceItemHeight

		// 最低高さと最大高さを設定
		minHeightWithChoices := 400
		maxHeightWithChoices := int(float64(w.world.Resources.ScreenDimensions.Height) * 0.8) // 画面高さの80%まで

		calculatedHeight := baseHeight + choiceHeight + 100 // +100 for padding and buttons

		if calculatedHeight < minHeightWithChoices {
			baseHeight = minHeightWithChoices
		} else if calculatedHeight > maxHeightWithChoices {
			baseHeight = maxHeightWithChoices
		} else {
			baseHeight = calculatedHeight
		}
	}

	return WindowSize{
		Width:  w.config.Size.Width,
		Height: baseHeight,
	}
}

// createWindowContainer はウィンドウコンテナを作成する
func (w *Window) createWindowContainer() *widget.Container {
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(w.world.Resources.UIResources.Panel.ImageTrans),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxHeight: 500,
			}),
		),
	)

	contentArea := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(&widget.Insets{
					Top:    20,
					Bottom: 60,
					Left:   10,
					Right:  10,
				}),
				widget.RowLayoutOpts.Spacing(2),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				StretchHorizontal:  true,
				StretchVertical:    true,
			}),
		),
	)

	messageText := styled.NewListItemText(
		w.content.Text,
		w.config.TextStyle.Color,
		false,
		w.world.Resources.UIResources,
	)
	contentArea.AddChild(messageText)

	if w.hasChoices && w.choiceMenu != nil {
		choicesContainer := w.createChoicesContainer()
		contentArea.AddChild(choicesContainer)
	}

	windowContainer.AddChild(contentArea)

	if !w.hasChoices || w.choiceMenu == nil {
		enterPrompt := w.createEnterPrompt()
		windowContainer.AddChild(enterPrompt)
	}

	return windowContainer
}

// createChoicesContainer は選択肢コンテナを作成する
func (w *Window) createChoicesContainer() *widget.Container {
	if w.choiceBuilder != nil && w.choiceMenu != nil {
		return w.choiceBuilder.BuildUI(w.choiceMenu)
	}

	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(5),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 10}),
			),
		),
	)

	for i, choice := range w.content.Choices {
		choiceText := styled.NewListItemText(
			choice.Text,
			w.config.TextStyle.Color,
			i == 0,
			w.world.Resources.UIResources,
		)
		container.AddChild(choiceText)
	}

	return container
}

// createEnterPrompt はEnterプロンプトを作成する
func (w *Window) createEnterPrompt() *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(0),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				Padding:            &widget.Insets{Bottom: 15},
			}),
		),
	)

	prompt := styled.NewListItemText(
		"Enter",
		color.RGBA{255, 255, 255, 255},
		true,
		w.world.Resources.UIResources,
	)

	container.AddChild(prompt)
	return container
}

// HandleInput はキーボード入力をActionに変換する
func (w *Window) HandleInput() (inputmapper.ActionID, bool) {
	keyboardInput := input.GetSharedKeyboardInput()

	for _, key := range w.config.SkippableKeys {
		if key == ebiten.KeyEnter {
			if keyboardInput.IsEnterJustPressedOnce() {
				return inputmapper.ActionConfirm, true
			}
		} else if keyboardInput.IsKeyJustPressed(key) {
			return inputmapper.ActionSkip, true
		}
	}

	return "", false
}

// DoAction はActionを実行する
func (w *Window) DoAction(action inputmapper.ActionID) {
	switch action {
	case inputmapper.ActionConfirm, inputmapper.ActionSkip:
		w.Close()
	}
}

// initChoiceMenu は選択肢メニューを初期化する
func (w *Window) initChoiceMenu() {
	if len(w.content.Choices) == 0 {
		return
	}

	items := make([]menu.Item, len(w.content.Choices))
	for i, choice := range w.content.Choices {
		items[i] = menu.Item{
			ID:       choice.Text,
			Label:    choice.Text,
			Disabled: false,
			UserData: i,
		}
	}

	itemsPerPage := w.calculateItemsPerPage(len(w.content.Choices))

	config := menu.Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
		ItemsPerPage:   itemsPerPage,
	}

	callbacks := menu.Callbacks{
		OnSelect: func(index int, _ menu.Item) {
			w.selectChoice(index)
		},
		OnCancel: func() {
			w.Close()
		},
		OnFocusChange: func(_, _ int) {
			if w.choiceMenu != nil {
				newPage := w.choiceMenu.GetCurrentPage()
				if newPage != w.currentMenuPage {
					w.currentMenuPage = newPage
					w.needsUIRebuild = true
				}
			}

			if w.choiceBuilder != nil && !w.needsUIRebuild {
				w.choiceBuilder.UpdateFocus(w.choiceMenu)
			}
		},
	}

	w.choiceMenu = menu.NewMenu(config, callbacks)
	w.choiceBuilder = menu.NewUIBuilder(w.world)
	w.currentMenuPage = 1
}

// rebuildUI はUIを再構築する
func (w *Window) rebuildUI() {
	w.ui = nil
	w.window = nil
	w.initUI()
}

// calculateItemsPerPage は1ページあたりのアイテム数を計算する
func (w *Window) calculateItemsPerPage(totalItems int) int {
	windowHeight := w.calculateWindowSize().Height

	textAreaHeight := 150
	pageIndicatorHeight := 30
	availableHeight := windowHeight - textAreaHeight - pageIndicatorHeight

	itemHeight := 35
	maxItemsPerPage := availableHeight / itemHeight

	if maxItemsPerPage < 3 {
		maxItemsPerPage = 3
	} else if maxItemsPerPage > 15 {
		maxItemsPerPage = 15
	}

	if totalItems <= maxItemsPerPage {
		return totalItems
	}

	return maxItemsPerPage
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
