package states

import (
	"fmt"
	"image"
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
	"github.com/kijimaD/ruins/lib/worldhelper/items"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type CraftMenuState struct {
	selection int
	craftMenu []ecs.Entity
	ui        *ebitenui.UI

	selectedItem       ecs.Entity        // 選択中のアイテム
	selectedItemButton *widget.Button    // 使用済みのアイテムのボタン
	items              []ecs.Entity      // 表示対象とするアイテム
	itemDesc           *widget.Text      // アイテムの概要
	itemList           *widget.Container // アイテムリストのコンテナ
	resultWindow       *widget.Window    // 合成結果ウィンドウ
	recipeList         *widget.Container // レシピリストのコンテナ
	winRect            image.Rectangle   // ウィンドウの開く位置
	specContainer      *widget.Container // 性能表示のコンテナ
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

func (st *CraftMenuState) initUI(world w.World) *ebitenui.UI {
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

	st.recipeList = st.newItemSpecContainer(world)
	st.specContainer = st.newItemSpecContainer(world)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("合成", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(toggleContainer)

		sc, v := eui.NewScrollContainer(st.itemList)
		rootContainer.AddChild(sc)
		rootContainer.AddChild(v)
		rootContainer.AddChild(st.newVSplitContainer(st.specContainer, st.recipeList, world))

		rootContainer.AddChild(itemDescContainer)
	}

	return &ebitenui.UI{Container: rootContainer}
}

// 新しいクエリを実行してitemsをセットする
func (st *CraftMenuState) queryMenuConsumable(world w.World) {
	st.itemList.RemoveChildren()
	st.items = []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Recipe,
		gameComponents.Consumable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.items = append(st.items, entity)
	}))
	st.generateList(world)
}

// 新しいクエリを実行してitemsをセットする
func (st *CraftMenuState) queryMenuWeapon(world w.World) {
	st.itemList.RemoveChildren()
	st.items = []ecs.Entity{}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.Recipe,
		gameComponents.Weapon,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		st.items = append(st.items, entity)
	}))
	st.generateList(world)
}

// itemsからUIを生成する
func (st *CraftMenuState) generateList(world world.World) {
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
			st.winRect = image.Rect(0, 0, x, y)
			st.winRect = st.winRect.Add(image.Point{x + 20, y + 20})
			actionWindow.SetLocation(st.winRect)
			st.ui.AddWindow(actionWindow)

			st.selectedItem = entity
			st.initWindowContainer(world, name.Name, windowContainer, actionWindow)
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.selectedItem != entity {
				st.selectedItem = entity
			}

			var description string
			world.Manager.Join(gameComponents.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
				if entity == st.selectedItem && entity.HasComponent(gameComponents.Description) {
					desc := gameComponents.Description.Get(entity).(*gc.Description)
					description = desc.Description
				}
			}))
			st.itemDesc.Label = description

			views.UpdateSpec(world, st.specContainer, []any{items.GetWeapon(world, entity)})

			st.updateRecipeList(world)
		})
		st.itemList.AddChild(itemButton)
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

func (st *CraftMenuState) specText(world w.World) *widget.Text {
	return widget.NewText(
		widget.TextOpts.Text("", eui.LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{}),
		),
	)
}

// 縦分割コンテナ
func (st *CraftMenuState) newVSplitContainer(top *widget.Container, bottom *widget.Container, world w.World) *widget.Container {
	split := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.DebugColor)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Spacing(2, 0),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true, true}),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:    2,
					Bottom: 2,
					Left:   2,
					Right:  2,
				}),
			)),
	)
	split.AddChild(top)
	split.AddChild(bottom)

	return split
}

func (st *CraftMenuState) initResultWindow(world w.World, entity ecs.Entity) {
	resultContainer := eui.NewWindowContainer()
	st.resultWindow = eui.NewSmallWindow(eui.NewWindowHeaderContainer("合成結果", world), resultContainer)

	views.UpdateSpec(world, resultContainer, []any{items.GetWeapon(world, entity)})

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
		if entity == st.selectedItem {
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
		st.resultWindow.SetLocation(st.winRect)
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
