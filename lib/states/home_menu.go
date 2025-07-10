// 拠点でのコマンド選択画面
package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/views"
	"github.com/kijimaD/ruins/lib/worldhelper"
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

var _ es.State = &HomeMenuState{}

func (st *HomeMenuState) OnPause(world w.World) {}

func (st *HomeMenuState) OnResume(world w.World) {
	st.updateMemberContainer(world)
}

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
		worldhelper.SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
		card1 := worldhelper.SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
		card2 := worldhelper.SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
		card3 := worldhelper.SpawnItem(world, "M72 LAW", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
		armor := worldhelper.SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
		ishihara := worldhelper.SpawnMember(world, "イシハラ", true)
		shirase := worldhelper.SpawnMember(world, "シラセ", true)
		worldhelper.SpawnMember(world, "タチバナ", true)
		worldhelper.SpawnAllMaterials(world)
		worldhelper.PlusAmount("鉄", 40, world)
		worldhelper.PlusAmount("鉄くず", 4, world)
		worldhelper.PlusAmount("緑ハーブ", 2, world)
		worldhelper.PlusAmount("フェライトコア", 30, world)
		worldhelper.SpawnAllRecipes(world)
		worldhelper.SpawnAllCards(world)

		worldhelper.Equip(world, card1, ishihara, gc.EquipmentSlotNumber(0))
		worldhelper.Equip(world, card2, ishihara, gc.EquipmentSlotNumber(0))
		worldhelper.Equip(world, card3, shirase, gc.EquipmentSlotNumber(0))
		worldhelper.Equip(world, armor, ishihara, gc.EquipmentSlotNumber(0))
	}

	bg := (*world.Resources.SpriteSheets)["bg_cup1"]
	st.bg = bg.Texture.Image

	st.ui = st.initUI(world)
}

func (st *HomeMenuState) OnStop(world w.World) {}

func (st *HomeMenuState) Update(world w.World) states.Transition {
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

	st.actionListContainer = eui.NewVerticalContainer()
	st.actionDescContainer = eui.NewRowContainer()
	st.actionDescContainer.AddChild(eui.NewMenuText(" ", world))

	rootContainer.AddChild(
		st.memberContainer,
		st.actionDescContainer,
		st.actionListContainer,
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
		trans: states.Transition{Type: states.TransPush, NewStates: []states.State{&CraftMenuState{}}},
		desc:  "アイテムを合成する",
	},
	{
		label: "入替",
		trans: states.Transition{Type: states.TransNone},
		desc:  "仲間を入れ替える(未実装)",
	},
	{
		label: "所持",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&InventoryMenuState{}}},
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
		btn := eui.NewButton(
			data.label,
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			}),
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
