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

	entries := make([]any, 0, 10)
	for _, itemEntity := range items {
		name := gameComponents.Name.Get(itemEntity).(*gc.Name)
		desc := gameComponents.Description.Get(itemEntity).(*gc.Description)

		ds := entryStruct{
			entity:      itemEntity,
			name:        name.Name,
			description: desc.Description,
		}
		entries = append(entries, ds)
	}

	list := widget.NewList(
		widget.ListOpts.ContainerOpts(widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(150, 300),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchVertical:    true,
			}),
		)),
		widget.ListOpts.Entries(entries),
		widget.ListOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle:     e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Disabled: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Mask:     e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(&widget.SliderTrackImage{
				Idle:  e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: e_image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}, buttonImage),
			widget.SliderOpts.MinHandleSize(5),
			widget.SliderOpts.TrackPadding(widget.NewInsetsSimple(2)),
		),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntryFontFace(face),
		widget.ListOpts.EntryColor(&widget.ListEntryColor{
			Selected:                   color.NRGBA{0, 255, 0, 255},
			Unselected:                 color.NRGBA{254, 255, 255, 255},
			SelectedBackground:         color.NRGBA{R: 130, G: 130, B: 200, A: 255},
			SelectedFocusedBackground:  color.NRGBA{R: 130, G: 130, B: 170, A: 255},
			FocusedBackground:          color.NRGBA{R: 170, G: 170, B: 180, A: 255},
			DisabledUnselected:         color.NRGBA{100, 100, 100, 255},
			DisabledSelected:           color.NRGBA{100, 100, 100, 255},
			DisabledSelectedBackground: color.NRGBA{100, 100, 100, 255},
		}),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			return e.(entryStruct).name
		}),
		widget.ListOpts.EntryTextPadding(widget.NewInsetsSimple(5)),
		widget.ListOpts.EntryTextPosition(widget.TextPositionStart, widget.TextPositionCenter),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			entry := args.Entry.(entryStruct)

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
					effects.ItemTrigger(nil, entry.entity, effects.Single{members[0]}, world)
					st.ui = st.initUI(world)
				}),
			))

			windowContainer.AddChild(widget.NewButton(
				widget.ButtonOpts.Image(buttonImage),
				widget.ButtonOpts.Text("捨てる", face, &widget.ButtonTextColor{
					Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
				}),
			))

			itemDesc := widget.NewText(
				widget.TextOpts.Text(entry.description, face, color.White),
				widget.TextOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.RowLayoutData{
						Position: widget.RowLayoutPositionCenter,
					}),
				),
			)
			windowContainer.AddChild(itemDesc)

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
				widget.WindowOpts.MinSize(240, 200),
				widget.WindowOpts.MaxSize(240, 400),
			)

			x, y := window.Contents.PreferredSize()
			r := image.Rect(0, 0, x, y)
			r = r.Add(image.Point{200, 50})
			window.SetLocation(r)
			ui.AddWindow(window)
			st.clickedItem = entry.entity
		}),
	)

	rootContainer.AddChild(list)

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
