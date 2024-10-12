package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/engine/world"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	ui *ebitenui.UI

	selectedItem       ecs.Entity        // 選択中のアイテム
	selectedItemButton *widget.Button    // 使用済みのアイテムのボタン
	items              []ecs.Entity      // 表示対象とするアイテム
	itemDesc           *widget.Text      // アイテムの概要
	actionContainer    *widget.Container // アクションの起点となるコンテナ
	specContainer      *widget.Container // 性能表示のコンテナ
	partyWindow        *widget.Window    // 仲間を選択するウィンドウ
	category           ItemCategoryType
}

func (st InventoryMenuState) String() string {
	return "InventoryMenu"
}

// State interface ================

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *InventoryMenuState) OnStop(world w.World) {}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	st.ui.Update()

	return states.Transition{Type: states.TransNone}
}

func (st *InventoryMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================
var _ haveCategory = &InventoryMenuState{}

func (st *InventoryMenuState) setCategoryReload(world w.World, category ItemCategoryType) {
	st.category = category
	st.categoryReload(world)
}

func (st *InventoryMenuState) categoryReload(world w.World) {
	st.actionContainer.RemoveChildren()
	st.items = []ecs.Entity{}

	switch st.category {
	case ItemCategoryTypeItem:
		st.items = st.queryMenuItem(world)
	case ItemCategoryTypeWearable:
		st.items = st.queryMenuWearable(world)
	case ItemCategoryTypeCard:
		st.items = st.queryMenuCard(world)
	case ItemCategoryTypeMaterial:
		st.items = st.queryMenuMaterial(world)
	default:
		log.Fatal("未定義のcategory")
	}

	st.generateList(world)
}

// TODO: 後で整理する
func (st *InventoryMenuState) SetCategory(category ItemCategoryType) {
	st.category = category
}

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	// 各アクションが入るコンテナ
	st.actionContainer = eui.NewScrollContentContainer()
	st.categoryReload(world)

	// 種類トグル
	toggleContainer := st.newToggleContainer(world)

	// アイテムの説明文
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空文字だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	sc, v := eui.NewScrollContainer(st.actionContainer)
	st.specContainer = st.newItemSpecContainer(world)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("インベントリ", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(toggleContainer)

		rootContainer.AddChild(sc)
		rootContainer.AddChild(v)
		rootContainer.AddChild(st.specContainer)

		rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: rootContainer}
}

func (st *InventoryMenuState) queryMenuItem(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.InBackpack,
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
		gameComponents.InBackpack,
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
		gameComponents.InBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *InventoryMenuState) queryMenuMaterial(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	simple.OwnedMaterial(func(entity ecs.Entity) {
		material := gameComponents.Material.Get(entity).(*gc.Material)
		// 0で初期化してるから、インスタンスは全て存在する。個数で判定する
		if material.Amount > 0 {
			items = append(items, entity)
		}
	}, world)

	return items
}

// itemsからUIを生成する
// 使用などでアイテム数が変動した場合は再実行する必要がある
func (st *InventoryMenuState) generateList(world world.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	count := fmt.Sprintf("合計 %02d個", len(st.items))
	st.actionContainer.AddChild(eui.NewWindowHeaderContainer(count, world))
	for _, entity := range st.items {
		entity := entity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		windowContainer := eui.NewWindowContainer()
		titleContainer := eui.NewWindowHeaderContainer("アクション", world)
		actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

		// アイテムの名前がラベルについたボタン
		itemButton := eui.NewItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			actionWindow.SetLocation(setWinRect())
			st.ui.AddWindow(actionWindow)
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.selectedItem != entity {
				st.selectedItem = entity
			}
			st.itemDesc.Label = simple.GetDescription(world, entity).Description
			views.UpdateSpec(world, st.specContainer, entity)
		})
		st.actionContainer.AddChild(itemButton)

		useButton := eui.NewItemButton("使う　", func(args *widget.ButtonClickedEventArgs) {
			consumable := gameComponents.Consumable.Get(entity).(*gc.Consumable)
			switch consumable.TargetType.TargetNum {
			case gc.TargetSingle:
				st.initPartyWindow(world)
				st.partyWindow.SetLocation(getWinRect())

				st.ui.AddWindow(st.partyWindow)
				actionWindow.Close()
				st.selectedItem = entity
				st.selectedItemButton = itemButton
			case gc.TargetAll:
				effects.ItemTrigger(nil, entity, effects.Party{}, world)
				actionWindow.Close()
				st.actionContainer.RemoveChild(itemButton)
				st.categoryReload(world)
			}
		}, world)
		if entity.HasComponent(gameComponents.Consumable) {
			windowContainer.AddChild(useButton)
		}

		dropButton := eui.NewItemButton("捨てる", func(args *widget.ButtonClickedEventArgs) {
			world.Manager.DeleteEntity(entity)
			actionWindow.Close()
			st.categoryReload(world)
		}, world)
		if !entity.HasComponent(gameComponents.Material) {
			windowContainer.AddChild(dropButton)
		}

		closeButton := eui.NewItemButton("閉じる", func(args *widget.ButtonClickedEventArgs) {
			actionWindow.Close()
		}, world)
		windowContainer.AddChild(closeButton)
	}
}

// メンバー選択画面を初期化する
func (st *InventoryMenuState) initPartyWindow(world w.World) {
	partyContainer := eui.NewWindowContainer()
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
		partyButton := eui.NewItemButton("使う", func(args *widget.ButtonClickedEventArgs) {
			effects.ItemTrigger(nil, st.selectedItem, effects.Single{entity}, world)
			st.partyWindow.Close()
			st.actionContainer.RemoveChild(st.selectedItemButton)
			st.categoryReload(world)
		}, world)
		memberContainer.AddChild(partyButton)
		views.AddMemberBar(world, memberContainer, entity)

		rowContainer.AddChild(memberContainer)
	}))
}

func (st *InventoryMenuState) newToggleContainer(world w.World) *widget.Container {
	toggleContainer := eui.NewRowContainer()
	toggleConsumableButton := eui.NewItemButton("道具", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, ItemCategoryTypeItem) }, world)
	toggleCardButton := eui.NewItemButton("手札", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, ItemCategoryTypeCard) }, world)
	toggleWearableButton := eui.NewItemButton("防具", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, ItemCategoryTypeWearable) }, world)
	toggleMaterialButton := eui.NewItemButton("素材", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, ItemCategoryTypeMaterial) }, world)
	toggleContainer.AddChild(toggleConsumableButton)
	toggleContainer.AddChild(toggleCardButton)
	toggleContainer.AddChild(toggleWearableButton)
	toggleContainer.AddChild(toggleMaterialButton)

	return toggleContainer
}

func (st *InventoryMenuState) newItemSpecContainer(world w.World) *widget.Container {
	itemSpecContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.ForegroundColor)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(4),
				widget.RowLayoutOpts.Padding(widget.Insets{
					Top:    10,
					Bottom: 10,
					Left:   10,
					Right:  10,
				}),
			)),
	)

	return itemSpecContainer
}
