package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/movement"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TileInputSystem はプレイヤーからのタイルベース入力を処理する。Actionシステムを使用して移動を実行する
// AIの移動・攻撃も将来的に同じActionシステムを使用予定
// TODO: 文脈に応じて発行アクションを判定する
func TileInputSystem(world w.World) {
	// ターン管理チェック - プレイヤーターンでない場合は入力を受け付けない
	if world.Resources.TurnManager != nil {
		turnManager := world.Resources.TurnManager.(*turns.TurnManager)
		if !turnManager.CanPlayerAct() {
			return
		}
	}
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		executeWaitAction(world)
		return
	}

	// 移動アクションを実行
	if direction != gc.DirectionNone {
		executeMoveAction(world, direction)
	}

	// Enterキー: 状況に応じたアクションを実行
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		executeEnterAction(world)
	}
}

// executeActivity はアクティビティ実行関数
func executeActivity(world w.World, activityType actions.ActivityType, params actions.ActionParams) {
	actionAPI := actions.NewActionAPI()

	result, err := actionAPI.Execute(activityType, params, world)
	if err != nil {
		_ = result // エラーの場合は結果を使用しない
		return
	}

	// 移動の場合は追加でタイルイベントをチェック
	if activityType == actions.ActivityMove && result != nil && result.Success && params.Destination != nil {
		// TODO: AI用と共通化したほうがよさそう? プレイヤーの場合だけログを出す、とかはありそうなものの
		checkTileEvents(world, params.Actor, int(params.Destination.X), int(params.Destination.Y))
	}
}

// executeMoveAction は移動アクションを実行する
func executeMoveAction(world w.World, direction gc.Direction) {
	// プレイヤーエンティティを取得
	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 現在位置を取得
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		currentX := int(gridElement.X)
		currentY := int(gridElement.Y)

		// 移動先を計算
		deltaX, deltaY := direction.GetDelta()
		newX := currentX + deltaX
		newY := currentY + deltaY

		// 移動可能かチェックして移動
		canMove := CanMoveTo(world, newX, newY, entity)

		if canMove {
			// 統一されたアクティビティ実行関数を使用
			destination := gc.Position{X: gc.Pixel(newX), Y: gc.Pixel(newY)}
			params := actions.ActionParams{
				Actor:       entity,
				Destination: &destination,
			}
			executeActivity(world, actions.ActivityMove, params)
		}
	}))
}

// executeWaitAction は待機アクションを実行する
func executeWaitAction(world w.World) {
	// プレイヤーエンティティを取得
	world.Manager.Join(
		world.Components.Player,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		params := actions.ActionParams{
			Actor:    entity,
			Duration: 1,
			Reason:   "プレイヤー待機",
		}
		executeActivity(world, actions.ActivityWait, params)
	}))
}

// executeEnterAction はEnterキーによる状況に応じたアクションを実行する
func executeEnterAction(world w.World) {
	// プレイヤーエンティティを取得
	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		tileX := int(gridElement.X)
		tileY := int(gridElement.Y)

		// ワープホールチェック
		if checkForWarp(world, entity, tileX, tileY) {
			params := actions.ActionParams{Actor: entity}
			executeActivity(world, actions.ActivityWarp, params)
			return
		}

		// アイテム拾得チェック
		if checkForItems(world, tileX, tileY) {
			params := actions.ActionParams{Actor: entity}
			executeActivity(world, actions.ActivityPickup, params)
			return
		}
	}))
}

// checkTileEvents はタイル上のイベントをチェックする
func checkTileEvents(world w.World, entity ecs.Entity, tileX, tileY int) {
	// プレイヤーの場合のみタイルイベントをチェック
	if entity.HasComponent(world.Components.Operator) {
		gameResources := world.Resources.Dungeon.(*resources.Dungeon)

		// プレイヤーが新しいタイルに移動した場合のみチェック
		if tileX != gameResources.PlayerTileState.LastTileX || tileY != gameResources.PlayerTileState.LastTileY {
			gridElement := &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}

			// ワープホールのチェック
			checkTileWarp(world, gridElement)

			// アイテムのチェック
			checkTileItemsForGridPlayer(world, gridElement)

			// 現在の位置を記録
			gameResources.PlayerTileState.LastTileX = tileX
			gameResources.PlayerTileState.LastTileY = tileY
		}
	}
}

// CanMoveTo は指定位置に移動可能かチェックする
func CanMoveTo(world w.World, tileX, tileY int, movingEntity ecs.Entity) bool {
	return movement.CanMoveTo(world, tileX, tileY, movingEntity)
}

// getWarpAtPlayerPosition はプレイヤーの現在位置のワープホールを取得する
func getWarpAtPlayerPosition(world w.World, playerGrid *gc.GridElement) *gc.Warp {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	pixelX := int(playerGrid.X) * 32
	pixelY := int(playerGrid.Y) * 32
	tileEntity := gameResources.Level.AtEntity(gc.Pixel(pixelX), gc.Pixel(pixelY))

	if tileEntity.HasComponent(world.Components.Warp) {
		return world.Components.Warp.Get(tileEntity).(*gc.Warp)
	}
	return nil
}

// checkTileWarp はプレイヤーがいるタイルのワープホールをチェックする
func checkTileWarp(world w.World, playerGrid *gc.GridElement) {
	warp := getWarpAtPlayerPosition(world, playerGrid)

	if warp != nil {
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

// checkForWarp はプレイヤー位置にワープホールがあるかチェック
func checkForWarp(world w.World, entity ecs.Entity, tileX, tileY int) bool {
	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
	return getWarpAtPlayerPosition(world, gridElement) != nil
}

// checkForItems はプレイヤー位置にアイテムがあるかチェック
func checkForItems(world w.World, tileX, tileY int) bool {
	hasItem := false
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationOnField,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(itemEntity ecs.Entity) {
		itemGrid := world.Components.GridElement.Get(itemEntity).(*gc.GridElement)
		if int(itemGrid.X) == tileX && int(itemGrid.Y) == tileY {
			hasItem = true
		}
	}))
	return hasItem
}
