package states

import (
	"fmt"

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

// ShopMenuState はショップメニューのゲームステート
type ShopMenuState struct {
	es.BaseState[w.World]
	ui *ebitenui.UI

	tabMenu             *tabmenu.TabMenu
	keyboardInput       input.KeyboardInput
	selectedItem        menu.Item         // 選択中のアイテム
	itemDesc            *widget.Text      // アイテムの概要
	specContainer       *widget.Container // 性能表示のコンテナ
	rootContainer       *widget.Container
	tabDisplayContainer *widget.Container // タブ表示のコンテナ
	categoryContainer   *widget.Container // カテゴリ一覧のコンテナ

	// アクション選択ウィンドウ用
	actionWindow     *widget.Window // アクション選択ウィンドウ
	actionFocusIndex int            // アクションウィンドウ内のフォーカス
	actionItems      []string       // アクション項目リスト
	isWindowMode     bool           // ウィンドウ操作モードかどうか
}

func (st ShopMenuState) String() string {
	return "ShopMenu"
}

// State interface ================

var _ es.State[w.World] = &ShopMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *ShopMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *ShopMenuState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *ShopMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *ShopMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *ShopMenuState) Update(world w.World) es.Transition[w.World] {
	// ウィンドウモードの場合はウィンドウ操作を優先
	if st.isWindowMode {
		if st.updateWindowMode(world) {
			return es.Transition[w.World]{Type: es.TransNone}
		}
	}

	st.tabMenu.Update()
	st.ui.Update()

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *ShopMenuState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *ShopMenuState) initUI(world w.World) *ebitenui.UI {
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
		// 3x3グリッドレイアウト
		// 1行目
		st.rootContainer.AddChild(styled.NewTitleText("店", world.Resources.UIResources))
		st.rootContainer.AddChild(st.categoryContainer) // カテゴリ一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())

		// 2行目
		st.rootContainer.AddChild(st.tabDisplayContainer)
		st.rootContainer.AddChild(widget.NewContainer())
		st.rootContainer.AddChild(st.specContainer)

		// 3行目
		st.rootContainer.AddChild(itemDescContainer)
		st.rootContainer.AddChild(widget.NewContainer())
		st.rootContainer.AddChild(widget.NewContainer())
	}

	return &ebitenui.UI{Container: st.rootContainer}
}

// createTabs はTabMenuで使用するタブを作成する
func (st *ShopMenuState) createTabs(world w.World) []tabmenu.TabItem {
	return []tabmenu.TabItem{
		{
			ID:    "buy",
			Label: "購入",
			Items: st.createBuyItems(world),
		},
		{
			ID:    "sell",
			Label: "売却",
			Items: st.createSellItems(world),
		},
	}
}

// createBuyItems は購入アイテムリストを作成
func (st *ShopMenuState) createBuyItems(world w.World) []menu.Item {
	shopInventory := worldhelper.GetShopInventory()
	items := make([]menu.Item, 0, len(shopInventory))

	for _, itemName := range shopInventory {
		price := st.getItemPrice(world, itemName, true)
		items = append(items, menu.Item{
			Label:    fmt.Sprintf("%s  ◆ %d", itemName, price),
			UserData: map[string]interface{}{"itemName": itemName, "price": price, "isBuy": true},
		})
	}

	return items
}

// createSellItems は売却アイテムリストを作成
func (st *ShopMenuState) createSellItems(world w.World) []menu.Item {
	var items []menu.Item

	worldhelper.QueryPlayer(world, func(_ ecs.Entity) {
		world.Manager.Join(
			world.Components.Item,
			world.Components.Name,
			world.Components.ItemLocationInBackpack,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			nameComp := world.Components.Name.Get(entity).(*gc.Name)
			itemName := nameComp.Name

			baseValue := worldhelper.GetItemValue(world, entity)
			price := worldhelper.CalculateSellPrice(baseValue)

			displayName := itemName
			if entity.HasComponent(world.Components.Stackable) {
				stackable := world.Components.Stackable.Get(entity).(*gc.Stackable)
				displayName = fmt.Sprintf("%s (%d個)", itemName, stackable.Count)
			}

			items = append(items, menu.Item{
				Label: fmt.Sprintf("%s  ◆ %d", displayName, price),
				UserData: map[string]interface{}{
					"itemName": itemName,
					"entity":   entity,
					"price":    price,
					"isBuy":    false,
				},
			})
		}))
	})

	if len(items) == 0 {
		items = append(items, menu.Item{
			Label:    "売却可能なアイテムがありません",
			UserData: map[string]interface{}{},
		})
	}

	return items
}

// getItemPrice はアイテム名から価格を計算
func (st *ShopMenuState) getItemPrice(world w.World, itemName string, isBuy bool) int {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	itemIdx, ok := rawMaster.ItemIndex[itemName]
	if !ok {
		return 0
	}
	itemDef := rawMaster.Raws.Items[itemIdx]
	if itemDef.Value == nil {
		return 0
	}

	baseValue := *itemDef.Value
	if isBuy {
		return worldhelper.CalculateBuyPrice(baseValue)
	}
	return worldhelper.CalculateSellPrice(baseValue)
}

// handleItemSelection はアイテム選択時の処理
func (st *ShopMenuState) handleItemSelection(world w.World, _ tabmenu.TabItem, item menu.Item) {
	if item.UserData == nil {
		return
	}

	st.selectedItem = item
	st.showActionWindow(world, item)
}

// handleItemChange はアイテムフォーカス変更時の処理
func (st *ShopMenuState) handleItemChange(world w.World, item menu.Item) {
	if item.UserData == nil {
		st.itemDesc.Label = " "
		st.specContainer.RemoveChildren()
		return
	}

	data := item.UserData.(map[string]interface{})
	itemName, _ := data["itemName"].(string)

	// 性能表示をクリア
	st.specContainer.RemoveChildren()

	// 一時的にエンティティを作成して性能表示とDescription取得
	tempEntity, err := worldhelper.SpawnItem(world, itemName, gc.ItemLocationInBackpack)
	if err != nil {
		st.itemDesc.Label = TextNoDescription
		return
	}
	defer world.Manager.DeleteEntity(tempEntity)

	// Descriptionコンポーネントの存在チェック
	if !tempEntity.HasComponent(world.Components.Description) {
		st.itemDesc.Label = TextNoDescription
	} else {
		desc := world.Components.Description.Get(tempEntity).(*gc.Description)
		if desc == nil {
			st.itemDesc.Label = TextNoDescription
		} else {
			st.itemDesc.Label = desc.Description
		}
	}

	views.UpdateSpec(world, st.specContainer, tempEntity)
}

// handlePurchase はアイテムの購入処理
func (st *ShopMenuState) handlePurchase(world w.World, item menu.Item) {
	data := item.UserData.(map[string]interface{})
	itemName, _ := data["itemName"].(string)

	worldhelper.QueryPlayer(world, func(playerEntity ecs.Entity) {
		err := worldhelper.BuyItem(world, playerEntity, itemName)
		if err != nil {
			// エラーの場合は何もしない（通貨不足など）
			return
		}
		// タブを再読み込み
		st.reloadTabs(world)
	})
}

// handleSell はアイテムの売却処理
func (st *ShopMenuState) handleSell(world w.World, item menu.Item) {
	data := item.UserData.(map[string]interface{})
	itemName, _ := data["itemName"].(string)
	entity, ok := data["entity"].(ecs.Entity)
	if !ok {
		return
	}
	price, _ := data["price"].(int)

	worldhelper.QueryPlayer(world, func(playerEntity ecs.Entity) {
		err := worldhelper.SellItem(world, playerEntity, entity)
		if err != nil {
			st.itemDesc.Label = fmt.Sprintf("売却失敗: %v", err)
		} else {
			st.itemDesc.Label = fmt.Sprintf("%sを売却しました（◆ %d）", itemName, price)
			// タブを再読み込み
			st.reloadTabs(world)
		}
	})
}

// reloadTabs はタブの内容を再読み込みする
func (st *ShopMenuState) reloadTabs(world w.World) {
	newTabs := st.createTabs(world)
	st.tabMenu.UpdateTabs(newTabs)
	// UpdateTabs後に表示を更新
	st.updateTabDisplay(world)
	st.updateCategoryDisplay(world)
}

// createTabDisplayUI はタブ表示UIを作成する
func (st *ShopMenuState) createTabDisplayUI(world w.World) {
	st.updateTabDisplay(world)
}

// createCategoryDisplayUI はカテゴリ表示UIを作成する
func (st *ShopMenuState) createCategoryDisplayUI(world w.World) {
	st.updateCategoryDisplay(world)
}

// updateCategoryDisplay はカテゴリ表示を更新する
func (st *ShopMenuState) updateCategoryDisplay(world w.World) {
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
func (st *ShopMenuState) updateTabDisplay(world w.World) {
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

// showActionWindow はアクションウィンドウを表示する
func (st *ShopMenuState) showActionWindow(world w.World, item menu.Item) {
	windowContainer := styled.NewWindowContainer(world.Resources.UIResources)
	titleContainer := styled.NewWindowHeaderContainer("アクション選択", world.Resources.UIResources)
	st.actionWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を準備
	st.actionItems = []string{}
	st.selectedItem = item

	data := item.UserData.(map[string]interface{})
	isBuy, _ := data["isBuy"].(bool)

	if isBuy {
		st.actionItems = append(st.actionItems, "購入する")
	} else {
		st.actionItems = append(st.actionItems, "売却する")
	}
	st.actionItems = append(st.actionItems, TextClose)

	st.actionFocusIndex = 0
	st.isWindowMode = true

	// UI要素を作成（表示のみ、操作はキーボードで行う）
	st.createActionWindowUI(world, windowContainer)

	st.actionWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.actionWindow)
}

// createActionWindowUI はアクションウィンドウのUI要素を作成する
func (st *ShopMenuState) createActionWindowUI(world w.World, _ *widget.Container) {
	st.updateActionWindowDisplay(world)
}

// updateActionWindowDisplay はアクションウィンドウの表示を更新する
func (st *ShopMenuState) updateActionWindowDisplay(world w.World) {
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
func (st *ShopMenuState) updateWindowMode(world w.World) bool {
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

// closeActionWindow はアクションウィンドウを閉じる
func (st *ShopMenuState) closeActionWindow() {
	if st.actionWindow != nil {
		st.actionWindow.Close()
		st.actionWindow = nil
	}
	st.isWindowMode = false
	st.actionFocusIndex = 0
	st.actionItems = nil
}

// executeActionItem は選択されたアクション項目を実行する
func (st *ShopMenuState) executeActionItem(world w.World) {
	if st.actionFocusIndex >= len(st.actionItems) {
		return
	}

	selectedAction := st.actionItems[st.actionFocusIndex]

	switch selectedAction {
	case "購入する":
		st.handlePurchase(world, st.selectedItem)
		st.closeActionWindow()
	case "売却する":
		st.handleSell(world, st.selectedItem)
		st.closeActionWindow()
	case TextClose:
		st.closeActionWindow()
	}
}
