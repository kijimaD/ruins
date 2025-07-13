package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
)

// MainMenuState は新しいメニューコンポーネントを使用するメインメニュー
type MainMenuState struct {
	states.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.MenuUIBuilder
	keyboardInput input.KeyboardInput
}

func (st MainMenuState) String() string {
	return "MainMenu"
}

// State interface ================

var _ es.State = &MainMenuState{}

func (st *MainMenuState) OnPause(world w.World) {}

func (st *MainMenuState) OnResume(world w.World) {}

func (st *MainMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.initMenu(world)
	st.ui = st.initUI(world)
}

func (st *MainMenuState) OnStop(world w.World) {}

func (st *MainMenuState) Update(world w.World) states.Transition {
	// Escapeキーでの終了処理はメニューのOnCancelで処理するため、ここでは削除

	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

func (st *MainMenuState) Draw(world w.World, screen *ebiten.Image) {
	bg := (*world.Resources.SpriteSheets)["bg_title1"]
	screen.DrawImage(bg.Texture.Image, nil)

	st.ui.Draw(screen)
}

// initMenu はメニューコンポーネントを初期化する
func (st *MainMenuState) initMenu(world w.World) {
	// メニュー項目の定義
	items := []menu.MenuItem{
		{
			ID:       "intro",
			Label:    "導入",
			UserData: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&IntroState{}}},
		},
		{
			ID:       "home",
			Label:    "拠点",
			UserData: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}},
		},
		{
			ID:       "explore",
			Label:    "探検",
			UserData: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&DungeonState{Depth: 1}}},
		},
		{
			ID:       "exit",
			Label:    "終了",
			UserData: states.Transition{Type: states.TransQuit},
		},
	}

	// メニューの設定
	config := menu.MenuConfig{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.MenuCallbacks{
		OnSelect: func(index int, item menu.MenuItem) {
			// 選択されたアイテムのUserDataからTransitionを取得
			if trans, ok := item.UserData.(states.Transition); ok {
				st.SetTransition(trans)
			}
		},
		OnCancel: func() {
			// Escapeキーが押された時の処理
			st.SetTransition(states.Transition{Type: states.TransQuit})
		},
		OnFocusChange: func(oldIndex, newIndex int) {
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewMenuUIBuilder(world)
}

// initUI はUIを初期化する
func (st *MainMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()

	// メニューのUIを構築してコンテナに追加
	menuContainer := st.uiBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}
