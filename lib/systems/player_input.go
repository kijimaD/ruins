package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// PlayerInputSystem はプレイヤーからの入力を処理する
func PlayerInputSystem(world w.World) {

	var playerVelocity *gc.Velocity
	var playerPos *gc.Position
	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.Operator,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerVelocity = world.Components.Velocity.Get(entity).(*gc.Velocity)
		playerPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	// デフォルト
	playerVelocity.ThrottleMode = gc.ThrottleModeNope
	// 同時押しがありうる
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		playerVelocity.ThrottleMode = gc.ThrottleModeFront
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		playerVelocity.ThrottleMode = gc.ThrottleModeBack
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		playerVelocity.Angle += math.Pi / 90
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		playerVelocity.Angle -= math.Pi / 90
	}

	{
		// カメラの追従
		var camera *gc.Camera
		var cPos *gc.Position
		world.Manager.Join(
			world.Components.Camera,
			world.Components.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			camera = world.Components.Camera.Get(entity).(*gc.Camera)
			cPos = world.Components.Position.Get(entity).(*gc.Position)
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
