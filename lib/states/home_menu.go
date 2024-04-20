// 拠点でのコマンド選択画面
package states

import (
	"fmt"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/loader"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/spawner"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type HomeMenuState struct {
	selection int
	homeMenu  []ecs.Entity
	ui        *ebitenui.UI

	memberContainer     *widget.Container
	actionListContainer *widget.Container
	actionDescContainer *widget.Container
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
		armor := spawner.SpawnItem(world, "西洋鎧", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "作業用ヘルメット", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "革のブーツ", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "ルビー原石", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "回復スプレー", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack)
		ishihara := spawner.SpawnMember(world, "イシハラ", true)
		spawner.SpawnMember(world, "シラセ", true)
		spawner.SpawnAllMaterials(world)
		material.PlusAmount("鉄", 40, world)
		material.PlusAmount("鉄くず", 4, world)
		material.PlusAmount("緑ハーブ", 2, world)
		material.PlusAmount("フェライトコア", 30, world)
		spawner.SpawnAllRecipes(world)

		equips.Equip(world, armor, ishihara, gc.EquipmentSlotZero)
	}

	st.ui = st.initUI(world)
}

func (st *HomeMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.homeMenu...)
}

// FIXME: 毎ループでやってるので重い。変更があったときだけ、でいい
func (st *HomeMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	effects.RunEffectQueue(world)
	_ = gs.EquipmentChangedSystem(world)

	// 完全回復
	effects.AddEffect(nil, effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
	effects.AddEffect(nil, effects.RecoveryStamina{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})

	labels := []string{
		"出発",
		"合成",
		"入替",
		"陣営",
		"終了",
	}
	st.actionListContainer.RemoveChildren()
	for i, label := range labels {
		btn := eui.NewMenuText(label, world)
		if st.selection == i {
			btn.Color = styles.ButtonHoverColor
		}
		st.actionListContainer.AddChild(btn)
	}

	descs := []string{
		"遺跡に出発する",
		"アイテムを合成する",
		"仲間を入れ替える",
		"キャンプメニューを開く",
		"終了する",
	}
	desc := descs[st.selection]

	st.actionDescContainer.RemoveChildren()
	st.actionDescContainer.AddChild(eui.NewMenuText(desc, world))

	st.memberContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))

	st.ui.Update()

	return updateMenu(st, world)
}

func (st *HomeMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
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
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&CraftMenuState{category: itemCategoryTypeItem}}}
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

// ================

func (st *HomeMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewItemGridContainer()
	st.memberContainer = eui.NewVerticalContainer()
	st.actionListContainer = eui.NewRowContainer()
	st.actionDescContainer = eui.NewRowContainer()
	{
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(st.memberContainer)

		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(eui.NewEmptyContainer())
		rootContainer.AddChild(eui.NewEmptyContainer())

		rootContainer.AddChild(eui.NewVSplitContainer(st.actionListContainer, st.actionDescContainer))
	}

	return &ebitenui.UI{Container: rootContainer}
}
