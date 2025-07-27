package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/eui"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/menu"
	w "github.com/kijimaD/ruins/lib/world"
)

// MainMenuState は新しいメニューコンポーネントを使用するメインメニュー
type MainMenuState struct {
	es.BaseState
	ui            *ebitenui.UI
	menu          *menu.Menu
	uiBuilder     *menu.UIBuilder
	keyboardInput input.KeyboardInput
}

func (st MainMenuState) String() string {
	return "MainMenu"
}

// State interface ================

var _ es.State = &MainMenuState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *MainMenuState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *MainMenuState) OnResume(_ w.World) {}

// OnStart はステート開始時の処理を行う
func (st *MainMenuState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}
	st.initMenu(world)
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *MainMenuState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *MainMenuState) Update(_ w.World) es.Transition {
	// Escapeキーでの終了処理はメニューのOnCancelで処理するため、ここでは削除

	// メニューの更新
	st.menu.Update(st.keyboardInput)

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はスクリーンに描画する
func (st *MainMenuState) Draw(world w.World, screen *ebiten.Image) {
	bg := (*world.Resources.SpriteSheets)["bg_title1"]
	screen.DrawImage(bg.Texture.Image, nil)

	st.ui.Draw(screen)
}

// initMenu はメニューコンポーネントを初期化する
func (st *MainMenuState) initMenu(world w.World) {
	// メニュー項目の定義
	items := []menu.Item{
		{
			ID:       "intro",
			Label:    "導入",
			UserData: es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewIntroState}},
		},
		{
			ID:       "home",
			Label:    "拠点",
			UserData: es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}},
		},
		{
			ID:       "explore",
			Label:    "探検",
			UserData: es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(1)}},
		},
		{
			ID:       "exit",
			Label:    "終了",
			UserData: es.Transition{Type: es.TransQuit},
		},
	}

	// メニューの設定
	config := menu.Config{
		Items:          items,
		InitialIndex:   0,
		WrapNavigation: true,
		Orientation:    menu.Vertical,
	}

	// コールバックの設定
	callbacks := menu.Callbacks{
		OnSelect: func(_ int, item menu.Item) {
			// 選択されたアイテムのUserDataからTransitionを取得
			if trans, ok := item.UserData.(es.Transition); ok {
				st.SetTransition(trans)
			}
		},
		OnCancel: func() {
			// Escapeキーが押された時の処理
			st.SetTransition(es.Transition{Type: es.TransQuit})
		},
		OnFocusChange: func(_, _ int) {
			// フォーカス変更時にUIを更新
			if st.uiBuilder != nil {
				st.uiBuilder.UpdateFocus(st.menu)
			}
		},
	}

	// メニューを作成
	st.menu = menu.NewMenu(config, callbacks)

	// UIビルダーを作成
	st.uiBuilder = menu.NewUIBuilder(world)
}

// initUI はUIを初期化する
func (st *MainMenuState) initUI(_ w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()

	// メニューのUIを構築してコンテナに追加
	menuContainer := st.uiBuilder.BuildUI(st.menu)
	rootContainer.AddChild(menuContainer)

	return &ebitenui.UI{Container: rootContainer}
}
