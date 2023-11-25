package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
)

type MainMenuState struct {
	selection int
}

// State interface ================

func (st *MainMenuState) OnPause(world w.World) {}

func (st *MainMenuState) OnResume(world w.World) {}

func (st *MainMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.MainMenu)
}

func (st *MainMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *MainMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}
	return updateMenu(st, world)
}

// Menu Interface ================

func (st *MainMenuState) getSelection() int {
	return st.selection
}

func (st *MainMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MainMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GamePlayState{}}}
	case 1:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&IntroState{}}}
	case 2:
		return states.Transition{Type: states.TransQuit}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MainMenuState) getMenuIDs() []string {
	return []string{"start", "intro", "exit"}
}

func (st *MainMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_start", "cursor_intro", "cursor_exit"}
}
