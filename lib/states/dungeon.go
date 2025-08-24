package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/systems"
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
	baseImage = ebiten.NewImage(screenWidth, screenHeight)
	baseImage.Fill(color.Black)

	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.Depth = st.Depth

	// seed が 0 の場合は NewLevel 内部でランダムシードが生成される
	gameResources.Level = mapbuilder.NewLevel(world, 50, 50, st.Seed)

	// フロア移動時に探索済みマップをリセット
	gameResources.ExploredTiles = make(map[string]bool)

	// 視界キャッシュをクリア（新しい階のために）
	gs.ClearVisionCaches()
}

// OnStop はステートが停止される際に呼ばれる
func (st *DungeonState) OnStop(world w.World) {
	world.Manager.Join(
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	world.Manager.Join(
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	world.Manager.Join(
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))

	// reset
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	gameResources.SetStateEvent(resources.StateEventNone)

	// 視界キャッシュをクリア
	gs.ClearVisionCaches()
}

// Update はゲームステートの更新処理を行う
func (st *DungeonState) Update(world w.World) es.Transition {
	gs.PlayerInputSystem(world)
	gs.AIInputSystem(world)
	gs.MoveSystem(world)
	gs.CollisionSystem(world)

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

// handleStateEvent はStateEventを処理し、対応する遷移を返す
func (st *DungeonState) handleStateEvent(world w.World) es.Transition {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	switch gameResources.ConsumeStateEvent() {
	case resources.StateEventWarpNext:
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewDungeonStateWithDepth(gameResources.Depth + 1)}}
	case resources.StateEventWarpEscape:
		return es.Transition{Type: es.TransSwitch, NewStateFuncs: []es.StateFactory{NewHomeMenuState}}
	case resources.StateEventBattleStart:
		// 戦闘開始
		battleStateFactory := func() es.State {
			// プレイヤーエンティティを検索
			var playerEntity ecs.Entity
			world.Manager.Join(
				world.Components.Position,
				world.Components.Operator,
			).Visit(ecs.Visit(func(entity ecs.Entity) {
				playerEntity = entity
			}))

			// 最も近い敵エンティティを検索
			var fieldEnemyEntity ecs.Entity
			var minDistance float64 = -1
			if playerEntity != ecs.Entity(0) {
				playerPos := world.Components.Position.Get(playerEntity).(*gc.Position)
				world.Manager.Join(
					world.Components.Position,
					world.Components.AIMoveFSM,
				).Visit(ecs.Visit(func(entity ecs.Entity) {
					enemyPos := world.Components.Position.Get(entity).(*gc.Position)
					dx := float64(playerPos.X - enemyPos.X)
					dy := float64(playerPos.Y - enemyPos.Y)
					distance := dx*dx + dy*dy // 距離の2乗で比較（平方根計算を省略）

					if minDistance < 0 || distance < minDistance {
						minDistance = distance
						fieldEnemyEntity = entity
					}
				}))
			}

			return &BattleState{
				FieldEnemyEntity: fieldEnemyEntity,
			}
		}
		return es.Transition{Type: es.TransPush, NewStateFuncs: []es.StateFactory{battleStateFactory}}
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
	gs.HUDSystem(world, screen)
}
