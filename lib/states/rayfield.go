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
	"github.com/kijimaD/ruins/lib/spawner"
	gs "github.com/kijimaD/ruins/lib/systems"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var (
	baseImage *ebiten.Image // 一番下にある黒背景
	bgImage   *ebiten.Image // 床を表現する
)

type RayFieldState struct{}

func (st RayFieldState) String() string {
	return "RayField"
}

// State interface ================

func (st *RayFieldState) OnPause(world w.World) {}

func (st *RayFieldState) OnResume(world w.World) {}

func (st *RayFieldState) OnStart(world w.World) {
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
}

func (st *RayFieldState) OnStop(world w.World) {}

func (st *RayFieldState) Update(world w.World) states.Transition {
	gs.MoveRaySystem(world)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{Type: states.TransPush, NewStates: []states.State{&FieldMenuState{}}}
	}

	return states.Transition{}
}

func (st *RayFieldState) Draw(world w.World, screen *ebiten.Image) {
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	screen.DrawImage(baseImage, nil)
	{
		tileWidth, tileHeight := bgImage.Size()
		// 背景画像を敷き詰める
		for i := 0; i < screenWidth; i += tileWidth {
			for j := 0; j < screenHeight; j += tileHeight {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i), float64(j))
				screen.DrawImage(bgImage, op)
			}
		}
	}

	gs.RenderObjectSystem(world, screen)
	gs.RenderVisionSystem(world, screen)
	gs.RenderShadowSystem(world, screen)
}
