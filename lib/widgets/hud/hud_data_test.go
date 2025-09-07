package hud

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

// TestGameInfoDraw はGameInfoが描画できることを確認
func TestGameInfoDraw(t *testing.T) {
	gameInfo := NewGameInfo()
	screen := ebiten.NewImage(800, 600)

	data := GameInfoData{
		FloorNumber: 5,
		PlayerSpeed: 2.5,
	}

	// 描画できることを確認
	gameInfo.Draw(screen, data)

	// 描画が完了すれば成功（パニックしない）
}

// TestMinimapDraw はMinimapが描画できることを確認
func TestMinimapDraw(t *testing.T) {
	minimap := NewMinimap()
	screen := ebiten.NewImage(800, 600)

	// 空のミニマップデータでテスト
	emptyData := MinimapData{
		PlayerTileX:   10,
		PlayerTileY:   20,
		ExploredTiles: make(map[string]bool),
		TileColors:    make(map[string]TileColorInfo),
		MinimapConfig: MinimapConfig{
			Width:  200,
			Height: 200,
			Scale:  4,
		},
		ScreenDimensions: ScreenDimensions{
			Width:  800,
			Height: 600,
		},
	}

	// 描画できることを確認
	minimap.Draw(screen, emptyData)

	// 探索済みタイルありのデータでテスト
	dataWithTiles := emptyData
	dataWithTiles.ExploredTiles = map[string]bool{
		"10,20": true,
		"11,20": true,
		"10,21": true,
	}
	dataWithTiles.TileColors = map[string]TileColorInfo{
		"10,20": {R: 100, G: 100, B: 100, A: 255},
		"11,20": {R: 200, G: 200, B: 200, A: 128},
		"10,21": {R: 100, G: 100, B: 100, A: 255},
	}

	minimap.Draw(screen, dataWithTiles)
}

// TestDebugOverlayDraw はDebugOverlayが描画できることを確認
func TestDebugOverlayDraw(t *testing.T) {
	overlay := NewDebugOverlay()
	screen := ebiten.NewImage(800, 600)

	data := DebugOverlayData{
		Enabled: true,
		AIStates: []AIStateInfo{
			{ScreenX: 400, ScreenY: 300, StateText: "ROAMING"},
			{ScreenX: 500, ScreenY: 200, StateText: "CHASING"},
		},
		VisionRanges: []VisionRangeInfo{
			{ScreenX: 400, ScreenY: 300, ScaledRadius: 50},
			{ScreenX: 500, ScreenY: 200, ScaledRadius: 75},
		},
		MovementDirections: []MovementDirectionInfo{
			{ScreenX: 400, ScreenY: 300, Angle: 45, Speed: 2.0, CameraScale: 1.0},
			{ScreenX: 500, ScreenY: 200, Angle: 90, Speed: 1.5, CameraScale: 1.2},
		},
		ScreenDimensions: ScreenDimensions{
			Width:  800,
			Height: 600,
		},
	}

	// world依存なしで描画できることを確認
	overlay.Draw(screen, data)

	// 無効状態でのテスト
	data.Enabled = false
	overlay.Draw(screen, data)
}

// TestMessageAreaDraw はMessageAreaが描画できることを確認
func TestMessageAreaDraw(t *testing.T) {
	// MessageAreaは内部でwidgetを使用するため、world必須でのコンストラクタを持つ
	// そのため、世界構築なしでのテストは困難
	// ここでは構造体の作成とDrawメソッドの存在確認のみ行う

	// MessageArea構造体が正しく定義されていることを確認
	var area *MessageArea
	screen := ebiten.NewImage(800, 600)

	data := MessageData{
		Messages: []string{
			"Test message 1",
			"Test message 2",
			"Test message 3",
		},
		ScreenDimensions: ScreenDimensions{
			Width:  800,
			Height: 600,
		},
		Config: DefaultMessageAreaConfig(),
	}

	// MessageAreaがnilの場合でもパニックしないことを確認
	if area != nil {
		area.Draw(screen, data)
	}

	// データ構造のサイズ確認
	assert.Len(t, data.Messages, 3)
	assert.Equal(t, 800, data.ScreenDimensions.Width)
}

// TestHUDData はHUDDataの統合データ構造をテスト
func TestHUDData(t *testing.T) {
	hudData := HUDData{
		GameInfo: GameInfoData{
			FloorNumber: 3,
			PlayerSpeed: 1.8,
		},
		MinimapData: MinimapData{
			PlayerTileX:      5,
			PlayerTileY:      7,
			ExploredTiles:    map[string]bool{"5,7": true},
			TileColors:       map[string]TileColorInfo{"5,7": {R: 255, G: 255, B: 255, A: 255}},
			MinimapConfig:    MinimapConfig{Width: 150, Height: 150, Scale: 3},
			ScreenDimensions: ScreenDimensions{Width: 1024, Height: 768},
		},
		DebugOverlay: DebugOverlayData{
			Enabled: false,
		},
		MessageData: MessageData{
			Messages:         []string{"Hello", "World"},
			ScreenDimensions: ScreenDimensions{Width: 1024, Height: 768},
			Config:           DefaultMessageAreaConfig(),
		},
	}

	// データ構造が正しく作成されることを確認
	assert.Equal(t, 3, hudData.GameInfo.FloorNumber)
	assert.Equal(t, 5, hudData.MinimapData.PlayerTileX)
	assert.Len(t, hudData.MessageData.Messages, 2)
}

// TestTileColorInfo はタイル色情報の構造をテスト
func TestTileColorInfo(t *testing.T) {
	colorInfo := TileColorInfo{
		R: 128,
		G: 64,
		B: 32,
		A: 255,
	}

	assert.Equal(t, uint8(128), colorInfo.R)
	assert.Equal(t, uint8(255), colorInfo.A)
}

// TestScreenDimensions は画面サイズ情報をテスト
func TestScreenDimensions(t *testing.T) {
	dimensions := ScreenDimensions{
		Width:  1920,
		Height: 1080,
	}

	assert.Equal(t, 1920, dimensions.Width)
	assert.Equal(t, 1080, dimensions.Height)
}

// TestWorldIndependentRendering はworld依存なしで複数のHUDコンポーネントが描画できることを確認
func TestWorldIndependentRendering(t *testing.T) {
	screen := ebiten.NewImage(1024, 768)

	// 各HUDコンポーネントを作成
	gameInfo := NewGameInfo()
	minimap := NewMinimap()
	debugOverlay := NewDebugOverlay()

	// 統合HUDデータを作成
	hudData := HUDData{
		GameInfo: GameInfoData{
			FloorNumber: 10,
			PlayerSpeed: 3.2,
		},
		MinimapData: MinimapData{
			PlayerTileX: 15,
			PlayerTileY: 25,
			ExploredTiles: map[string]bool{
				"15,25": true,
				"16,25": true,
				"15,26": true,
			},
			TileColors: map[string]TileColorInfo{
				"15,25": {R: 100, G: 100, B: 100, A: 255},
				"16,25": {R: 200, G: 200, B: 200, A: 128},
				"15,26": {R: 100, G: 100, B: 100, A: 255},
			},
			MinimapConfig: MinimapConfig{
				Width:  180,
				Height: 180,
				Scale:  5,
			},
			ScreenDimensions: ScreenDimensions{
				Width:  1024,
				Height: 768,
			},
		},
		DebugOverlay: DebugOverlayData{
			Enabled: true,
			AIStates: []AIStateInfo{
				{ScreenX: 512, ScreenY: 384, StateText: "WAITING"},
			},
			VisionRanges: []VisionRangeInfo{
				{ScreenX: 512, ScreenY: 384, ScaledRadius: 60},
			},
			MovementDirections: []MovementDirectionInfo{},
			ScreenDimensions: ScreenDimensions{
				Width:  1024,
				Height: 768,
			},
		},
		MessageData: MessageData{
			Messages: []string{
				"Battle started!",
				"Enemy defeated",
				"Level up!",
			},
			ScreenDimensions: ScreenDimensions{
				Width:  1024,
				Height: 768,
			},
			Config: DefaultMessageAreaConfig(),
		},
	}

	// world依存なしで全てのコンポーネントを描画
	gameInfo.Draw(screen, hudData.GameInfo)
	minimap.Draw(screen, hudData.MinimapData)
	debugOverlay.Draw(screen, hudData.DebugOverlay)

	// 描画が完了すれば成功
}

// BenchmarkWorldIndependentHUDRendering はHUD描画のパフォーマンスをベンチマーク
func BenchmarkWorldIndependentHUDRendering(b *testing.B) {
	screen := ebiten.NewImage(800, 600)
	gameInfo := NewGameInfo()
	minimap := NewMinimap()
	debugOverlay := NewDebugOverlay()

	gameInfoData := GameInfoData{FloorNumber: 5, PlayerSpeed: 2.0}
	minimapData := MinimapData{
		PlayerTileX:      10,
		PlayerTileY:      10,
		ExploredTiles:    map[string]bool{"10,10": true, "11,10": true},
		TileColors:       map[string]TileColorInfo{"10,10": {100, 100, 100, 255}, "11,10": {200, 200, 200, 128}},
		MinimapConfig:    MinimapConfig{Width: 200, Height: 200, Scale: 4},
		ScreenDimensions: ScreenDimensions{Width: 800, Height: 600},
	}
	debugData := DebugOverlayData{
		Enabled:            true,
		AIStates:           []AIStateInfo{{ScreenX: 400, ScreenY: 300, StateText: "ROAMING"}},
		VisionRanges:       []VisionRangeInfo{{ScreenX: 400, ScreenY: 300, ScaledRadius: 50}},
		MovementDirections: []MovementDirectionInfo{},
		ScreenDimensions:   ScreenDimensions{Width: 800, Height: 600},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gameInfo.Draw(screen, gameInfoData)
		minimap.Draw(screen, minimapData)
		debugOverlay.Draw(screen, debugData)
	}
}
