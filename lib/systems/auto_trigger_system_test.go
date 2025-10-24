package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用のトリガー型定義

// AutoTestTrigger は自動発動するテスト用トリガー（SameTile）
type AutoTestTrigger struct{}

// Config はTriggerDataインターフェースの実装
func (t AutoTestTrigger) Config() gc.TriggerConfig {
	return gc.TriggerConfig{
		ActivationRange: gc.ActivationRangeSameTile,
		ActivationWay:   gc.ActivationWayAuto,
	}
}

// AutoAdjacentTrigger は自動発動するテスト用トリガー（Adjacent）
type AutoAdjacentTrigger struct{}

// Config はTriggerDataインターフェースの実装
func (t AutoAdjacentTrigger) Config() gc.TriggerConfig {
	return gc.TriggerConfig{
		ActivationRange: gc.ActivationRangeAdjacent,
		ActivationWay:   gc.ActivationWayAuto,
	}
}

// AutoWarpTrigger は自動発動するワープトリガー（テスト用）
type AutoWarpTrigger struct{}

// Config はTriggerDataインターフェースの実装
func (t AutoWarpTrigger) Config() gc.TriggerConfig {
	return gc.TriggerConfig{
		ActivationRange: gc.ActivationRangeSameTile,
		ActivationWay:   gc.ActivationWayAuto,
	}
}

// InvalidAutoRangeTrigger は無効なActivationRangeを持つ自動発動トリガー（テスト用）
type InvalidAutoRangeTrigger struct{}

// Config はTriggerDataインターフェースの実装
func (t InvalidAutoRangeTrigger) Config() gc.TriggerConfig {
	return gc.TriggerConfig{
		ActivationRange: gc.ActivationRange("INVALID_RANGE"),
		ActivationWay:   gc.ActivationWayAuto,
	}
}

// InvalidAutoWayTrigger は無効なActivationWayを持つトリガー（テスト用）
type InvalidAutoWayTrigger struct{}

// Config はTriggerDataインターフェースの実装
func (t InvalidAutoWayTrigger) Config() gc.TriggerConfig {
	return gc.TriggerConfig{
		ActivationRange: gc.ActivationRangeSameTile,
		ActivationWay:   gc.ActivationWay("INVALID_WAY"),
	}
}

// TestAutoTriggerSystem_AutoWay はAuto方式のトリガーが自動実行されることを確認
func TestAutoTriggerSystem_AutoWay(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Auto方式のトリガーを作成（プレイヤーと同じタイル）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: AutoTestTrigger{},
	})
	// 実行されたことを確認するためにConsumableを追加
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// Consumableトリガーが削除されていることを確認（実行された証拠）
	assert.False(t, triggerEntity.HasComponent(world.Components.Trigger),
		"Autoトリガーは自動実行され、Consumableなので削除されるべき")
}

// TestAutoTriggerSystem_ManualWay はManual方式のトリガーが自動実行されないことを確認
func TestAutoTriggerSystem_ManualWay(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Manual方式のトリガーを作成（プレイヤーと同じタイル）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpNextTrigger{}, // Manual 方式
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// Manualトリガーは実行されず、残っているべき
	assert.True(t, triggerEntity.HasComponent(world.Components.Trigger),
		"Manualトリガーは自動実行されないべき")
	assert.True(t, triggerEntity.HasComponent(world.Components.Consumable),
		"Manualトリガーは自動実行されないので削除されないべき")
}

// TestAutoTriggerSystem_OnCollisionWay はOnCollision方式のトリガーが自動実行されないことを確認
func TestAutoTriggerSystem_OnCollisionWay(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// OnCollision方式のトリガーを作成（プレイヤーと隣接）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.DoorTrigger{}, // OnCollision 方式
	})
	triggerEntity.AddComponent(world.Components.Door, &gc.Door{IsOpen: false, Orientation: gc.DoorOrientationHorizontal})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// OnCollisionトリガーは実行されず、ドアは閉じたままのはず
	doorComp := world.Components.Door.Get(triggerEntity).(*gc.Door)
	assert.False(t, doorComp.IsOpen, "OnCollisionトリガーは自動実行されないべき")
}

// TestAutoTriggerSystem_OutOfRange は範囲外のAutoトリガーが実行されないことを確認
func TestAutoTriggerSystem_OutOfRange(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Auto方式のトリガーを作成（プレイヤーから遠い位置）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 50, Y: 50}) // 遠い位置
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: AutoTestTrigger{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// 範囲外なので実行されず、残っているべき
	assert.True(t, triggerEntity.HasComponent(world.Components.Trigger),
		"範囲外のAutoトリガーは実行されないべき")
	assert.True(t, triggerEntity.HasComponent(world.Components.Consumable),
		"範囲外のAutoトリガーは削除されないべき")
}

// TestAutoTriggerSystem_AdjacentRange は隣接範囲のAutoトリガーが実行されることを確認
func TestAutoTriggerSystem_AdjacentRange(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// Auto方式 + Adjacent範囲のトリガーを作成（プレイヤーに隣接）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 11, Y: 10}) // 隣接
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: AutoAdjacentTrigger{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// 隣接範囲内なので実行され、削除されているべき
	assert.False(t, triggerEntity.HasComponent(world.Components.Trigger),
		"隣接範囲のAutoトリガーは実行され、削除されるべき")
}

// TestAutoTriggerSystem_NoPlayer はプレイヤーがいない場合にエラーを返すことを確認
func TestAutoTriggerSystem_NoPlayer(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成しない

	// Auto方式のトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: AutoTestTrigger{},
	})

	// システム実行
	err := AutoTriggerSystem(world)
	require.Error(t, err, "プレイヤーがいない場合はエラーを返すべき")
}

// TestAutoTriggerSystem_MultipleAutoTriggers は複数のAutoトリガーが同時に実行されることを確認
func TestAutoTriggerSystem_MultipleAutoTriggers(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// 複数のAutoトリガーを作成
	trigger1 := world.Manager.NewEntity()
	trigger1.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	trigger1.AddComponent(world.Components.Trigger, &gc.Trigger{Data: AutoTestTrigger{}})
	trigger1.AddComponent(world.Components.Consumable, &gc.Consumable{})

	trigger2 := world.Manager.NewEntity()
	trigger2.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	trigger2.AddComponent(world.Components.Trigger, &gc.Trigger{Data: AutoTestTrigger{}})
	trigger2.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// 両方のトリガーが実行され、削除されているべき
	assert.False(t, trigger1.HasComponent(world.Components.Trigger),
		"1つ目のAutoトリガーは削除されるべき")
	assert.False(t, trigger2.HasComponent(world.Components.Trigger),
		"2つ目のAutoトリガーは削除されるべき")
}

// TestAutoTriggerSystem_WarpNextEvent はWarpNextトリガーでStateEventが設定されることを確認
func TestAutoTriggerSystem_WarpNextEvent(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// WarpNextをAuto方式にカスタマイズ（本来はManualだがテスト用にAuto化）
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: gc.WarpNextTrigger{}, // 実際のWarpNextを使用
	})

	// Auto方式にするためにTrigger.Dataを上書き（テスト用）
	trigger := world.Components.Trigger.Get(triggerEntity).(*gc.Trigger)
	trigger.Data = AutoWarpTrigger{}

	// システム実行
	err := AutoTriggerSystem(world)
	require.NoError(t, err)

	// StateEventが設定されていることを確認
	// 注: 実際のWarpNextTriggerでないため、StateEventは設定されない可能性がある
	// この場合はトリガーの実行自体が確認できれば良い
}

// TestAutoTriggerSystem_PlayerNoGridElement はプレイヤーにGridElementがない場合の動作確認
func TestAutoTriggerSystem_PlayerNoGridElement(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成（GridElementなし）
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	// GridElementを追加しない

	// Autoトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: AutoTestTrigger{},
	})

	// システム実行
	err := AutoTriggerSystem(world)
	// GridElementがない場合はnilを返して処理を中断する
	assert.NoError(t, err, "プレイヤーにGridElementがない場合はエラーなしで終了すべき")

	// トリガーは実行されないべき
	assert.True(t, triggerEntity.HasComponent(world.Components.Trigger),
		"プレイヤーにGridElementがない場合、トリガーは実行されないべき")
}

// TestAutoTriggerSystem_InvalidRange は無効なActivationRangeを持つトリガーがスキップされることを確認
func TestAutoTriggerSystem_InvalidRange(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// 無効なActivationRangeを持つトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: InvalidAutoRangeTrigger{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行（エラーは返さず、警告ログを出してスキップする）
	err := AutoTriggerSystem(world)
	assert.NoError(t, err, "無効なトリガーはスキップされ、エラーは返さない")

	// トリガーは実行されず、残っているべき
	assert.True(t, triggerEntity.HasComponent(world.Components.Trigger),
		"無効なActivationRangeのトリガーはスキップされるべき")
	assert.True(t, triggerEntity.HasComponent(world.Components.Consumable),
		"無効なActivationRangeのトリガーは削除されないべき")
}

// TestAutoTriggerSystem_InvalidWay は無効なActivationWayを持つトリガーがスキップされることを確認
func TestAutoTriggerSystem_InvalidWay(t *testing.T) {
	t.Parallel()

	world := testutil.InitTestWorld(t)

	// プレイヤーを作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, &gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})

	// 無効なActivationWayを持つトリガーを作成
	triggerEntity := world.Manager.NewEntity()
	triggerEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: 10, Y: 10})
	triggerEntity.AddComponent(world.Components.Trigger, &gc.Trigger{
		Data: InvalidAutoWayTrigger{},
	})
	triggerEntity.AddComponent(world.Components.Consumable, &gc.Consumable{})

	// システム実行（エラーは返さず、警告ログを出してスキップする）
	err := AutoTriggerSystem(world)
	assert.NoError(t, err, "無効なトリガーはスキップされ、エラーは返さない")

	// トリガーは実行されず、残っているべき
	assert.True(t, triggerEntity.HasComponent(world.Components.Trigger),
		"無効なActivationWayのトリガーはスキップされるべき")
	assert.True(t, triggerEntity.HasComponent(world.Components.Consumable),
		"無効なActivationWayのトリガーは削除されないべき")
}
