package states

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/raw"
	"github.com/kijimaD/sokotwo/lib/resources"
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
		componentList := loader.EntityComponentList{}
		rawMaster := world.Resources.RawMaster.(raw.RawMaster)
		componentList.Game = append(componentList.Game, rawMaster.GenerateItem("回復薬"))
		componentList.Engine = append(componentList.Engine, loader.EngineComponentList{})
		loader.AddEntities(world, componentList)
		log.Println("スポーンした")
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *DebugMenuState) getMenuIDs() []string {
	return []string{"spawn"}
}

func (st *DebugMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_spawn"}
}
