package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/tabmenu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// EquipMenuState は装備メニューのゲームステート
type EquipMenuState struct {
	es.BaseState
	ui *ebitenui.UI

	tabMenu             *tabmenu.TabMenu
	keyboardInput       input.KeyboardInput
	itemDesc            *widget.Text      // アイテムの説明
	specContainer       *widget.Container // 性能コンテナ
	abilityContainer    *widget.Container // メンバーの能力表示コンテナ
	rootContainer       *widget.Container
	tabDisplayContainer *widget.Container // タブ表示のコンテナ
	categoryContainer   *widget.Container // カテゴリ一覧のコンテナ

	// 選択中の味方
	curMemberIdx int
	// 現在のタブインデックス（装備選択時の復元用）
	previousTabIndex int

	// アクション選択ウィンドウ用
	actionWindow     *widget.Window // アクション選択ウィンドウ
	actionFocusIndex int            // アクションウィンドウ内のフォーカス
	actionItems      []string       // アクション項目リスト
	isWindowMode     bool           // ウィンドウ操作モードかどうか

	// 装備選択状態管理
	isEquipMode       bool                   // 装備選択モードかどうか
	equipSlotNumber   gc.EquipmentSlotNumber // 装備スロット番号
	previousEquipment *ecs.Entity            // 前の装備
	equipTargetMember ecs.Entity             // 装備対象のメンバー
}

func (st EquipMenuState) String() string {
	return "EquipMenu"
}

// 装備対象
type equipTarget int

const (
	equipTargetWear equipTarget = iota
	equipTargetCard
)

// State interface ================

var _ es.State = &EquipMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *EquipMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *EquipMenuState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *EquipMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *EquipMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *EquipMenuState) Update(world w.World) es.Transition {
	changed := gs.EquipmentChangedSystem(world)
	if changed {
		st.reloadAbilityContainer(world)
	}
	effects.RunEffectQueue(world)

	if st.keyboardInput.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition{Type: es.TransPush, NewStates: []es.State{&DebugMenuState{}}}
	}

	// ウィンドウモードの場合はウィンドウ操作を優先
	if st.isWindowMode {
		if st.updateWindowMode(world) {
			return es.Transition{Type: es.TransNone}
		}
	}

	st.tabMenu.Update()
	st.ui.Update()

	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *EquipMenuState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *EquipMenuState) initUI(world w.World) *ebitenui.UI {
	res := world.Resources.UIResources

	// TabMenuの設定
	tabs := st.createTabs(world)
	config := tabmenu.TabMenuConfig{
		Tabs:             tabs,
		InitialTabIndex:  0,
		InitialItemIndex: 0,
		WrapNavigation:   true,
	}

	callbacks := tabmenu.TabMenuCallbacks{
		OnSelectItem: func(_ int, _ int, tab tabmenu.TabItem, item menu.MenuItem) {
			st.handleItemSelection(world, tab, item)
		},
		OnCancel: func() {
			// Escapeでホームメニューに戻る
			st.SetTransition(es.Transition{Type: es.TransSwitch, NewStates: []es.State{&HomeMenuState{}}})
		},
		OnTabChange: func(_, _ int, _ tabmenu.TabItem) {
			st.updateTabDisplay(world)
			st.updateCategoryDisplay(world)
			st.updateAbilityDisplay(world)
		},
		OnItemChange: func(_ int, _, _ int, item menu.MenuItem) {
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
	st.abilityContainer = eui.NewVerticalContainer(
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
		st.rootContainer.AddChild(eui.NewTitleText("装備", world))
		st.rootContainer.AddChild(st.categoryContainer) // カテゴリ一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())

		st.rootContainer.AddChild(st.tabDisplayContainer) // タブとアイテム一覧の表示
		st.rootContainer.AddChild(widget.NewContainer())
		st.rootContainer.AddChild(eui.NewWSplitContainer(st.specContainer, st.abilityContainer))

		st.rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: st.rootContainer}
}

// createTabs はTabMenuで使用するタブを作成する
func (st *EquipMenuState) createTabs(world w.World) []tabmenu.TabItem {
	members := []ecs.Entity{}
	worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})

	tabs := []tabmenu.TabItem{}

	// 各メンバーごとにタブを作成（防具と手札のスロットを統合）
	for memberIdx, member := range members {
		memberName := world.Components.Game.Name.Get(member).(*gc.Name).Name

		// 防具と手札の両方のスロットを統合して作成
		allItems := st.createAllSlotItems(world, member, memberIdx)

		tabs = append(tabs, tabmenu.TabItem{
			ID:    fmt.Sprintf("member_%d", memberIdx),
			Label: memberName,
			Items: allItems,
		})
	}

	return tabs
}

// createAllSlotItems は防具と手札の全スロットのMenuItemを作成する
func (st *EquipMenuState) createAllSlotItems(world w.World, member ecs.Entity, _ int) []menu.MenuItem {
	items := []menu.MenuItem{}

	// 防具スロットを追加
	wearSlots := worldhelper.GetWearEquipments(world, member)
	for i, slot := range wearSlots {
		var name string
		if slot != nil {
			name = fmt.Sprintf("防具%d: %s", i+1, world.Components.Game.Name.Get(*slot).(*gc.Name).Name)
		} else {
			name = fmt.Sprintf("防具%d: -", i+1)
		}

		items = append(items, menu.MenuItem{
			ID:    fmt.Sprintf("wear_slot_%d", i),
			Label: name,
			UserData: map[string]interface{}{
				"member":     member,
				"slotNumber": i,
				"entity":     slot,
				"equipType":  equipTargetWear,
			},
		})
	}

	// 手札スロットを追加
	cardSlots := worldhelper.GetCardEquipments(world, member)
	for i, slot := range cardSlots {
		var name string
		if slot != nil {
			name = fmt.Sprintf("手札%d: %s", i+1, world.Components.Game.Name.Get(*slot).(*gc.Name).Name)
		} else {
			name = fmt.Sprintf("手札%d: -", i+1)
		}

		items = append(items, menu.MenuItem{
			ID:    fmt.Sprintf("card_slot_%d", i),
			Label: name,
			UserData: map[string]interface{}{
				"member":     member,
				"slotNumber": i,
				"entity":     slot,
				"equipType":  equipTargetCard,
			},
		})
	}

	return items
}

// handleItemSelection はアイテム選択時の処理
func (st *EquipMenuState) handleItemSelection(world w.World, _ tabmenu.TabItem, item menu.MenuItem) {
	if st.isEquipMode {
		// 装備選択モードの場合
		st.handleEquipItemSelection(world, item)
	} else {
		// スロット選択モードの場合
		userData, ok := item.UserData.(map[string]interface{})
		if !ok {
			log.Fatal("unexpected item UserData")
		}

		st.showActionWindow(world, userData)
	}
}

// handleItemChange はアイテム変更時の処理（カーソル移動）
func (st *EquipMenuState) handleItemChange(world w.World, item menu.MenuItem) {
	// 無効なアイテムの場合は何もしない
	if item.UserData == nil {
		st.itemDesc.Label = " "
		st.specContainer.RemoveChildren()
		return
	}

	if st.isEquipMode {
		// 装備選択モードの場合
		entity, ok := item.UserData.(ecs.Entity)
		if !ok {
			log.Fatal("unexpected item UserData")
		}

		if entity.HasComponent(world.Components.Game.Description) {
			desc := world.Components.Game.Description.Get(entity).(*gc.Description)
			st.itemDesc.Label = desc.Description
		}
		views.UpdateSpec(world, st.specContainer, entity)
	} else {
		// スロット選択モードの場合
		userData, ok := item.UserData.(map[string]interface{})
		if !ok {
			log.Fatal("unexpected item UserData")
		}

		slotEntity := userData["entity"].(*ecs.Entity)
		if slotEntity != nil {
			if (*slotEntity).HasComponent(world.Components.Game.Description) {
				desc := world.Components.Game.Description.Get(*slotEntity).(*gc.Description)
				st.itemDesc.Label = desc.Description
			}
			views.UpdateSpec(world, st.specContainer, *slotEntity)
		} else {
			st.itemDesc.Label = " "
			st.specContainer.RemoveChildren()
		}

		// 現在の選択に基づいてメンバー情報を更新
		if member, ok := userData["member"].(ecs.Entity); ok {
			members := []ecs.Entity{}
			worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
				members = append(members, entity)
			})

			// メンバーインデックスを更新
			for i, m := range members {
				if m == member {
					st.curMemberIdx = i
					break
				}
			}

			st.updateAbilityDisplay(world)
		}
	}
}

// createTabDisplayUI はタブ表示UIを作成する
func (st *EquipMenuState) createTabDisplayUI(world w.World) {
	st.updateTabDisplay(world)
}

// createCategoryDisplayUI はカテゴリ表示UIを作成する
func (st *EquipMenuState) createCategoryDisplayUI(world w.World) {
	st.updateCategoryDisplay(world)
}

// updateCategoryDisplay はカテゴリ表示を更新する
func (st *EquipMenuState) updateCategoryDisplay(world w.World) {
	// 既存の子要素をクリア
	st.categoryContainer.RemoveChildren()

	// 装備選択モードの場合は、全メンバーを表示して装備対象をハイライト
	// 選択中かをエンティティIDで比較している
	if st.isEquipMode {
		members := []ecs.Entity{}
		worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})

		for _, member := range members {
			memberName := world.Components.Game.Name.Get(member).(*gc.Name).Name
			isTargetMember := member == st.equipTargetMember

			if isTargetMember {
				// 装備対象のメンバーは背景色付きで明るい文字色
				categoryWidget := eui.NewListItemText(memberName, styles.TextColor, true, world)
				st.categoryContainer.AddChild(categoryWidget)
			} else {
				// その他のメンバーは背景なしでグレー文字色
				categoryWidget := eui.NewListItemText(memberName, styles.ForegroundColor, false, world)
				st.categoryContainer.AddChild(categoryWidget)
			}
		}
		return
	}

	// 通常モード: 全カテゴリを横並びで表示
	// 選択中かをタブ番号で比較している
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
func (st *EquipMenuState) updateTabDisplay(world w.World) {
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
func (st *EquipMenuState) updateInitialItemDisplay(world w.World) {
	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	if len(currentTab.Items) > 0 && currentItemIndex >= 0 && currentItemIndex < len(currentTab.Items) {
		currentItem := currentTab.Items[currentItemIndex]
		st.handleItemChange(world, currentItem)
	}
}

// updateAbilityDisplay はメンバー能力表示を更新する
func (st *EquipMenuState) updateAbilityDisplay(world w.World) {
	st.reloadAbilityContainer(world)
}

// メンバーの能力表示コンテナを更新する
func (st *EquipMenuState) reloadAbilityContainer(world w.World) {
	st.abilityContainer.RemoveChildren()

	currentTab := st.tabMenu.GetCurrentTab()
	// タブIDからメンバーIDを取得（新しいパターン: "member_0"）
	var memberIdx int
	if _, err := fmt.Sscanf(currentTab.ID, "member_%d", &memberIdx); err != nil {
		log.Printf("Failed to parse tab ID %s: %v", currentTab.ID, err)
		return
	}

	members := []ecs.Entity{}
	worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})

	if memberIdx >= len(members) {
		return
	}

	targetMember := members[memberIdx]

	views.AddMemberBar(world, st.abilityContainer, targetMember)

	attrs := world.Components.Game.Attributes.Get(targetMember).(*gc.Attributes)
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.VitalityLabel, attrs.Vitality.Total, attrs.Vitality.Modifier), styles.TextColor, world))
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.StrengthLabel, attrs.Strength.Total, attrs.Strength.Modifier), styles.TextColor, world))
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.SensationLabel, attrs.Sensation.Total, attrs.Sensation.Modifier), styles.TextColor, world))
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.DexterityLabel, attrs.Dexterity.Total, attrs.Dexterity.Modifier), styles.TextColor, world))
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.AgilityLabel, attrs.Agility.Total, attrs.Agility.Modifier), styles.TextColor, world))
	st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.DefenseLabel, attrs.Defense.Total, attrs.Defense.Modifier), styles.TextColor, world))
}

// 装備可能な防具を取得する
func (st *EquipMenuState) queryMenuWear(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	world.Manager.Join(
		world.Components.Game.Item,
		world.Components.Game.ItemLocationInBackpack,
		world.Components.Game.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// 装備可能な手札を取得する
func (st *EquipMenuState) queryMenuCard(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	world.Manager.Join(
		world.Components.Game.Item,
		world.Components.Game.ItemLocationInBackpack,
		world.Components.Game.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// showActionWindow はアクションウィンドウを表示する
func (st *EquipMenuState) showActionWindow(world w.World, userData map[string]interface{}) {
	windowContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("アクション選択", world)
	st.actionWindow = eui.NewSmallWindow(titleContainer, windowContainer)

	// アクション項目を準備
	st.actionItems = []string{}

	// スロットに装備されているかチェック
	slotEntity, hasEquipment := userData["entity"].(*ecs.Entity)
	if hasEquipment && slotEntity != nil {
		st.actionItems = append(st.actionItems, "外す")
	}
	st.actionItems = append(st.actionItems, "装備する")
	st.actionItems = append(st.actionItems, TextClose)

	st.actionFocusIndex = 0
	st.isWindowMode = true

	// UI要素を作成（表示のみ、操作はキーボードで行う）
	st.createActionWindowUI(world, windowContainer)

	st.actionWindow.SetLocation(getCenterWinRect())
	st.ui.AddWindow(st.actionWindow)
}

// createActionWindowUI はアクションウィンドウのUI要素を作成する
func (st *EquipMenuState) createActionWindowUI(world w.World, _ *widget.Container) {
	st.updateActionWindowDisplay(world)
}

// updateActionWindowDisplay はアクションウィンドウの表示を更新する
func (st *EquipMenuState) updateActionWindowDisplay(world w.World) {
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

	st.actionWindow.SetLocation(getCenterWinRect())
	st.ui.AddWindow(st.actionWindow)
}

// updateWindowMode はウィンドウモード時の操作を処理する
func (st *EquipMenuState) updateWindowMode(world w.World) bool {
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
	if st.keyboardInput.IsEnterJustPressedOnce() || st.keyboardInput.IsKeyJustPressed(ebiten.KeySpace) {
		st.executeActionItem(world)
		return true
	}

	return true
}

// closeActionWindow はアクションウィンドウを閉じる
func (st *EquipMenuState) closeActionWindow() {
	if st.actionWindow != nil {
		st.actionWindow.Close()
		st.actionWindow = nil
	}
	st.isWindowMode = false
	st.actionFocusIndex = 0
	st.actionItems = nil
}

// executeActionItem は選択されたアクション項目を実行する
func (st *EquipMenuState) executeActionItem(world w.World) {
	if st.actionFocusIndex >= len(st.actionItems) {
		return
	}

	selectedAction := st.actionItems[st.actionFocusIndex]
	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	if currentItemIndex < 0 || currentItemIndex >= len(currentTab.Items) {
		st.closeActionWindow()
		return
	}

	userData, ok := currentTab.Items[currentItemIndex].UserData.(map[string]interface{})
	if !ok {
		st.closeActionWindow()
		return
	}

	switch selectedAction {
	case "装備する":
		st.startEquipMode(world, userData)
	case "外す":
		st.unequipItem(world, userData)
	case TextClose:
		st.closeActionWindow()
	}
}

// startEquipMode は装備選択モードを開始する
func (st *EquipMenuState) startEquipMode(world w.World, userData map[string]interface{}) {
	member := userData["member"].(ecs.Entity)
	slotNumber := userData["slotNumber"].(int)
	equipType := userData["equipType"].(equipTarget)
	previousEquipment := userData["entity"].(*ecs.Entity)

	// 現在のタブインデックスを保存
	st.previousTabIndex = st.tabMenu.GetCurrentTabIndex()

	st.isEquipMode = true
	st.equipSlotNumber = gc.EquipmentSlotNumber(slotNumber)
	st.previousEquipment = previousEquipment
	st.equipTargetMember = member // 装備対象のメンバーを保存

	// 装備選択用のタブを作成
	var items []ecs.Entity
	switch equipType {
	case equipTargetWear:
		items = st.queryMenuWear(world)
	case equipTargetCard:
		items = st.queryMenuCard(world)
	}

	equipItems := st.createEquipMenuItems(world, items, member)

	newTabs := []tabmenu.TabItem{
		{
			ID:    "equip_selection",
			Label: "装備選択",
			Items: equipItems,
		},
	}

	st.tabMenu.UpdateTabs(newTabs)
	st.updateTabDisplay(world)
	st.updateCategoryDisplay(world)
	st.closeActionWindow()
}

// createEquipMenuItems は装備選択用のMenuItemを作成する
func (st *EquipMenuState) createEquipMenuItems(world w.World, entities []ecs.Entity, _ ecs.Entity) []menu.MenuItem {
	items := make([]menu.MenuItem, len(entities))

	for i, entity := range entities {
		name := world.Components.Game.Name.Get(entity).(*gc.Name).Name
		items[i] = menu.MenuItem{
			ID:       fmt.Sprintf("equip_entity_%d", entity),
			Label:    name,
			UserData: entity,
		}
	}

	return items
}

// handleEquipItemSelection は装備選択時の処理
func (st *EquipMenuState) handleEquipItemSelection(world w.World, item menu.MenuItem) {
	entity, ok := item.UserData.(ecs.Entity)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	// 前の装備を外す
	if st.previousEquipment != nil {
		worldhelper.Disarm(world, *st.previousEquipment)
	}

	// 保存されたメンバーに新しい装備を装着
	worldhelper.Equip(world, entity, st.equipTargetMember, st.equipSlotNumber)

	// 装備モードを終了して元の表示に戻る
	st.exitEquipMode(world)
}

// unequipItem は装備を外す
func (st *EquipMenuState) unequipItem(world w.World, userData map[string]interface{}) {
	slotEntity, hasEquipment := userData["entity"].(*ecs.Entity)
	if hasEquipment && slotEntity != nil {
		worldhelper.Disarm(world, *slotEntity)
		st.reloadTabs(world)
		st.updateTabDisplay(world)
		st.updateAbilityDisplay(world)
	}
	st.closeActionWindow()
}

// exitEquipMode は装備選択モードを終了する
func (st *EquipMenuState) exitEquipMode(world w.World) {
	st.isEquipMode = false
	st.equipSlotNumber = 0
	st.previousEquipment = nil
	st.equipTargetMember = 0 // メンバー情報をクリア

	// 元のタブに戻る
	newTabs := st.createTabs(world)
	st.tabMenu.UpdateTabs(newTabs)

	// 保存されたタブインデックスに復元
	if st.previousTabIndex >= 0 && st.previousTabIndex < len(newTabs) {
		st.tabMenu.SetTabIndex(st.previousTabIndex)
	}

	st.updateTabDisplay(world)
	st.updateCategoryDisplay(world)
	st.updateAbilityDisplay(world)
}

// reloadTabs はタブの内容を再読み込みする
func (st *EquipMenuState) reloadTabs(world w.World) {
	newTabs := st.createTabs(world)
	st.tabMenu.UpdateTabs(newTabs)
	st.updateTabDisplay(world)
}
