package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TileMoveSystem はタイルベース移動処理を行う
func TileMoveSystem(world w.World) {
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.TurnBased,
		world.Components.WantsToMove,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		wants := world.Components.WantsToMove.Get(entity).(*gc.WantsToMove)

		// 移動意図がない場合はスキップ
		if wants.Direction == gc.DirectionNone {
			return
		}

		// 移動先を計算
		deltaX, deltaY := wants.Direction.GetDelta()
		newX := int(gridElement.X) + deltaX
		newY := int(gridElement.Y) + deltaY

		// 移動可能かチェック
		if canMoveTo(world, newX, newY, entity) {
			gridElement.X = gc.Tile(newX)
			gridElement.Y = gc.Tile(newY)
		}

		// 移動意図をクリア
		wants.Direction = gc.DirectionNone
	}))

	// プレイヤーの移動後にタイルイベントをチェック
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Operator, // プレイヤーであることを示す
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		currentTileX := int(gridElement.X)
		currentTileY := int(gridElement.Y)

		// プレイヤーが新しいタイルに移動した場合のみチェック
		if currentTileX != gameResources.PlayerTileState.LastTileX || currentTileY != gameResources.PlayerTileState.LastTileY {
			// ワープホールのチェック
			checkTileWarp(world, gridElement)

			// アイテムのチェック
			checkTileItemsForGridPlayer(world, gridElement)

			// 現在の位置を記録
			gameResources.PlayerTileState.LastTileX = currentTileX
			gameResources.PlayerTileState.LastTileY = currentTileY
		}
	}))
}

// canMoveTo は指定位置に移動可能かチェックする
func canMoveTo(world w.World, tileX, tileY int, movingEntity ecs.Entity) bool {
	// 他のエンティティとの衝突チェック
	canMove := true
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.BlockPass,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 自分自身は除外
		if entity == movingEntity {
			return
		}

		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(gridElement.X) == tileX && int(gridElement.Y) == tileY {
			canMove = false
		}
	}))

	// TODO: マップの境界チェックやタイルの通行可否チェックを追加
	return canMove
}

// checkTileWarp はプレイヤーがいるタイルのワープホールをチェックする
func checkTileWarp(world w.World, playerGrid *gc.GridElement) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	pixelX := int(playerGrid.X) * 32
	pixelY := int(playerGrid.Y) * 32
	tileEntity := gameResources.Level.AtEntity(gc.Pixel(pixelX), gc.Pixel(pixelY))

	if tileEntity.HasComponent(world.Components.Warp) {
		warp := world.Components.Warp.Get(tileEntity).(*gc.Warp)
		gameResources.PlayerTileState.CurrentWarp = warp // 現在のワープホールを記録

		switch warp.Mode {
		case gc.WarpModeNext:
			gamelog.New(gamelog.FieldLog).
				Append("階段を発見した。Enterキーで移動").
				Log()
		case gc.WarpModeEscape:
			gamelog.New(gamelog.FieldLog).
				Append("出口を発見した。Enterキーで移動").
				Log()
		}
	} else {
		// ワープホールから離れた場合はリセット
		gameResources.PlayerTileState.CurrentWarp = nil
	}
}

// HandleWarpInput はワープホール上でのEnterキー入力を処理する
func HandleWarpInput(world w.World) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	if gameResources.PlayerTileState.CurrentWarp == nil {
		return // ワープホール上にいない
	}

	switch gameResources.PlayerTileState.CurrentWarp.Mode {
	case gc.WarpModeNext:
		gameResources.SetStateEvent(resources.StateEventWarpNext)
	case gc.WarpModeEscape:
		gameResources.SetStateEvent(resources.StateEventWarpEscape)
	}
}

// checkTileItemsForGridPlayer はグリッドベースプレイヤーのタイルアイテムをチェックする
func checkTileItemsForGridPlayer(world w.World, playerGrid *gc.GridElement) {
	playerTileX := int(playerGrid.X)
	playerTileY := int(playerGrid.Y)

	// GridElementベースのアイテムをチェック
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Item,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		nameComp := world.Components.Name.Get(entity).(*gc.Name)

		itemTileX := int(gridElement.X)
		itemTileY := int(gridElement.Y)

		// GridElementアイテムをチェック

		if itemTileX == playerTileX && itemTileY == playerTileY {
			// アイテムを発見したメッセージを表示
			gamelog.New(gamelog.FieldLog).
				ItemName(nameComp.Name).
				Append("を発見した。").
				Log()
		}
	}))
}
