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

	// 移動先にOnCollision方式のTriggerがある場合は自動実行
	targetGrid := &gc.GridElement{X: gc.Tile(newX), Y: gc.Tile(newY)}
	trigger, triggerEntity := getTriggerAtSameTile(world, targetGrid)
	if trigger != nil && trigger.Data.Config().ActivationWay == gc.ActivationWayOnCollision {
		// DoorTriggerの場合は、閉じている場合のみ実行（開いている場合は通過）
		if _, isDoorTrigger := trigger.Data.(gc.DoorTrigger); isDoorTrigger {
			if triggerEntity.HasComponent(world.Components.Door) {
				door := world.Components.Door.Get(triggerEntity).(*gc.Door)
				if !door.IsOpen {
					// 閉じているドアは開くトリガーを実行
					params := actions.ActionParams{Actor: entity}
					executeActivity(world, &actions.TriggerActivateActivity{TriggerEntity: triggerEntity}, params)
					return
				}
				// 開いているドアは通過可能なので、トリガーを実行せずに下の移動処理に進む
			}
		} else {
			// ドア以外のOnCollisionトリガーは常に実行
			params := actions.ActionParams{Actor: entity}
			executeActivity(world, &actions.TriggerActivateActivity{TriggerEntity: triggerEntity}, params)
			return
		}
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

// ExecuteEnterAction は直上タイルのTriggerを実行する
func ExecuteEnterAction(world w.World) {
	entity, err := worldhelper.GetPlayerEntity(world)
	if err != nil {
		return
	}

	if !entity.HasComponent(world.Components.GridElement) {
		return
	}

	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

	trigger, triggerEntity := getTriggerAtSameTile(world, gridElement)
	if trigger != nil && trigger.Data.Config().ActivationRange == gc.ActivationRangeSameTile {
		params := actions.ActionParams{Actor: entity}
		executeActivity(world, &actions.TriggerActivateActivity{TriggerEntity: triggerEntity}, params)
	}
}

// checkTileEvents はタイル上のイベントをチェックする
func checkTileEvents(world w.World, entity ecs.Entity, tileX, tileY int) {
	// プレイヤーの場合のみタイルイベントをチェック
	if entity.HasComponent(world.Components.Player) {
		gridElement := &gc.GridElement{X: gc.Tile(tileX), Y: gc.Tile(tileY)}

		// 手動トリガーのメッセージ表示
		showTileTriggerMessage(world, gridElement)
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
		if worldhelper.IsInActivationRange(playerGrid, ge, t.Data.Config().ActivationRange) {
			trigger = t
			triggerEntity = entity
		}
	}))
	return trigger, triggerEntity
}

// getAllInteractiveTriggersInRange はプレイヤーの範囲内の全てのインタラクティブなTriggerエンティティを取得する
// Manual と OnCollision 方式のTriggerが対象
func getAllInteractiveTriggersInRange(world w.World, playerGrid *gc.GridElement) []ecs.Entity {
	var results []ecs.Entity

	world.Manager.Join(
		world.Components.GridElement,
		world.Components.Trigger,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		trigger := world.Components.Trigger.Get(entity).(*gc.Trigger)
		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)

		way := trigger.Data.Config().ActivationWay
		// ManualまたはOnCollision方式で、範囲内にあるものを取得
		if (way == gc.ActivationWayManual || way == gc.ActivationWayOnCollision) &&
			worldhelper.IsInActivationRange(playerGrid, gridElement, trigger.Data.Config().ActivationRange) {
			results = append(results, entity)
		}
	}))

	return results
}

// getDirectionLabel はプレイヤーからターゲットへの方向ラベルを取得する
func getDirectionLabel(playerGrid, targetGrid *gc.GridElement) string {
	dx := int(targetGrid.X) - int(playerGrid.X)
	dy := int(targetGrid.Y) - int(playerGrid.Y)

	// 同じタイル
	if dx == 0 && dy == 0 {
		return "直上"
	}

	// 8方向を判定
	if dy < 0 {
		if dx < 0 {
			return "左上"
		} else if dx > 0 {
			return "右上"
		}
		return "上"
	} else if dy > 0 {
		if dx < 0 {
			return "左下"
		} else if dx > 0 {
			return "右下"
		}
		return "下"
	}
	if dx < 0 {
		return "左"
	}
	return "右"
}

// showTileTriggerMessage は手動トリガーのメッセージを表示する
func showTileTriggerMessage(world w.World, playerGrid *gc.GridElement) {
	trigger, triggerEntity := getTriggerInRange(world, playerGrid)
	if trigger == nil {
		return
	}

	if trigger.Data.Config().ActivationWay != gc.ActivationWayManual {
		return
	}

	switch trigger.Data.(type) {
	case gc.WarpNextTrigger:
		gamelog.New(gamelog.FieldLog).
			Append("転移ゲートがある。Enterキーで移動。").
			Log()
	case gc.WarpEscapeTrigger:
		gamelog.New(gamelog.FieldLog).
			Append("脱出ゲートがある。Enterキーで移動。").
			Log()
	case gc.ItemTrigger:
		// アイテムの名前を取得して表示
		if triggerEntity.HasComponent(world.Components.Name) {
			nameComp := world.Components.Name.Get(triggerEntity).(*gc.Name)
			gamelog.New(gamelog.FieldLog).
				ItemName(nameComp.Name).
				Append(" がある。").
				Log()
		}
	}
}

// InteractionAction はインタラクション可能なアクション情報
type InteractionAction struct {
	Label    string                    // 表示ラベル（例："開く(上)"）
	Activity actions.ActivityInterface // 実行するアクティビティ
	Target   ecs.Entity                // ターゲットエンティティ
}

// getTriggerActions はTriggerに対応するアクションを取得する
func getTriggerActions(world w.World, trigger *gc.Trigger, triggerEntity ecs.Entity, dirLabel string) []InteractionAction {
	var result []InteractionAction

	switch trigger.Data.(type) {
	case gc.WarpNextTrigger:
		result = append(result, InteractionAction{
			Label:    "転移(" + dirLabel + ")",
			Activity: &actions.TriggerActivateActivity{TriggerEntity: triggerEntity},
			Target:   triggerEntity,
		})
	case gc.WarpEscapeTrigger:
		result = append(result, InteractionAction{
			Label:    "脱出(" + dirLabel + ")",
			Activity: &actions.TriggerActivateActivity{TriggerEntity: triggerEntity},
			Target:   triggerEntity,
		})
	case gc.DoorTrigger:
		// ドアの状態に応じたアクションを生成
		if triggerEntity.HasComponent(world.Components.Door) {
			door := world.Components.Door.Get(triggerEntity).(*gc.Door)
			if door.IsOpen {
				result = append(result, InteractionAction{
					Label:    "閉じる(" + dirLabel + ")",
					Activity: &actions.CloseDoorActivity{},
					Target:   triggerEntity,
				})
			} else {
				result = append(result, InteractionAction{
					Label:    "開く(" + dirLabel + ")",
					Activity: &actions.OpenDoorActivity{},
					Target:   triggerEntity,
				})
			}
		}
	case gc.TalkTrigger:
		// 会話アクションを生成
		if triggerEntity.HasComponent(world.Components.Name) {
			name := world.Components.Name.Get(triggerEntity).(*gc.Name)
			result = append(result, InteractionAction{
				Label:    "話しかける(" + name.Name + ")",
				Activity: &actions.TalkActivity{},
				Target:   triggerEntity,
			})
		}
	case gc.ItemTrigger:
		// アイテム拾得アクションを生成
		if triggerEntity.HasComponent(world.Components.Name) {
			name := world.Components.Name.Get(triggerEntity).(*gc.Name)
			result = append(result, InteractionAction{
				Label:    "拾う(" + name.Name + ")",
				Activity: &actions.PickupActivity{},
				Target:   triggerEntity,
			})
		}
	}

	return result
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

	var interactionActions []InteractionAction

	// インタラクティブなTriggerを全て取得してアクションを生成
	triggerEntities := getAllInteractiveTriggersInRange(world, gridElement)
	for _, triggerEntity := range triggerEntities {
		if !triggerEntity.HasComponent(world.Components.GridElement) {
			continue
		}
		if !triggerEntity.HasComponent(world.Components.Trigger) {
			continue
		}

		triggerGrid := world.Components.GridElement.Get(triggerEntity).(*gc.GridElement)
		trigger := world.Components.Trigger.Get(triggerEntity).(*gc.Trigger)
		dirLabel := getDirectionLabel(gridElement, triggerGrid)
		triggerActions := getTriggerActions(world, trigger, triggerEntity, dirLabel)
		interactionActions = append(interactionActions, triggerActions...)
	}

	return interactionActions
}
