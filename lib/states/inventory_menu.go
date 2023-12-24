package states

import (
	"image"
	"image/color"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/effects"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	er "github.com/kijimaD/sokotwo/lib/engine/resources"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
	"golang.org/x/image/font"
)

type InventoryMenuState struct {
	selection     int
	inventoryMenu []ecs.Entity
	menuLen       int
	ui            *ebitenui.UI
	clickedItem   ecs.Entity
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
	return []string{}
}

// ================

func (st *InventoryMenuState) initUI(world w.World) *ebitenui.UI {
	ui := ebitenui.UI{}
	buttonImage, _ := loadButtonImage()
	face, _ := loadFont((*world.Resources.Fonts)["kappa"])
	titleFace := face

	gameComponents := world.Components.Game.(*gc.Components)
	var members []ecs.Entity
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		members = append(members, entity)
	}))

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

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Left:  25,
				Right: 25,
			}),
		)),
	)

	title := widget.NewText(
		widget.TextOpts.Text("インベントリ", face, color.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	rootContainer.AddChild(title)

	for _, itemEntity := range items {
		entity := itemEntity
		name := gameComponents.Name.Get(entity).(*gc.Name)

		windowContainer := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
			widget.ContainerOpts.Layout(
				widget.NewGridLayout(
					widget.GridLayoutOpts.Columns(1),
					widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, false, false}),
					widget.GridLayoutOpts.Padding(widget.Insets{
						Top:    20,
						Bottom: 20,
						Left:   10,
						Right:  10,
					}),
					widget.GridLayoutOpts.Spacing(0, 15),
				),
			),
		)

		windowContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),
			widget.ButtonOpts.Text("使う", face, &widget.ButtonTextColor{
				Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			}),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				effects.ItemTrigger(nil, entity, effects.Single{members[0]}, world)
				st.ui = st.initUI(world)
			}),
		))

		windowContainer.AddChild(widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),
			widget.ButtonOpts.Text("捨てる", face, &widget.ButtonTextColor{
				Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			}),
		))

		titleContainer := widget.NewContainer(
			widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{150, 150, 150, 255})),
			widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
		)
		titleContainer.AddChild(widget.NewText(
			widget.TextOpts.Text("アクション", titleFace, color.NRGBA{254, 255, 255, 255}),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			})),
		))

		window := widget.NewWindow(
			widget.WindowOpts.Contents(windowContainer),
			widget.WindowOpts.TitleBar(titleContainer, 25),
			widget.WindowOpts.Modal(),
			widget.WindowOpts.CloseMode(widget.CLICK_OUT),
			widget.WindowOpts.Draggable(),
			widget.WindowOpts.Resizeable(),
			widget.WindowOpts.MinSize(200, 200),
			widget.WindowOpts.MaxSize(300, 400),
		)

		button := widget.NewButton(
			widget.ButtonOpts.Image(buttonImage),
			widget.ButtonOpts.Text(name.Name, face, &widget.ButtonTextColor{
				Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			}),
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:   30,
				Right:  30,
				Top:    5,
				Bottom: 5,
			}),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				x, y := window.Contents.PreferredSize()
				r := image.Rect(0, 0, x, y)
				r = r.Add(image.Point{100, 50})
				window.SetLocation(r)
				ui.AddWindow(window)

				st.clickedItem = entity
			}),
		)
		button.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			if st.clickedItem != entity {
				st.clickedItem = entity
				st.ui = st.initUI(world)
			}
		})

		rootContainer.AddChild(button)
	}

	var description string
	world.Manager.Join(gameComponents.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
		switch {
		case entity.HasComponent(gameComponents.Description):
			if entity == st.clickedItem {
				c := gameComponents.Description.Get(entity).(*gc.Description)
				description = c.Description
			}
		}
	}))

	itemDesc := widget.NewText(
		widget.TextOpts.Text(description, face, color.White),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	rootContainer.AddChild(itemDesc)

	ui.Container = rootContainer

	return &ui
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := e_image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := e_image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := e_image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(font er.Font) (font.Face, error) {
	return truetype.NewFace(font.Font, &truetype.Options{
		Size: 24,
		DPI:  72,
	}), nil
}
