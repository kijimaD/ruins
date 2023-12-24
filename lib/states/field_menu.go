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

type FieldMenuState struct {
	selection int
	fieldMenu []ecs.Entity
}

// State interface ================

func (st *FieldMenuState) OnPause(world w.World) {}

func (st *FieldMenuState) OnResume(world w.World) {}

func (st *FieldMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.fieldMenu = append(st.fieldMenu, loader.AddEntities(world, prefabs.Menu.FieldMenu)...)
}

func (st *FieldMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.fieldMenu...)
}

func (st *FieldMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}

	return updateMenu(st, world)
}

func (st *FieldMenuState) Draw(world w.World, screen *ebiten.Image) {}

// Menu Interface ================

func (st *FieldMenuState) getSelection() int {
	return st.selection
}

func (st *FieldMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *FieldMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransReplace, NewStates: []states.State{&MainMenuState{}}}
	case 1:
		return states.Transition{Type: states.TransPop}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *FieldMenuState) getMenuIDs() []string {
	return []string{"item", "close"}
}

func (st *FieldMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_item", "cursor_close"}
}
