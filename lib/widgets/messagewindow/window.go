package messagewindow

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/messagedata"
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

	// メッセージキュー管理
	queueManager   *QueueManager
	currentMessage *messagedata.MessageData
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

	// 現在のメッセージの完了コールバック実行
	if w.currentMessage != nil && w.currentMessage.OnComplete != nil {
		w.currentMessage.OnComplete()
	}

	// キューに次のメッセージがある場合は次を表示
	if w.queueManager != nil && w.queueManager.HasNext() {
		w.showNextMessage()
		return
	}

	w.isOpen = false

	// コールバック実行
	if w.onClose != nil {
		w.onClose()
	}
}

// showNextMessage は次のメッセージを表示する
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

	// UIを再初期化
	w.initialized = false
	w.ui = nil
	w.window = nil
	w.choiceMenu = nil
	w.choiceBuilder = nil
}

// updateContentFromMessage はMessageDataからcontentを更新する
func (w *Window) updateContentFromMessage(msg *messagedata.MessageData) {
	w.content.Text = msg.Text
	w.content.SpeakerName = msg.Speaker

	// 選択肢をconvertする
	w.content.Choices = make([]Choice, len(msg.Choices))
	for i, choice := range msg.Choices {
		choiceCopy := choice // クロージャのキャプチャ問題を回避
		w.content.Choices[i] = Choice{
			Text: choice.Text,
			Action: func() {
				if choiceCopy.Action != nil {
					choiceCopy.Action(w.world)
				}
				// 選択肢にメッセージが関連付けられている場合はキューに追加
				if choiceCopy.MessageData != nil {
					if w.queueManager == nil {
						w.queueManager = NewQueueManager()
					}
					w.queueManager.EnqueueFront(choiceCopy.MessageData)
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

	// 選択肢に応じて位置を調整
	screenWidth := w.world.Resources.ScreenDimensions.Width
	screenHeight := w.world.Resources.ScreenDimensions.Height
	x := (screenWidth - windowSize.Width) / 2

	var y int
	numChoices := len(w.content.Choices)

	if w.hasChoices && numChoices > 0 {
		// 選択肢の数に応じて位置を調整
		if numChoices <= 3 {
			// 少ない選択肢は画面中央
			y = (screenHeight - windowSize.Height) / 2
		} else if numChoices <= 8 {
			// 中程度の選択肢は上寄り中央
			y = screenHeight / 3
		} else {
			// 多い選択肢は上寄りに配置
			y = screenHeight / 5
		}
	} else {
		// 選択肢がない場合は画面中央に配置
		y = (screenHeight - windowSize.Height) / 2
	}

	// 画面からはみ出さないように調整
	margin := 30
	if y+windowSize.Height > screenHeight-margin {
		y = screenHeight - windowSize.Height - margin
	}
	if y < margin {
		y = margin
	}

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
	// AnchorLayoutを使用してEnterプロンプトを絶対位置に配置
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(w.world.Resources.UIResources.Panel.ImageTrans),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.GridLayoutData{
				MaxHeight: 500,
			}),
		),
	)

	// メッセージコンテンツエリア（上部全体）
	contentArea := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    20,
					Bottom: 60, // Enterプロンプト用のスペースを確保
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

	// 選択肢表示エリア
	if w.hasChoices && w.choiceMenu != nil {
		choicesContainer := w.createChoicesContainer()
		contentArea.AddChild(choicesContainer)
	}

	windowContainer.AddChild(contentArea)

	// 選択肢がない場合のみ Enter プロンプトを固定下部に表示
	if !w.hasChoices || w.choiceMenu == nil {
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

// createEnterPrompt はウィンドウ下部に固定されたEnterプロンプトを作成する
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
				HorizontalPosition: widget.AnchorLayoutPositionCenter, // 水平中央
				VerticalPosition:   widget.AnchorLayoutPositionEnd,    // 下部に固定
				Padding:            widget.Insets{Bottom: 15},         // 下端からの余白
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
			Disabled: false,
			UserData: i, // インデックスを保存
		}
	}

	// Menu設定
	// 選択肢の数に応じてページサイズを調整
	itemsPerPage := w.calculateItemsPerPage(len(w.content.Choices))

	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      itemsPerPage,
		ShowPageIndicator: true, // ページインジケーターを表示
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

// calculateItemsPerPage は選択肢の数とウィンドウサイズに応じて1ページあたりのアイテム数を計算する
func (w *Window) calculateItemsPerPage(totalItems int) int {
	// ウィンドウの高さからメッセージテキストとパディング、ページインジケーターを除いた利用可能な高さを計算
	windowHeight := w.calculateWindowSize().Height

	// メッセージテキスト部分とパディングを除く
	textAreaHeight := 150     // メッセージテキスト + パディング
	pageIndicatorHeight := 30 // ページインジケーターの高さ
	availableHeight := windowHeight - textAreaHeight - pageIndicatorHeight

	// 1つのアイテムあたりの高さ（35px）
	itemHeight := 35
	maxItemsPerPage := availableHeight / itemHeight

	// 最低3つ、最大15つの範囲で制限
	if maxItemsPerPage < 3 {
		maxItemsPerPage = 3
	} else if maxItemsPerPage > 15 {
		maxItemsPerPage = 15
	}

	// 総アイテム数がページサイズより少ない場合は、総アイテム数を返す
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
