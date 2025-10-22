package states

import (
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/movement"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ExecuteMoveAction は移動アクションを実行する
func ExecuteMoveAction(world w.World, direction gc.Direction) {
	entity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return
	}

	if !entity.HasComponent(world.Components.GridElement) {
		return
	}

	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
	currentX := int(gridElement.X)
	currentY := int(gridElement.Y)

	deltaX, deltaY := direction.GetDelta()
	newX := currentX + deltaX
	newY := currentY + deltaY

	// 移動先に会話NPCがいる場合は会話アクション
	npc := findNeutralNPCAtPosition(world, newX, newY)
	if npc != nil {
		params := actions.ActionParams{
			Actor:  entity,
			Target: npc,
		}
		executeActivity(world, &actions.TalkActivity{}, params)
		return
	}

	// 移動先に敵がいる場合は攻撃アクション
	enemy := findEnemyAtPosition(world, entity, newX, newY)
	if enemy != nil {
		params := actions.ActionParams{
			Actor:  entity,
			Target: enemy,
		}
		executeActivity(world, &actions.AttackActivity{}, params)
		return
	}

	// 移動先に閉じたドアがある場合はドアを開くアクション
	door := findClosedDoorAtPosition(world, newX, newY)
	if door != nil {
		params := actions.ActionParams{
			Actor:  entity,
			Target: door,
		}
		executeActivity(world, &actions.OpenDoorActivity{}, params)
		return
	}

	canMove := movement.CanMoveTo(world, newX, newY, entity)
	if canMove {
		destination := gc.Position{X: gc.Pixel(newX), Y: gc.Pixel(newY)}
		params := actions.ActionParams{
			Actor:       entity,
			Destination: &destination,
		}
		executeActivity(world, &actions.MoveActivity{}, params)
	}
}

// executeActivity はアクティビティ実行関数
func executeActivity(world w.World, actorImpl actions.ActivityInterface, params actions.ActionParams) {
	manager := actions.NewActivityManager(logger.New(logger.CategoryAction))

	result, err := manager.Execute(actorImpl, params, world)
	if err != nil {
		_ = result // エラーの場合は結果を使用しない
		return
	}

	// 会話の場合は会話メッセージを表示するStateEventを設定
	if _, isTalkActivity := actorImpl.(*actions.TalkActivity); isTalkActivity && result != nil && result.Success && params.Target != nil {
		targetEntity := *params.Target
		if targetEntity.HasComponent(world.Components.Dialog) {
			dialog := world.Components.Dialog.Get(targetEntity).(*gc.Dialog)
			world.Resources.Dungeon.SetStateEvent(resources.ShowDialogEvent{
				MessageKey:    dialog.MessageKey,
				SpeakerEntity: targetEntity,
			})
		}
	}

	// 移動の場合は追加でタイルイベントをチェック
	if _, isMoveActivity := actorImpl.(*actions.MoveActivity); isMoveActivity && result != nil && result.Success && params.Destination != nil {
		checkTileEvents(world, params.Actor, int(params.Destination.X), int(params.Destination.Y))
	}
}

// ExecuteWaitAction は待機アクションを実行する
func ExecuteWaitAction(world w.World) {
	entity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return
	}

	params := actions.ActionParams{
		Actor:    entity,
		Duration: 1,
		Reason:   "プレイヤー待機",
	}
	executeActivity(world, &actions.WaitActivity{}, params)
}

// ExecuteEnterAction はEnterキーによる状況に応じたアクションを実行する
func ExecuteEnterAction(world w.World) {
	entity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return
	}

	if !entity.HasComponent(world.Components.GridElement) {
		return
	}

	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
	tileX := int(gridElement.X)
	tileY := int(gridElement.Y)

	// 手動実行型かつ直上タイル型のTriggerを実行する
	trigger, triggerEntity := getTriggerAtSameTile(world, gridElement)
	if trigger != nil && trigger.ActivationMode == gc.ActivationModeManual {
		params := actions.ActionParams{Actor: entity}
		executeActivity(world, &actions.TriggerActivateActivity{TriggerEntity: triggerEntity}, params)
		return
	}

	if checkForItems(world, tileX, tileY) {
		params := actions.ActionParams{Actor: entity}
		executeActivity(world, &actions.PickupActivity{}, params)
		return
	}
}

// checkTileEvents はタイル上のイベントをチェックする
func checkTileEvents(world w.World, entity ecs.Entity, tileX, tileY int) {
	// プレイヤーの場合のみタイルイベントをチェック
	if entity.HasComponent(world.Components.Player) {
		gridElement := &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}

		// 手動トリガーのメッセージ表示
		showTileTriggerMessage(world, gridElement)

		// アイテムのメッセージ表示
		showTileItemsForGridPlayer(world, gridElement)
	}
}

// getTriggerAtSameTile はプレイヤーの直上タイルのTriggerとエンティティを取得する
// 複数ある場合は最初に見つかったものを返す
func getTriggerAtSameTile(world w.World, playerGrid *gc.GridElement) (*gc.Trigger, ecs.Entity) {
	var trigger *gc.Trigger
	var triggerEntity ecs.Entity
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Trigger,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if trigger != nil {
			return // 既に見つかっている
		}
		ge := world.Components.GridElement.Get(entity).(*gc.GridElement)
		// 直上タイルのみ
		if ge.X == playerGrid.X && ge.Y == playerGrid.Y {
			trigger = world.Components.Trigger.Get(entity).(*gc.Trigger)
			triggerEntity = entity
		}
	}))
	return trigger, triggerEntity
}

// getTriggerInRange はプレイヤーの範囲内のTriggerとエンティティを取得する
// 複数ある場合は最初に見つかったものを返す
func getTriggerInRange(world w.World, playerGrid *gc.GridElement) (*gc.Trigger, ecs.Entity) {
	var trigger *gc.Trigger
	var triggerEntity ecs.Entity
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Trigger,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if trigger != nil {
			return // 既に見つかっている
		}
		t := world.Components.Trigger.Get(entity).(*gc.Trigger)
		ge := world.Components.GridElement.Get(entity).(*gc.GridElement)

		// ActivationRangeに応じた範囲チェック
		if isInRange(playerGrid, ge, t.ActivationRange) {
			trigger = t
			triggerEntity = entity
		}
	}))
	return trigger, triggerEntity
}

// isInRange はプレイヤーがトリガーの発動範囲内にいるかを判定する
func isInRange(playerGrid, triggerGrid *gc.GridElement, activationRange gc.ActivationRange) bool {
	switch activationRange {
	case gc.ActivationRangeSameTile:
		// 直上（同じタイル）
		return playerGrid.X == triggerGrid.X && playerGrid.Y == triggerGrid.Y
	case gc.ActivationRangeAdjacent:
		// 隣接タイル（近傍8タイル）
		diffX := int(playerGrid.X) - int(triggerGrid.X)
		diffY := int(playerGrid.Y) - int(triggerGrid.Y)
		dx := max(diffX, -diffX)
		dy := max(diffY, -diffY)
		return dx <= 1 && dy <= 1
	default:
		return false
	}
}

// showTileTriggerMessage は手動トリガーのメッセージを表示する
func showTileTriggerMessage(world w.World, playerGrid *gc.GridElement) {
	trigger, _ := getTriggerInRange(world, playerGrid)
	if trigger == nil {
		return
	}

	// ActivationMode=MANUALのTriggerのみメッセージ表示
	if trigger.ActivationMode != gc.ActivationModeManual {
		return
	}

	// Triggerの種類に応じてメッセージ表示
	switch trigger.Detail.(type) {
	case gc.WarpNextTrigger:
		gamelog.New(gamelog.FieldLog).
			Append("転移ゲートがある。Enterキーで移動。").
			Log()
	case gc.WarpEscapeTrigger:
		gamelog.New(gamelog.FieldLog).
			Append("脱出ゲートがある。Enterキーで移動。").
			Log()
	}
}

// showTileItemsForGridPlayer はグリッドベースプレイヤーのタイル上のアイテムメッセージを表示する
func showTileItemsForGridPlayer(world w.World, playerGrid *gc.GridElement) {
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
				Append(" を発見した。").
				Log()
		}
	}))
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
func findEnemyAtPosition(world w.World, movingEntity ecs.Entity, tileX, tileY int) *ecs.Entity {
	var foundEnemy *ecs.Entity

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
				foundEnemy = &entity
				return
			}
		}
	}))

	return foundEnemy
}

// findNeutralNPCAtPosition は指定位置にある中立NPC（会話可能）を検索する
func findNeutralNPCAtPosition(world w.World, tileX, tileY int) *ecs.Entity {
	var foundNPC *ecs.Entity

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.FactionNeutral,
		world.Components.Dialog,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(gridElement.X) == tileX && int(gridElement.Y) == tileY {
			foundNPC = &entity
		}
	}))

	return foundNPC
}

// findClosedDoorAtPosition は指定位置にある閉じたドアエンティティを検索する
func findClosedDoorAtPosition(world w.World, tileX, tileY int) *ecs.Entity {
	var foundDoor *ecs.Entity

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Door,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(gridElement.X) == tileX && int(gridElement.Y) == tileY {
			door := world.Components.Door.Get(entity).(*gc.Door)
			// 閉じているドアのみを対象
			if !door.IsOpen {
				foundDoor = &entity
			}
		}
	}))

	return foundDoor
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

// InteractionAction はインタラクション可能なアクション情報
type InteractionAction struct {
	Label    string                    // 表示ラベル（例："開く(上)"）
	Activity actions.ActivityInterface // 実行するアクティビティ
	Target   ecs.Entity                // ターゲットエンティティ
}

// GetInteractionActions はプレイヤー周辺の実行可能なアクションを取得する
func GetInteractionActions(world w.World) []InteractionAction {
	playerEntity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return nil
	}

	if !playerEntity.HasComponent(world.Components.GridElement) {
		return nil
	}

	gridElement := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)
	playerX := int(gridElement.X)
	playerY := int(gridElement.Y)

	var interactionActions []InteractionAction

	// 近傍8マスをスキャン（プレイヤー位置を除く周辺）
	directions := []struct {
		dx    int
		dy    int
		label string
	}{
		{0, -1, "上"},
		{0, 1, "下"},
		{-1, 0, "左"},
		{1, 0, "右"},
		{-1, -1, "左上"},
		{1, -1, "右上"},
		{-1, 1, "左下"},
		{1, 1, "右下"},
	}

	for _, dir := range directions {
		tileX := playerX + dir.dx
		tileY := playerY + dir.dy

		// ドアをチェック
		world.Manager.Join(
			world.Components.GridElement,
			world.Components.Door,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			ge := world.Components.GridElement.Get(entity).(*gc.GridElement)
			if int(ge.X) == tileX && int(ge.Y) == tileY {
				door := world.Components.Door.Get(entity).(*gc.Door)

				if door.IsOpen {
					// 開いているドアは閉じる
					interactionActions = append(interactionActions, InteractionAction{
						Label:    "閉じる(" + dir.label + ")",
						Activity: &actions.CloseDoorActivity{},
						Target:   entity,
					})
				} else {
					// 閉じているドアは開く
					interactionActions = append(interactionActions, InteractionAction{
						Label:    "開く(" + dir.label + ")",
						Activity: &actions.OpenDoorActivity{},
						Target:   entity,
					})
				}
			}
		}))

		// 会話可能NPCをチェック
		world.Manager.Join(
			world.Components.GridElement,
			world.Components.FactionNeutral,
			world.Components.Dialog,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			ge := world.Components.GridElement.Get(entity).(*gc.GridElement)
			if int(ge.X) == tileX && int(ge.Y) == tileY {
				name := world.Components.Name.Get(entity).(*gc.Name)
				interactionActions = append(interactionActions, InteractionAction{
					Label:    "話しかける(" + name.Name + ")",
					Activity: &actions.TalkActivity{},
					Target:   entity,
				})
			}
		}))
	}

	return interactionActions
}
