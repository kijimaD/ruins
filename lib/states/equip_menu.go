package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/utils/consts"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type EquipMenuState struct {
	ui *ebitenui.UI

	slots            []*ecs.Entity     // スロット一覧
	items            []ecs.Entity      // インベントリにあるアイテム一覧
	subMenuContainer *widget.Container // 右上のサブメニュー
	actionContainer  *widget.Container // 操作の起点となるメインメニュー
	specContainer    *widget.Container // 性能コンテナ
	abilityContainer *widget.Container // メンバーの能力表示コンテナ
	itemDesc         *widget.Text      // アイテムの説明
	curMemberIdx     int               // 選択中の味方
}

func (st EquipMenuState) String() string {
	return "EquipMenu"
}

// State interface ================

func (st *EquipMenuState) OnPause(world w.World) {}

func (st *EquipMenuState) OnResume(world w.World) {}

func (st *EquipMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *EquipMenuState) OnStop(world w.World) {}

func (st *EquipMenuState) Update(world w.World) states.Transition {
	changed := gs.EquipmentChangedSystem(world)
	if changed {
		st.reloadAbilityContainer(world)
	}
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

func (st *EquipMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *EquipMenuState) initUI(world w.World) *ebitenui.UI {
	st.actionContainer = eui.NewVerticalContainer()
	st.generateActionContainer(world)
	st.specContainer = eui.NewVerticalContainer()
	st.abilityContainer = eui.NewVerticalContainer()
	st.reloadAbilityContainer(world)

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
		rootContainer.AddChild(eui.NewWSplitContainer(st.specContainer, st.abilityContainer))

		rootContainer.AddChild(st.itemDesc)
	}

	return &ebitenui.UI{Container: rootContainer}
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
				// 該当スロットに装備がある場合はその性能を表示する
				views.UpdateSpec(world, st.specContainer, *v)
			} else {
				// 非表示
				st.specContainer.RemoveChildren()
			}
		})
		equipButton := eui.NewItemButton("装備する", func(args *widget.ButtonClickedEventArgs) {
			st.items = st.queryMenuWearable(world)
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
			views.UpdateSpec(world, st.specContainer, entity)
		})
		st.actionContainer.AddChild(itemButton)
	}
}

// サブメニューコンテナの表示を切り替える
func (st *EquipMenuState) toggleSubMenu(world w.World, isInventory bool, reloadFunc func()) {
	st.subMenuContainer.RemoveChildren()

	if !isInventory {
		members := []ecs.Entity{}
		simple.InPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})

		prevMemberButton := eui.NewItemButton("前", func(args *widget.ButtonClickedEventArgs) {
			// Goでは負の剰余は負のままになる
			result := st.curMemberIdx - 1
			if result < 0 {
				result += len(members)
			}
			st.curMemberIdx = result
			st.generateActionContainer(world)

			st.reloadAbilityContainer(world)
		}, world)
		nextMemberButton := eui.NewItemButton("次", func(args *widget.ButtonClickedEventArgs) {
			st.curMemberIdx = (st.curMemberIdx + 1) % len(members)
			st.generateActionContainer(world)

			st.reloadAbilityContainer(world)
		}, world)
		st.subMenuContainer.AddChild(prevMemberButton)
		st.subMenuContainer.AddChild(nextMemberButton)
	}
}

// メンバーの能力表示コンテナを更新する
func (st *EquipMenuState) reloadAbilityContainer(world w.World) {
	st.abilityContainer.RemoveChildren()

	members := []ecs.Entity{}
	simple.InPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})
	targetMember := members[st.curMemberIdx]

	gameComponents := world.Components.Game.(*gc.Components)
	for _, entity := range st.queryAbility(world) {
		if entity != targetMember {
			continue
		}
		views.AddMemberBar(world, st.abilityContainer, entity)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %3d/%3d", consts.SPLabel, pools.SP.Current, pools.SP.Max), styles.TextColor, world))

		attrs := gameComponents.Attributes.Get(entity).(*gc.Attributes)
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.VitalityLabel, attrs.Vitality.Total, attrs.Vitality.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.StrengthLabel, attrs.Strength.Total, attrs.Strength.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.SensationLabel, attrs.Sensation.Total, attrs.Sensation.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.DexterityLabel, attrs.Dexterity.Total, attrs.Dexterity.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.AgilityLabel, attrs.Agility.Total, attrs.Agility.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", consts.DefenseLabel, attrs.Defense.Total, attrs.Defense.Modifier), styles.TextColor, world))
	}
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

func (st *EquipMenuState) queryAbility(world w.World) []ecs.Entity {
	entities := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entities = append(entities, entity)
	}))

	return entities
}
