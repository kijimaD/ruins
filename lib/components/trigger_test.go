package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWarpNextTrigger_Config はWarpNextTriggerの設定が正しいことを確認
func TestWarpNextTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := WarpNextTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeSameTile, config.ActivationRange,
		"WarpNextは直上タイルで発動する")
	assert.Equal(t, ActivationWayManual, config.ActivationWay,
		"WarpNextは手動発動する")
}

// TestWarpEscapeTrigger_Config はWarpEscapeTriggerの設定が正しいことを確認
func TestWarpEscapeTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := WarpEscapeTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeSameTile, config.ActivationRange,
		"WarpEscapeは直上タイルで発動する")
	assert.Equal(t, ActivationWayManual, config.ActivationWay,
		"WarpEscapeは手動発動する")
}

// TestDoorTrigger_Config はDoorTriggerの設定が正しいことを確認
func TestDoorTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := DoorTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeAdjacent, config.ActivationRange,
		"Doorは隣接タイルで発動する")
	assert.Equal(t, ActivationWayOnCollision, config.ActivationWay,
		"Doorは衝突時に自動発動する")
}

// TestTalkTrigger_Config はTalkTriggerの設定が正しいことを確認
func TestTalkTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := TalkTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeAdjacent, config.ActivationRange,
		"Talkは隣接タイルで発動する")
	assert.Equal(t, ActivationWayOnCollision, config.ActivationWay,
		"Talkは衝突時に自動発動する")
}

// TestItemTrigger_Config はItemTriggerの設定が正しいことを確認
func TestItemTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := ItemTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeSameTile, config.ActivationRange,
		"Itemは直上タイルで発動する")
	assert.Equal(t, ActivationWayManual, config.ActivationWay,
		"Itemは手動発動する")
}

// TestMeleeTrigger_Config はMeleeTriggerの設定が正しいことを確認
func TestMeleeTrigger_Config(t *testing.T) {
	t.Parallel()

	trigger := MeleeTrigger{}
	config := trigger.Config()

	assert.Equal(t, ActivationRangeAdjacent, config.ActivationRange,
		"Meleeは隣接タイルで発動する")
	assert.Equal(t, ActivationWayOnCollision, config.ActivationWay,
		"Meleeは衝突時に自動発動する")
}

// TestActivationRange_Valid は有効なActivationRangeの検証
func TestActivationRange_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		activationRange ActivationRange
		expectValid     bool
	}{
		{
			name:            "SameTile は有効",
			activationRange: ActivationRangeSameTile,
			expectValid:     true,
		},
		{
			name:            "Adjacent は有効",
			activationRange: ActivationRangeAdjacent,
			expectValid:     true,
		},
		{
			name:            "空文字列は無効",
			activationRange: ActivationRange(""),
			expectValid:     false,
		},
		{
			name:            "未定義の値は無効",
			activationRange: ActivationRange("INVALID_RANGE"),
			expectValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.activationRange.Valid()
			if tt.expectValid {
				assert.NoError(t, err, "Valid()はエラーを返さないべき")
			} else {
				assert.Error(t, err, "Valid()はエラーを返すべき")
			}
		})
	}
}

// TestActivationWay_Valid は有効なActivationWayの検証
func TestActivationWay_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		activationWay ActivationWay
		expectValid   bool
	}{
		{
			name:          "Auto は有効",
			activationWay: ActivationWayAuto,
			expectValid:   true,
		},
		{
			name:          "Manual は有効",
			activationWay: ActivationWayManual,
			expectValid:   true,
		},
		{
			name:          "OnCollision は有効",
			activationWay: ActivationWayOnCollision,
			expectValid:   true,
		},
		{
			name:          "空文字列は無効",
			activationWay: ActivationWay(""),
			expectValid:   false,
		},
		{
			name:          "未定義の値は無効",
			activationWay: ActivationWay("INVALID_MODE"),
			expectValid:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.activationWay.Valid()
			if tt.expectValid {
				assert.NoError(t, err, "Valid()はエラーを返さないべき")
			} else {
				assert.Error(t, err, "Valid()はエラーを返すべき")
			}
		})
	}
}

// TestTriggerInterfaceImplementation は全てのトリガーがTriggerDataインターフェースを実装していることを確認
func TestTriggerInterfaceImplementation(t *testing.T) {
	t.Parallel()

	// 全てのトリガータイプがTriggerDataインターフェースを実装していることを確認
	var _ TriggerData = WarpNextTrigger{}
	var _ TriggerData = WarpEscapeTrigger{}
	var _ TriggerData = DoorTrigger{}
	var _ TriggerData = TalkTrigger{}
	var _ TriggerData = ItemTrigger{}
	var _ TriggerData = MeleeTrigger{}
}

// TestTriggerConfigConsistency は全トリガーの設定が一貫していることを確認
func TestTriggerConfigConsistency(t *testing.T) {
	t.Parallel()

	// 全トリガータイプ
	triggers := []struct {
		name    string
		trigger TriggerData
	}{
		{"WarpNext", WarpNextTrigger{}},
		{"WarpEscape", WarpEscapeTrigger{}},
		{"Door", DoorTrigger{}},
		{"Talk", TalkTrigger{}},
		{"Item", ItemTrigger{}},
		{"Melee", MeleeTrigger{}},
	}

	for _, tt := range triggers {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			config := tt.trigger.Config()

			// ActivationRangeが有効な値であることを確認
			err := config.ActivationRange.Valid()
			require.NoError(t, err, "%s のActivationRangeは有効でなければならない", tt.name)

			// ActivationWayが有効な値であることを確認
			err = config.ActivationWay.Valid()
			require.NoError(t, err, "%s のActivationWayは有効でなければならない", tt.name)
		})
	}
}

// TestTriggerDesignConstraints は設計上の制約をテスト（仕様確認）
func TestTriggerDesignConstraints(t *testing.T) {
	t.Parallel()

	t.Run("SameTileトリガーはManual方式", func(t *testing.T) {
		t.Parallel()
		// 仕様: 直上タイルトリガー（WarpNext, WarpEscape, Item）は手動発動
		sameTileTriggers := []TriggerData{
			WarpNextTrigger{},
			WarpEscapeTrigger{},
			ItemTrigger{},
		}

		for _, trigger := range sameTileTriggers {
			config := trigger.Config()
			assert.Equal(t, ActivationRangeSameTile, config.ActivationRange)
			assert.Equal(t, ActivationWayManual, config.ActivationWay,
				"直上タイルトリガーは手動発動である")
		}
	})

	t.Run("AdjacentトリガーはOnCollision方式", func(t *testing.T) {
		t.Parallel()
		// 仕様: 隣接タイルトリガー（Door, Talk）は衝突時自動発動
		adjacentTriggers := []TriggerData{
			DoorTrigger{},
			TalkTrigger{},
		}

		for _, trigger := range adjacentTriggers {
			config := trigger.Config()
			assert.Equal(t, ActivationRangeAdjacent, config.ActivationRange)
			assert.Equal(t, ActivationWayOnCollision, config.ActivationWay,
				"隣接タイルトリガーは衝突時自動発動である")
		}
	})

	t.Run("MeleeトリガーはAdjacent+OnCollision方式", func(t *testing.T) {
		t.Parallel()
		// 仕様: 近接攻撃トリガーは隣接タイルで衝突時自動発動
		trigger := MeleeTrigger{}
		config := trigger.Config()
		assert.Equal(t, ActivationRangeAdjacent, config.ActivationRange,
			"Meleeは隣接タイルで発動する")
		assert.Equal(t, ActivationWayOnCollision, config.ActivationWay,
			"Meleeは衝突時自動発動する")
	})
}
