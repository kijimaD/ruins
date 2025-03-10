package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils/mathutil"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func PlayerInputSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	var playerPos *gc.Position
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Operator,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerPos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	const maxSpeed = 2.0
	const minSpeed = -1.0
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW):
		playerPos.Speed += 0.1
		playerPos.Speed = mathutil.Min(maxSpeed, playerPos.Speed)
	case ebiten.IsKeyPressed(ebiten.KeyS):
		playerPos.Speed -= 0.1
		playerPos.Speed = mathutil.Max(minSpeed, playerPos.Speed)
	case ebiten.IsKeyPressed(ebiten.KeyD):
		playerPos.Angle += math.Pi / 90
	case ebiten.IsKeyPressed(ebiten.KeyA):
		playerPos.Angle -= math.Pi / 90
	}

	{
		// カメラの追従
		var camera *gc.Camera
		var cPos *gc.Position
		world.Manager.Join(
			gameComponents.Camera,
			gameComponents.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			camera = gameComponents.Camera.Get(entity).(*gc.Camera)
			cPos = gameComponents.Position.Get(entity).(*gc.Position)
		}))
		cPos.X = playerPos.X
		cPos.Y = playerPos.Y

		// ズーム率変更
		// 参考: https://ebitengine.org/ja/examples/isometric.html
		var scrollY float64
		if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
			scrollY = -0.25
		} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
			scrollY = 0.25
		} else {
			_, scrollY = ebiten.Wheel()
			if scrollY < -1 {
				scrollY = -1
			} else if scrollY > 1 {
				scrollY = 1
			}
		}
		camera.ScaleTo += scrollY * (camera.ScaleTo / 7)

		// Clamp target zoom level.
		if camera.ScaleTo < 0.8 {
			camera.ScaleTo = 0.8
		} else if camera.ScaleTo > 10 {
			camera.ScaleTo = 10
		}

		// Smooth zoom transition.
		div := 10.0
		if camera.ScaleTo > camera.Scale {
			camera.Scale += (camera.ScaleTo - camera.Scale) / div
		} else if camera.ScaleTo < camera.Scale {
			camera.Scale -= (camera.Scale - camera.ScaleTo) / div
		}
	}
}
