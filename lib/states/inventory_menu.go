package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/tabmenu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// InventoryMenuState はインベントリメニューのゲームステート
type InventoryMenuState struct {
	es.BaseState
	ui *ebitenui.UI

	tabMenu             *tabmenu.TabMenu
	keyboardInput       input.KeyboardInput
	selectedItem        ecs.Entity        // 選択中のアイテム
	itemDesc            *widget.Text      // アイテムの概要
	specContainer       *widget.Container // 性能表示のコンテナ
	partyWindow         *widget.Window    // 仲間を選択するウィンドウ
	rootContainer       *widget.Container
	tabDisplayContainer *widget.Container // タブ表示のコンテナ
	categoryContainer   *widget.Container // カテゴリ一覧のコンテナ

	// アクション選択ウィンドウ用
	actionWindow     *widget.Window // アクション選択ウィンドウ
	actionFocusIndex int            // アクションウィンドウ内のフォーカス
	actionItems      []string       // アクション項目リスト
	isWindowMode     bool           // ウィンドウ操作モードかどうか

	// パーティ選択ウィンドウ用
	partyFocusIndex int          // パーティウィンドウ内のフォーカス
	partyMembers    []ecs.Entity // パーティメンバーのエンティティリスト
	isPartyMode     bool         // パーティ選択モードかどうか
}

func (st InventoryMenuState) String() string {
	return "InventoryMenu"
}

// State interface ================

var _ es.State = &InventoryMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *InventoryMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *InventoryMenuState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *InventoryMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *InventoryMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *InventoryMenuState) Update(world w.World) es.Transition {

	if st.keyboardInput.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewDebugMenuState}}
	}

	// ウィンドウモードの場合はウィンドウ操作を優先
	if st.isWindowMode {
		if st.updateWindowMode(world) {
			return es.Transition{Type: es.TransNone}
		}
	}

	// パーティ選択モードの場合はパーティ操作を優先
	if st.isPartyMode {
		if st.updatePartyMode(world) {
			return es.Transition{Type: es.TransNone}
		}
	}

	st.tabMenu.Update()
	st.ui.Update()

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *InventoryMenuState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	res := world.Resources.UIResources

	// TabMenuの設定
	tabs := st.createTabs(world)
	config := tabmenu.Config{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
	}

	callbacks := tabmenu.Callbacks{
		OnSelectItem: func(_ int, _ int, tab tabmenu.TabItem, item menu.Item) {
			st.handleItemSelection(world, tab, item)
		},
		OnCancel: func() {
			// Escapeでホームメニューに戻る
			st.SetTransition(es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}})
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
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空文字だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.specContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)

	// 初期状態の表示を更新
	st.updateInitialItemDisplay(world)

	// タブ表示のコンテナを作成
	st.tabDisplayContainer = eui.NewVerticalContainer()
	st.createTabDisplayUI(world)

	// カテゴリ一覧のコンテナを作成（横並び）
	st.categoryContainer = eui.NewRowContainer()
	st.createCategoryDisplayUI(world)

	st.rootContainer = eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		st.rootContainer.AddChild(eui.NewTitleText("インベントリ", world))
		st.rootContainer.AddChild(st.categoryContainer) // カテゴリ一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())

		st.rootContainer.AddChild(st.tabDisplayContainer) // タブとアイテム一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())
		st.rootContainer.AddChild(st.specContainer)

		st.rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: st.rootContainer}
}

// createTabs はTabMenuで使用するタブを作成する
func (st *InventoryMenuState) createTabs(world w.World) []tabmenu.TabItem {
	tabs := []tabmenu.TabItem{
		{
			ID:    "items",
			Label: "道具",
			Items: st.createMenuItems(world, st.queryMenuItem(world)),
		},
		{
			ID:    "cards",
			Label: "手札",
			Items: st.createMenuItems(world, st.queryMenuCard(world)),
		},
		{
			ID:    "wearables",
			Label: "防具",
			Items: st.createMenuItems(world, st.queryMenuWearable(world)),
		},
		{
			ID:    "materials",
			Label: "素材",
			Items: st.createMenuItems(world, st.queryMenuMaterial(world)),
		},
	}

	return tabs
}

// createMenuItems はECSエンティティをMenuItemに変換する
func (st *InventoryMenuState) createMenuItems(world w.World, entities []ecs.Entity) []menu.Item {
	items := make([]menu.Item, len(entities))

	for i, entity := range entities {
		name := world.Components.Name.Get(entity).(*gc.Name).Name
		items[i] = menu.Item{
			ID:       fmt.Sprintf("entity_%d", entity),
			Label:    name,
			UserData: entity,
		}
	}

	return items
}

// handleItemSelection はアイテム選択時の処理
func (st *InventoryMenuState) handleItemSelection(world w.World, _ tabmenu.TabItem, item menu.Item) {
	entity, ok := item.UserData.(ecs.Entity)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	st.selectedItem = entity
	st.showActionWindow(world, entity)
}

// handleItemChange はアイテム変更時の処理（カーソル移動）
func (st *InventoryMenuState) handleItemChange(world w.World, item menu.Item) {
	// 無効なアイテムの場合は何もしない
	if item.UserData == nil {
		st.itemDesc.Label = " "
		st.specContainer.RemoveChildren()
		return
	}

	entity, ok := item.UserData.(ecs.Entity)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	// Descriptionコンポーネントの存在チェック
	if !entity.HasComponent(world.Components.Description) {
		st.itemDesc.Label = TextNoDescription
		st.specContainer.RemoveChildren()
		return
	}

	desc := world.Components.Description.Get(entity).(*gc.Description)
	if desc == nil {
		st.itemDesc.Label = TextNoDescription
		st.specContainer.RemoveChildren()
		return
	}

	st.itemDesc.Label = desc.Description
	views.UpdateSpec(world, st.specContainer, entity)
}

// updateWindowMode はウィンドウモード時の操作を処理する
func (st *InventoryMenuState) updateWindowMode(world w.World) bool {
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
func (st *InventoryMenuState) closeActionWindow() {
	if st.actionWindow != nil {
		st.actionWindow.Close()
		st.actionWindow = nil
	}
	st.isWindowMode = false
	st.actionFocusIndex = 0
	st.actionItems = nil
}

// closePartyWindow はパーティウィンドウを閉じる
func (st *InventoryMenuState) closePartyWindow() {
	if st.partyWindow != nil {
		st.partyWindow.Close()
		st.partyWindow = nil
	}
	st.isPartyMode = false
	st.partyFocusIndex = 0
	st.partyMembers = nil
}

// updatePartyMode はパーティ選択モード時の操作を処理する
func (st *InventoryMenuState) updatePartyMode(world w.World) bool {
	// Escapeでパーティモードを終了
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyEscape) {
		st.closePartyWindow()
		return false
	}

	memberCount := len(st.partyMembers)

	// 2x2グリッドでの移動
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if st.partyFocusIndex == memberCount { // キャンセル項目から上移動
			// キャンセルから上に移動する場合、最下段のメンバーに移動
			if memberCount >= 3 {
				st.partyFocusIndex = 2 // 下段左
			} else if memberCount >= 1 {
				st.partyFocusIndex = 0 // 上段左
			}
		} else if st.partyFocusIndex >= 2 { // 下段から上段へ
			st.partyFocusIndex -= 2
		} else { // 上段からキャンセルへ
			st.partyFocusIndex = memberCount
		}
		st.updatePartyWindowDisplay(world)
		return true
	}
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if st.partyFocusIndex == memberCount { // キャンセル項目から下移動
			// キャンセルから下に移動する場合、最上段のメンバーに移動
			st.partyFocusIndex = 0
		} else if st.partyFocusIndex < 2 { // 上段から下段へ
			if st.partyFocusIndex+2 < memberCount {
				st.partyFocusIndex += 2
			} else {
				st.partyFocusIndex = memberCount // キャンセルへ
			}
		} else { // 下段からキャンセルへ
			st.partyFocusIndex = memberCount
		}
		st.updatePartyWindowDisplay(world)
		return true
	}
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if st.partyFocusIndex == memberCount { // キャンセル項目は左右移動なし
			return true
		}
		if st.partyFocusIndex%2 == 0 { // 左列から右列へ（循環）
			if st.partyFocusIndex+1 < memberCount {
				st.partyFocusIndex++
			}
		} else { // 右列から左列へ
			st.partyFocusIndex--
		}
		st.updatePartyWindowDisplay(world)
		return true
	}
	if st.keyboardInput.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if st.partyFocusIndex == memberCount { // キャンセル項目は左右移動なし
			return true
		}
		if st.partyFocusIndex%2 == 0 { // 左列から右列へ
			if st.partyFocusIndex+1 < memberCount {
				st.partyFocusIndex++
			}
		} else { // 右列から左列へ（循環）
			st.partyFocusIndex--
		}
		st.updatePartyWindowDisplay(world)
		return true
	}

	// Enterでメンバー選択実行（押下-押上ワンセット）
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.executePartySelection(world)
		return true
	}

	return true
}

// showActionWindow はアクションウィンドウを表示する
func (st *InventoryMenuState) showActionWindow(world w.World, entity ecs.Entity) {
	windowContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("アクション選択", world)
	st.actionWindow = eui.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を準備
	st.actionItems = []string{}
	st.selectedItem = entity

	// 使用可能なアクションを登録
	if entity.HasComponent(world.Components.Consumable) {
		st.actionItems = append(st.actionItems, "使う")
	}
	if !entity.HasComponent(world.Components.Material) {
		st.actionItems = append(st.actionItems, "捨てる")
	}
	st.actionItems = append(st.actionItems, TextClose)

	st.actionFocusIndex = 0
	st.isWindowMode = true

	// UI要素を作成（表示のみ、操作はキーボードで行う）
	st.createActionWindowUI(world, windowContainer, entity)

	st.actionWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.actionWindow)
}

// createActionWindowUI はアクションウィンドウのUI要素を作成する
func (st *InventoryMenuState) createActionWindowUI(world w.World, _ *widget.Container, _ ecs.Entity) {
	st.updateActionWindowDisplay(world)
}

// updateActionWindowDisplay はアクションウィンドウの表示を更新する
func (st *InventoryMenuState) updateActionWindowDisplay(world w.World) {
	if st.actionWindow == nil {
		return
	}

	// 既存のウィンドウを閉じて新しく作成
	st.actionWindow.Close()

	windowContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("アクション選択", world)
	st.actionWindow = eui.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を表示
	for i, action := range st.actionItems {
		isSelected := i == st.actionFocusIndex
		actionWidget := eui.NewListItemText(action, styles.TextColor, isSelected, world)
		windowContainer.AddChild(actionWidget)
	}

	st.actionWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.actionWindow)
}

// executeActionItem は選択されたアクション項目を実行する
func (st *InventoryMenuState) executeActionItem(world w.World) {
	if st.actionFocusIndex >= len(st.actionItems) {
		return
	}

	selectedAction := st.actionItems[st.actionFocusIndex]

	switch selectedAction {
	case "使う":
		consumable := world.Components.Consumable.Get(st.selectedItem).(*gc.Consumable)
		switch consumable.TargetType.TargetNum {
		case gc.TargetSingle:
			st.closeActionWindow()
			st.initPartyWindowWithKeyboard(world)
		case gc.TargetAll:
			processor := effects.NewProcessor()
			useItemEffect := effects.UseItem{Item: st.selectedItem}
			partySelector := effects.TargetParty{}
			if err := processor.AddTargetedEffect(useItemEffect, nil, partySelector, world); err != nil {
				log.Printf("アイテムエフェクト追加エラー: %v", err)
			}
			if err := processor.Execute(world); err != nil {
				log.Printf("アイテムエフェクト実行エラー: %v", err)
			}
			st.closeActionWindow()
			st.reloadTabs(world)
			st.updateTabDisplay(world)
			st.updateCategoryDisplay(world)
		}
	case "捨てる":
		world.Manager.DeleteEntity(st.selectedItem)
		st.closeActionWindow()
		st.reloadTabs(world)
		st.updateTabDisplay(world)
		st.updateCategoryDisplay(world)
	case TextClose:
		st.closeActionWindow()
	}
}

// reloadTabs はタブの内容を再読み込みする
func (st *InventoryMenuState) reloadTabs(world w.World) {
	newTabs := st.createTabs(world)
	st.tabMenu.UpdateTabs(newTabs)
	// UpdateTabs後に表示を更新
	st.updateTabDisplay(world)
}

func (st *InventoryMenuState) queryMenuItem(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
		world.Components.Wearable.Not(),
		world.Components.Card.Not(),
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuCard(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	world.Manager.Join(
		world.Components.Item,
		world.Components.Card,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuWearable(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	world.Manager.Join(
		world.Components.Item,
		world.Components.Wearable,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuMaterial(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	worldhelper.QueryOwnedMaterial(func(entity ecs.Entity) {
		material := world.Components.Material.Get(entity).(*gc.Material)
		// 0で初期化してるから、インスタンスは全て存在する。個数で判定する
		if material.Amount > 0 {
			items = append(items, entity)
		}
	}, world)

	return items
}

// createTabDisplayUI はタブ表示UIを作成する
func (st *InventoryMenuState) createTabDisplayUI(world w.World) {
	st.updateTabDisplay(world)
}

// createCategoryDisplayUI はカテゴリ表示UIを作成する
func (st *InventoryMenuState) createCategoryDisplayUI(world w.World) {
	st.updateCategoryDisplay(world)
}

// updateCategoryDisplay はカテゴリ表示を更新する
func (st *InventoryMenuState) updateCategoryDisplay(world w.World) {
	// 既存の子要素をクリア
	st.categoryContainer.RemoveChildren()

	// 全カテゴリを横並びで表示
	currentTabIndex := st.tabMenu.GetCurrentTabIndex()
	tabs := st.createTabs(world) // 最新のタブ情報を取得

	for i, tab := range tabs {
		isSelected := i == currentTabIndex
		if isSelected {
			// 選択中のカテゴリは背景色付きで明るい文字色
			categoryWidget := eui.NewListItemText(tab.Label, styles.TextColor, true, world)
			st.categoryContainer.AddChild(categoryWidget)
		} else {
			// 非選択のカテゴリは背景なしでグレー文字色
			categoryWidget := eui.NewListItemText(tab.Label, styles.ForegroundColor, false, world)
			st.categoryContainer.AddChild(categoryWidget)
		}
	}
}

// updateTabDisplay はタブ表示を更新する
func (st *InventoryMenuState) updateTabDisplay(world w.World) {
	// 既存の子要素をクリア
	st.tabDisplayContainer.RemoveChildren()

	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	// タブ名を表示（サブタイトルとして）
	tabNameText := eui.NewSubtitleText(fmt.Sprintf("【%s】", currentTab.Label), world)
	st.tabDisplayContainer.AddChild(tabNameText)

	// アイテム一覧を表示
	for i, item := range currentTab.Items {
		isSelected := i == currentItemIndex && currentItemIndex >= 0
		if isSelected {
			// 選択中のアイテムは背景色付きで明るい文字色
			itemWidget := eui.NewListItemText(item.Label, styles.TextColor, true, world)
			st.tabDisplayContainer.AddChild(itemWidget)
		} else {
			// 非選択のアイテムは背景なしでグレー文字色
			itemWidget := eui.NewListItemText(item.Label, styles.ForegroundColor, false, world)
			st.tabDisplayContainer.AddChild(itemWidget)
		}
	}

	// アイテムがない場合の表示
	if len(currentTab.Items) == 0 {
		emptyText := eui.NewDescriptionText("(アイテムなし)", world)
		st.tabDisplayContainer.AddChild(emptyText)
	}
}

// updateInitialItemDisplay は初期状態のアイテム表示を更新する
func (st *InventoryMenuState) updateInitialItemDisplay(world w.World) {
	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	if len(currentTab.Items) > 0 && currentItemIndex >= 0 && currentItemIndex < len(currentTab.Items) {
		currentItem := currentTab.Items[currentItemIndex]
		st.handleItemChange(world, currentItem)
	}
}

// メンバー選択画面を初期化する（キーボード操作版）
func (st *InventoryMenuState) initPartyWindowWithKeyboard(world w.World) {
	partyContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("ターゲット選択", world)
	st.partyWindow = eui.NewSmallWindow(titleContainer, partyContainer)

	// パーティメンバーリストを作成
	st.partyMembers = []ecs.Entity{}
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
		world.Components.Name,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.partyMembers = append(st.partyMembers, entity)
	}))

	st.partyFocusIndex = 0
	st.isPartyMode = true

	// UI要素を作成
	st.updatePartyWindowDisplay(world)

	st.partyWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.partyWindow)
}

// updatePartyWindowDisplay はパーティウィンドウの表示を更新する
func (st *InventoryMenuState) updatePartyWindowDisplay(world w.World) {
	if st.partyWindow == nil {
		return
	}

	// 既存のウィンドウを閉じて新しく作成
	st.partyWindow.Close()

	partyContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("対象選択", world)
	st.partyWindow = eui.NewSmallWindow(titleContainer, partyContainer)

	// 2x2グリッドコンテナを作成
	gridContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(8, 8),
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true, true}),
		)),
	)

	// パーティメンバーを2x2グリッドで表示
	for i, memberEntity := range st.partyMembers {
		isSelected := i == st.partyFocusIndex

		// 選択状態に応じた背景色のコンテナを作成
		var memberContainer *widget.Container
		if isSelected {
			memberContainer = eui.NewVerticalContainer(
				widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(styles.ButtonHoverColor)),
			)
		} else {
			memberContainer = eui.NewVerticalContainer()
		}

		// メンバーバーを追加（キャラ名も含む）
		views.AddMemberBar(world, memberContainer, memberEntity)

		gridContainer.AddChild(memberContainer)
	}

	// グリッドをメインコンテナに追加
	partyContainer.AddChild(gridContainer)

	// キャンセル項目を別途追加（グリッドの下）
	cancelIndex := len(st.partyMembers)
	isSelected := st.partyFocusIndex == cancelIndex
	cancelWidget := eui.NewListItemText("キャンセル", styles.TextColor, isSelected, world)
	partyContainer.AddChild(cancelWidget)

	st.partyWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.partyWindow)
}

// executePartySelection は選択されたパーティメンバーでアイテムを使用する
func (st *InventoryMenuState) executePartySelection(world w.World) {
	// キャンセル項目の場合
	if st.partyFocusIndex >= len(st.partyMembers) {
		st.closePartyWindow()
		return
	}

	// 選択されたメンバーでアイテムを使用
	selectedMember := st.partyMembers[st.partyFocusIndex]
	processor := effects.NewProcessor()
	useItemEffect := effects.UseItem{Item: st.selectedItem}
	processor.AddEffect(useItemEffect, nil, selectedMember)
	if err := processor.Execute(world); err != nil {
		log.Printf("アイテムエフェクト実行エラー: %v", err)
	}

	st.closePartyWindow()
	st.reloadTabs(world)
	st.updateTabDisplay(world)
	st.updateCategoryDisplay(world)
}
