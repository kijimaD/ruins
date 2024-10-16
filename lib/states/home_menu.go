// 拠点でのコマンド選択画面
package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper/equips"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type HomeMenuState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	// 背景
	bg *ebiten.Image

	// メンバー一覧コンテナ
	memberContainer *widget.Container
	// 選択肢アクション一覧コンテナ
	actionListContainer *widget.Container
	// 選択肢アクションの説明コンテナ
	actionDescContainer *widget.Container
}

func (st HomeMenuState) String() string {
	return "HomeMenu"
}

// State interface ================

func (st *HomeMenuState) OnPause(world w.World) {}

func (st *HomeMenuState) OnResume(world w.World) {}

func (st *HomeMenuState) OnStart(world w.World) {
	// デバッグ用
	// 初回のみ追加する
	memberCount := 0
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		memberCount++
	}))
	if memberCount == 0 {
		spawner.SpawnItem(world, "体当たり", gc.ItemLocationNone)

		spawner.SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
		card1 := spawner.SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
		card2 := spawner.SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
		armor := spawner.SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		spawner.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		ishihara := spawner.SpawnMember(world, "イシハラ", true)
		spawner.SpawnMember(world, "シラセ", true)
		spawner.SpawnMember(world, "タチバナ", true)
		spawner.SpawnAllMaterials(world)
		material.PlusAmount("鉄", 40, world)
		material.PlusAmount("鉄くず", 4, world)
		material.PlusAmount("緑ハーブ", 2, world)
		material.PlusAmount("フェライトコア", 30, world)
		spawner.SpawnAllRecipes(world)

		equips.Equip(world, card1, ishihara, gc.EquipmentSlotNumber(0))
		equips.Equip(world, card2, ishihara, gc.EquipmentSlotNumber(0))
		equips.Equip(world, armor, ishihara, gc.EquipmentSlotNumber(0))
	}

	{
		// ステータス反映(最大HP)
		_ = gs.EquipmentChangedSystem(world)
		// 完全回復
		effects.AddEffect(nil, effects.Healing{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
		effects.AddEffect(nil, effects.RecoveryStamina{Amount: gc.RatioAmount{Ratio: float64(1.0)}}, effects.Party{})
		effects.RunEffectQueue(world)
	}

	bg := (*world.Resources.SpriteSheets)["bg_cup1"]
	st.bg = bg.Texture.Image

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

	return states.Transition{Type: states.TransNone}
}

func (st *HomeMenuState) Draw(world w.World, screen *ebiten.Image) {
	if st.bg != nil {
		screen.DrawImage(st.bg, &ebiten.DrawImageOptions{})
	}

	st.ui.Draw(screen)
}

// ================

func (st *HomeMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.memberContainer = eui.NewRowContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)

	st.actionListContainer = eui.NewRowContainer()
	st.actionDescContainer = eui.NewRowContainer()
	st.actionDescContainer.AddChild(eui.NewMenuText(" ", world))

	actionContainer := eui.NewVSplitContainer(st.actionListContainer, st.actionDescContainer)
	rootContainer.AddChild(
		st.memberContainer,
		actionContainer,
	)

	st.updateActionList(world)
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

	for _, data := range homeMenuTrans {
		data := data
		btn := eui.NewItemButton(
			data.label,
			func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			},
			world,
		)
		btn.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.actionDescContainer.RemoveChildren()
			st.actionDescContainer.AddChild(eui.NewMenuText(data.desc, world))

			st.updateMemberContainer(world)
		})
		st.actionListContainer.AddChild(btn)
	}
}

// メンバー一覧を更新する
func (st *HomeMenuState) updateMemberContainer(world w.World) {
	st.memberContainer.RemoveChildren()
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		views.AddMemberBar(world, st.memberContainer, entity)
	}))
}
