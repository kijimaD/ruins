package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/sokotwo/lib/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type InventoryMenuState struct {
	selection int
}

// State interface ================

func (st *InventoryMenuState) OnPause(world w.World) {}

func (st *InventoryMenuState) OnResume(world w.World) {}

func (st *InventoryMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.InventoryMenu)
}

func (st *InventoryMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *InventoryMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(gameComponents.Item).Visit(ecs.Visit(func(entity ecs.Entity) {
		item := gameComponents.Item.Get(entity).(*gc.Item)
		fmt.Println(item)
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
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonSelectState{}}}

	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *InventoryMenuState) getMenuIDs() []string {
	return []string{"dungeon"}
}

func (st *InventoryMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_dungeon"}
}
