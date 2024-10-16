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
	"github.com/kijimaD/ruins/lib/engine/world"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/euiext"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	ui *ebitenui.UI

	selectedItem    ecs.Entity        // 選択中のアイテム
	items           []ecs.Entity      // 表示対象とするアイテム
	itemDesc        *widget.Text      // アイテムの概要
	actionContainer *widget.Container // アクションの起点となるコンテナ
	specContainer   *widget.Container // 性能表示のコンテナ
	partyWindow     *widget.Window    // 仲間を選択するウィンドウ
	category        ItemCategoryType
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

func (st *InventoryMenuState) setCategory(world w.World, category ItemCategoryType) {
	st.category = category
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

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	res := world.Resources.UIResources

	// 各アクションが入るコンテナ
	st.actionContainer = eui.NewVerticalContainer()
	st.categoryReload(world)

	// 種類トグル
	toggleContainer := st.newToggleContainer(world)

	// アイテムの説明文
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空文字だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.specContainer = eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)

	rootContainer := eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		rootContainer.AddChild(eui.NewMenuText("インベントリ", world))
		rootContainer.AddChild(widget.NewContainer())
		rootContainer.AddChild(toggleContainer)

		rootContainer.AddChild(st.actionContainer)
		rootContainer.AddChild(widget.NewContainer())
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

	entities := []any{}
	for _, entity := range st.items {
		entities = append(entities, entity)
	}

	res := world.Resources.UIResources
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
			if st.selectedItem != entity {
				st.selectedItem = entity
			}
			desc := gameComponents.Description.Get(entity).(*gc.Description)
			st.itemDesc.Label = desc.Description
			views.UpdateSpec(world, st.specContainer, entity)
		}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			windowContainer := eui.NewWindowContainer(world)
			titleContainer := eui.NewWindowHeaderContainer("アクション", world)
			actionWindow := eui.NewSmallWindow(titleContainer, windowContainer)

			entity, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatalf("unexpected entry: %#v", entity)
			}

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
						st.generateList(world)
					case gc.TargetAll:
						effects.ItemTrigger(nil, entity, effects.Party{}, world)
						actionWindow.Close()
						st.categoryReload(world)
					}
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
					st.categoryReload(world)
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
		}),
		euiext.ListOpts.EntryTextPadding(widget.NewInsetsSimple(10)),
		euiext.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(440, 520),
		)),
		euiext.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(res.List.ImageTrans)),
	}
	list := eui.NewList(
		entities,
		opts,
		world,
	)
	st.actionContainer.AddChild(list)

	count := fmt.Sprintf("合計 %02d個", len(st.items))
	st.actionContainer.AddChild(eui.NewMenuText(count, world))
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
				st.categoryReload(world)
			}),
		)
		memberContainer.AddChild(partyButton)
		views.AddMemberBar(world, memberContainer, entity)

		rowContainer.AddChild(memberContainer)
	}))
}

func (st *InventoryMenuState) newToggleContainer(world w.World) *widget.Container {
	toggleContainer := eui.NewRowContainer()
	toggleConsumableButton := eui.NewButton("道具",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeItem)
		}),
	)
	toggleCardButton := eui.NewButton("手札",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeCard)
		}),
	)
	toggleWearableButton := eui.NewButton("防具",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeWearable)
		}),
	)
	toggleMaterialButton := eui.NewButton("素材",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeMaterial)
		}),
	)
	toggleContainer.AddChild(toggleConsumableButton)
	toggleContainer.AddChild(toggleCardButton)
	toggleContainer.AddChild(toggleWearableButton)
	toggleContainer.AddChild(toggleMaterialButton)

	return toggleContainer
}
