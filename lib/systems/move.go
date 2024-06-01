package systems

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

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

	const speed = 3

	// 元の画像を0度(時計の12時の位置スタート)として、何度回転させるか
	switch {
	// Right
	case ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		pos.X += speed
		pos.Angle = math.Pi / 2
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			pos.Y -= speed
			pos.Angle = math.Pi / 4
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			pos.Y += speed
			pos.Angle = 3 * math.Pi / 4
		}
	// Down
	case ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		pos.Y += speed
		pos.Angle = math.Pi
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			pos.X -= speed
			pos.Angle = 5 * math.Pi / 4
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			pos.X += speed
			pos.Angle = 3 * math.Pi / 4
		}
	// Left
	case ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		pos.X -= speed
		pos.Angle = 3 * math.Pi / 2
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			pos.Y -= speed
			pos.Angle = 7 * math.Pi / 4
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			pos.Y += speed
			pos.Angle = 5 * math.Pi / 4
		}
	// Up
	case ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		pos.Y -= speed
		pos.Angle = math.Pi * 2
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			pos.X -= speed
			pos.Angle = 7 * math.Pi / 4
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			pos.X += speed
			pos.Angle = math.Pi / 4
		}
	}

	// 移動した場合、衝突判定して衝突していれば元の位置に戻す
	// 2つの矩形を比較して、重複する部分があれば衝突とみなす
	if pos.X != originalX || pos.Y != originalY {
		sprite := spriteRender.SpriteSheet.Sprites[spriteRender.SpriteNumber]
		padding := 4 // 1マスの道を進みやすくする
		x1 := float64(pos.X - sprite.Width/2 + padding)
		x2 := float64(pos.X + sprite.Width/2 - padding)
		y1 := float64(pos.Y - sprite.Height/2 + padding)
		y2 := float64(pos.Y + sprite.Height/2 - padding)

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
				if (math.Max(x1, objectx1) < math.Min(x2, objectx2)) && (math.Max(y1, objecty1) < math.Min(y2, objecty2)) {
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
				if (math.Max(x1, objectx1) < math.Min(x2, objectx2)) && (math.Max(y1, objecty1) < math.Min(y2, objecty2)) {
					// 衝突していれば元の位置に戻す
					pos.X = originalX
					pos.Y = originalY
				}
			}
		}))
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
		cPos.X = pos.X
		cPos.Y = pos.Y

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

	padding := 20
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
				gameResources.StateEvent = resources.StateEventWarpNext
			case gc.WarpModeEscape:
				gameResources.StateEvent = resources.StateEventWarpEscape
			}

		}
	}
}
