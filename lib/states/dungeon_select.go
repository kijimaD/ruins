package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type DungeonSelectState struct {
	ui                     *ebitenui.UI
	trans                  *states.Transition
	dungeonSelectContainer *widget.Container
	dungeonDescContainer   *widget.Container
}

func (st DungeonSelectState) String() string {
	return "DungeonSelect"
}

// State interface ================

func (st *DungeonSelectState) OnPause(world w.World) {}

func (st *DungeonSelectState) OnResume(world w.World) {}

func (st *DungeonSelectState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *DungeonSelectState) OnStop(world w.World) {}

func (st *DungeonSelectState) Update(world w.World) states.Transition {
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

func (st *DungeonSelectState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

func (st *DungeonSelectState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()

	st.dungeonSelectContainer = eui.NewRowContainer()
	rootContainer.AddChild(st.dungeonSelectContainer)

	st.dungeonDescContainer = eui.NewRowContainer()
	rootContainer.AddChild(st.dungeonDescContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *DungeonSelectState) updateMenuContainer(world w.World) {
	st.dungeonSelectContainer.RemoveChildren()

	for _, data := range dungeonSelectTrans {
		data := data
		btn := eui.NewItemButton(
			data.label,
			func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			},
			world,
		)
		btn.GetWidget().CursorEnterEvent.AddHandler(func(args interface{}) {
			st.dungeonDescContainer.RemoveChildren()
			st.dungeonDescContainer.AddChild(eui.NewMenuText(data.desc, world))
		})
		st.dungeonSelectContainer.AddChild(btn)
	}
}

// MEMO: まだtransは全部同じ
var dungeonSelectTrans = []struct {
	label string
	desc  string
	trans states.Transition
}{
	{
		label: "森の遺跡",
		desc:  "鬱蒼とした森の奥地にある遺跡",
		trans: states.Transition{Type: states.TransReplace, NewStates: []states.State{&DungeonState{Depth: 1}}},
	},
	{
		label: "山の遺跡",
		desc:  "切り立った山の洞窟にある遺跡",
		trans: states.Transition{Type: states.TransReplace, NewStates: []states.State{&DungeonState{Depth: 1}}},
	},
	{
		label: "塔の遺跡",
		desc:  "雲にまで届く塔を持つ遺跡",
		trans: states.Transition{Type: states.TransReplace, NewStates: []states.State{&DungeonState{Depth: 1}}},
	},
	{
		label: "拠点メニューに戻る",
		desc:  "",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}},
	},
}
