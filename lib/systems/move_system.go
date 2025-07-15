package systems

import (
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// MoveSystem はエンティティの移動処理を行う
func MoveSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	maxFrontSpeed := 2.0
	maxBackSpeed := -1.0
	accelerationSpeed := 0.05
	world.Manager.Join(
		gameComponents.Velocity,
		gameComponents.Position,
		gameComponents.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := gameComponents.Velocity.Get(entity).(*gc.Velocity)
		switch velocity.ThrottleMode {
		case gc.ThrottleModeFront:
			velocity.Speed = utils.Min(maxFrontSpeed, velocity.Speed+accelerationSpeed)
		case gc.ThrottleModeBack:
			velocity.Speed = utils.Max(maxBackSpeed, velocity.Speed-accelerationSpeed)
		case gc.ThrottleModeNope:
			// 何もしない
		}
		tryMove(world, entity, velocity.Angle, velocity.Speed)
	}))

	// 操作キャラに対してタイルイベントを発行する
	world.Manager.Join(
		gameComponents.Position,
		gameComponents.SpriteRender,
		gameComponents.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := gameComponents.Position.Get(entity).(*gc.Position)
		gameResources := world.Resources.Game.(*resources.Game)
		tileEntity := gameResources.Level.AtEntity(pos.X, pos.Y)

		if tileEntity.HasComponent(gameComponents.Warp) {
			warp := gameComponents.Warp.Get(tileEntity).(*gc.Warp)
			switch warp.Mode {
			case gc.WarpModeNext:
				effects.AddEffect(nil, effects.WarpNext{}, effects.None{})
			case gc.WarpModeEscape:
				effects.AddEffect(nil, effects.WarpEscape{}, effects.None{})
			}
		}
	}))
}

// 角度と距離を指定して相対移動させる
func tryMove(world w.World, entity ecs.Entity, angle float64, distance float64) {
	gameComponents := world.Components.Game.(*gc.Components)

	pos := gameComponents.Position.Get(entity).(*gc.Position)
	spriteRender := gameComponents.SpriteRender.Get(entity).(*ec.SpriteRender)

	originalX := pos.X
	originalY := pos.Y
	radian := angle + math.Pi/2 // 度数法 -> 弧度法
	pos.X = gc.Pixel(float64(pos.X) - math.Cos(radian)*distance)
	pos.Y = gc.Pixel(float64(pos.Y) - math.Sin(radian)*distance)

	{
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		padding := 4 // 1マスの道を進みやすくする
		x1 := float64(int(pos.X) - sprite.Width/2 + padding)
		x2 := float64(int(pos.X) + sprite.Width/2 - padding)
		y1 := float64(int(pos.Y) - sprite.Height/2 + padding)
		y2 := float64(int(pos.Y) + sprite.Height/2 - padding)

		world.Manager.Join(
			gameComponents.SpriteRender,
			gameComponents.BlockPass,
		).Visit(ecs.Visit(func(entityAnother ecs.Entity) {
			if entity == entityAnother {
				return
			}
			switch {
			case entityAnother.HasComponent(gameComponents.Position):
				objectPos := gameComponents.Position.Get(entityAnother).(*gc.Position)
				objectSpriteRender := gameComponents.SpriteRender.Get(entityAnother).(*ec.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]

				objectx1 := float64(int(objectPos.X) - objectSprite.Width/2)
				objectx2 := float64(int(objectPos.X) + objectSprite.Width/2)
				objecty1 := float64(int(objectPos.Y) - objectSprite.Height/2)
				objecty2 := float64(int(objectPos.Y) + objectSprite.Height/2)
				if (math.Max(x1, objectx1) < math.Min(x2, objectx2)) && (math.Max(y1, objecty1) < math.Min(y2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			case entityAnother.HasComponent(gameComponents.GridElement):
				objectGrid := gameComponents.GridElement.Get(entityAnother).(*gc.GridElement)
				objectSpriteRender := gameComponents.SpriteRender.Get(entityAnother).(*ec.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]
				x := int(objectGrid.Row) * sprite.Width
				y := int(objectGrid.Col) * sprite.Height
				objectx1 := float64(x)
				objectx2 := float64(x + objectSprite.Width)
				objecty1 := float64(y)
				objecty2 := float64(y + objectSprite.Height)
				if (math.Max(x1, objectx1) < math.Min(x2, objectx2)) && (math.Max(y1, objecty1) < math.Min(y2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			}
		}))
	}

	padding := gc.Pixel(10)
	gameResources := world.Resources.Game.(*resources.Game)
	levelWidth := gameResources.Level.Width()
	levelHeight := gameResources.Level.Height()

	if pos.X >= gc.Pixel(levelWidth-padding) {
		pos.X = gc.Pixel(levelWidth - padding - 1)
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
}
