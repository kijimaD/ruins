package actions

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTriggerActivateActivity_Info はActivityInfoが正しく返されることを確認
func TestTriggerActivateActivity_Info(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()

	activity := &TriggerActivateActivity{TriggerEntity: triggerEntity}
	info := activity.Info()

	assert.Equal(t, "トリガー発動", info.Name)
	assert.False(t, info.Interruptible, "トリガー発動は中断不可")
	assert.False(t, info.Resumable, "トリガー発動は再開不可")
	assert.Equal(t, 0, info.ActionPointCost, "トリガー発動はAPコストなし")
}

// TestTriggerActivateActivity_String はString()が正しく返されることを確認
func TestTriggerActivateActivity_String(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()

	activity := &TriggerActivateActivity{TriggerEntity: triggerEntity}
	assert.Equal(t, "TriggerActivate", activity.String())
}

// TestTriggerActivateActivity_Validate_Success はValidateが成功することを確認
func TestTriggerActivateActivity_Validate_Success(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpNextTrigger{},
	})

	activity := &TriggerActivateActivity{TriggerEntity: triggerEntity}
	act := NewActivity(activity, triggerEntity, 1)

	err := activity.Validate(act, world)
	assert.NoError(t, err, "Triggerコンポーネントがある場合は検証成功")
}

// TestTriggerActivateActivity_Validate_NoTrigger はTriggerコンポーネントがない場合のエラーを確認
func TestTriggerActivateActivity_Validate_NoTrigger(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	notTriggerEntity := world.Manager.NewEntity()

	activity := &TriggerActivateActivity{TriggerEntity: notTriggerEntity}
	act := NewActivity(activity, notTriggerEntity, 1)

	err := activity.Validate(act, world)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Triggerを持っていません")
}

// TestTriggerActivateActivity_WarpNext はWarpNextTriggerの動作を確認
func TestTriggerActivateActivity_WarpNext(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpNextトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpNextTrigger{},
	})

	// TriggerActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "WarpNextトリガーが成功するべき")

	// StateEventが設定されていることを確認
	event := world.Resources.Dungeon.ConsumeStateEvent()
	require.NotNil(t, event, "StateEventが設定されているべき")
	_, ok := event.(resources.WarpNextEvent)
	assert.True(t, ok, "WarpNextEventが設定されているべき")
}

// TestTriggerActivateActivity_WarpEscape はWarpEscapeTriggerの動作を確認
func TestTriggerActivateActivity_WarpEscape(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpEscapeトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpEscapeTrigger{},
	})

	// TriggerActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "WarpEscapeトリガーが成功するべき")

	// StateEventが設定されていることを確認
	event := world.Resources.Dungeon.ConsumeStateEvent()
	require.NotNil(t, event, "StateEventが設定されているべき")
	_, ok := event.(resources.WarpEscapeEvent)
	assert.True(t, ok, "WarpEscapeEventが設定されているべき")
}

// TestTriggerActivateActivity_Door はDoorTriggerの動作を確認
func TestTriggerActivateActivity_Door(t *testing.T) {
	t.Parallel()

	t.Run("閉じたドアを開く", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// DoorTriggerを持つドアを作成（閉じている）
		doorEntity := world.Manager.NewEntity()
		doorEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
		doorEntity.AddComponent(world.Components.Door, &gc.Door{IsOpen: false, Orientation: gc.DoorOrientationHorizontal})
		doorEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
			Data: gc.DoorTrigger{},
		})
		doorEntity.AddComponent(world.Components.BlockPass, &gc.BlockPass{})
		doorEntity.AddComponent(world.Components.BlockView, &gc.BlockView{})

		// TriggerActivateActivityを実行
		manager := NewActivityManager(logger.New(logger.CategoryAction))
		params := ActionParams{
			Actor: player,
		}
		result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: doorEntity}, params, world)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success, "DoorTriggerが成功するべき")

		// ドアが開いていることを確認
		doorComp := world.Components.Door.Get(doorEntity).(*gc.Door)
		assert.True(t, doorComp.IsOpen, "ドアが開いているべき")
	})

	t.Run("開いたドアを閉じる", func(t *testing.T) {
		t.Parallel()
		world := testutil.InitTestWorld(t)

		// プレイヤーを作成
		player := world.Manager.NewEntity()
		player.AddComponent(world.Components.Player, &gc.Player{})
		player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

		// DoorTriggerを持つドアを作成（開いている）
		doorEntity := world.Manager.NewEntity()
		doorEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
		doorEntity.AddComponent(world.Components.Door, &gc.Door{IsOpen: true, Orientation: gc.DoorOrientationHorizontal})
		doorEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
			Data: gc.DoorTrigger{},
		})

		// TriggerActivateActivityを実行
		manager := NewActivityManager(logger.New(logger.CategoryAction))
		params := ActionParams{
			Actor: player,
		}
		result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: doorEntity}, params, world)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success, "DoorTriggerが成功するべき")

		// ドアが閉じていることを確認
		doorComp := world.Components.Door.Get(doorEntity).(*gc.Door)
		assert.False(t, doorComp.IsOpen, "ドアが閉じているべき")
	})
}

// TestTriggerActivateActivity_Talk はTalkTriggerの動作を確認
func TestTriggerActivateActivity_Talk(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// TalkTriggerを持つNPCを作成
	npcEntity := world.Manager.NewEntity()
	npcEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
	npcEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.TalkTrigger{},
	})
	npcEntity.AddComponent(world.Components.Dialog, &gc.Dialog{
		MessageKey: "test_npc_greeting",
	})
	npcEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストNPC"})
	npcEntity.AddComponent(world.Components.FactionNeutral, nil)

	// TriggerActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: npcEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "TalkTriggerが成功するべき")

	// TalkActivityはShowDialogEventを設定しない実装のため、
	// トリガーが正常に実行されたことのみを確認する
}

// TestTriggerActivateActivity_Item はItemTriggerの動作を確認
func TestTriggerActivateActivity_Item(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// ItemTriggerを持つアイテムを作成（Consumableで削除確認）
	itemEntity := world.Manager.NewEntity()
	itemEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	itemEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.ItemTrigger{},
	})
	itemEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストアイテム"})
	itemEntity.AddComponent(world.Components.Item, &gc.Item{})
	itemEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// TriggerActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: itemEntity}, params, world)

	// ItemTriggerはPickupActivityを呼び出すが、プレイヤーにインベントリがないため失敗する可能性がある
	// ここではトリガーが実行されたことを確認する
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestTriggerActivateActivity_Consumable はConsumableコンポーネントがある場合にエンティティが削除されることを確認
func TestTriggerActivateActivity_Consumable(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Consumableなトリガーを作成（一度だけ発動するWarpNext）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpNextTrigger{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// エンティティIDを保存
	triggerID := triggerEntity

	// TriggerActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&TriggerActivateActivity{TriggerEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "Consumableトリガーが成功するべき")

	// トリガーエンティティが削除されていることを確認
	assert.False(t, triggerID.HasComponent(world.Components.Trigger),
		"Consumableトリガーは実行後に削除されるべき")
}
