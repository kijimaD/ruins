package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func PlayerMoveSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	var playerEntity ecs.Entity
	var playerPos *gc.Position // player position
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Player,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = entity
		playerPos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	const distance = 3
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyD):
		tryMove(world, playerEntity, 45, distance)
	case ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyA):
		tryMove(world, playerEntity, 135, distance)
	case ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyA):
		tryMove(world, playerEntity, 225, distance)
	case ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyD):
		tryMove(world, playerEntity, 315, distance)
	case ebiten.IsKeyPressed(ebiten.KeyD):
		tryMove(world, playerEntity, 0, distance)
	case ebiten.IsKeyPressed(ebiten.KeyW):
		tryMove(world, playerEntity, 90, distance)
	case ebiten.IsKeyPressed(ebiten.KeyA):
		tryMove(world, playerEntity, 180, distance)
	case ebiten.IsKeyPressed(ebiten.KeyS):
		tryMove(world, playerEntity, 270, distance)
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

// 角度と距離を指定して相対移動させる
func tryMove(world w.World, entity ecs.Entity, angle float64, distance float64) {
	gameComponents := world.Components.Game.(*gc.Components)

	pos := gameComponents.Position.Get(entity).(*gc.Position) // player pos
	spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

	originalX := pos.X
	originalY := pos.Y

	radians := angle * math.Pi / 180.0
	pos.X += int(math.Cos(radians) * distance)
	pos.Y -= int(math.Sin(radians) * distance) // ゲーム画面におけるy軸は上がマイナスであるので逆転させる

	{
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		padding := 4 // 1マスの道を進みやすくする
		playerx1 := float64(pos.X - sprite.Width/2 + padding)
		playerx2 := float64(pos.X + sprite.Width/2 - padding)
		playery1 := float64(pos.Y - sprite.Height/2 + padding)
		playery2 := float64(pos.Y + sprite.Height/2 - padding)

		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.BlockPass,
			gameComponents.Player.Not(),
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			switch {
			case entity.HasComponent(gameComponents.Position):
				objectPos := gameComponents.Position.Get(entity).(*gc.Position)
				objectSpriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]

				objectx1 := float64(objectPos.X - objectSprite.Width/2)
				objectx2 := float64(objectPos.X + objectSprite.Width/2)
				objecty1 := float64(objectPos.Y - objectSprite.Height/2)
				objecty2 := float64(objectPos.Y + objectSprite.Height/2)
				if (math.Max(playerx1, objectx1) < math.Min(playerx2, objectx2)) && (math.Max(playery1, objecty1) < math.Min(playery2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			case entity.HasComponent(gameComponents.GridElement):
				objectGrid := gameComponents.GridElement.Get(entity).(*gc.GridElement)
				objectSpriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]
				x := int(objectGrid.Row) * sprite.Width
				y := int(objectGrid.Col) * sprite.Height
				objectx1 := float64(x)
				objectx2 := float64(x + objectSprite.Width)
				objecty1 := float64(y)
				objecty2 := float64(y + objectSprite.Height)
				if (math.Max(playerx1, objectx1) < math.Min(playerx2, objectx2)) && (math.Max(playery1, objecty1) < math.Min(playery2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			}
		}))
	}

	padding := 10
	gameResources := world.Resources.Game.(*resources.Game)
	levelWidth := gameResources.Level.Width()
	levelHeight := gameResources.Level.Height()

	// +1/-1 is to stop player before it reaches the border
	if pos.X >= levelWidth-padding {
		pos.X = levelWidth - padding - 1
	}

	if pos.X <= padding {
		pos.X = padding + 1
	}

	if pos.Y >= levelHeight-padding {
		pos.Y = levelHeight - padding - 1
	}

	if pos.Y <= padding {
		pos.Y = padding + 1
	}

	// タイルイベントを発行する
	{
		gameResources := world.Resources.Game.(*resources.Game)
		entity := gameResources.Level.AtEntity(pos.X, pos.Y)

		gameComponents := world.Components.Game.(*gc.Components)
		if entity.HasComponent(gameComponents.Warp) {
			warp := gameComponents.Warp.Get(entity).(*gc.Warp)
			switch warp.Mode {
			case gc.WarpModeNext:
				effects.AddEffect(nil, effects.WarpNext{}, effects.None{})
			case gc.WarpModeEscape:
				effects.AddEffect(nil, effects.WarpEscape{}, effects.None{})
			}
		}
	}
}
