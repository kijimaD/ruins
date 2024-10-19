package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"
)

type DebugMenuState struct {
	ui                 *ebitenui.UI
	trans              *states.Transition
	debugMenuContainer *widget.Container
}

func (st DebugMenuState) String() string {
	return "DebugMenu"
}

// State interface ================

func (st *DebugMenuState) OnPause(world w.World) {}

func (st *DebugMenuState) OnResume(world w.World) {}

func (st *DebugMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *DebugMenuState) OnStop(world w.World) {}

func (st *DebugMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPop}
	}

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	st.ui.Update()

	return states.Transition{Type: states.TransNone}
}

func (st *DebugMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *DebugMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)
	st.debugMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(st.debugMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *DebugMenuState) updateMenuContainer(world w.World) {
	st.debugMenuContainer.RemoveChildren()

	for _, data := range debugMenuTrans {
		data := data
		btn := eui.NewButton(
			data.label,
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				data.f(world)
				st.trans = &data.trans
			}),
		)
		st.debugMenuContainer.AddChild(btn)
	}
}

var debugMenuTrans = []struct {
	label string
	f     func(world w.World)
	trans states.Transition
}{
	{
		label: "回復薬スポーン(インベントリ)",
		f:     func(world w.World) { spawner.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack) },
		trans: states.Transition{Type: states.TransNone},
	},
	{
		label: "手榴弾スポーン(インベントリ)",
		f:     func(world w.World) { spawner.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack) },
		trans: states.Transition{Type: states.TransNone},
	},
	{
		label: "戦闘開始",
		f:     func(world w.World) {},
		trans: states.Transition{Type: states.TransPush, NewStates: []states.State{&BattleState{}}},
	},
	{
		label: "ゲームオーバー",
		f:     func(world w.World) {},
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&GameOverState{}}},
	},
	{
		label: "閉じる",
		f:     func(world w.World) {},
		trans: states.Transition{Type: states.TransPop},
	},
}
