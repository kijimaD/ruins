package states

import (
	"fmt"
	"image"
	"strconv"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/engine/world"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/styles"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	selection     int
	inventoryMenu []ecs.Entity
	ui            *ebitenui.UI

	selectedItem       ecs.Entity        // 選択中のアイテム
	selectedItemButton *widget.Button    // 使用済みのアイテムのボタン
	items              []ecs.Entity      // 表示対象とするアイテム
	itemDesc           *widget.Text      // アイテムの概要
	itemList           *widget.Container // アイテムリストのコンテナ
	partyWindow        *widget.Window    // 仲間を選択するウィンドウ
	weaponAccuracy     *widget.Text      // 武器の命中率
	weaponBaseDamage   *widget.Text      // 武器の攻撃力
	weaponConsumption  *widget.Text      // 武器の消費エネルギー
}

// State interface ================

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.inventoryMenu = append(st.inventoryMenu, loader.AddEntities(world, prefabs.Menu.InventoryMenu)...)
	st.ui = st.initUI(world)
}

func (st *InventoryMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.inventoryMenu...)
}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CampMenuState{}}}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	st.ui.Update()

	return updateMenu(st, world)
}

func (st *InventoryMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// Menu Interface ================

func (st *InventoryMenuState) getSelection() int {
	return st.selection
}

func (st *InventoryMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *InventoryMenuState) confirmSelection(world w.World) states.Transition {
	return states.Transition{Type: states.TransNone}
}

func (st *InventoryMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *InventoryMenuState) getCursorMenuIDs() []string {
	return []string{""}
}

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	// 各アイテムが入るコンテナ
	st.itemList = eui.NewScrollContentContainer()

	// アイテムの説明文
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空白だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.queryMenuConsumable(world)
	toggleContainer := eui.NewRowContainer()
	toggleConsumableButton := eui.NewItemButton("アイテム", func(args *widget.ButtonClickedEventArgs) { st.queryMenuConsumable(world) }, world)
	toggleWeaponButton := eui.NewItemButton("武器", func(args *widget.ButtonClickedEventArgs) { st.queryMenuWeapon(world) }, world)
	toggleContainer.AddChild(toggleConsumableButton)
	toggleContainer.AddChild(toggleWeaponButton)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("インベントリ", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(toggleContainer)

		sc, v := eui.NewScrollContainer(st.itemList)
		rootContainer.AddChild(sc)
		rootContainer.AddChild(v)
		rootContainer.AddChild(st.newItemSpecContainer(world))

		rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: rootContainer}
}

// 新しいクエリを実行してitemsをセットする
func (st *InventoryMenuState) queryMenuConsumable(world w.World) {
	st.itemList.RemoveChildren()
	st.items = []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Name,
		gameComponents.Description,
		gameComponents.InBackpack,
		gameComponents.Consumable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.items = append(st.items, entity)
	}))
	st.generateList(world)
}

// 新しいクエリを実行してitemsをセットする
func (st *InventoryMenuState) queryMenuWeapon(world w.World) {
	st.itemList.RemoveChildren()
	st.items = []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Name,
		gameComponents.Description,
		gameComponents.InBackpack,
		gameComponents.Weapon,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.items = append(st.items, entity)
	}))
	st.generateList(world)
}

// itemsからUIを生成する
func (st *InventoryMenuState) generateList(world world.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	for _, entity := range st.items {
		entity := entity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		windowContainer := eui.NewWindowContainer()
		titleContainer := eui.NewWindowHeaderContainer("アクション", world)
		actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

		// アイテムの名前がラベルについたボタン
		itemButton := eui.NewItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			x, y := ebiten.CursorPosition()
			r := image.Rect(0, 0, x, y)
			r = r.Add(image.Point{x + 20, y + 20})
			actionWindow.SetLocation(r)
			st.ui.AddWindow(actionWindow)

			st.selectedItem = entity
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.selectedItem != entity {
				st.selectedItem = entity
			}

			var description string
			world.Manager.Join(gameComponents.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
				if entity == st.selectedItem && entity.HasComponent(gameComponents.Description) {
					c := gameComponents.Description.Get(entity).(*gc.Description)
					description = c.Description
				}
			}))
			st.itemDesc.Label = description

			var accuracy string
			var baseDamage string
			var consumption string
			world.Manager.Join(gameComponents.Weapon).Visit(ecs.Visit(func(entity ecs.Entity) {
				if entity == st.selectedItem && entity.HasComponent(gameComponents.Weapon) {
					weapon := gameComponents.Weapon.Get(entity).(*gc.Weapon)
					accuracy = fmt.Sprintf("命中率 %s", strconv.Itoa(weapon.Accuracy))
					baseDamage = fmt.Sprintf("攻撃力 %s", strconv.Itoa(weapon.BaseDamage))
					consumption = fmt.Sprintf("消費SP %s", strconv.Itoa(weapon.EnergyConsumption))
				}
			}))
			st.weaponAccuracy.Label = accuracy
			st.weaponBaseDamage.Label = baseDamage
			st.weaponConsumption.Label = consumption
		})
		st.itemList.AddChild(itemButton)

		useButton := eui.NewItemButton("使う　", func(args *widget.ButtonClickedEventArgs) {
			x, y := ebiten.CursorPosition()
			r := image.Rect(0, 0, x, y)
			r = r.Add(image.Point{x + 20, y + 20})
			st.initPartyWindow(world)
			st.partyWindow.SetLocation(r)

			consumable := gameComponents.Consumable.Get(entity).(*gc.Consumable)
			switch consumable.TargetType.TargetNum {
			case gc.TargetSingle:
				st.ui.AddWindow(st.partyWindow)
				actionWindow.Close()
				st.selectedItem = entity
				st.selectedItemButton = itemButton
			case gc.TargetAll:
				effects.ItemTrigger(nil, entity, effects.Party{}, world)
				actionWindow.Close()
				st.itemList.RemoveChild(itemButton)
			}
		}, world)
		if gameComponents.Consumable.Get(entity) != nil {
			windowContainer.AddChild(useButton)
		}

		dropButton := eui.NewItemButton("捨てる", func(args *widget.ButtonClickedEventArgs) {
			world.Manager.DeleteEntity(entity)
			st.itemList.RemoveChild(itemButton)
			actionWindow.Close()
		}, world)
		windowContainer.AddChild(dropButton)

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
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		partyButton := eui.NewItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			effects.ItemTrigger(nil, st.selectedItem, effects.Single{entity}, world)
			st.partyWindow.Close()
			st.itemList.RemoveChild(st.selectedItemButton)
		}, world)
		partyContainer.AddChild(partyButton)
	}))
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
	st.weaponAccuracy = st.specText(world)
	st.weaponBaseDamage = st.specText(world)
	st.weaponConsumption = st.specText(world)
	itemSpecContainer.AddChild(st.weaponAccuracy)
	itemSpecContainer.AddChild(st.weaponBaseDamage)
	itemSpecContainer.AddChild(st.weaponConsumption)
	return itemSpecContainer
}

func (st *InventoryMenuState) specText(world w.World) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text("", eui.LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}
