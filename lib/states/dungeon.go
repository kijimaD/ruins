package states

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/input"
	"github.com/kijimaD/ruins/lib/inputmapper"
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/mapspawner"
	"github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/systems"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
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
	BuilderType mapplanner.PlannerType
}

func (st DungeonState) String() string {
	return "Dungeon"
}

// State interface ================

var _ es.State[w.World] = &DungeonState{}
var _ es.ActionHandler[w.World] = &DungeonState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *DungeonState) OnPause(_ w.World) error { return nil }

// OnResume はステートが再開される際に呼ばれる
func (st *DungeonState) OnResume(_ w.World) error { return nil }

// OnStart はステートが開始される際に呼ばれる
func (st *DungeonState) OnStart(world w.World) error {
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

	// 計画作成する
	plan, err := mapplanner.Plan(world, consts.MapTileWidth, consts.MapTileHeight, st.Seed, st.BuilderType)
	if err != nil {
		return err
	}
	// スポーンする
	level, err := mapspawner.Spawn(world, plan)
	if err != nil {
		return err
	}
	world.Resources.Dungeon.Level = level

	// プレイヤー位置を取得する
	playerX, playerY, hasPlayerPos := plan.GetPlayerStartPosition()
	if !hasPlayerPos {
		return fmt.Errorf("プレイヤー開始位置が設定されていません")
	}
	// プレイヤーを配置する
	if err := worldhelper.MovePlayerToPosition(world, playerX, playerY); err != nil {
		return err
	}

	// フロア移動時に探索済みマップをリセット
	world.Resources.Dungeon.ExploredTiles = make(map[gc.GridElement]bool)

	// 視界キャッシュをクリア（新しい階のために）
	gs.ClearVisionCaches()

	// 初回の冒険開始時のみ操作ガイドを表示
	if st.BuilderType.Name == mapplanner.PlannerTypeTown.Name {
		gamelog.New(gamelog.FieldLog).
			System("WASD: 移動する。").
			Log()
		gamelog.New(gamelog.FieldLog).
			System("Mキー: メニューを開く。").
			Log()
	}
	return nil
}

// OnStop はステートが停止される際に呼ばれる
func (st *DungeonState) OnStop(world w.World) error {
	world.Manager.Join(
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// プレイヤーエンティティ、バックパック内アイテム、装備中アイテムは次のフロアでも必要なので削除しない
		if !entity.HasComponent(world.Components.Player) &&
			!entity.HasComponent(world.Components.ItemLocationInBackpack) &&
			!entity.HasComponent(world.Components.ItemLocationEquipped) {
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
	return nil
}

// Update はゲームステートの更新処理を行う
func (st *DungeonState) Update(world w.World) (es.Transition[w.World], error) {
	// キー入力をActionに変換
	if action, ok := st.HandleInput(); ok {
		if transition, err := st.DoAction(world, action); err != nil {
			return es.Transition[w.World]{}, err
		} else if transition.Type != es.TransNone {
			return transition, nil
		}
	}

	if err := gs.TurnSystem(world); err != nil {
		return es.Transition[w.World]{}, err
	}
	// 移動処理の後にカメラ更新
	gs.CameraSystem(world)

	// プレイヤー死亡チェック
	if st.checkPlayerDeath(world) {
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewGameOverMessageState}}, nil
	}

	// StateEvent処理をチェック
	if transition := st.handleStateEvent(world); transition.Type != es.TransNone {
		return transition, nil
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition(), nil
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) error {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.VisionSystem(world, screen)
	gs.HUDSystem(world, screen) // HUD systemでメッセージも描画
	return nil
}

// ================

// HandleInput はキー入力をActionに変換する
func (st *DungeonState) HandleInput() (inputmapper.ActionID, bool) {
	keyboardInput := input.GetSharedKeyboardInput()

	// メニューキー（M）でダンジョンメニューを開く
	if keyboardInput.IsKeyJustPressed(ebiten.KeyM) {
		return inputmapper.ActionOpenDungeonMenu, true
	}

	cfg := config.MustGet()
	if cfg.Debug && keyboardInput.IsKeyJustPressed(ebiten.KeySlash) {
		return inputmapper.ActionOpenDebugMenu, true
	}

	// 8方向移動キー入力
	if keyboardInput.IsKeyJustPressed(ebiten.KeyW) || keyboardInput.IsKeyJustPressed(ebiten.KeyUp) {
		if keyboardInput.IsKeyJustPressed(ebiten.KeyA) || keyboardInput.IsKeyJustPressed(ebiten.KeyLeft) {
			return inputmapper.ActionMoveNorthWest, true
		}
		if keyboardInput.IsKeyJustPressed(ebiten.KeyD) || keyboardInput.IsKeyJustPressed(ebiten.KeyRight) {
			return inputmapper.ActionMoveNorthEast, true
		}
		return inputmapper.ActionMoveNorth, true
	}
	if keyboardInput.IsKeyJustPressed(ebiten.KeyS) || keyboardInput.IsKeyJustPressed(ebiten.KeyDown) {
		if keyboardInput.IsKeyJustPressed(ebiten.KeyA) || keyboardInput.IsKeyJustPressed(ebiten.KeyLeft) {
			return inputmapper.ActionMoveSouthWest, true
		}
		if keyboardInput.IsKeyJustPressed(ebiten.KeyD) || keyboardInput.IsKeyJustPressed(ebiten.KeyRight) {
			return inputmapper.ActionMoveSouthEast, true
		}
		return inputmapper.ActionMoveSouth, true
	}
	if keyboardInput.IsKeyJustPressed(ebiten.KeyA) || keyboardInput.IsKeyJustPressed(ebiten.KeyLeft) {
		return inputmapper.ActionMoveWest, true
	}
	if keyboardInput.IsKeyJustPressed(ebiten.KeyD) || keyboardInput.IsKeyJustPressed(ebiten.KeyRight) {
		return inputmapper.ActionMoveEast, true
	}

	// 待機キー
	if keyboardInput.IsKeyJustPressed(ebiten.KeyPeriod) {
		return inputmapper.ActionWait, true
	}

	// 相互作用キー（Enter）
	if keyboardInput.IsKeyJustPressed(ebiten.KeyEnter) {
		return inputmapper.ActionInteract, true
	}

	return "", false
}

// DoAction はActionを実行する
func (st *DungeonState) DoAction(world w.World, action inputmapper.ActionID) (es.Transition[w.World], error) {
	// UI系アクションは常に実行可能
	switch action {
	case inputmapper.ActionOpenDungeonMenu, inputmapper.ActionOpenDebugMenu, inputmapper.ActionOpenInventory:
		// UI系はターンチェック不要
	default:
		// ゲーム内アクション（移動、攻撃など）はターンチェックが必要
		if world.Resources.TurnManager != nil {
			turnManager := world.Resources.TurnManager.(*turns.TurnManager)
			if !turnManager.CanPlayerAct() {
				return es.Transition[w.World]{Type: es.TransNone}, nil
			}
		}
	}

	switch action {
	// UI系アクション（ステート遷移）
	case inputmapper.ActionOpenDungeonMenu:
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonMenuState}}, nil
	case inputmapper.ActionOpenDebugMenu:
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewDebugMenuState}}, nil
	case inputmapper.ActionOpenInventory:
		return es.Transition[w.World]{Type: es.TransPush, NewStateFuncs: []es.StateFactory[w.World]{NewInventoryMenuState}}, nil

	// 移動系アクション（World状態を変更）
	case inputmapper.ActionMoveNorth:
		gs.ExecuteMoveAction(world, gc.DirectionUp)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveSouth:
		gs.ExecuteMoveAction(world, gc.DirectionDown)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveEast:
		gs.ExecuteMoveAction(world, gc.DirectionRight)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveWest:
		gs.ExecuteMoveAction(world, gc.DirectionLeft)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveNorthEast:
		gs.ExecuteMoveAction(world, gc.DirectionUpRight)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveNorthWest:
		gs.ExecuteMoveAction(world, gc.DirectionUpLeft)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveSouthEast:
		gs.ExecuteMoveAction(world, gc.DirectionDownRight)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionMoveSouthWest:
		gs.ExecuteMoveAction(world, gc.DirectionDownLeft)
		return es.Transition[w.World]{Type: es.TransNone}, nil
	case inputmapper.ActionWait:
		gs.ExecuteWaitAction(world)
		return es.Transition[w.World]{Type: es.TransNone}, nil

	// 相互作用系アクション
	case inputmapper.ActionInteract:
		gs.ExecuteEnterAction(world)
		return es.Transition[w.World]{Type: es.TransNone}, nil

	default:
		// 未知のActionの場合は何もしない
		return es.Transition[w.World]{Type: es.TransNone}, nil
	}
}

// ================

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
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonState(world.Resources.Dungeon.Depth + 1)}}
	case resources.StateEventWarpEscape:
		return es.Transition[w.World]{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory[w.World]{NewDungeonState(1, WithBuilderType(mapplanner.PlannerTypeTown))}}
	default:
		// StateEventNoneまたは未知のイベントの場合は何もしない
		return es.Transition[w.World]{Type: es.TransNone}
	}
}
