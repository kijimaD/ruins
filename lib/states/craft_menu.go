package states

import (
	"fmt"
	"image/color"
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
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/craft"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type CraftMenuState struct {
	ui *ebitenui.UI

	hoveredItem        ecs.Entity        // ホバー中のアイテム
	selectedItemButton *widget.Button    // 使用済みのアイテムのボタン
	items              []ecs.Entity      // 表示対象とするアイテム
	itemDesc           *widget.Text      // アイテムの概要
	actionContainer    *widget.Container // アクションの起点となるコンテナ
	specContainer      *widget.Container // 性能表示のコンテナ
	resultWindow       *widget.Window    // 合成結果ウィンドウ
	recipeList         *widget.Container // レシピリストのコンテナ
	category           ItemCategoryType
}

func (st CraftMenuState) String() string {
	return "CraftMenu"
}

// State interface ================

func (st *CraftMenuState) OnPause(world w.World) {}

func (st *CraftMenuState) OnResume(world w.World) {}

func (st *CraftMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *CraftMenuState) OnStop(world w.World) {}

func (st *CraftMenuState) Update(world w.World) states.Transition {
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

func (st *CraftMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

var _ haveCategory = &CraftMenuState{}

func (st *CraftMenuState) setCategoryReload(world w.World, category ItemCategoryType) {
	st.category = category
	st.categoryReload(world)
}

func (st *CraftMenuState) categoryReload(world w.World) {
	st.actionContainer.RemoveChildren()
	st.items = []ecs.Entity{}

	switch st.category {
	case ItemCategoryTypeItem:
		st.items = st.queryMenuConsumable(world)
	case ItemCategoryTypeCard:
		st.items = st.queryMenuCard(world)
	case ItemCategoryTypeWearable:
		st.items = st.queryMenuWearable(world)
	default:
		log.Fatal("未定義のcategory")
	}

	st.generateActionContainer(world)
}

// TODO: あとで整理する
func (st *CraftMenuState) SetCategory(category ItemCategoryType) {
	st.category = category
}

// ================

func (st *CraftMenuState) initUI(world w.World) *ebitenui.UI {
	// 各アイテムが入るコンテナ
	st.actionContainer = eui.NewRowContainer()
	st.categoryReload(world)

	// アイテムの説明文
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空白だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.queryMenuConsumable(world)
	toggleContainer := eui.NewRowContainer()
	toggleConsumableButton := eui.NewItemButton("道具",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeItem)
		}),
	)
	toggleWearableButton := eui.NewItemButton("装備",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeWearable)
		}),
	)
	toggleCardButton := eui.NewItemButton("手札",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.setCategoryReload(world, ItemCategoryTypeCard)
		}),
	)
	toggleContainer.AddChild(toggleConsumableButton)
	toggleContainer.AddChild(toggleWearableButton)
	toggleContainer.AddChild(toggleCardButton)

	st.recipeList = eui.NewVerticalContainer()
	st.specContainer = eui.NewVerticalContainer()

	res := world.Resources.UIResources
	rootContainer := eui.NewItemGridContainer(
		widget.ContainerOpts.BackgroundImage(res.Panel.ImageTrans),
	)
	{
		rootContainer.AddChild(eui.NewMenuText("合成", world))
		rootContainer.AddChild(widget.NewContainer())
		rootContainer.AddChild(toggleContainer)

		rootContainer.AddChild(st.actionContainer)
		rootContainer.AddChild(widget.NewContainer())
		rootContainer.AddChild(eui.NewVSplitContainer(st.specContainer, st.recipeList))

		rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: rootContainer}
}

func (st *CraftMenuState) queryMenuConsumable(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Recipe,
		gameComponents.Consumable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity.HasComponent(gameComponents.Card) {
			return
		}

		items = append(items, entity)
	}))

	return items
}

func (st *CraftMenuState) queryMenuCard(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Recipe,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

func (st *CraftMenuState) queryMenuWearable(world w.World) []ecs.Entity {
	items := []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Recipe,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// itemsからactionContainerを生成する
func (st *CraftMenuState) generateActionContainer(world world.World) {
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
			name := gameComponents.Name.Get(entity).(*gc.Name)

			return string(name.Name)
		}),
		euiext.ListOpts.EntryEnterFunc(func(e any) {
			entity, ok := e.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			if st.hoveredItem != entity {
				st.hoveredItem = entity
			}
			desc := gameComponents.Description.Get(entity).(*gc.Description)
			st.itemDesc.Label = desc.Description
			views.UpdateSpec(world, st.specContainer, entity)
			st.updateRecipeList(world)
		}),
		euiext.ListOpts.EntryButtonOpts(),
		euiext.ListOpts.EntrySelectedHandler(func(args *euiext.ListEntrySelectedEventArgs) {
			entity, ok := args.Entry.(ecs.Entity)
			if !ok {
				log.Fatal("unexpected entry detect!")
			}
			name := gameComponents.Name.Get(entity).(*gc.Name)
			windowContainer := eui.NewWindowContainer(world)
			actionWindow := eui.NewSmallWindow(
				eui.NewWindowHeaderContainer("アクション", world),
				windowContainer,
			)

			actionWindow.SetLocation(setWinRect())
			st.initWindowContainer(world, name.Name, windowContainer, actionWindow)
			st.ui.AddWindow(actionWindow)
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
}

func (st *CraftMenuState) initResultWindow(world w.World, entity ecs.Entity) {
	resultContainer := eui.NewWindowContainer(world)
	st.resultWindow = eui.NewSmallWindow(eui.NewWindowHeaderContainer("合成結果", world), resultContainer)

	views.UpdateSpec(world, resultContainer, entity)

	closeButton := eui.NewItemButton("閉じる",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			st.resultWindow.Close()
		}),
	)
	resultContainer.AddChild(closeButton)
}

func (st *CraftMenuState) updateRecipeList(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	st.recipeList.RemoveChildren()
	world.Manager.Join(
		gameComponents.Recipe,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if entity == st.hoveredItem {
			recipe := gameComponents.Recipe.Get(entity).(*gc.Recipe)
			for _, input := range recipe.Inputs {
				str := fmt.Sprintf("%s %d pcs\n    所持: %d pcs", input.Name, input.Amount, material.GetAmount(input.Name, world))
				var color color.RGBA
				if material.GetAmount(input.Name, world) >= input.Amount {
					color = styles.SuccessColor
				} else {
					color = styles.DangerColor
				}

				st.recipeList.AddChild(eui.NewBodyText(str, color, world))
			}
		}
	}))
}

// アクションウィンドウはクリックのたびに毎回中身を作り直す
// useButton.GetWidget().Disabled = true を使ってボタンを非活性にする方が楽でよさそうなのだが、非活性にすると描画の色まわりでヌルポになる。色は設定しているのに...
func (st *CraftMenuState) initWindowContainer(world w.World, name string, windowContainer *widget.Container, actionWindow *widget.Window) {
	windowContainer.RemoveChildren()
	useButton := eui.NewItemButton("合成する",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			resultEntity, err := craft.Craft(world, name)
			if err != nil {
				log.Fatal(err)
			}
			st.updateRecipeList(world)

			actionWindow.Close()
			st.initResultWindow(world, *resultEntity)
			st.resultWindow.SetLocation(getWinRect())
			st.ui.AddWindow(st.resultWindow)
		}),
	)
	if craft.CanCraft(world, name) {
		windowContainer.AddChild(useButton)
	}

	closeButton := eui.NewItemButton("閉じる",
		world,
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			actionWindow.Close()
		}),
	)
	windowContainer.AddChild(closeButton)
}
