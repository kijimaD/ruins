package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/sokotwo/lib/components"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	selection     int
	inventoryMenu []ecs.Entity
}

// State interface ================

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.inventoryMenu = append(st.inventoryMenu, loader.AddEntities(world, prefabs.Menu.InventoryMenu)...)
}

func (st *InventoryMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.inventoryMenu...)
}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CampMenuState{}}}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	itemTexts := ""
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Name,
		gameComponents.Description,
		gameComponents.InBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		description := gameComponents.Description.Get(entity).(*gc.Description)
		itemTexts += fmt.Sprintf("%s %s\n", name.Name, description.Description)
	}))

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "description" {
			text.Text = itemTexts
		}
	}))

	return updateMenu(st, world)
}

// Menu Interface ================

func (st *InventoryMenuState) getSelection() int {
	return st.selection
}

func (st *InventoryMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *InventoryMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		// TODO: 実装
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *InventoryMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *InventoryMenuState) getCursorMenuIDs() []string {
	return []string{""}
}
