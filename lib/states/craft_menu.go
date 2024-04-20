package states

import (
	"fmt"
	"image/color"
	"log"

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
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/craft"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type CraftMenuState struct {
	selection int
	craftMenu []ecs.Entity
	ui        *ebitenui.UI

	hoveredItem        ecs.Entity        // ホバー中のアイテム
	selectedItemButton *widget.Button    // 使用済みのアイテムのボタン
	items              []ecs.Entity      // 表示対象とするアイテム
	itemDesc           *widget.Text      // アイテムの概要
	actionContainer    *widget.Container // アクションの起点となるコンテナ
	specContainer      *widget.Container // 性能表示のコンテナ
	resultWindow       *widget.Window    // 合成結果ウィンドウ
	recipeList         *widget.Container // レシピリストのコンテナ
	category           itemCategoryType
}

// State interface ================

func (st *CraftMenuState) OnPause(world w.World) {}

func (st *CraftMenuState) OnResume(world w.World) {}

func (st *CraftMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.craftMenu = append(st.craftMenu, loader.AddEntities(world, prefabs.Menu.CraftMenu)...)
	st.ui = st.initUI(world)
}

func (st *CraftMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.craftMenu...)
}

func (st *CraftMenuState) Update(world w.World) states.Transition {
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

func (st *CraftMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// Menu Interface ================

func (st *CraftMenuState) getSelection() int {
	return st.selection
}

func (st *CraftMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *CraftMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *CraftMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *CraftMenuState) getCursorMenuIDs() []string {
	return []string{""}
}

// ================
var _ haveCategory = &CraftMenuState{}

func (st *CraftMenuState) setCategoryReload(world w.World, category itemCategoryType) {
	st.category = category
	st.categoryReload(world)
}

func (st *CraftMenuState) categoryReload(world w.World) {
	st.actionContainer.RemoveChildren()
	st.items = []ecs.Entity{}

	switch st.category {
	case itemCategoryTypeItem:
		st.items = st.queryMenuConsumable(world)
	case itemCategoryTypeCard:
		st.items = st.queryMenuCard(world)
	case itemCategoryTypeWearable:
		st.items = st.queryMenuWearable(world)
	default:
		log.Fatal("未定義のcategory")
	}

	st.generateActionContainer(world)
}

// ================

func (st *CraftMenuState) initUI(world w.World) *ebitenui.UI {
	// 各アイテムが入るコンテナ
	st.actionContainer = eui.NewScrollContentContainer()
	st.categoryReload(world)

	// アイテムの説明文
	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空白だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	st.queryMenuConsumable(world)
	toggleContainer := eui.NewRowContainer()
	toggleConsumableButton := eui.NewItemButton("道具", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, itemCategoryTypeItem) }, world)
	toggleWearableButton := eui.NewItemButton("装備", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, itemCategoryTypeWearable) }, world)
	toggleCardButton := eui.NewItemButton("手札", func(args *widget.ButtonClickedEventArgs) { st.setCategoryReload(world, itemCategoryTypeCard) }, world)
	toggleContainer.AddChild(toggleConsumableButton)
	toggleContainer.AddChild(toggleWearableButton)
	toggleContainer.AddChild(toggleCardButton)

	st.recipeList = st.newItemSpecContainer(world)
	st.specContainer = st.newItemSpecContainer(world)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("合成", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(toggleContainer)

		sc, v := eui.NewScrollContainer(st.actionContainer)
		rootContainer.AddChild(sc)
		rootContainer.AddChild(v)
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

	for _, entity := range st.items {
		entity := entity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		windowContainer := eui.NewWindowContainer()
		actionWindow := eui.NewSmallWindow(
			eui.NewWindowHeaderContainer("アクション", world),
			windowContainer,
		)

		itemButton := eui.NewItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			actionWindow.SetLocation(setWinRect())
			st.initWindowContainer(world, name.Name, windowContainer, actionWindow)
			st.ui.AddWindow(actionWindow)
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.hoveredItem != entity {
				st.hoveredItem = entity
			}
			st.itemDesc.Label = simple.GetDescription(world, entity).Description
			views.UpdateSpec(world, st.specContainer, entity)
			st.updateRecipeList(world)
		})
		st.actionContainer.AddChild(itemButton)
	}
}

func (st *CraftMenuState) newItemSpecContainer(world w.World) *widget.Container {
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

func (st *CraftMenuState) initResultWindow(world w.World, entity ecs.Entity) {
	resultContainer := eui.NewWindowContainer()
	st.resultWindow = eui.NewSmallWindow(eui.NewWindowHeaderContainer("合成結果", world), resultContainer)

	views.UpdateSpec(world, resultContainer, entity)

	closeButton := eui.NewItemButton("閉じる", func(args *widget.ButtonClickedEventArgs) {
		st.resultWindow.Close()
	}, world)
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
	useButton := eui.NewItemButton("合成する", func(args *widget.ButtonClickedEventArgs) {
		resultEntity, err := craft.Craft(world, name)
		if err != nil {
			log.Fatal(err)
		}
		st.updateRecipeList(world)

		actionWindow.Close()
		st.initResultWindow(world, *resultEntity)
		st.resultWindow.SetLocation(getWinRect())
		st.ui.AddWindow(st.resultWindow)
	}, world)
	if craft.CanCraft(world, name) {
		windowContainer.AddChild(useButton)
	}

	closeButton := eui.NewItemButton("閉じる", func(args *widget.ButtonClickedEventArgs) {
		actionWindow.Close()
	}, world)
	windowContainer.AddChild(closeButton)
}
