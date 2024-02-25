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

	slots           []*ecs.Entity
	actionContainer *widget.Container
	specContainer   *widget.Container
	itemDesc        *widget.Text
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
	gameComponents := world.Components.Game.(*gc.Components)
	members := []ecs.Entity{}
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		members = append(members, entity)
	}))

	st.actionContainer = st.newItemSpecContainer(world)
	st.specContainer = st.newItemSpecContainer(world)
	st.generateActionContainer(world, members[0])

	toggleContainer := eui.NewRowContainer()
	toggleEquipButton := eui.NewItemButton("装備", func(args *widget.ButtonClickedEventArgs) {}, world)
	toggleSkillButton := eui.NewItemButton("技能", func(args *widget.ButtonClickedEventArgs) {}, world)
	toggleContainer.AddChild(toggleEquipButton)
	toggleContainer.AddChild(toggleSkillButton)

	itemDescContainer := eui.NewRowContainer()
	st.itemDesc = eui.NewMenuText(" ", world) // 空白だと初期状態の縦サイズがなくなる
	itemDescContainer.AddChild(st.itemDesc)

	rootContainer := eui.NewItemGridContainer()
	{
		rootContainer.AddChild(eui.NewMenuText("装備", world))
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(toggleContainer)

		rootContainer.AddChild(st.actionContainer)
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(st.specContainer)

		rootContainer.AddChild(st.itemDesc)
	}

	return &ebitenui.UI{Container: rootContainer}
}

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
func (st *EquipMenuState) generateActionContainer(world w.World, member ecs.Entity) {
	st.actionContainer.RemoveChildren()
	st.slots = equips.GetEquipments(world, member)

	gameComponents := world.Components.Game.(*gc.Components)
	for _, slot := range st.slots {
		slot := slot
		var name = "_"
		var desc = " "
		if slot != nil {
			name = gameComponents.Name.Get(*slot).(*gc.Name).Name
			desc = gameComponents.Description.Get(*slot).(*gc.Description).Description
		}
		itemButton := eui.NewItemButton(name, func(args *widget.ButtonClickedEventArgs) {
		}, world)

		itemButton.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.itemDesc.Label = desc
			if slot != nil {
				views.UpdateSpec(world, st.specContainer, []any{
					simple.GetWeapon(world, *slot),
					simple.GetWearable(world, *slot),
				})
			} else {
				st.specContainer.RemoveChildren()
			}
		})
		st.actionContainer.AddChild(itemButton)
	}
}
