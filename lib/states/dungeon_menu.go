package states

import (
	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/styles"
	"github.com/kijimaD/ruins/lib/widgets/menu"
)

type DungeonMenuState struct {
	es.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.MenuUIBuilder
	keyboardInput input.KeyboardInput
}

func (st DungeonMenuState) String() string {
	return "DungeonMenu"
}

// State interface ================

var _ es.State = &DungeonMenuState{}

func (st *DungeonMenuState) OnPause(world w.World) {}

func (st *DungeonMenuState) OnResume(world w.World) {
	// フォーカス状態を更新
	if st.uiBuilder != nil && st.menu != nil {
		st.uiBuilder.UpdateFocus(st.menu)
	}
}

func (st *DungeonMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	st.initMenu(world)
	st.ui = st.initUI(world)
}

func (st *DungeonMenuState) OnStop(world w.World) {}

func (st *DungeonMenuState) Update(world w.World) es.Transition {
	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

func (st *DungeonMenuState) Draw(world w.World, screen *ebiten.Image) {
	st.ui.Draw(screen)
}

// ================

// initMenu はメニューコンポーネントを初期化する
func (st *DungeonMenuState) initMenu(world w.World) {
	// メニュー項目の定義
	items := []menu.MenuItem{
		{
			ID:          "close",
			Label:       "閉じる",
			Description: "ダンジョンメニューを閉じる",
			UserData:    es.Transition{Type: es.TransPop},
		},
		{
			ID:          "exit",
			Label:       "終了",
			Description: "メインメニューに戻る",
			UserData:    es.Transition{Type: es.TransSwitch, NewStates: []es.State{&MainMenuState{}}},
		},
	}

	// メニュー設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバック設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(index int, item menu.MenuItem) {
			if trans, ok := item.UserData.(es.Transition); ok {
				st.SetTransition(trans)
			}
		},
		OnCancel: func() {
			st.SetTransition(es.Transition{Type: es.TransPop})
		},
		OnFocusChange: func(oldIndex, newIndex int) {
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	st.menu = menu.NewMenu(config, callbacks)
}

func (st *DungeonMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer(
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(styles.TransBlackColor)),
	)

	// UIビルダーを使用してメニューUIを構築
	st.uiBuilder = menu.NewMenuUIBuilder(world)
	menuContainer := st.uiBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}
