package states

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/colors"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
)

// GameOverState はゲームオーバーのゲームステート
type GameOverState struct {
	es.BaseState
	ui            *ebitenui.UI
	keyboardInput input.KeyboardInput
}

func (st GameOverState) String() string {
	return "GameOver"
}

// State interface ================

var _ es.State = &GameOverState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *GameOverState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *GameOverState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *GameOverState) OnStart(world w.World) {
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *GameOverState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *GameOverState) Update(_ w.World) es.Transition {
	if st.keyboardInput.IsEnterJustPressedOnce() {
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
	}

	st.ui.Update()

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *GameOverState) Draw(_ w.World, screen *ebiten.Image) {
	// 半透明の黒い背景を描画してダンジョン画面を暗くする
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(colors.TransBlackColor)
	screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

	st.ui.Draw(screen)
}

// ================

func (st *GameOverState) initUI(world w.World) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	windowContainer := styled.NewWindowContainer(world)
	titleContainer := styled.NewWindowHeaderContainer("GAME OVER", world)
	gameOverWindow := styled.NewSmallWindow(
		titleContainer,
		windowContainer,
		widget.WindowOpts.CloseMode(widget.NONE),
	)

	// コンテンツを追加
	gameOverText := styled.NewTitleText("死亡した。", world)
	windowContainer.AddChild(gameOverText)

	instructionText := styled.NewDescriptionText("Enterキーを押してメインメニューに戻る", world)
	windowContainer.AddChild(instructionText)

	// ウィンドウを中央に配置
	gameOverWindow.SetLocation(getCenterWinRect(world))

	ui := &ebitenui.UI{Container: rootContainer}
	ui.AddWindow(gameOverWindow)

	return ui
}
