package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// raycast用move
func MoveSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	var pos *gc.Position
	var spriteRender *ec.SpriteRender
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.Player,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos = gameComponents.Position.Get(entity).(*gc.Position)
		spriteRender = gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
	}))

	// オブジェクトの位置関係によっては、進めないこともある。そのときに戻す用
	originalX := pos.X
	originalY := pos.Y

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		pos.X += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		pos.Y += 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		pos.X -= 2
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		pos.Y -= 2
	}

	// 衝突判定して、衝突していれば元の位置に戻す
	// 2つの矩形を比較して、重複する部分があれば衝突とみなす
	{
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		x1 := float64(pos.X - sprite.Width/2)
		x2 := float64(pos.X + sprite.Width/2)
		y1 := float64(pos.Y - sprite.Height/2)
		y2 := float64(pos.Y + sprite.Height/2)

		world.Manager.Join(
			gameComponents.Position,
			gameComponents.SpriteRender,
			gameComponents.BlockPass,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			if !entity.HasComponent(gameComponents.Player) {
				objectPos := gameComponents.Position.Get(entity).(*gc.Position)
				objectSpriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]

				objectx1 := float64(objectPos.X - objectSprite.Width/2)
				objectx2 := float64(objectPos.X + objectSprite.Width/2)
				objecty1 := float64(objectPos.Y - objectSprite.Height/2)
				objecty2 := float64(objectPos.Y + objectSprite.Height/2)

				if (math.Max(x1, objectx1) < math.Min(x2, objectx2)) && (math.Max(y1, objecty1) < math.Min(y2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			}
		}))
	}

	// カメラの追従
	{
		var cPos *gc.Position
		world.Manager.Join(
			gameComponents.Camera,
			gameComponents.Position,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			cPos = gameComponents.Position.Get(entity).(*gc.Position)
		}))
		cPos.X = pos.X
		cPos.Y = pos.Y
	}

	padding := 20
	screenWidth := world.Resources.ScreenDimensions.Width
	screenHeight := world.Resources.ScreenDimensions.Height

	// // +1/-1 is to stop player before it reaches the border
	if pos.X >= screenWidth-padding {
		pos.X = screenWidth - padding - 1
	}

	if pos.X <= padding {
		pos.X = padding + 1
	}

	if pos.Y >= screenHeight-padding {
		pos.Y = screenHeight - padding - 1
	}

	if pos.Y <= padding {
		pos.Y = padding + 1
	}
}
