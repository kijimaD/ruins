package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type EquipMenuState struct {
	selection int
	equipMenu []ecs.Entity
	ui        *ebitenui.UI
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
	rootContainer := eui.NewItemGridContainer()

	return &ebitenui.UI{Container: rootContainer}
}
