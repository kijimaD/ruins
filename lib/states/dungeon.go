package states

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/turns"
	"github.com/kijimaD/ruins/lib/widgets/styled"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	baseImage *ebiten.Image // 一番下にある黒背景
)

// DungeonState はダンジョン探索中のゲームステート
type DungeonState struct {
	es.BaseState
	Depth int
	// Seed はマップ生成用のシード値（0の場合はDungeonリソースのシード値を使用）
	Seed uint64
	// BuilderType は使用するマップビルダーのタイプ（BuilderTypeRandomの場合はランダム選択）
	BuilderType mapbuilder.BuilderType
	// ゲームオーバー状態
	gameOver bool
	// UI関連
	ui             *ebitenui.UI
	gameOverWindow *widget.Window
	keyboardInput  input.KeyboardInput
}

func (st DungeonState) String() string {
	return "Dungeon"
}

// State interface ================

var _ es.State = &DungeonState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *DungeonState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *DungeonState) OnResume(_ w.World) {}

// OnStart はステートが開始される際に呼ばれる
func (st *DungeonState) OnStart(world w.World) {
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height
	if screenWidth > 0 && screenHeight > 0 {
		baseImage = ebiten.NewImage(screenWidth, screenHeight)
		baseImage.Fill(color.Black)
	}

	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.Depth = st.Depth

	// キーボード入力を初期化
	if st.keyboardInput == nil {
		st.keyboardInput = input.GetSharedKeyboardInput()
	}

	// ターンマネージャーを初期化
	if world.Resources.TurnManager == nil {
		world.Resources.TurnManager = turns.NewTurnManager()
	}

	// seed が 0 の場合は NewLevel 内部でランダムシードが生成される
	level, err := mapbuilder.NewLevel(world, consts.MapTileWidth, consts.MapTileHeight, st.Seed, st.BuilderType)
	if err != nil {
		panic(err)
	}
	gameResources.Level = level

	// フロア移動時に探索済みマップをリセット
	gameResources.ExploredTiles = make(map[string]bool)

	// プレイヤーのタイル状態をリセット（新しい階のために）
	gameResources.ResetPlayerTileState()

	// 視界キャッシュをクリア（新しい階のために）
	gs.ClearVisionCaches()

	// UI初期化
	st.ui = st.initUI(world)
}

// OnStop はステートが停止される際に呼ばれる
func (st *DungeonState) OnStop(world w.World) {
	world.Manager.Join(
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// プレイヤーエンティティは次のフロアでも必要なので削除しない
		if !entity.HasComponent(world.Components.Player) {
			world.Manager.DeleteEntity(entity)
		}
	}))
	world.Manager.Join(
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// プレイヤーエンティティは次のフロアでも必要なので削除しない
		if !entity.HasComponent(world.Components.Player) {
			world.Manager.DeleteEntity(entity)
		}
	}))

	// reset
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.SetStateEvent(resources.StateEventNone)

	// 視界キャッシュをクリア
	gs.ClearVisionCaches()
}

// Update はゲームステートの更新処理を行う
func (st *DungeonState) Update(world w.World) es.Transition {
	// UI更新
	if st.ui != nil {
		st.ui.Update()
	}

	// ゲームオーバー状態でない場合のみゲームロジックを実行
	if st.gameOver {
		// ゲームオーバー状態：IsEnterJustPressedOnce()でメインメニューに戻る
		if st.keyboardInput.IsEnterJustPressedOnce() {
			return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewMainMenuState}}
		}
		return es.Transition{Type: es.TransNone}
	}

	gs.TurnSystem(world)
	// 移動処理の後にカメラ更新
	gs.CameraSystem(world)

	// プレイヤー死亡チェック
	if st.checkPlayerDeath(world) {
		st.gameOver = true
		st.showGameOverWindow(world) // ウィンドウを表示
		return es.Transition{Type: es.TransNone}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewDungeonMenuState}}
	}

	cfg := config.MustGet()
	if cfg.Debug && inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{NewDebugMenuState}}
	}

	// StateEvent処理をチェック
	if transition := st.handleStateEvent(world); transition.Type != es.TransNone {
		return transition
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// checkPlayerDeath はプレイヤーの死亡状態をチェックする
func (st *DungeonState) checkPlayerDeath(world w.World) bool {
	playerDead := false
	world.Manager.Join(
		world.Components.Player,
		world.Components.Dead,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerDead = true
	}))
	return playerDead
}

// handleStateEvent はStateEventを処理し、対応する遷移を返す
func (st *DungeonState) handleStateEvent(world w.World) es.Transition {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	switch gameResources.ConsumeStateEvent() {
	case resources.StateEventWarpNext:
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(gameResources.Depth + 1)}}
	case resources.StateEventWarpEscape:
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}}
	default:
		// StateEventNoneまたは未知のイベントの場合は何もしない
		return es.Transition{Type: es.TransNone}
	}
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.VisionSystem(world, screen)
	gs.HUDSystem(world, screen) // HUD systemでメッセージも描画

	// UI描画（ゲームオーバーウィンドウなどを含む）
	if st.ui != nil {
		st.ui.Draw(screen)
	}
}

// initUI はUIを初期化する
func (st *DungeonState) initUI(world w.World) *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}

// showGameOverWindow はゲームオーバーウィンドウを表示する（craft_menuスタイル）
func (st *DungeonState) showGameOverWindow(world w.World) {
	windowContainer := styled.NewWindowContainer(world)
	titleContainer := styled.NewWindowHeaderContainer("GAME OVER", world)
	st.gameOverWindow = styled.NewSmallWindow(titleContainer, windowContainer)

	// コンテンツを追加
	gameOverText := styled.NewTitleText("あなたは死んでしまった...", world)
	windowContainer.AddChild(gameOverText)

	instructionText := styled.NewDescriptionText("Enterキーを押してメインメニューに戻る", world)
	windowContainer.AddChild(instructionText)

	// ウィンドウを中央に配置して表示
	st.gameOverWindow.SetLocation(getCenterWinRect(world))
	st.ui.AddWindow(st.gameOverWindow)
}
