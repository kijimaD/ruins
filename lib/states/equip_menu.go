package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type EquipMenuState struct {
	selection int
	equipMenu []ecs.Entity
	ui        *ebitenui.UI

	slots            []*ecs.Entity     // スロット一覧
	items            []ecs.Entity      // インベントリにあるアイテム一覧
	subMenuContainer *widget.Container // 右上のサブメニュー
	actionContainer  *widget.Container // 操作の起点となるメインメニュー
	specContainer    *widget.Container // 性能コンテナ
	itemDesc         *widget.Text      // アイテムの説明
	curMemberIdx     int               // 選択中の味方
}

// State interface ================

func (st *EquipMenuState) OnPause(world w.World) {}

func (st *EquipMenuState) OnResume(world w.World) {}

func (st *EquipMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.equipMenu = append(st.equipMenu, loader.AddEntities(world, prefabs.Menu.EquipMenu)...)
	st.ui = st.initUI(world)
}

func (st *EquipMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.equipMenu...)
}

func (st *EquipMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	st.ui.Update()

	return updateMenu(st, world)
}

func (st *EquipMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// Menu Interface ================

func (st *EquipMenuState) getSelection() int {
	return st.selection
}

func (st *EquipMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *EquipMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *EquipMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *EquipMenuState) getCursorMenuIDs() []string {
	return []string{""}
}

// ================

func (st *EquipMenuState) initUI(world w.World) *ebitenui.UI {
	st.actionContainer = st.newItemSpecContainer(world)
	st.specContainer = st.newItemSpecContainer(world)
	st.generateActionContainer(world)

	st.subMenuContainer = eui.NewRowContainer()
	st.toggleSubMenu(world, false, func() {})

	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空白だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("装備", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(st.subMenuContainer)

		sc, v := eui.NewScrollContainer(st.actionContainer)
		rootContainer.AddChild(sc)
		rootContainer.AddChild(v)
		rootContainer.AddChild(st.specContainer)

		rootContainer.AddChild(st.itemDesc)
	}

	return &ebitenui.UI{Container: rootContainer}
}

// アイテムの性能表示コンテナ
func (st *EquipMenuState) newItemSpecContainer(world w.World) *widget.Container {
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

// スロットコンテナを生成する
func (st *EquipMenuState) generateActionContainer(world w.World) {
	st.actionContainer.RemoveChildren()

	members := []ecs.Entity{}
	simple.InPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})
	member := members[st.curMemberIdx]

	st.slots = equips.GetEquipments(world, member)

	gameComponents := world.Components.Game.(*gc.Components)
	for i, v := range st.slots {
		windowContainer := eui.NewWindowContainer()
		titleContainer := eui.NewWindowHeaderContainer("アクション", world)
		actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

		v := v
		i := i
		var name = ""
		var desc = " "
		if v != nil {
			name = gameComponents.Name.Get(*v).(*gc.Name).Name
			desc = gameComponents.Description.Get(*v).(*gc.Description).Description
		}

		slotButton := eui.NewItemButton(fmt.Sprintf("[ %s ]", name), func(args *widget.ButtonClickedEventArgs) {
			actionWindow.SetLocation(setWinRect())
			st.ui.AddWindow(actionWindow)
		}, world)
		slotButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.itemDesc.Label = desc
			if v != nil {
				views.UpdateSpec(world, st.specContainer, []any{
					simple.GetWeapon(world, *v),
					simple.GetWearable(world, *v),
				})
			} else {
				st.specContainer.RemoveChildren()
			}
		})
		equipButton := eui.NewItemButton("装備する", func(args *widget.ButtonClickedEventArgs) {
			st.items = st.queryMenuWeapon(world)
			f := func() { st.generateActionContainerEquip(world, member, gc.EquipmentSlotNumber(i), v) }
			f()
			st.toggleSubMenu(world, true, f)
			actionWindow.Close()
		}, world)
		windowContainer.AddChild(equipButton)

		if v != nil {
			disarmButton := eui.NewItemButton("外す", func(args *widget.ButtonClickedEventArgs) {
				equips.Disarm(world, *v)
				st.generateActionContainer(world)
				actionWindow.Close()
			}, world)
			windowContainer.AddChild(disarmButton)
		}

		closeButton := eui.NewItemButton("閉じる", func(args *widget.ButtonClickedEventArgs) {
			actionWindow.Close()
		}, world)
		windowContainer.AddChild(closeButton)

		st.actionContainer.AddChild(slotButton)
	}
}

// インベントリにある装備選択を生成する
func (st *EquipMenuState) generateActionContainerEquip(world w.World, member ecs.Entity, targetSlot gc.EquipmentSlotNumber, previousEquipment *ecs.Entity) {
	st.actionContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	for _, entity := range st.items {
		entity := entity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		itemButton := eui.NewItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			if previousEquipment != nil {
				equips.Disarm(world, *previousEquipment)
			}
			equips.Equip(world, entity, member, targetSlot)

			// 画面を戻す
			st.generateActionContainer(world)
			st.toggleSubMenu(world, false, func() {})
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.itemDesc.Label = simple.GetDescription(world, entity).Description
			views.UpdateSpec(world, st.specContainer, []any{
				simple.GetWeapon(world, entity),
				simple.GetWearable(world, entity),
				simple.GetMaterial(world, entity),
			})
		})
		st.actionContainer.AddChild(itemButton)
	}
}

// サブメニューコンテナの表示を切り替える
func (st *EquipMenuState) toggleSubMenu(world w.World, isInventory bool, reloadFunc func()) {
	st.subMenuContainer.RemoveChildren()

	if isInventory {
		toggleWeaponButton := eui.NewItemButton("武器", func(args *widget.ButtonClickedEventArgs) { st.items = st.queryMenuWeapon(world); reloadFunc() }, world)
		toggleWearableButton := eui.NewItemButton("防具", func(args *widget.ButtonClickedEventArgs) { st.items = st.queryMenuWearable(world); reloadFunc() }, world)
		st.subMenuContainer.AddChild(toggleWeaponButton)
		st.subMenuContainer.AddChild(toggleWearableButton)
	} else {
		members := []ecs.Entity{}
		simple.InPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})

		getName := func() string { return simple.GetName(world, members[st.curMemberIdx]).Name }
		memberName := eui.NewMenuText(getName(), world)

		prevMemberButton := eui.NewItemButton("前", func(args *widget.ButtonClickedEventArgs) {
			// Goでは負の剰余は負のままになる
			result := st.curMemberIdx - 1
			if result < 0 {
				result += len(members)
			}
			st.curMemberIdx = result
			st.generateActionContainer(world)

			memberName.Label = getName()
		}, world)
		nextMemberButton := eui.NewItemButton("次", func(args *widget.ButtonClickedEventArgs) {
			st.curMemberIdx = (st.curMemberIdx + 1) % len(members)
			st.generateActionContainer(world)

			memberName.Label = getName()
		}, world)
		st.subMenuContainer.AddChild(prevMemberButton)
		st.subMenuContainer.AddChild(memberName)
		st.subMenuContainer.AddChild(nextMemberButton)
	}
}

// 装備可能な武器を取得する
func (st *EquipMenuState) queryMenuWeapon(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.InBackpack,
		gameComponents.Weapon,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// 装備可能な防具を取得する
func (st *EquipMenuState) queryMenuWearable(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.InBackpack,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}
