package states

import (
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	"github.com/kijimaD/ruins/lib/widgets/tabmenu"
	"github.com/kijimaD/ruins/lib/widgets/views"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// CraftMenuState はクラフトメニューのゲームステート
type CraftMenuState struct {
	es.BaseState[w.World]
	ui *ebitenui.UI

	tabMenu             *tabmenu.TabMenu
	keyboardInput       input.KeyboardInput
	selectedItem        string            // 選択中のレシピ名
	itemDesc            *widget.Text      // アイテムの概要
	specContainer       *widget.Container // 性能表示のコンテナ
	recipeList          *widget.Container // レシピリストのコンテナ
	resultWindow        *widget.Window    // 合成結果ウィンドウ
	rootContainer       *widget.Container
	tabDisplayContainer *widget.Container // タブ表示のコンテナ
	categoryContainer   *widget.Container // カテゴリ一覧のコンテナ

	// アクション選択ウィンドウ用
	actionWindow     *widget.Window // アクション選択ウィンドウ
	actionFocusIndex int            // アクションウィンドウ内のフォーカス
	actionItems      []string       // アクション項目リスト
	isWindowMode     bool           // ウィンドウ操作モードかどうか

	// 結果ウィンドウ用
	resultFocusIndex int        // 結果ウィンドウ内のフォーカス
	resultItems      []string   // 結果ウィンドウの項目リスト
	resultEntity     ecs.Entity // 生成された結果アイテムのエンティティ
	isResultMode     bool       // 結果ウィンドウ操作モードかどうか
}

func (st CraftMenuState) String() string {
	return "CraftMenu"
}

// State interface ================

var _ es.State[w.World] = &CraftMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *CraftMenuState) OnPause(_ w.World) error { return nil }

// OnResume はステートが再開される際に呼ばれる
func (st *CraftMenuState) OnResume(_ w.World) error { return nil }

// OnStart はステートが開始される際に呼ばれる
func (st *CraftMenuState) OnStart(world w.World) error {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
	return nil
}

// OnStop はステートが停止される際に呼ばれる
func (st *CraftMenuState) OnStop(_ w.World) error { return nil }

// Update はゲームステートの更新処理を行う
func (st *CraftMenuState) Update(world w.World) (es.Transition[w.World], error) {

	if st.keyboardInput.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewDebugMenuState}}, nil
	}

	// ウィンドウモードの場合はウィンドウ操作を優先
	if st.isWindowMode {
		if st.updateWindowMode(world) {
			return es.Transition[w.World]{Type: es.TransNone}, nil
		}
	}

	// 結果ウィンドウモードの場合は結果ウィンドウ操作を優先
	if st.isResultMode {
		if st.updateResultMode(world) {
			return es.Transition[w.World]{Type: es.TransNone}, nil
		}
	}

	st.tabMenu.Update()
	st.ui.Update()

	return st.ConsumeTransition(), nil
}

// Draw はゲームステートの描画処理を行う
func (st *CraftMenuState) Draw(_ w.World, screen *ebiten.Image) error {
	st.ui.Draw(screen)
	return nil
}

// ================

func (st *CraftMenuState) initUI(world w.World) *ebitenui.UI {
	res := world.Resources.UIResources

	// TabMenuの設定
	tabs := st.createTabs(world)
	config := tabmenu.Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
		ItemsPerPage:     20,
	}

	callbacks := tabmenu.Callbacks{
		OnSelectItem: func(_ int, _ int, tab tabmenu.TabItem, item menu.Item) {
			st.handleItemSelection(world, tab, item)
		},
		OnCancel: func() {
			// Escapeで前の画面に戻る
			st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
		},
		OnTabChange: func(_, _ int, _ tabmenu.TabItem) {
			st.updateTabDisplay(world)
			st.updateCategoryDisplay(world)
		},
		OnItemChange: func(_ int, _, _ int, item menu.Item) {
			st.handleItemChange(world, item)
			st.updateTabDisplay(world)
		},
	}

	st.tabMenu = tabmenu.NewTabMenu(config, callbacks, st.keyboardInput)

	// アイテムの説明文
	itemDescContainer := styled.NewRowContainer()
	st.itemDesc = styled.NewMenuText(" ", world.Resources.UIResources) // 空文字だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.specContainer = styled.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	st.recipeList = styled.NewVerticalContainer()

	// 初期状態の表示を更新
	st.updateInitialItemDisplay(world)

	// タブ表示のコンテナを作成
	st.tabDisplayContainer = styled.NewVerticalContainer()
	st.createTabDisplayUI(world)

	// カテゴリ一覧のコンテナを作成（横並び）
	st.categoryContainer = styled.NewRowContainer()
	st.createCategoryDisplayUI(world)

	st.rootContainer = styled.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		// 3x3グリッドレイアウト: 9個の要素が必要
		// 1行目
		st.rootContainer.AddChild(styled.NewTitleText("合成", world.Resources.UIResources))
		st.rootContainer.AddChild(st.categoryContainer) // カテゴリ一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())

		// 2行目
		st.rootContainer.AddChild(st.tabDisplayContainer) // タブとアイテム一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())
		st.rootContainer.AddChild(styled.NewVSplitContainer(st.specContainer, st.recipeList))

		// 3行目
		st.rootContainer.AddChild(itemDescContainer)
		st.rootContainer.AddChild(widget.NewContainer()) // 空
		st.rootContainer.AddChild(widget.NewContainer()) // 空
	}

	return &ebitenui.UI{Container: st.rootContainer}
}

// createTabs はTabMenuで使用するタブを作成する
func (st *CraftMenuState) createTabs(world w.World) []tabmenu.TabItem {
	tabs := []tabmenu.TabItem{
		{
			ID:    "consumables",
			Label: "道具",
			Items: st.createMenuItems(world, st.queryMenuConsumable(world)),
		},
		{
			ID:    "cards",
			Label: "手札",
			Items: st.createMenuItems(world, st.queryMenuCard(world)),
		},
		{
			ID:    "wearables",
			Label: "装備",
			Items: st.createMenuItems(world, st.queryMenuWearable(world)),
		},
	}

	return tabs
}

// createMenuItems はレシピ名リストをMenuItemに変換する
func (st *CraftMenuState) createMenuItems(_ w.World, recipeNames []string) []menu.Item {
	items := make([]menu.Item, len(recipeNames))

	for i, recipeName := range recipeNames {
		items[i] = menu.Item{
			ID:       fmt.Sprintf("recipe_%s", recipeName),
			Label:    recipeName,
			UserData: recipeName,
		}
	}

	return items
}

// handleItemSelection はアイテム選択時の処理
func (st *CraftMenuState) handleItemSelection(world w.World, _ tabmenu.TabItem, item menu.Item) {
	recipeName, ok := item.UserData.(string)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	st.selectedItem = recipeName
	st.showActionWindow(world, recipeName)
}

// handleItemChange はアイテム変更時の処理（カーソル移動）
func (st *CraftMenuState) handleItemChange(world w.World, item menu.Item) {
	// 無効なアイテムの場合は何もしない
	if item.UserData == nil {
		st.itemDesc.Label = " "
		st.specContainer.RemoveChildren()
		st.recipeList.RemoveChildren()
		return
	}

	recipeName, ok := item.UserData.(string)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	// RawMasterからEntitySpecを取得
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	spec, err := rawMaster.NewRecipeSpec(recipeName)
	if err != nil {
		st.itemDesc.Label = TextNoDescription
		st.specContainer.RemoveChildren()
		st.recipeList.RemoveChildren()
		return
	}

	// Descriptionを取得
	if spec.Description != nil {
		st.itemDesc.Label = spec.Description.Description
	} else {
		st.itemDesc.Label = TextNoDescription
	}

	// EntitySpecから性能表示を更新
	views.UpdateSpecFromSpec(world, st.specContainer, spec)
	if err := st.updateRecipeList(world, spec.Recipe); err != nil {
		log.Fatal(err)
	}
}

func (st *CraftMenuState) queryMenuConsumable(world w.World) []string {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	var items []string

	// 全レシピから消耗品（カード以外）を抽出
	for recipeName := range rawMaster.RecipeIndex {
		spec, err := rawMaster.NewRecipeSpec(recipeName)
		if err != nil {
			continue
		}
		// 消耗品でカード以外
		if spec.Consumable != nil && spec.Card == nil {
			items = append(items, recipeName)
		}
	}

	sort.Strings(items)
	return items
}

func (st *CraftMenuState) queryMenuCard(world w.World) []string {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	var items []string

	// 全レシピからカードを抽出
	for recipeName := range rawMaster.RecipeIndex {
		spec, err := rawMaster.NewRecipeSpec(recipeName)
		if err != nil {
			continue
		}
		if spec.Card != nil {
			items = append(items, recipeName)
		}
	}

	sort.Strings(items)
	return items
}

func (st *CraftMenuState) queryMenuWearable(world w.World) []string {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	var items []string

	// 全レシピから装備品を抽出
	for recipeName := range rawMaster.RecipeIndex {
		spec, err := rawMaster.NewRecipeSpec(recipeName)
		if err != nil {
			continue
		}
		if spec.Wearable != nil {
			items = append(items, recipeName)
		}
	}

	sort.Strings(items)
	return items
}

// showResultWindow は合成結果ウィンドウを表示する
func (st *CraftMenuState) showResultWindow(world w.World, entity ecs.Entity) {
	windowContainer := styled.NewWindowContainer(world.Resources.UIResources)
	titleContainer := styled.NewWindowHeaderContainer("合成結果", world.Resources.UIResources)
	st.resultWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// 結果項目を準備
	st.resultItems = []string{TextClose}
	st.resultFocusIndex = 0
	st.resultEntity = entity // 生成されたアイテムのエンティティを保存
	st.isResultMode = true

	// UI要素を作成（表示のみ、操作はキーボードで行う）
	st.createResultWindowUI(world, windowContainer, entity)

	st.resultWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.resultWindow)
}

// createResultWindowUI は結果ウィンドウのUI要素を作成する
func (st *CraftMenuState) createResultWindowUI(world w.World, container *widget.Container, entity ecs.Entity) {
	// アイテム詳細を表示
	views.UpdateSpec(world, container, entity)

	st.updateResultWindowDisplay(world)
}

// updateResultWindowDisplay は結果ウィンドウの表示を更新する
func (st *CraftMenuState) updateResultWindowDisplay(world w.World) {
	if st.resultWindow == nil {
		return
	}

	// 既存のウィンドウを閉じて新しく作成
	st.resultWindow.Close()

	windowContainer := styled.NewWindowContainer(world.Resources.UIResources)
	titleContainer := styled.NewWindowHeaderContainer("合成結果", world.Resources.UIResources)
	st.resultWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// アイテム詳細を表示（生成されたアイテムの値を使用）
	views.UpdateSpec(world, windowContainer, st.resultEntity)

	// ボタン項目を表示
	for i, action := range st.resultItems {
		isSelected := i == st.resultFocusIndex
		actionWidget := styled.NewListItemText(action, consts.TextColor, isSelected, world.Resources.UIResources)
		windowContainer.AddChild(actionWidget)
	}

	st.resultWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.resultWindow)
}

func (st *CraftMenuState) updateRecipeList(world w.World, recipe *gc.Recipe) error {
	st.recipeList.RemoveChildren()

	if recipe == nil {
		return fmt.Errorf("recipeがnilです")
	}

	for _, input := range recipe.Inputs {
		var currentAmount int
		if stackableEntity, found := worldhelper.FindStackableInInventory(world, input.Name); found {
			stackable := world.Components.Stackable.Get(stackableEntity).(*gc.Stackable)
			currentAmount = stackable.Count
		}
		str := fmt.Sprintf("%s %d pcs\n    所持: %d pcs", input.Name, input.Amount, currentAmount)
		var color color.RGBA
		if currentAmount >= input.Amount {
			color = consts.SuccessColor
		} else {
			color = consts.DangerColor
		}

		st.recipeList.AddChild(styled.NewBodyText(str, color, world.Resources.UIResources))
	}
	return nil
}

// showActionWindow はアクションウィンドウを表示する
func (st *CraftMenuState) showActionWindow(world w.World, recipeName string) {
	windowContainer := styled.NewWindowContainer(world.Resources.UIResources)
	titleContainer := styled.NewWindowHeaderContainer("アクション選択", world.Resources.UIResources)
	st.actionWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を準備
	st.actionItems = []string{}
	st.selectedItem = recipeName

	// 合成可能かチェック
	if canCraft, _ := worldhelper.CanCraft(world, recipeName); canCraft {
		st.actionItems = append(st.actionItems, "合成する")
	}
	st.actionItems = append(st.actionItems, TextClose)

	st.actionFocusIndex = 0
	st.isWindowMode = true

	// UI要素を作成（表示のみ、操作はキーボードで行う）
	st.updateActionWindowDisplay(world)

	st.actionWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.actionWindow)
}

// updateActionWindowDisplay はアクションウィンドウの表示を更新する
func (st *CraftMenuState) updateActionWindowDisplay(world w.World) {
	if st.actionWindow == nil {
		return
	}

	// 既存のウィンドウを閉じて新しく作成
	st.actionWindow.Close()

	windowContainer := styled.NewWindowContainer(world.Resources.UIResources)
	titleContainer := styled.NewWindowHeaderContainer("アクション選択", world.Resources.UIResources)
	st.actionWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を表示
	for i, action := range st.actionItems {
		isSelected := i == st.actionFocusIndex
		actionWidget := styled.NewListItemText(action, consts.TextColor, isSelected, world.Resources.UIResources)
		windowContainer.AddChild(actionWidget)
	}

	st.actionWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.actionWindow)
}

// updateWindowMode はウィンドウモード時の操作を処理する
func (st *CraftMenuState) updateWindowMode(world w.World) bool {
	// Escapeでウィンドウモードを終了
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		st.closeActionWindow()
		return false
	}

	// 上下矢印でフォーカス移動
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		st.actionFocusIndex--
		if st.actionFocusIndex < 0 {
			st.actionFocusIndex = len(st.actionItems) - 1
		}
		st.updateActionWindowDisplay(world)
		return true
	}
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		st.actionFocusIndex++
		if st.actionFocusIndex >= len(st.actionItems) {
			st.actionFocusIndex = 0
		}
		st.updateActionWindowDisplay(world)
		return true
	}

	// Enterで選択実行（押下-押上ワンセット）
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.executeActionItem(world)
		return true
	}

	return true
}

// updateResultMode は結果ウィンドウモード時の操作を処理する
func (st *CraftMenuState) updateResultMode(world w.World) bool {
	// Escapeで結果ウィンドウモードを終了
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		st.closeResultWindow()
		return false
	}

	// 上下矢印でフォーカス移動
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		st.resultFocusIndex--
		if st.resultFocusIndex < 0 {
			st.resultFocusIndex = len(st.resultItems) - 1
		}
		st.updateResultWindowDisplay(world)
		return true
	}
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		st.resultFocusIndex++
		if st.resultFocusIndex >= len(st.resultItems) {
			st.resultFocusIndex = 0
		}
		st.updateResultWindowDisplay(world)
		return true
	}

	// Enterで選択実行（押下-押上ワンセット）
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.executeResultItem(world)
		return true
	}

	return true
}

// closeActionWindow はアクションウィンドウを閉じる
func (st *CraftMenuState) closeActionWindow() {
	if st.actionWindow != nil {
		st.actionWindow.Close()
		st.actionWindow = nil
	}
	st.isWindowMode = false
	st.actionFocusIndex = 0
	st.actionItems = nil
}

// closeResultWindow は結果ウィンドウを閉じる
func (st *CraftMenuState) closeResultWindow() {
	if st.resultWindow != nil {
		st.resultWindow.Close()
		st.resultWindow = nil
	}
	st.isResultMode = false
	st.resultFocusIndex = 0
	st.resultItems = nil
	st.resultEntity = 0 // エンティティIDをリセット
}

// executeResultItem は選択された結果項目を実行する
func (st *CraftMenuState) executeResultItem(_ w.World) {
	if st.resultFocusIndex >= len(st.resultItems) {
		return
	}

	selectedAction := st.resultItems[st.resultFocusIndex]

	switch selectedAction {
	case TextClose:
		st.closeResultWindow()
	}
}

// executeActionItem は選択されたアクション項目を実行する
func (st *CraftMenuState) executeActionItem(world w.World) {
	if st.actionFocusIndex >= len(st.actionItems) {
		return
	}

	selectedAction := st.actionItems[st.actionFocusIndex]
	recipeName := st.selectedItem

	switch selectedAction {
	case "合成する":
		resultEntity, err := worldhelper.Craft(world, recipeName)
		if err != nil {
			log.Fatal(err)
		}

		// レシピリストを更新
		rawMaster := world.Resources.RawMaster.(*raw.Master)
		var spec gc.EntitySpec
		spec, err = rawMaster.NewRecipeSpec(recipeName)
		if err != nil {
			log.Fatal(err)
		}
		if err = st.updateRecipeList(world, spec.Recipe); err != nil {
			log.Fatal(err)
		}

		st.closeActionWindow()
		st.showResultWindow(world, *resultEntity)
		st.reloadTabs(world)
		st.updateTabDisplay(world)
		st.updateCategoryDisplay(world)
	case TextClose:
		st.closeActionWindow()
	}
}

// reloadTabs はタブの内容を再読み込みする
func (st *CraftMenuState) reloadTabs(world w.World) {
	newTabs := st.createTabs(world)
	st.tabMenu.UpdateTabs(newTabs)
	// UpdateTabs後に表示を更新
	st.updateTabDisplay(world)
}

// createTabDisplayUI はタブ表示UIを作成する
func (st *CraftMenuState) createTabDisplayUI(world w.World) {
	st.updateTabDisplay(world)
}

// createCategoryDisplayUI はカテゴリ表示UIを作成する
func (st *CraftMenuState) createCategoryDisplayUI(world w.World) {
	st.updateCategoryDisplay(world)
}

// updateCategoryDisplay はカテゴリ表示を更新する
func (st *CraftMenuState) updateCategoryDisplay(world w.World) {
	// 既存の子要素をクリア
	st.categoryContainer.RemoveChildren()

	// 全カテゴリを横並びで表示
	currentTabIndex := st.tabMenu.GetCurrentTabIndex()
	tabs := st.createTabs(world) // 最新のタブ情報を取得

	for i, tab := range tabs {
		isSelected := i == currentTabIndex
		if isSelected {
			// 選択中のカテゴリは背景色付きで明るい文字色
			categoryWidget := styled.NewListItemText(tab.Label, consts.TextColor, true, world.Resources.UIResources)
			st.categoryContainer.AddChild(categoryWidget)
		} else {
			// 非選択のカテゴリは背景なしでグレー文字色
			categoryWidget := styled.NewListItemText(tab.Label, consts.ForegroundColor, false, world.Resources.UIResources)
			st.categoryContainer.AddChild(categoryWidget)
		}
	}
}

// updateTabDisplay はタブ表示を更新する
func (st *CraftMenuState) updateTabDisplay(world w.World) {
	// 既存の子要素をクリア
	st.tabDisplayContainer.RemoveChildren()

	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	// タブ名を表示（サブタイトルとして）
	tabNameText := styled.NewSubtitleText(fmt.Sprintf("【%s】", currentTab.Label), world.Resources.UIResources)
	st.tabDisplayContainer.AddChild(tabNameText)

	// ページインジケーターを表示
	pageText := st.tabMenu.GetPageIndicatorText()
	if pageText != "" {
		pageIndicator := styled.NewPageIndicator(pageText, world.Resources.UIResources)
		st.tabDisplayContainer.AddChild(pageIndicator)
	}

	// 現在のページで表示されるアイテムとインデックスを取得
	visibleItems, indices := st.tabMenu.GetVisibleItems()

	// アイテム一覧を表示（ページ内のアイテムのみ）
	for i, item := range visibleItems {
		actualIndex := indices[i]
		isSelected := actualIndex == currentItemIndex && currentItemIndex >= 0
		if isSelected {
			// 選択中のアイテムは背景色付きで明るい文字色
			itemWidget := styled.NewListItemText(item.Label, consts.TextColor, true, world.Resources.UIResources)
			st.tabDisplayContainer.AddChild(itemWidget)
		} else {
			// 非選択のアイテムは背景なしでグレー文字色
			itemWidget := styled.NewListItemText(item.Label, consts.ForegroundColor, false, world.Resources.UIResources)
			st.tabDisplayContainer.AddChild(itemWidget)
		}
	}

	// アイテムがない場合の表示
	if len(currentTab.Items) == 0 {
		emptyText := styled.NewDescriptionText("(アイテムなし)", world.Resources.UIResources)
		st.tabDisplayContainer.AddChild(emptyText)
	}
}

// updateInitialItemDisplay は初期状態のアイテム表示を更新する
func (st *CraftMenuState) updateInitialItemDisplay(world w.World) {
	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	if len(currentTab.Items) > 0 && currentItemIndex >= 0 && currentItemIndex < len(currentTab.Items) {
		currentItem := currentTab.Items[currentItemIndex]
		st.handleItemChange(world, currentItem)
	}
}
