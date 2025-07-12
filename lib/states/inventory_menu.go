package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/tabmenu"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	ui *ebitenui.UI

	tabMenu             *tabmenu.TabMenu
	keyboardInput       input.KeyboardInput
	selectedItem        ecs.Entity        // 選択中のアイテム
	itemDesc            *widget.Text      // アイテムの概要
	specContainer       *widget.Container // 性能表示のコンテナ
	partyWindow         *widget.Window    // 仲間を選択するウィンドウ
	rootContainer       *widget.Container
	tabDisplayContainer *widget.Container  // タブ表示のコンテナ
	trans               *states.Transition // 状態遷移
}

func (st InventoryMenuState) String() string {
	return "InventoryMenu"
}

// State interface ================

var _ es.State = &InventoryMenuState{}

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

func (st *InventoryMenuState) OnStop(world w.World) {}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)

	if st.keyboardInput.IsKeyJustPressedIfDifferent(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	// 遷移が設定されていたら返す
	if st.trans != nil {
		trans := *st.trans
		st.trans = nil
		return trans
	}

	st.tabMenu.Update()
	st.ui.Update()

	return states.Transition{Type: states.TransNone}
}

func (st *InventoryMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	res := world.Resources.UIResources

	// TabMenuの設定
	tabs := st.createTabs(world)
	config := tabmenu.TabMenuConfig{
		Tabs:              tabs,
		InitialTabIndex:   0,
		InitialItemIndex:  0,
		WrapNavigation:    true,
		OnlyDifferentKeys: false, // 一時的にfalseにしてテスト
	}

	callbacks := tabmenu.TabMenuCallbacks{
		OnSelectItem: func(tabIndex int, itemIndex int, tab tabmenu.TabItem, item menu.MenuItem) {
			st.handleItemSelection(world, tab, item)
		},
		OnCancel: func() {
			// Escapeでホームメニューに戻る
			st.trans = &states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
		},
		OnTabChange: func(oldTabIndex, newTabIndex int, tab tabmenu.TabItem) {
			st.updateTabDisplay(world)
		},
		OnItemChange: func(tabIndex int, oldItemIndex, newItemIndex int, item menu.MenuItem) {
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

	st.rootContainer = eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		st.rootContainer.AddChild(eui.NewMenuText("インベントリ", world))
		st.rootContainer.AddChild(widget.NewContainer())
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
func (st *InventoryMenuState) createMenuItems(world w.World, entities []ecs.Entity) []menu.MenuItem {
	gameComponents := world.Components.Game.(*gc.Components)
	items := make([]menu.MenuItem, len(entities))

	for i, entity := range entities {
		name := gameComponents.Name.Get(entity).(*gc.Name).Name
		items[i] = menu.MenuItem{
			ID:       fmt.Sprintf("entity_%d", entity),
			Label:    string(name),
			UserData: entity,
		}
	}

	return items
}

// handleItemSelection はアイテム選択時の処理
func (st *InventoryMenuState) handleItemSelection(world w.World, tab tabmenu.TabItem, item menu.MenuItem) {
	entity, ok := item.UserData.(ecs.Entity)
	if !ok {
		log.Fatal("unexpected item UserData")
	}

	st.selectedItem = entity
	st.showActionWindow(world, entity)
}

// handleItemChange はアイテム変更時の処理（カーソル移動）
func (st *InventoryMenuState) handleItemChange(world w.World, item menu.MenuItem) {
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

	gameComponents := world.Components.Game.(*gc.Components)

	// Descriptionコンポーネントの存在チェック
	if !entity.HasComponent(gameComponents.Description) {
		st.itemDesc.Label = "説明なし"
		st.specContainer.RemoveChildren()
		return
	}

	desc := gameComponents.Description.Get(entity).(*gc.Description)
	if desc == nil {
		st.itemDesc.Label = "説明なし"
		st.specContainer.RemoveChildren()
		return
	}

	st.itemDesc.Label = desc.Description
	views.UpdateSpec(world, st.specContainer, entity)
}

// showActionWindow はアクションウィンドウを表示する
func (st *InventoryMenuState) showActionWindow(world w.World, entity ecs.Entity) {
	windowContainer := eui.NewWindowContainer(world)
	titleContainer := eui.NewWindowHeaderContainer("アクション", world)
	actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

	gameComponents := world.Components.Game.(*gc.Components)

	useButton := eui.NewButton("使う　",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			consumable := gameComponents.Consumable.Get(entity).(*gc.Consumable)
			switch consumable.TargetType.TargetNum {
			case gc.TargetSingle:
				st.initPartyWindow(world)
				st.partyWindow.SetLocation(getWinRect())
				st.ui.AddWindow(st.partyWindow)
				actionWindow.Close()
				st.selectedItem = entity
			case gc.TargetAll:
				effects.ItemTrigger(nil, entity, effects.Party{}, world)
				actionWindow.Close()
			}
			st.reloadTabs(world)
		}),
	)
	if entity.HasComponent(gameComponents.Consumable) {
		windowContainer.AddChild(useButton)
	}

	dropButton := eui.NewButton("捨てる",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			world.Manager.DeleteEntity(entity)
			actionWindow.Close()
			st.reloadTabs(world)
		}),
	)
	if !entity.HasComponent(gameComponents.Material) {
		windowContainer.AddChild(dropButton)
	}

	closeButton := eui.NewButton("閉じる",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			actionWindow.Close()
		}),
	)
	windowContainer.AddChild(closeButton)

	actionWindow.SetLocation(setWinRect())
	st.ui.AddWindow(actionWindow)
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

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationInBackpack,
		gameComponents.Wearable.Not(),
		gameComponents.Card.Not(),
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuCard(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Card,
		gameComponents.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuWearable(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Wearable,
		gameComponents.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuMaterial(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	worldhelper.QueryOwnedMaterial(func(entity ecs.Entity) {
		material := gameComponents.Material.Get(entity).(*gc.Material)
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

// updateTabDisplay はタブ表示を更新する
func (st *InventoryMenuState) updateTabDisplay(world w.World) {
	// 既存の子要素をクリア
	st.tabDisplayContainer.RemoveChildren()

	currentTab := st.tabMenu.GetCurrentTab()
	currentItemIndex := st.tabMenu.GetCurrentItemIndex()

	// タブ名を表示
	tabNameText := eui.NewMenuText(fmt.Sprintf("【%s】", currentTab.Label), world)
	st.tabDisplayContainer.AddChild(tabNameText)

	// アイテム一覧を表示
	for i, item := range currentTab.Items {
		itemText := item.Label
		if i == currentItemIndex && currentItemIndex >= 0 {
			itemText = "-> " + itemText // 選択中のアイテムにマーカーを追加
		} else {
			itemText = "   " + itemText
		}
		itemWidget := eui.NewMenuText(itemText, world)
		st.tabDisplayContainer.AddChild(itemWidget)
	}

	// アイテムがない場合の表示
	if len(currentTab.Items) == 0 {
		emptyText := eui.NewMenuText("  (アイテムなし)", world)
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

// メンバー選択画面を初期化する
func (st *InventoryMenuState) initPartyWindow(world w.World) {
	partyContainer := eui.NewWindowContainer(world)
	st.partyWindow = eui.NewSmallWindow(eui.NewWindowHeaderContainer("選択", world), partyContainer)
	rowContainer := eui.NewRowContainer()
	partyContainer.AddChild(rowContainer)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		memberContainer := eui.NewVerticalContainer()
		partyButton := eui.NewButton("使う",
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				effects.ItemTrigger(nil, st.selectedItem, effects.Single{entity}, world)
				st.partyWindow.Close()
				st.reloadTabs(world)
			}),
		)
		memberContainer.AddChild(partyButton)
		views.AddMemberBar(world, memberContainer, entity)

		rowContainer.AddChild(memberContainer)
	}))
}
