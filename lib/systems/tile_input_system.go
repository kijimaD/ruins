package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TileInputSystem はプレイヤーからのタイルベース入力を処理する
func TileInputSystem(world w.World) {
	// プレイヤーエンティティの移動意図をクリア
	world.Manager.Join(
		world.Components.Player,
		world.Components.TurnBased,
		world.Components.WantsToMove,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		wants := world.Components.WantsToMove.Get(entity).(*gc.WantsToMove)
		wants.Direction = gc.DirectionNone
	}))

	// キー入力を方向に変換
	var direction gc.Direction

	// 8方向キー入力
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			direction = gc.DirectionUpLeft
		} else if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			direction = gc.DirectionUpRight
		} else {
			direction = gc.DirectionUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			direction = gc.DirectionDownLeft
		} else if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			direction = gc.DirectionDownRight
		} else {
			direction = gc.DirectionDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		direction = gc.DirectionLeft
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		direction = gc.DirectionRight
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		// スペースキーまたはピリオドで待機
		direction = gc.DirectionNone
	}

	// 移動意図を設定
	if direction != gc.DirectionNone {
		world.Manager.Join(
			world.Components.Player,
			world.Components.TurnBased,
			world.Components.WantsToMove,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			wants := world.Components.WantsToMove.Get(entity).(*gc.WantsToMove)
			wants.Direction = direction
		}))
	}
}
