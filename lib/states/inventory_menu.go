package states

import (
	"image"

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
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	selection     int
	inventoryMenu []ecs.Entity
	menuLen       int
	ui            *ebitenui.UI
	clickedItem   ecs.Entity
}

var selectedItem ecs.Entity
var selectedItemButton *widget.Button // 選択中のアイテム

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

type entryStruct struct {
	entity      ecs.Entity
	name        string
	description string
}

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	ui := ebitenui.UI{}
	gameComponents := world.Components.Game.(*gc.Components)

	var members []ecs.Entity
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		members = append(members, entity)
	}))

	// TODO: わかりにくいのでグリッドをやめてrowにしたい
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// アイテム, スクロール, アイテムspecで3列になっている
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Spacing(2, 0),
			widget.GridLayoutOpts.Stretch([]bool{true, false, true}, []bool{false, true, false}),
		)),
	)

	rootContainer.AddChild(eui.NewMenuTitle("インベントリ", world))
	rootContainer.AddChild(eui.EmptyContainer())
	rootContainer.AddChild(eui.EmptyContainer())

	// 各アイテムが入るコンテナ
	itemList := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(2),
	)))

	// アイテムの説明文コンテナ
	itemDescContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(0, 40),
		),
	)

	// アイテムの説明文
	itemDesc := widget.NewText(
		widget.TextOpts.Text(" ", eui.LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
	)
	itemDescContainer.AddChild(itemDesc)

	var items []ecs.Entity
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Name,
		gameComponents.Description,
		gameComponents.InBackpack,
		gameComponents.Consumable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	newWindowContainer := func() *widget.Container {
		return widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.WindowBodyColor)),
			widget.ContainerOpts.Layout(
				// TODO: gridである必要はなさそう。RowContainerを使えばよい?
				widget.NewGridLayout(
					widget.GridLayoutOpts.Columns(1),
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false}),
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    10,
						Bottom: 10,
						Left:   10,
						Right:  10,
					}),
					widget.GridLayoutOpts.Spacing(0, 2),
				),
			),
		)
	}

	newTitleContainer := func(title string) *widget.Container {
		container := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.WindowHeaderColor)),
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		container.AddChild(widget.NewText(
			widget.TextOpts.Text(title, eui.LoadFont(world), styles.TextColor),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
		))
		return container
	}

	newWindow := func(title *widget.Container, content *widget.Container) *widget.Window {
		return widget.NewWindow(
			widget.WindowOpts.Contents(content),
			widget.WindowOpts.TitleBar(title, 25),
			widget.WindowOpts.Modal(),
			widget.WindowOpts.CloseMode(widget.CLICK_OUT),
			widget.WindowOpts.Draggable(),
			widget.WindowOpts.Resizeable(),
			widget.WindowOpts.MinSize(200, 200),
			widget.WindowOpts.MaxSize(300, 400),
		)
	}

	newItemButton := func(text string, f func(args *widget.ButtonClickedEventArgs)) *widget.Button {
		return widget.NewButton(
			widget.ButtonOpts.Image(eui.LoadButtonImage()),
			widget.ButtonOpts.Text(text, eui.LoadFont(world), &widget.ButtonTextColor{
				Idle: styles.TextColor,
			}),
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			widget.ButtonOpts.ClickedHandler(f),
		)
	}

	partyContainer := newWindowContainer()
	partyWindow := newWindow(newTitleContainer("選択"), partyContainer)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		partyButton := newItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			effects.ItemTrigger(nil, selectedItem, effects.Single{entity}, world)
			partyWindow.Close()
			itemList.RemoveChild(selectedItemButton)
		})
		partyContainer.AddChild(partyButton)
	}))

	for _, entity := range items {
		entity := entity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		windowContainer := newWindowContainer()
		titleContainer := newTitleContainer("アクション")

		actionWindow := newWindow(titleContainer, windowContainer)

		// アイテムの名前がラベルについたボタン
		itemButton := newItemButton(name.Name, func(args *widget.ButtonClickedEventArgs) {
			x, y := ebiten.CursorPosition()
			r := image.Rect(0, 0, x, y)
			r = r.Add(image.Point{x + 20, y + 20})
			actionWindow.SetLocation(r)
			ui.AddWindow(actionWindow)

			st.clickedItem = entity
		})

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.clickedItem != entity {
				st.clickedItem = entity
			}

			var description string
			world.Manager.Join(gameComponents.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
				if entity == st.clickedItem && entity.HasComponent(gameComponents.Description) {
					c := gameComponents.Description.Get(entity).(*gc.Description)
					description = c.Description
				}
			}))
			itemDesc.Label = description
		})
		itemList.AddChild(itemButton)

		useButton := newItemButton("使う", func(args *widget.ButtonClickedEventArgs) {
			x, y := ebiten.CursorPosition()
			r := image.Rect(0, 0, x, y)
			r = r.Add(image.Point{x + 20, y + 20})
			partyWindow.SetLocation(r)

			consumable := gameComponents.Consumable.Get(entity).(*gc.Consumable)
			switch consumable.TargetType.TargetNum {
			case gc.TargetSingle:
				ui.AddWindow(partyWindow)
				actionWindow.Close()
				selectedItem = entity
				selectedItemButton = itemButton
			case gc.TargetAll:
				effects.ItemTrigger(nil, entity, effects.Party{}, world)
				actionWindow.Close()
				itemList.RemoveChild(itemButton)
			}
		})
		windowContainer.AddChild(useButton)

		dropButton := newItemButton("捨てる", func(args *widget.ButtonClickedEventArgs) {
			world.Manager.DeleteEntity(entity)
			itemList.RemoveChild(itemButton)
			actionWindow.Close()
		})
		windowContainer.AddChild(dropButton)

		closeButton := newItemButton("閉じる", func(args *widget.ButtonClickedEventArgs) {
			actionWindow.Close()
		})
		windowContainer.AddChild(closeButton)
	}

	sc, v := eui.NewScrollContainer(itemList)
	rootContainer.AddChild(sc)
	rootContainer.AddChild(v)

	itemSpec := widget.NewText(
		widget.TextOpts.Text("性能", eui.LoadFont(world), styles.TextColor),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	rootContainer.AddChild(itemSpec)
	rootContainer.AddChild(itemDescContainer)

	ui = ebitenui.UI{
		Container: rootContainer,
	}
	ui.Container = rootContainer

	return &ui
}
