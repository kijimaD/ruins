package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type CampMenuState struct {
	selection int
	campMenu  []ecs.Entity
}

// State interface ================

func (st *CampMenuState) OnPause(world w.World) {}

func (st *CampMenuState) OnResume(world w.World) {}

func (st *CampMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.campMenu = append(st.campMenu, loader.AddEntities(world, prefabs.Menu.CampMenu)...)
}

func (st *CampMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.campMenu...)
}

func (st *CampMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "description" {
			switch st.selection {
			case 0:
				text.Text = "アイテムを使う"
			case 1:
				text.Text = "装備を変更する"
			}
		}
	}))

	return updateMenu(st, world)
}

// Menu Interface ================

func (st *CampMenuState) getSelection() int {
	return st.selection
}

func (st *CampMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *CampMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&InventoryMenuState{}}}
	case 1:
		// TODO: 実装する
		return states.Transition{Type: states.TransNone}
	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *CampMenuState) getMenuIDs() []string {
	return []string{"item", "equip"}
}

func (st *CampMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_item", "cursor_equip"}
}
