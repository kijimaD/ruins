// 拠点でのコマンド選択画面
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

type HomeMenuState struct {
	selection int
	homeMenu  []ecs.Entity
}

// State interface ================

func (st *HomeMenuState) OnPause(world w.World) {
	st.OnStop(world)
}

func (st *HomeMenuState) OnResume(world w.World) {
	st.OnStart(world)
}

func (st *HomeMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	st.homeMenu = append(st.homeMenu, loader.AddEntities(world, prefabs.Menu.HomeMenu)...)
}

func (st *HomeMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.homeMenu...)
}

func (st *HomeMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		if text.ID == "description" {
			switch st.selection {
			case 0:
				text.Text = "遺跡に出発する"
			case 1:
				text.Text = "アイテムを合成する"
			case 2:
				text.Text = "アイテムを使う"
			case 3:
				text.Text = "仲間を入れ替える"
			case 4:
				text.Text = "設定を変更する"
			case 5:
				text.Text = "終了する"
			}
		}
	}))

	return updateMenu(st, world)
}

// Menu Interface ================

func (st *HomeMenuState) getSelection() int {
	return st.selection
}

func (st *HomeMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *HomeMenuState) confirmSelection(world w.World) states.Transition {
	switch st.selection {
	case 0:
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonSelectState{}}}
	case 1:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MixMenuState{}}}
	case 2:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&InventoryMenuState{}}}
	case 3:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&FieldState{}}}
	case 4:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&FieldState{}}}
	case 5:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}

	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *HomeMenuState) getMenuIDs() []string {
	return []string{"dungeon", "buy", "equip", "party", "system", "exit"}
}

func (st *HomeMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_dungeon", "cursor_buy", "cursor_equip", "cursor_party", "cursor_system", "cursor_exit"}
}
