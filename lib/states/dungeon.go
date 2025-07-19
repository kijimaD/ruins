package states

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.Depth = st.Depth
	gameResources.Level = mapbuilder.NewLevel(world, 50, 50)
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
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventNone
}

// Update はゲームステートの更新処理を行う
func (st *DungeonState) Update(world w.World) es.Transition {
	gs.PlayerInputSystem(world)
	gs.AIInputSystem(world)
	gs.MoveSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return es.Transition{Type: es.TransPush, NewStates: []es.State{&DungeonMenuState{}}}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventWarpNext:
		return es.Transition{Type: es.TransSwitch, NewStates: []es.State{&DungeonState{Depth: gameResources.Depth + 1}}}
	case resources.StateEventWarpEscape:
		return es.Transition{Type: es.TransSwitch, NewStates: []es.State{&HomeMenuState{}}}
	}

	// BaseStateの共通処理を使用
	return st.ConsumeTransition()
}

// Draw はゲームステートの描画処理を行う
func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.VisionSystem(world, screen)
	gs.HUDSystem(world, screen)
}
