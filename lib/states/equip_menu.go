package states

import (
	"fmt"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type EquipMenuState struct {
	states.BaseState
	ui *ebitenui.UI

	// スロット一覧
	slots []*ecs.Entity
	// インベントリにあるアイテム一覧
	items []ecs.Entity
	// 右上のサブメニュー
	subMenuContainer *widget.Container
	// 操作の起点となるメインメニュー
	actionContainer *widget.Container
	// 性能コンテナ
	specContainer *widget.Container
	// メンバーの能力表示コンテナ
	abilityContainer *widget.Container
	// 装備対象の切り替えコンテナ
	equipTargetContainer *widget.Container
	// アイテムの説明
	itemDesc *widget.Text
	// 選択中の味方
	curMemberIdx int
	// 装備対象。防具もしくは手札
	equipTarget equipTarget
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

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
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
	st.equipTargetContainer = eui.NewRowContainer()
	st.reloadEquipTargetContainer(world, true)

	st.subMenuContainer = eui.NewRowContainer()
	st.reloadSubMenu(world, true, func() {})

	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空文字だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	res := world.Resources.UIResources
	rootContainer := eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		rootContainer.AddChild(st.equipTargetContainer)
		rootContainer.AddChild(widget.NewContainer())
		rootContainer.AddChild(st.subMenuContainer)

		rootContainer.AddChild(st.actionContainer)
		rootContainer.AddChild(widget.NewContainer())
		rootContainer.AddChild(eui.NewWSplitContainer(st.specContainer, st.abilityContainer))

		rootContainer.AddChild(st.itemDesc)
	}

	return &ebitenui.UI{Container: rootContainer}
}

type equipActionEntry struct {
	entity     *ecs.Entity
	slotNumber int
}

// アクションコンテナを生成する
func (st *EquipMenuState) generateActionContainer(world w.World) {
	st.actionContainer.RemoveChildren()

	members := []ecs.Entity{}
	worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})
	member := members[st.curMemberIdx]
	st.setSlots(world, member)

	gameComponents := world.Components.Game.(*gc.Components)
	slots := []any{}
	for i, v := range st.slots {
		entry := equipActionEntry{
			entity:     v,
			slotNumber: i,
		}
		slots = append(slots, entry)
	}

	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			v, ok := e.(equipActionEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			var name = "-"
			if v.entity != nil {
				name = gameComponents.Name.Get(*v.entity).(*gc.Name).Name
			}

			return string(name)
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {
			if e == nil {
				return
			}

			v, ok := e.(equipActionEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			var desc = " "
			if v.entity != nil {
				desc = gameComponents.Description.Get(*v.entity).(*gc.Description).Description
			}

			st.itemDesc.Label = desc
			if v.entity != nil {
				// 該当スロットに装備がある場合はその性能を表示する
				views.UpdateSpec(world, st.specContainer, *v.entity)
			} else {
				// 非表示
				st.specContainer.RemoveChildren()
			}
		}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			windowContainer := eui.NewWindowContainer(world)
			titleContainer := eui.NewWindowHeaderContainer("アクション", world)
			actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

			v, ok := args.Entry.(equipActionEntry)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			equipButton := eui.NewButton("装備する",
				world,
				widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
					st.setItems(world)
					f := func() {
						st.generateActionContainerEquip(world, member, gc.EquipmentSlotNumber(v.slotNumber), v.entity)
					}
					f()
					st.reloadSubMenu(world, false, f)
					st.reloadEquipTargetContainer(world, false)
					actionWindow.Close()
				}),
			)
			windowContainer.AddChild(equipButton)

			if v.entity != nil {
				disarmButton := eui.NewButton("外す",
					world,
					widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
						worldhelper.Disarm(world, *v.entity)
						st.generateActionContainer(world)
						actionWindow.Close()
					}),
				)
				windowContainer.AddChild(disarmButton)
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
		}),
		euiext.ListOpts.EntryTextPadding(widget.NewInsetsSimple(10)),
		euiext.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(440, 400),
		)),
	}
	list := eui.NewList(
		slots,
		opts,
		world,
	)

	st.actionContainer.AddChild(list)
}

// インベントリにある装備選択コンテナを生成する
func (st *EquipMenuState) generateActionContainerEquip(world w.World, member ecs.Entity, targetSlot gc.EquipmentSlotNumber, previousEquipment *ecs.Entity) {
	// 切り替えのために消す
	st.actionContainer.RemoveChildren()

	gameComponents := world.Components.Game.(*gc.Components)
	entities := []any{}
	for _, entity := range st.items {
		entities = append(entities, entity)
	}

	opts := []euiext.ListOpt{
		euiext.ListOpts.EntryLabelFunc(func(e any) string {
			entity, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := gameComponents.Name.Get(entity).(*gc.Name).Name

			return string(name)
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {
			entity, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			desc := gameComponents.Description.Get(entity).(*gc.Description).Description
			st.itemDesc.Label = desc
			views.UpdateSpec(world, st.specContainer, entity)
		}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			entity, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}

			// 装備
			if previousEquipment != nil {
				worldhelper.Disarm(world, *previousEquipment)
			}
			worldhelper.Equip(world, entity, member, targetSlot)

			// 画面を戻す
			st.generateActionContainer(world)
			st.reloadSubMenu(world, true, func() {})
			st.reloadEquipTargetContainer(world, true)
		}),
		euiext.ListOpts.EntryTextPadding(widget.NewInsetsSimple(10)),
		euiext.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(440, 400),
		)),
	}

	list := eui.NewList(
		entities,
		opts,
		world,
	)
	st.actionContainer.AddChild(list)
	st.actionContainer.AddChild(eui.NewButton(
		"戻る",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.generateActionContainer(world)
			st.reloadSubMenu(world, true, func() {})
			st.reloadEquipTargetContainer(world, true)
		}),
	))
}

func (st *EquipMenuState) reloadEquipTargetContainer(world w.World, visible bool) {
	st.equipTargetContainer.RemoveChildren()

	st.equipTargetContainer.AddChild(eui.NewMenuText("装備", world))
	if !visible {
		return
	}
	toggleTargetWearButton := eui.NewButton("防具",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.equipTarget = equipTargetWear
			st.generateActionContainer(world)
			st.reloadEquipTargetContainer(world, visible)
		}),
	)
	toggleTargetCardButton := eui.NewButton("手札",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.equipTarget = equipTargetCard
			st.generateActionContainer(world)
			st.reloadEquipTargetContainer(world, visible)
		}),
	)

	switch st.equipTarget {
	case equipTargetWear:
		toggleTargetWearButton.GetWidget().Disabled = true
	case equipTargetCard:
		toggleTargetCardButton.GetWidget().Disabled = true
	}
	st.equipTargetContainer.AddChild(toggleTargetWearButton)
	st.equipTargetContainer.AddChild(toggleTargetCardButton)
}

// サブメニューコンテナの表示を切り替える
func (st *EquipMenuState) reloadSubMenu(world w.World, visible bool, reloadFunc func()) {
	st.subMenuContainer.RemoveChildren()

	if visible {
		members := []ecs.Entity{}
		worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
			members = append(members, entity)
		})

		prevMemberButton := eui.NewButton("前",
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				// Goでは負の剰余は負のままになる
				result := st.curMemberIdx - 1
				if result < 0 {
					result += len(members)
				}
				st.curMemberIdx = result
				st.generateActionContainer(world)

				st.reloadAbilityContainer(world)
			}),
		)
		nextMemberButton := eui.NewButton("次",
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				st.curMemberIdx = (st.curMemberIdx + 1) % len(members)
				st.generateActionContainer(world)

				st.reloadAbilityContainer(world)
			}),
		)
		st.subMenuContainer.AddChild(prevMemberButton)
		st.subMenuContainer.AddChild(nextMemberButton)
	}
}

// メンバーの能力表示コンテナを更新する
func (st *EquipMenuState) reloadAbilityContainer(world w.World) {
	st.abilityContainer.RemoveChildren()

	members := []ecs.Entity{}
	worldhelper.QueryInPartyMember(world, func(entity ecs.Entity) {
		members = append(members, entity)
	})
	targetMember := members[st.curMemberIdx]

	gameComponents := world.Components.Game.(*gc.Components)
	for _, entity := range st.queryAbility(world) {
		if entity != targetMember {
			continue
		}
		views.AddMemberBar(world, st.abilityContainer, entity)

		attrs := gameComponents.Attributes.Get(entity).(*gc.Attributes)
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.VitalityLabel, attrs.Vitality.Total, attrs.Vitality.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.StrengthLabel, attrs.Strength.Total, attrs.Strength.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.SensationLabel, attrs.Sensation.Total, attrs.Sensation.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.DexterityLabel, attrs.Dexterity.Total, attrs.Dexterity.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.AgilityLabel, attrs.Agility.Total, attrs.Agility.Modifier), styles.TextColor, world))
		st.abilityContainer.AddChild(eui.NewBodyText(fmt.Sprintf("%s %2d(%+d)", utils.DefenseLabel, attrs.Defense.Total, attrs.Defense.Modifier), styles.TextColor, world))
	}
}

func (st *EquipMenuState) setSlots(world w.World, member ecs.Entity) {
	switch st.equipTarget {
	case equipTargetWear:
		st.slots = worldhelper.GetWearEquipments(world, member)
	case equipTargetCard:
		st.slots = worldhelper.GetCardEquipments(world, member)
	}
}

func (st *EquipMenuState) setItems(world w.World) {
	switch st.equipTarget {
	case equipTargetWear:
		st.items = st.queryMenuWear(world)
	case equipTargetCard:
		st.items = st.queryMenuCard(world)
	default:
		log.Fatal(fmt.Sprintf("invalid equipTarget type: %d", st.equipTarget))
	}
}

// 装備可能な防具を取得する
func (st *EquipMenuState) queryMenuWear(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationInBackpack,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// 装備可能な手札を取得する
func (st *EquipMenuState) queryMenuCard(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationInBackpack,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *EquipMenuState) queryAbility(world w.World) []ecs.Entity {
	entities := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entities = append(entities, entity)
	}))

	return entities
}
