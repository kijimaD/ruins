// 拠点でのコマンド選択画面
package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/spawner"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type HomeMenuState struct {
	selection int
	ui        *ebitenui.UI
	trans     *states.Transition

	memberContainer     *widget.Container // メンバー一覧コンテナ
	actionListContainer *widget.Container // 選択肢アクション一覧コンテナ
	actionDescContainer *widget.Container // 選択肢アクションの説明コンテナ
}

func (st HomeMenuState) String() string {
	return "HomeMenu"
}

// State interface ================

func (st *HomeMenuState) OnPause(world w.World) {
	st.OnStop(world)
}

func (st *HomeMenuState) OnResume(world w.World) {
	st.OnStart(world)
}

func (st *HomeMenuState) OnStart(world w.World) {
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
		spawner.SpawnMember(world, "ヒラヤマ", true)
		spawner.SpawnAllMaterials(world)
		material.PlusAmount("鉄", 40, world)
		material.PlusAmount("鉄くず", 4, world)
		material.PlusAmount("緑ハーブ", 2, world)
		material.PlusAmount("フェライトコア", 30, world)
		spawner.SpawnAllRecipes(world)

		equips.Equip(world, armor, ishihara, gc.EquipmentSlotZero)
	}

	// ステータス反映(最大HP)
	_ = gs.EquipmentChangedSystem(world)

	// 完全回復
	effects.AddEffect(nil, effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
	effects.AddEffect(nil, effects.RecoveryStamina{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
	effects.RunEffectQueue(world)

	st.ui = st.initUI(world)
}

func (st *HomeMenuState) OnStop(world w.World) {}

func (st *HomeMenuState) Update(world w.World) states.Transition {
	effects.RunEffectQueue(world)
	_ = gs.EquipmentChangedSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DebugMenuState{}}}
	}

	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return updateMenu(st, world)
}

func (st *HomeMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// Menu Interface ================
// 使ってないので削除していく

func (st *HomeMenuState) getSelection() int {
	return st.selection
}

func (st *HomeMenuState) setSelection(selection int) {
	st.selection = selection
}

func (st *HomeMenuState) confirmSelection(world w.World) states.Transition {
	return states.Transition{Type: states.TransNone}
}

func (st *HomeMenuState) getMenuIDs() []string {
	return []string{""}
}

func (st *HomeMenuState) getCursorMenuIDs() []string {
	return []string{""}
}

// ================

func (st *HomeMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.memberContainer = eui.NewRowContainer()

	st.actionListContainer = eui.NewRowContainer()
	st.actionDescContainer = eui.NewRowContainer()
	actionContainer := eui.NewVSplitContainer(st.actionListContainer, st.actionDescContainer)
	{
		rootContainer.AddChild(st.memberContainer)
		rootContainer.AddChild(actionContainer)
	}

	st.updateActionList(world)
	st.updateActionDesc(world)
	st.updateMemberContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

var homeMenuTrans = []struct {
	label string
	trans states.Transition
	desc  string
}{
	{
		label: "出発",
		trans: states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonSelectState{}}},
		desc:  "遺跡に出発する",
	},
	{
		label: "合成",
		trans: states.Transition{Type: states.TransPush, NewStates: []states.State{&CraftMenuState{category: ItemCategoryTypeItem}}},
		desc:  "アイテムを合成する",
	},
	{
		label: "入替",
		trans: states.Transition{Type: states.TransNone},
		desc:  "仲間を入れ替える(未実装)",
	},
	{
		label: "所持",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&InventoryMenuState{category: ItemCategoryTypeItem}}},
		desc:  "所持品を確認する",
	},
	{
		label: "装備",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&EquipMenuState{}}},
		desc:  "装備を変更する",
	},
	{
		label: "終了",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&MainMenuState{}}},
		desc:  "タイトル画面に戻る",
	},
}

// 選択肢を更新する
func (st *HomeMenuState) updateActionList(world w.World) {
	st.actionListContainer.RemoveChildren()

	for i, data := range homeMenuTrans {
		i := i
		label := data.label
		trans := data.trans
		btn := eui.NewItemButton(
			label,
			func(args *widget.ButtonClickedEventArgs) {
				st.trans = &trans
			},
			world,
		)
		btn.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.selection = i
			st.updateActionDesc(world)
			st.updateMemberContainer(world)
		})
		st.actionListContainer.AddChild(btn)
	}
}

// 選択肢の解説を更新する
func (st *HomeMenuState) updateActionDesc(world w.World) {
	st.actionDescContainer.RemoveChildren()
	st.actionDescContainer.AddChild(eui.NewMenuText(homeMenuTrans[st.selection].desc, world))
}

// メンバー一覧を更新する
func (st *HomeMenuState) updateMemberContainer(world w.World) {
	st.memberContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Member,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}
