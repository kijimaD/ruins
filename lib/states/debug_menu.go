package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	"github.com/kijimaD/sokotwo/lib/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type DebugMenuState struct {
	selection int
	debugMenu []ecs.Entity
}

// State interface ================

func (st *DebugMenuState) OnPause(world w.World) {}

func (st *DebugMenuState) OnResume(world w.World) {}

func (st *DebugMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.debugMenu = append(st.debugMenu, loader.AddEntities(world, prefabs.Menu.DebugMenu)...)
}

func (st *DebugMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.debugMenu...)
}

func (st *DebugMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}
	return updateMenu(st, world)
}

func (st *DebugMenuState) Draw(world w.World, screen *ebiten.Image) {}

// Menu Interface ================

func (st *DebugMenuState) getSelection() int {
	return st.selection
}

func (st *DebugMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *DebugMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		spawner.SpawnItem(world, "回復薬")

		return states.Transition{Type: states.TransNone}
	case 1:
		spawner.SpawnItem(world, "手榴弾")

		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *DebugMenuState) getMenuIDs() []string {
	return []string{"spawn_item_potion", "spawn_item_grenade"}
}

func (st *DebugMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_spawn_item_potion", "cursor_spawn_item_grenade"}
}
