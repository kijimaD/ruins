package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// DebugMenuState はデバッグメニューのゲームステート
type DebugMenuState struct {
	es.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	menuBuilder   *menu.MenuUIBuilder
	keyboardInput input.KeyboardInput
}

func (st DebugMenuState) String() string {
	return "DebugMenu"
}

// State interface ================

var _ es.State = &DebugMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *DebugMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *DebugMenuState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *DebugMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *DebugMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *DebugMenuState) Update(_ w.World) es.Transition {
	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *DebugMenuState) Draw(_ w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

func (st *DebugMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)

	// Menuコンポーネントを作成
	st.createDebugMenu(world)

	// MenuのUIを構築
	st.menuBuilder = menu.NewMenuUIBuilder(world)
	menuContainer := st.menuBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}

// createDebugMenu はデバッグメニューを作成する
func (st *DebugMenuState) createDebugMenu(world w.World) {
	items := make([]menu.MenuItem, len(debugMenuTrans))
	for i, data := range debugMenuTrans {
		items[i] = menu.MenuItem{
			ID:       data.label,
			Label:    data.label,
			UserData: i, // debugMenuTransのインデックスを保存
		}
	}

	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
		Columns:        0,
	}

	callbacks := menu.MenuCallbacks{
		OnSelect: func(index int, _ menu.MenuItem) {
			st.executeDebugMenuItem(world, index)
		},
		OnCancel: func() {
			st.SetTransition(es.Transition{Type: es.TransPop})
		},
		OnFocusChange: func(_, _ int) {
			// フォーカス変更時にUIを更新
			if st.menuBuilder != nil {
				st.menuBuilder.UpdateFocus(st.menu)
			}
		},
	}

	st.menu = menu.NewMenu(config, callbacks)
}

// executeDebugMenuItem は選択されたデバッグメニュー項目を実行する
func (st *DebugMenuState) executeDebugMenuItem(world w.World, index int) {
	if index < 0 || index >= len(debugMenuTrans) {
		return
	}

	data := debugMenuTrans[index]
	data.f(world)

	// 遅延評価でTransitionを取得（毎回新しいインスタンスが作成される）
	st.SetTransition(data.getTransFunc())
}

var debugMenuTrans = []struct {
	label        string
	f            func(world w.World)
	getTransFunc func() es.Transition // 遅延評価で毎回新しいTransitionを取得
}{
	{
		label:        "回復薬スポーン(インベントリ)",
		f:            func(world w.World) { worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack) },
		getTransFunc: func() es.Transition { return es.Transition{Type: es.TransNone} },
	},
	{
		label:        "手榴弾スポーン(インベントリ)",
		f:            func(world w.World) { worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack) },
		getTransFunc: func() es.Transition { return es.Transition{Type: es.TransNone} },
	},
	{
		label: "戦闘開始",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition {
			return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewBattleState}}
		},
	},
	{
		label: "汎用戦闘イベント開始",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition {
			return es.Transition{Type: es.TransPush, NewStateFuncs: GetRaidEvent1Factories()}
		},
	},
	{
		label: "汎用アイテム入手イベント開始",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition {
			return es.Transition{Type: es.TransPush, NewStateFuncs: GetItemGetEvent1Factories()}
		},
	},
	{
		label: "ゲームオーバー",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition {
			return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewGameOverState}}
		},
	},
	{
		label:        TextClose,
		f:            func(_ w.World) {},
		getTransFunc: func() es.Transition { return es.Transition{Type: es.TransPop} },
	},
}
