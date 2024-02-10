package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type MixMenuState struct {
	selection int
	mixMenu   []ecs.Entity
}

// State interface ================

func (st *MixMenuState) OnPause(world w.World) {}

func (st *MixMenuState) OnResume(world w.World) {}

func (st *MixMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.mixMenu = append(st.mixMenu, loader.AddEntities(world, prefabs.Menu.MixMenu)...)
}

func (st *MixMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.mixMenu...)
}

func (st *MixMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}
	return updateMenu(st, world)
}

func (st *MixMenuState) Draw(world w.World, screen *ebiten.Image) {}

// Menu Interface ================

func (st *MixMenuState) getSelection() int {
	return st.selection
}

func (st *MixMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *MixMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *MixMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *MixMenuState) getCursorMenuIDs() []string {
	return []string{""}
}
