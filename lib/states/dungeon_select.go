package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type DungeonSelectState struct {
	selection     int
	dungeonSelect []ecs.Entity
}

// State interface ================

func (st *DungeonSelectState) OnPause(world w.World) {}

func (st *DungeonSelectState) OnResume(world w.World) {}

func (st *DungeonSelectState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.dungeonSelect = append(st.dungeonSelect, loader.AddEntities(world, prefabs.Menu.DungeonSelect)...)
}

func (st *DungeonSelectState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.dungeonSelect...)
}

func (st *DungeonSelectState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}
	return updateMenu(st, world)
}

// Menu Interface ================

func (st *DungeonSelectState) getSelection() int {
	return st.selection
}

func (st *DungeonSelectState) setSelection(selection int) {
	st.selection = selection
}

func (st *DungeonSelectState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&FieldState{}}}
	case 1:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&FieldState{}}}
	case 2:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&FieldState{}}}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *DungeonSelectState) getMenuIDs() []string {
	return []string{"forest", "mountain", "tower"}
}

func (st *DungeonSelectState) getCursorMenuIDs() []string {
	return []string{"cursor_forest", "cursor_mountain", "cursor_tower"}
}
