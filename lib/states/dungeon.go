package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/levelgen"
	"github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	baseImage *ebiten.Image // 一番下にある黒背景
)

// DungeonState はダンジョン探索中のゲームステート
type DungeonState struct {
	es.BaseState[w.World]
	Depth int
	// Seed はマップ生成用のシード値（0の場合はDungeonリソースのシード値を使用）
	Seed uint64
	// BuilderType は使用するマップビルダーのタイプ（BuilderTypeRandom の場合はランダム選択）
	BuilderType mapplaner.PlannerType
}

func (st DungeonState) String() string {
	return "Dungeon"
}

// State interface ================

var _ es.State[w.World] = &DungeonState{}

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

	world.Resources.Dungeon.Depth = st.Depth

	// ターンマネージャーを初期化
	if world.Resources.TurnManager == nil {
		world.Resources.TurnManager = turns.NewTurnManager()
	}

	// seed が 0 の場合は NewLevel 内部でランダムシードが生成される
	level, err := levelgen.NewLevel(world, consts.MapTileWidth, consts.MapTileHeight, st.Seed, st.BuilderType)
	if err != nil {
		panic(err)
	}
	world.Resources.Dungeon.Level = level

	// フロア移動時に探索済みマップをリセット
	world.Resources.Dungeon.ExploredTiles = make(map[gc.GridElement]bool)

	// 視界キャッシュをクリア（新しい階のために）
	gs.ClearVisionCaches()
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
	world.Resources.Dungeon.SetStateEvent(resources.StateEventNone)

	// 視界キャッシュをクリア
	gs.ClearVisionCaches()
}

// Update はゲームステートの更新処理を行う
func (st *DungeonState) Update(world w.World) es.Transition[w.World] {
	gs.TurnSystem(world)
	// 移動処理の後にカメラ更新
	gs.CameraSystem(world)

	// プレイヤー死亡チェック
	if st.checkPlayerDeath(world) {
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewGameOverMessageState}}
	}

	cfg := config.MustGet()
	if cfg.Debug && inpututil.IsKeyJustPressed(ebiten.KeySlash) {
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewDebugMenuState}}
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
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		playerDead = true
	}))
	return playerDead
}

// handleStateEvent はStateEventを処理し、対応する遷移を返す
func (st *DungeonState) handleStateEvent(world w.World) es.Transition[w.World] {

	switch world.Resources.Dungeon.ConsumeStateEvent() {
	case resources.StateEventWarpNext:
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonStateWithDepth(world.Resources.Dungeon.Depth + 1)}}
	case resources.StateEventWarpEscape:
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewHomeMenuState}}
	default:
		// StateEventNoneまたは未知のイベントの場合は何もしない
		return es.Transition[w.World]{Type: es.TransNone}
	}
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.VisionSystem(world, screen)
	gs.HUDSystem(world, screen) // HUD systemでメッセージも描画
}
