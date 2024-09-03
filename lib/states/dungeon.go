package states

import (
	"bytes"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/systems"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	baseImage *ebiten.Image // 一番下にある黒背景
	bgImage   *ebiten.Image // 床を表現する
)

type DungeonState struct {
	Depth int
}

func (st DungeonState) String() string {
	return "Dungeon"
}

// State interface ================

func (st *DungeonState) OnPause(world w.World) {}

func (st *DungeonState) OnResume(world w.World) {}

func (st *DungeonState) OnStart(world w.World) {
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height
	baseImage = ebiten.NewImage(screenWidth, screenHeight)
	img, _, err := image.Decode(bytes.NewReader(images.Tile_png))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(img)
	baseImage.Fill(color.Black)

	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.Depth = st.Depth
	gameResources.Level = mapbuilder.NewLevel(world, 50, 50)
}

func (st *DungeonState) OnStop(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	world.Manager.Join(
		gameComponents.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))
	world.Manager.Join(
		gameComponents.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		world.Manager.DeleteEntity(entity)
	}))

	// reset
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventNone
}

func (st *DungeonState) Update(world w.World) states.Transition {
	gs.PlayerMoveSystem(world)
	effects.RunEffectQueue(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonMenuState{}}}
	}

	gameResources := world.Resources.Game.(*resources.Game)
	switch gameResources.StateEvent {
	case resources.StateEventWarpNext:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&DungeonState{Depth: gameResources.Depth + 1}}}
	case resources.StateEventWarpEscape:
		return states.Transition{Type: states.TransSwitch, NewStates: []states.State{&HomeMenuState{}}}
	}

	return states.Transition{}
}

func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.DarknessSystem(world, screen)
	gs.BlindSpotSystem(world, screen)
	gs.HUDSystem(world, screen)
}
