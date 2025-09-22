package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/movement"
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
// 複数プレイヤーエンティティが存在する場合は最初のエンティティのみを処理する
func executeMoveAction(world w.World, direction gc.Direction) {
	var firstPlayerEntity ecs.Entity
	var hasFirstPlayer bool

	// 最初のプレイヤーエンティティを取得
	world.Manager.Join(
		world.Components.Player,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if !hasFirstPlayer {
			firstPlayerEntity = entity
			hasFirstPlayer = true
		}
	}))

	// 最初のプレイヤーエンティティのみを処理
	if hasFirstPlayer {
		entity := firstPlayerEntity
		// 現在位置を取得
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		currentX := int(gridElement.X)
		currentY := int(gridElement.Y)

		// 移動先を計算
		deltaX, deltaY := direction.GetDelta()
		newX := currentX + deltaX
		newY := currentY + deltaY

		// 移動先に敵がいるかチェック
		enemy := findEnemyAtPosition(world, entity, newX, newY)
		if enemy != ecs.Entity(0) {
			// 敵がいる場合は攻撃アクション
			params := actions.ActionParams{
				Actor:  entity,
				Target: &enemy,
			}
			executeActivity(world, actions.ActivityAttack, params)
			return
		}

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
	}
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
		if checkForWarp(world, entity) {
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
	if entity.HasComponent(world.Components.Player) {
		gridElement := &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}

		// ワープホールのチェック
		checkTileWarp(world, gridElement)

		// アイテムのチェック
		checkTileItemsForGridPlayer(world, gridElement)
	}
}

// CanMoveTo は指定位置に移動可能かチェックする
func CanMoveTo(world w.World, tileX, tileY int, movingEntity ecs.Entity) bool {
	return movement.CanMoveTo(world, tileX, tileY, movingEntity)
}

// getWarpAtPlayerPosition はプレイヤーの現在位置のワープホールを取得する
func getWarpAtPlayerPosition(world w.World, playerGrid *gc.GridElement) *gc.Warp {
	pixelX := int(playerGrid.X) * 32
	pixelY := int(playerGrid.Y) * 32
	tileEntity := world.Resources.Dungeon.Level.AtEntity(gc.Pixel(pixelX), gc.Pixel(pixelY))

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
func checkForWarp(world w.World, entity ecs.Entity) bool {
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

// findEnemyAtPosition は指定位置にいる敵エンティティを検索する
func findEnemyAtPosition(world w.World, movingEntity ecs.Entity, tileX, tileY int) ecs.Entity {
	var foundEnemy ecs.Entity

	// 指定位置にいる全エンティティをチェック
	world.Manager.Join(
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 自分自身は除外
		if entity == movingEntity {
			return
		}

		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(gridElement.X) == tileX && int(gridElement.Y) == tileY {
			// 死亡しているエンティティは除外
			if entity.HasComponent(world.Components.Dead) {
				return
			}

			// 敵対関係かチェック
			if isHostileFaction(world, movingEntity, entity) {
				foundEnemy = entity
				return // 最初に見つかった敵を返す
			}
		}
	}))

	return foundEnemy
}

// isHostileFaction は2つのエンティティが敵対関係にあるかを判定する
func isHostileFaction(world w.World, entity1, entity2 ecs.Entity) bool {
	// プレイヤー側(Ally)と敵(Enemy)は敵対関係
	entity1IsAlly := entity1.HasComponent(world.Components.FactionAlly)
	entity1IsEnemy := entity1.HasComponent(world.Components.FactionEnemy)
	entity2IsAlly := entity2.HasComponent(world.Components.FactionAlly)
	entity2IsEnemy := entity2.HasComponent(world.Components.FactionEnemy)

	// プレイヤー側 vs 敵側
	if (entity1IsAlly && entity2IsEnemy) || (entity1IsEnemy && entity2IsAlly) {
		return true
	}

	// その他の組み合わせは敵対関係ではない
	return false
}
