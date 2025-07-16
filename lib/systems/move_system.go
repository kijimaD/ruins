package systems

import (
	"math"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/mathutil"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// MoveSystem はエンティティの移動処理を行う
func MoveSystem(world w.World) {

	maxFrontSpeed := 2.0
	maxBackSpeed := -1.0
	accelerationSpeed := 0.05
	world.Manager.Join(
		world.Components.Velocity,
		world.Components.Position,
		world.Components.SpriteRender,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		velocity := world.Components.Velocity.Get(entity).(*gc.Velocity)
		switch velocity.ThrottleMode {
		case gc.ThrottleModeFront:
			velocity.Speed = mathutil.Min(maxFrontSpeed, velocity.Speed+accelerationSpeed)
		case gc.ThrottleModeBack:
			velocity.Speed = mathutil.Max(maxBackSpeed, velocity.Speed-accelerationSpeed)
		case gc.ThrottleModeNope:
			// 何もしない
		}
		tryMove(world, entity, velocity.Angle, velocity.Speed)
	}))

	// 操作キャラに対してタイルイベントを発行する
	world.Manager.Join(
		world.Components.Position,
		world.Components.SpriteRender,
		world.Components.Operator,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pos := world.Components.Position.Get(entity).(*gc.Position)
		gameResources := world.Resources.Game.(*resources.Game)
		tileEntity := gameResources.Level.AtEntity(pos.X, pos.Y)

		if tileEntity.HasComponent(world.Components.Warp) {
			warp := world.Components.Warp.Get(tileEntity).(*gc.Warp)
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

	pos := world.Components.Position.Get(entity).(*gc.Position)
	spriteRender := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)

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
			world.Components.SpriteRender,
			world.Components.BlockPass,
		).Visit(ecs.Visit(func(entityAnother ecs.Entity) {
			if entity == entityAnother {
				return
			}
			switch {
			case entityAnother.HasComponent(world.Components.Position):
				objectPos := world.Components.Position.Get(entityAnother).(*gc.Position)
				objectSpriteRender := world.Components.SpriteRender.Get(entityAnother).(*gc.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]

				objectx1 := float64(int(objectPos.X) - objectSprite.Width/2)
				objectx2 := float64(int(objectPos.X) + objectSprite.Width/2)
				objecty1 := float64(int(objectPos.Y) - objectSprite.Height/2)
				objecty2 := float64(int(objectPos.Y) + objectSprite.Height/2)
				if (max(x1, objectx1) < min(x2, objectx2)) && (max(y1, objecty1) < min(y2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			case entityAnother.HasComponent(world.Components.GridElement):
				objectGrid := world.Components.GridElement.Get(entityAnother).(*gc.GridElement)
				objectSpriteRender := world.Components.SpriteRender.Get(entityAnother).(*gc.SpriteRender)
				objectSprite := spriteRender.SpriteSheet.Sprites[objectSpriteRender.SpriteNumber]
				x := int(objectGrid.Row) * sprite.Width
				y := int(objectGrid.Col) * sprite.Height
				objectx1 := float64(x)
				objectx2 := float64(x + objectSprite.Width)
				objecty1 := float64(y)
				objecty2 := float64(y + objectSprite.Height)
				if (max(x1, objectx1) < min(x2, objectx2)) && (max(y1, objecty1) < min(y2, objecty2)) {
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
