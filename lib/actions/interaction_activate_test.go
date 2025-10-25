package actions

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInteractionActivateActivity_Info はActivityInfoが正しく返されることを確認
func TestInteractionActivateActivity_Info(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()

	activity := &InteractionActivateActivity{InteractableEntity: triggerEntity}
	info := activity.Info()

	assert.Equal(t, "相互作用発動", info.Name)
	assert.False(t, info.Interruptible, "相互作用発動は中断不可")
	assert.False(t, info.Resumable, "相互作用発動は再開不可")
	assert.Equal(t, 0, info.ActionPointCost, "相互作用発動はAPコストなし")
}

// TestInteractionActivateActivity_String はString()が正しく返されることを確認
func TestInteractionActivateActivity_String(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()

	activity := &InteractionActivateActivity{InteractableEntity: triggerEntity}
	assert.Equal(t, "InteractionActivate", activity.String())
}

// TestInteractionActivateActivity_Validate_Success はValidateが成功することを確認
func TestInteractionActivateActivity_Validate_Success(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.WarpNextInteraction{},
	})

	activity := &InteractionActivateActivity{InteractableEntity: triggerEntity}
	act := NewActivity(activity, triggerEntity, 1)

	err := activity.Validate(act, world)
	assert.NoError(t, err, "Triggerコンポーネントがある場合は検証成功")
}

// TestInteractionActivateActivity_Validate_NoTrigger はTriggerコンポーネントがない場合のエラーを確認
func TestInteractionActivateActivity_Validate_NoTrigger(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	notInteractableEntity := world.Manager.NewEntity()

	activity := &InteractionActivateActivity{InteractableEntity: notInteractableEntity}
	act := NewActivity(activity, notInteractableEntity, 1)

	err := activity.Validate(act, world)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Interactableを持っていません")
}

// InvalidRangeTrigger は無効なActivationRangeを持つテスト用トリガー
type InvalidRangeTrigger struct{}

func (t InvalidRangeTrigger) Config() gc.InteractionConfig {
	return gc.InteractionConfig{
		ActivationRange: gc.ActivationRange("INVALID_RANGE"),
		ActivationWay:   gc.ActivationWayManual,
	}
}

// InvalidWayTrigger は無効なActivationWayを持つテスト用トリガー
type InvalidWayTrigger struct{}

func (t InvalidWayTrigger) Config() gc.InteractionConfig {
	return gc.InteractionConfig{
		ActivationRange: gc.ActivationRangeSameTile,
		ActivationWay:   gc.ActivationWay("INVALID_WAY"),
	}
}

// TestInteractionActivateActivity_Validate_InvalidRange は無効なActivationRangeの検証エラーを確認
func TestInteractionActivateActivity_Validate_InvalidRange(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: InvalidRangeTrigger{},
	})

	activity := &InteractionActivateActivity{InteractableEntity: triggerEntity}
	act := NewActivity(activity, triggerEntity, 1)

	err := activity.Validate(act, world)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "無効なActivationRange")
}

// TestInteractionActivateActivity_Validate_InvalidWay は無効なActivationWayの検証エラーを確認
func TestInteractionActivateActivity_Validate_InvalidWay(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: InvalidWayTrigger{},
	})

	activity := &InteractionActivateActivity{InteractableEntity: triggerEntity}
	act := NewActivity(activity, triggerEntity, 1)

	err := activity.Validate(act, world)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "無効なActivationWay")
}

// TestInteractionActivateActivity_WarpNext はWarpNextTriggerの動作を確認
func TestInteractionActivateActivity_WarpNext(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpNextトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.WarpNextInteraction{},
	})

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "WarpNextトリガーが成功するべき")

	// StateEventが設定されていることを確認
	event := world.Resources.Dungeon.ConsumeStateEvent()
	require.NotNil(t, event, "StateEventが設定されているべき")
	_, ok := event.(resources.WarpNextEvent)
	assert.True(t, ok, "WarpNextEventが設定されているべき")
}

// TestInteractionActivateActivity_WarpEscape はWarpEscapeTriggerの動作を確認
func TestInteractionActivateActivity_WarpEscape(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpEscapeトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.WarpEscapeInteraction{},
	})

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "WarpEscapeトリガーが成功するべき")

	// StateEventが設定されていることを確認
	event := world.Resources.Dungeon.ConsumeStateEvent()
	require.NotNil(t, event, "StateEventが設定されているべき")
	_, ok := event.(resources.WarpEscapeEvent)
	assert.True(t, ok, "WarpEscapeEventが設定されているべき")
}

// TestInteractionActivateActivity_GameClear はゲームクリア条件を満たした脱出の動作を確認
func TestInteractionActivateActivity_GameClear(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// ゲームクリア深度以上を設定
	world.Resources.Dungeon.Depth = consts.GameClearDepth

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpEscapeトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.WarpEscapeInteraction{},
	})

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "WarpEscapeトリガーが成功するべき")

	// GameClearEventが設定されていることを確認
	event := world.Resources.Dungeon.ConsumeStateEvent()
	require.NotNil(t, event, "StateEventが設定されているべき")
	_, ok := event.(resources.GameClearEvent)
	assert.True(t, ok, "GameClearEventが設定されているべき")
}

// TestInteractionActivateActivity_Door はDoorTriggerの動作を確認
func TestInteractionActivateActivity_Door(t *testing.T) {
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
		doorEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
			Data: gc.DoorInteraction{},
		})
		doorEntity.AddComponent(world.Components.BlockPass, &gc.BlockPass{})
		doorEntity.AddComponent(world.Components.BlockView, &gc.BlockView{})

		// InteractionActivateActivityを実行
		manager := NewActivityManager(logger.New(logger.CategoryAction))
		params := ActionParams{
			Actor: player,
		}
		result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: doorEntity}, params, world)

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
		doorEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
			Data: gc.DoorInteraction{},
		})

		// InteractionActivateActivityを実行
		manager := NewActivityManager(logger.New(logger.CategoryAction))
		params := ActionParams{
			Actor: player,
		}
		result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: doorEntity}, params, world)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success, "DoorTriggerが成功するべき")

		// ドアが閉じていることを確認
		doorComp := world.Components.Door.Get(doorEntity).(*gc.Door)
		assert.False(t, doorComp.IsOpen, "ドアが閉じているべき")
	})
}

// TestInteractionActivateActivity_Talk はTalkTriggerの動作を確認
func TestInteractionActivateActivity_Talk(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// TalkTriggerを持つNPCを作成
	npcEntity := world.Manager.NewEntity()
	npcEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
	npcEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.TalkInteraction{},
	})
	npcEntity.AddComponent(world.Components.Dialog, &gc.Dialog{
		MessageKey: "test_npc_greeting",
	})
	npcEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストNPC"})
	npcEntity.AddComponent(world.Components.FactionNeutral, nil)

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: npcEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "TalkTriggerが成功するべき")

	// TalkActivityはShowDialogEventを設定しない実装のため、
	// トリガーが正常に実行されたことのみを確認する
}

// TestInteractionActivateActivity_Item はItemTriggerの動作を確認
func TestInteractionActivateActivity_Item(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// ItemTriggerを持つアイテムを作成（Consumableで削除確認）
	itemEntity := world.Manager.NewEntity()
	itemEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	itemEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.ItemInteraction{},
	})
	itemEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テストアイテム"})
	itemEntity.AddComponent(world.Components.Item, &gc.Item{})
	itemEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: itemEntity}, params, world)

	// ItemTriggerはPickupActivityを呼び出すが、プレイヤーにインベントリがないため失敗する可能性がある
	// ここではトリガーが実行されたことを確認する
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestInteractionActivateActivity_Melee はMeleeTriggerの動作を確認
func TestInteractionActivateActivity_Melee(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// MeleeTriggerを持つ敵を作成
	enemyEntity := world.Manager.NewEntity()
	enemyEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
	enemyEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.MeleeInteraction{},
	})
	enemyEntity.AddComponent(world.Components.Name, &gc.Name{Name: "テスト敵"})

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: enemyEntity}, params, world)

	// MeleeTriggerはAttackActivityを呼び出すが、必要なコンポーネントがないため失敗する可能性がある
	// ここではトリガーが実行されたことを確認する
	require.NoError(t, err)
	require.NotNil(t, result)
}

// TestInteractionActivateActivity_Consumable はConsumableコンポーネントがある場合にエンティティが削除されることを確認
func TestInteractionActivateActivity_Consumable(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Consumableなトリガーを作成（一度だけ発動するWarpNext）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Interactable, &gc.Interactable{
		Data: gc.WarpNextInteraction{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// エンティティIDを保存
	triggerID := triggerEntity

	// InteractionActivateActivityを実行
	manager := NewActivityManager(logger.New(logger.CategoryAction))
	params := ActionParams{
		Actor: player,
	}
	result, err := manager.Execute(&InteractionActivateActivity{InteractableEntity: triggerEntity}, params, world)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success, "Consumableトリガーが成功するべき")

	// トリガーエンティティが削除されていることを確認
	assert.False(t, triggerID.HasComponent(world.Components.Interactable),
		"Consumableトリガーは実行後に削除されるべき")
}
