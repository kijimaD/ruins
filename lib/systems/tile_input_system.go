package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
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
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyPeriod) {
		// スペースキーまたはピリオドで待機
		executeAction(world, actions.ActionWait, nil)
		return
	}

	// 移動アクションを実行
	if direction != gc.DirectionNone {
		executeMoveAction(world, direction)
	}

	// Enterキー: 状況に応じたアクションを実行
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		actionID := determineEnterAction(world)
		if actionID != actions.ActionNull {
			executeAction(world, actionID, nil)
		}
	}
}

// executeAction は統一されたアクション実行関数
func executeAction(world w.World, actionID actions.ActionID, position *gc.Position) {
	executor := actions.NewExecutor()

	// プレイヤーエンティティを取得
	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		ctx := actions.Context{
			Actor:    entity,
			Position: position,
			World:    world,
		}

		result, err := executor.Execute(actionID, ctx)
		if err != nil {
			// エラーログ（必要に応じて）
			_ = result // 現時点では結果を使用しない
		}

		// アクション成功時、ターンマネージャーで移動ポイントを消費
		if result != nil && result.Success && world.Resources.TurnManager != nil {
			if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
				turnManager.ConsumePlayerMoves(actionID)
			}
		}

		// 移動の場合は追加でタイルイベントをチェック
		if actionID == actions.ActionMove && result != nil && result.Success && position != nil {
			checkTileEvents(world, entity, int(position.X), int(position.Y))
		}
	}))
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
		if canMoveTo(world, newX, newY, entity) {
			// 統一されたアクション実行関数を使用
			position := &gc.Position{X: gc.Pixel(newX), Y: gc.Pixel(newY)}
			executeAction(world, actions.ActionMove, position)
		}
	}))
}

// checkTileEvents はタイル上のイベントをチェックする
// 既存のTileMoveSystemから移植
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

// canMoveTo は指定位置に移動可能かチェックする
// TileMoveSystemから移植
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
// TileMoveSystemから移植
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

// checkTileItemsForGridPlayer はグリッドベースプレイヤーのタイルアイテムをチェックする
// TileMoveSystemから移植
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

// determineEnterAction はEnterキー押下時のアクションを状況に応じて決定する
// CDDAのコンテキスト判断ロジックを参考にした実装
func determineEnterAction(world w.World) actions.ActionID {
	// プレイヤーエンティティを取得
	var playerEntity ecs.Entity
	var playerFound bool
	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = entity
		playerFound = true
	}))

	if !playerFound {
		return actions.ActionNull
	}

	gridElement := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)
	playerTileX := int(gridElement.X)
	playerTileY := int(gridElement.Y)

	// 優先順位1: ワープホールのチェック
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	if gameResources.PlayerTileState.CurrentWarp != nil {
		return actions.ActionWarp
	}

	// 優先順位2: アイテムのチェック
	hasItem := false
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationOnField,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(itemEntity ecs.Entity) {
		itemGrid := world.Components.GridElement.Get(itemEntity).(*gc.GridElement)
		if int(itemGrid.X) == playerTileX && int(itemGrid.Y) == playerTileY {
			hasItem = true
		}
	}))

	if hasItem {
		return actions.ActionPickupItem
	}

	// 何もアクションがない場合
	return actions.ActionNull
}
