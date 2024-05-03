package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/spawner"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type DebugMenuState struct {
	selection int
	debugMenu []ecs.Entity

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

func (st *DebugMenuState) OnStop(world w.World) {
	world.Manager.DeleteEntities(st.debugMenu...)
}

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
	rootContainer := eui.NewVerticalContainer()
	st.debugMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(st.debugMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *DebugMenuState) updateMenuContainer(world w.World) {
	st.debugMenuContainer.RemoveChildren()

	for _, data := range debugMenuTrans {
		data := data
		btn := eui.NewItemButton(
			data.label,
			func(args *widget.ButtonClickedEventArgs) {
				data.f(world)
				st.trans = &data.trans
			},
			world,
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
		f:     func(world w.World) { spawner.SpawnItem(world, "回復薬", raw.SpawnInBackpack) },
		trans: states.Transition{Type: states.TransNone},
	},
	{
		label: "手榴弾スポーン(インベントリ)",
		f:     func(world w.World) { spawner.SpawnItem(world, "手榴弾", raw.SpawnInBackpack) },
		trans: states.Transition{Type: states.TransNone},
	},
	{
		label: "閉じる",
		f:     func(world w.World) {},
		trans: states.Transition{Type: states.TransPop},
	},
}
