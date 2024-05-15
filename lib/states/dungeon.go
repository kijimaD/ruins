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
	"github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/spawner"
	gs "github.com/kijimaD/ruins/lib/systems"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	baseImage *ebiten.Image // 一番下にある黒背景
	bgImage   *ebiten.Image // 床を表現する
)

type DungeonState struct{}

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

	// debug
	// 終了したとき、現状は前回からの続きになっている
	// 終了や階層移動したときはフィールドにあるものを削除したい
	playerCount := 0
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerCount++
	}))
	if playerCount == 0 {
		spawner.SpawnPlayer(world, 200, 200)
		spawner.SpawnFieldWall(world, 240, 200)
		spawner.SpawnFieldWall(world, 320, 200)
		spawner.SpawnFieldWall(world, 352, 200)
		spawner.SpawnFieldWarpNext(world, 300, 300)
	}

	world.Resources.Game = &resources.Game{}
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.Level = resources.NewLevel(1, 50, 50)
}

func (st *DungeonState) OnStop(world w.World) {}

func (st *DungeonState) Update(world w.World) states.Transition {
	gs.MoveSystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&DungeonMenuState{}}}
	}

	return states.Transition{}
}

func (st *DungeonState) Draw(world w.World, screen *ebiten.Image) {
	screen.DrawImage(baseImage, nil)

	gs.RenderSpriteSystem(world, screen)
	gs.DarknessSystem(world, screen)
	gs.BlindSpotSystem(world, screen)
}
