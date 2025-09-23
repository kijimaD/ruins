package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// DebugMenuState はデバッグメニューのゲームステート
type DebugMenuState struct {
	es.BaseState[w.World]
	ui            *ebitenui.UI
	menu          *menu.Menu
	menuBuilder   *menu.UIBuilder
	keyboardInput input.KeyboardInput
}

func (st DebugMenuState) String() string {
	return "DebugMenu"
}

// State interface ================

var _ es.State[w.World] = &DebugMenuState{}

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
func (st *DebugMenuState) Update(_ w.World) es.Transition[w.World] {
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
	rootContainer := styled.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(consts.BlackColor)),
	)

	// Menuコンポーネントを作成
	st.createDebugMenu(world)

	// MenuのUIを構築
	st.menuBuilder = menu.NewUIBuilder(world)
	menuContainer := st.menuBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}

// createDebugMenu はデバッグメニューを作成する
func (st *DebugMenuState) createDebugMenu(world w.World) {
	items := make([]menu.Item, len(debugMenuTrans))
	for i, data := range debugMenuTrans {
		items[i] = menu.Item{
			ID:       data.label,
			Label:    data.label,
			UserData: i, // debugMenuTransのインデックスを保存
		}
	}

	config := menu.Config{
		Items:             items,
		InitialIndex:      0,
		WrapNavigation:    true,
		Orientation:       menu.Vertical,
		ItemsPerPage:      20,
		ShowPageIndicator: true,
	}

	callbacks := menu.Callbacks{
		OnSelect: func(index int, _ menu.Item) {
			st.executeDebugMenuItem(world, index)
		},
		OnCancel: func() {
			st.SetTransition(es.Transition[w.World]{Type: es.TransPop})
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
	getTransFunc func() es.Transition[w.World] // 遅延評価で毎回新しいTransitionを取得
}{
	{
		label: "回復薬スポーン(インベントリ)",
		f: func(world w.World) {
			_, err := worldhelper.SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
			if err != nil {
				panic(err)
			}
		},
		getTransFunc: func() es.Transition[w.World] { return es.Transition[w.World]{Type: es.TransNone} },
	},
	{
		label: "手榴弾スポーン(インベントリ)",
		f: func(world w.World) {
			_, err := worldhelper.SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
			if err != nil {
				panic(err)
			}
		},
		getTransFunc: func() es.Transition[w.World] { return es.Transition[w.World]{Type: es.TransNone} },
	},
	{
		label: "汎用アイテム入手イベント開始",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: GetItemGetEvent1Factories()}
		},
	},
	{
		label: "ゲームオーバー",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewGameOverState}}
		},
	},
	{
		label: "ダンジョン開始(大部屋)",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeBigRoom),
			}}
		},
	},
	{
		label: "ダンジョン開始(小部屋)",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeSmallRoom),
			}}
		},
	},
	{
		label: "ダンジョン開始(洞窟)",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeCave),
			}}
		},
	},
	{
		label: "ダンジョン開始(廃墟)",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeRuins),
			}}
		},
	},
	{
		label: "ダンジョン開始(森)",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransReplace, NewStateFuncs: []es.StateFactory[w.World]{
				NewDungeonStateWithBuilder(1, mapbuilder.BuilderTypeForest),
			}}
		},
	},
	{
		label: "メッセージウィンドウテスト",
		f:     func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] {
			return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewMessageWindowState}}
		},
	},
	{
		label:        TextClose,
		f:            func(_ w.World) {},
		getTransFunc: func() es.Transition[w.World] { return es.Transition[w.World]{Type: es.TransPop} },
	},
}
