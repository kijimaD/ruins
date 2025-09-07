package systems

import (
	"testing"

	"github.com/kijimaD/ruins/lib/widgets/hud"
	"github.com/stretchr/testify/assert"
)

// TestTileColorInfo はTileColorInfoの型エイリアスをテスト
func TestTileColorInfo(t *testing.T) {
	colorInfo := TileColorInfo{
		R: 255,
		G: 128,
		B: 64,
		A: 200,
	}

	// hud.TileColorInfoと同じ構造であることを確認
	var hudColorInfo hud.TileColorInfo = colorInfo

	assert.Equal(t, uint8(255), hudColorInfo.R)
	assert.Equal(t, uint8(128), hudColorInfo.G)
	assert.Equal(t, uint8(64), hudColorInfo.B)
	assert.Equal(t, uint8(200), hudColorInfo.A)
}

// TestExtractMessagesFromGameLog はメッセージ抽出関数をテスト
func TestExtractMessagesFromGameLog(t *testing.T) {
	// 現在の実装では空のスライスを返す（簡略化版）
	messages := extractMessagesFromGameLog()

	assert.NotNil(t, messages)
	assert.Empty(t, messages)
}

// TestHUDDataStructureConsistency はHUDDataの構造が一貫していることを確認
func TestHUDDataStructureConsistency(t *testing.T) {
	// hud.HUDDataがextractorで使用される各種データ型と一致することを確認
	var hudData hud.HUDData

	// GameInfoData
	hudData.GameInfo = hud.GameInfoData{
		FloorNumber: 1,
		PlayerSpeed: 1.0,
	}

	// MinimapData
	hudData.MinimapData = hud.MinimapData{
		PlayerTileX:   0,
		PlayerTileY:   0,
		ExploredTiles: make(map[string]bool),
		TileColors:    make(map[string]hud.TileColorInfo),
		MinimapConfig: hud.MinimapConfig{
			Width:  100,
			Height: 100,
			Scale:  2,
		},
		ScreenDimensions: hud.ScreenDimensions{
			Width:  800,
			Height: 600,
		},
	}

	// DebugOverlayData
	hudData.DebugOverlay = hud.DebugOverlayData{
		Enabled:            false,
		AIStates:           []hud.AIStateInfo{},
		VisionRanges:       []hud.VisionRangeInfo{},
		MovementDirections: []hud.MovementDirectionInfo{},
		ScreenDimensions: hud.ScreenDimensions{
			Width:  800,
			Height: 600,
		},
	}

	// MessageData
	hudData.MessageData = hud.MessageData{
		Messages:         []string{},
		ScreenDimensions: hud.ScreenDimensions{Width: 800, Height: 600},
		Config:           hud.DefaultMessageAreaConfig(),
	}

	// 構造体が正しく作成できることを確認
	assert.Equal(t, 1, hudData.GameInfo.FloorNumber)
	assert.Equal(t, 100, hudData.MinimapData.MinimapConfig.Width)
	assert.False(t, hudData.DebugOverlay.Enabled)
	assert.Equal(t, 800, hudData.MessageData.ScreenDimensions.Width)
}

// TestDefaultMessageAreaConfig はデフォルト設定との互換性を確認
func TestDefaultMessageAreaConfig(t *testing.T) {
	config := hud.DefaultMessageAreaConfig()

	// デフォルト値の確認
	expectedValues := map[string]int{
		"LogAreaHeight": 120,
		"MaxLogLines":   5,
		"LogAreaMargin": 8,
		"LineHeight":    20,
		"YPadding":      8,
	}

	assert.Equal(t, expectedValues["LogAreaHeight"], config.LogAreaHeight)
	assert.Equal(t, expectedValues["MaxLogLines"], config.MaxLogLines)
	assert.Equal(t, expectedValues["LogAreaMargin"], config.LogAreaMargin)
	assert.Equal(t, expectedValues["LineHeight"], config.LineHeight)
	assert.Equal(t, expectedValues["YPadding"], config.YPadding)
}

// BenchmarkHUDDataCreation はHUDData作成のパフォーマンスをベンチマーク
func BenchmarkHUDDataCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hudData := hud.HUDData{
			GameInfo: hud.GameInfoData{
				FloorNumber: 1,
				PlayerSpeed: 1.5,
			},
			MinimapData: hud.MinimapData{
				PlayerTileX:      10,
				PlayerTileY:      10,
				ExploredTiles:    map[string]bool{"10,10": true},
				TileColors:       map[string]hud.TileColorInfo{"10,10": {100, 100, 100, 255}},
				MinimapConfig:    hud.MinimapConfig{Width: 200, Height: 200, Scale: 4},
				ScreenDimensions: hud.ScreenDimensions{Width: 800, Height: 600},
			},
			DebugOverlay: hud.DebugOverlayData{
				Enabled:            true,
				AIStates:           []hud.AIStateInfo{{ScreenX: 400, ScreenY: 300, StateText: "ROAMING"}},
				VisionRanges:       []hud.VisionRangeInfo{{ScreenX: 400, ScreenY: 300, ScaledRadius: 50}},
				MovementDirections: []hud.MovementDirectionInfo{},
				ScreenDimensions:   hud.ScreenDimensions{Width: 800, Height: 600},
			},
			MessageData: hud.MessageData{
				Messages:         []string{"Test message"},
				ScreenDimensions: hud.ScreenDimensions{Width: 800, Height: 600},
				Config:           hud.DefaultMessageAreaConfig(),
			},
		}
		_ = hudData // 未使用変数警告を回避
	}
}
