// 拠点でのコマンド選択画面
package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/materialhelper"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/spawner"
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

	// デバッグ用
	// 初回のみ追加する
	count := 0
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		count++
	}))
	if count == 0 {
		spawner.SpawnItem(world, "木刀", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "ハンドガン", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "レイガン", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnMember(world, "村上", true)
		spawner.SpawnMember(world, "白瀬", true)
		spawner.SpawnAllMaterials(world)
		materialhelper.PlusAmount("鉄", 40, world)
		materialhelper.PlusAmount("鉄くず", 4, world)
		materialhelper.PlusAmount("緑ハーブ", 2, world)
		materialhelper.PlusAmount("フェライトコア", 30, world)
		spawner.SpawnAllRecipes(world)
	}
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
				text.Text = "仲間を入れ替える"
			case 3:
				text.Text = "キャンプメニューを開く"
			case 4:
				text.Text = "終了する"
			}
		}
	}))

	names := []string{}
	hps := []string{}
	sps := []string{}
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
		gameComponents.Name,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)

		names = append(names, fmt.Sprintf("%-4s Lv.%d", name.Name, pools.Level))
		hps = append(hps, fmt.Sprintf("HP %3d / %3d", pools.HP.Current, pools.HP.Max))
		sps = append(sps, fmt.Sprintf("SP %3d / %3d", pools.SP.Current, pools.SP.Max))
	}))

	world.Manager.Join(
		world.Components.Engine.Text,
		world.Components.Engine.UITransform,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*ec.Text)
		switch text.ID {
		case "party_1_name":
			if len(names) > 0 {
				text.Text = names[0]
			}
		case "party_2_name":
			if len(names) > 1 {
				text.Text = names[1]
			}
		case "party_3_name":
			if len(names) > 2 {
				text.Text = names[2]
			}
		case "party_4_name":
			if len(names) > 3 {
				text.Text = names[3]
			}
		case "party_1_hp_label":
			text.Text = hps[0]
		case "party_1_sp_label":
			text.Text = sps[0]
		case "party_2_hp_label":
			text.Text = hps[1]
		case "party_2_sp_label":
			text.Text = sps[1]
		}
	}))

	return updateMenu(st, world)
}

func (st *HomeMenuState) Draw(world w.World, screen *ebiten.Image) {}

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
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CraftMenuState{}}}
	case 2:
		// TODO: 実装する
		return states.Transition{Type: states.TransNone}
	case 3:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CampMenuState{}}}
	case 4:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}}

	}
	panic(fmt.Errorf("unknown selection: %d", st.selection))
}

func (st *HomeMenuState) getMenuIDs() []string {
	return []string{"dungeon", "mix", "party", "camp", "exit"}
}

func (st *HomeMenuState) getCursorMenuIDs() []string {
	return []string{"cursor_dungeon", "cursor_mix", "cursor_party", "cursor_camp", "cursor_exit"}
}
