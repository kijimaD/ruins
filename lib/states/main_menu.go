package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/eui"
)

type MainMenuState struct {
	ui    *ebitenui.UI
	trans *states.Transition

	mainMenuContainer *widget.Container
	menuButtons       []*widget.Button
	focusIndex        int
}

func (st MainMenuState) String() string {
	return "MainMenu"
}

// State interface ================

var _ es.State = &MainMenuState{}

func (st *MainMenuState) OnPause(world w.World) {}

func (st *MainMenuState) OnResume(world w.World) {}

func (st *MainMenuState) OnStart(world w.World) {
	st.ui = st.initUI(world)
}

func (st *MainMenuState) OnStop(world w.World) {}

func (st *MainMenuState) Update(world w.World) states.Transition {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransQuit}
	}

	// キーボードナビゲーション処理
	st.handleKeyboardNavigation()

	st.ui.Update()

	if st.trans != nil {
		next := *st.trans
		st.trans = nil
		return next
	}

	return states.Transition{Type: states.TransNone}
}

func (st *MainMenuState) Draw(world w.World, screen *ebiten.Image) {
	bg := (*world.Resources.SpriteSheets)["bg_title1"]
	screen.DrawImage(bg.Texture.Image, nil)

	st.ui.Draw(screen)
}

func (st *MainMenuState) initUI(world w.World) *ebitenui.UI {
	rootContainer := eui.NewVerticalContainer()
	st.mainMenuContainer = eui.NewVerticalContainer()
	rootContainer.AddChild(st.mainMenuContainer)

	st.updateMenuContainer(world)

	return &ebitenui.UI{Container: rootContainer}
}

func (st *MainMenuState) updateMenuContainer(world w.World) {
	st.mainMenuContainer.RemoveChildren()
	st.menuButtons = make([]*widget.Button, 0, len(mainMenuTrans))

	for i, data := range mainMenuTrans {
		data := data
		index := i
		btn := eui.NewButton(
			data.label,
			world,
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				st.trans = &data.trans
			}),
		)
		
		// 初期フォーカス設定
		if index == 0 {
			st.focusIndex = 0
			st.setButtonFocus(btn, true)
		} else {
			st.setButtonFocus(btn, false)
		}
		
		st.menuButtons = append(st.menuButtons, btn)
		st.mainMenuContainer.AddChild(btn)
	}
}

// キーボードナビゲーション処理
func (st *MainMenuState) handleKeyboardNavigation() {
	if len(st.menuButtons) == 0 {
		return
	}

	// Shift+Tabの判定
	isShiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift)
	
	// 矢印キー上下、Tab/Shift+Tabでフォーカス移動
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || (inpututil.IsKeyJustPressed(ebiten.KeyTab) && !isShiftPressed) {
		st.moveFocusDown()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || (inpututil.IsKeyJustPressed(ebiten.KeyTab) && isShiftPressed) {
		st.moveFocusUp()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// 選択されているメニューを実行
		st.executeCurrentSelection()
	}
}

// フォーカスを下に移動
func (st *MainMenuState) moveFocusDown() {
	if len(st.menuButtons) == 0 {
		return
	}

	// 現在のフォーカスを解除
	st.setButtonFocus(st.menuButtons[st.focusIndex], false)
	
	// 次のインデックスに移動（循環）
	st.focusIndex = (st.focusIndex + 1) % len(st.menuButtons)
	
	// 新しいフォーカスを設定
	st.setButtonFocus(st.menuButtons[st.focusIndex], true)
}

// フォーカスを上に移動
func (st *MainMenuState) moveFocusUp() {
	if len(st.menuButtons) == 0 {
		return
	}

	// 現在のフォーカスを解除
	st.setButtonFocus(st.menuButtons[st.focusIndex], false)
	
	// 前のインデックスに移動（循環）
	st.focusIndex = (st.focusIndex - 1 + len(st.menuButtons)) % len(st.menuButtons)
	
	// 新しいフォーカスを設定
	st.setButtonFocus(st.menuButtons[st.focusIndex], true)
}

// 現在選択されているメニューを実行
func (st *MainMenuState) executeCurrentSelection() {
	if len(st.menuButtons) > 0 && st.focusIndex >= 0 && st.focusIndex < len(mainMenuTrans) {
		st.trans = &mainMenuTrans[st.focusIndex].trans
	}
}

// ボタンのフォーカス状態を設定
func (st *MainMenuState) setButtonFocus(btn *widget.Button, focused bool) {
	if btn == nil {
		return
	}
	
	// EbitenUIのフォーカス機能を使用
	btn.Focus(focused)
}

var mainMenuTrans = []struct {
	label string
	trans states.Transition
}{
	{
		label: "導入",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&IntroState{}}},
	},
	{
		label: "拠点",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}},
	},
	{
		label: "探検",
		trans: states.Transition{Type: states.TransSwitch, NewStates: []states.State{&DungeonState{Depth: 1}}},
	},
	{
		label: "終了",
		trans: states.Transition{Type: states.TransQuit},
	},
}
