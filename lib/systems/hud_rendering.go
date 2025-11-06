package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/widgets/hud"
	w "github.com/kijimaD/ruins/lib/world"
)

// HUDRenderingSystem はクエリを実行し、UIを描画するシステム
type HUDRenderingSystem struct {
	gameInfo        *hud.GameInfo
	minimap         *hud.Minimap
	debugOverlay    *hud.DebugOverlay
	messageArea     *hud.MessageArea
	currencyDisplay *hud.CurrencyDisplay
	enabled         bool
}

// NewHUDRenderingSystem は新しいHUD描画システムを作成する
func NewHUDRenderingSystem(world w.World) *HUDRenderingSystem {
	return &HUDRenderingSystem{
		gameInfo:        hud.NewGameInfo(world),
		minimap:         hud.NewMinimap(world),
		debugOverlay:    hud.NewDebugOverlay(world),
		messageArea:     hud.NewMessageArea(world),
		currencyDisplay: hud.NewCurrencyDisplay(world),
		enabled:         true,
	}
}

// Update はHUDシステム全体を更新する
func (sys *HUDRenderingSystem) Update(world w.World) {
	if sys == nil || !sys.enabled {
		return
	}

	// 各ウィジェットの更新処理
	if sys.gameInfo != nil {
		sys.gameInfo.Update(world)
	}
	if sys.minimap != nil {
		sys.minimap.Update(world)
	}
	if sys.debugOverlay != nil {
		sys.debugOverlay.Update(world)
	}
	if sys.messageArea != nil {
		sys.messageArea.Update(world)
	}
	if sys.currencyDisplay != nil {
		sys.currencyDisplay.Update(world)
	}
}

// Run はHUD描画システムのメイン処理を実行する
// worldからデータを抽出し、HUDを描画する
func (sys *HUDRenderingSystem) Run(world w.World, screen *ebiten.Image) {
	if sys == nil || !sys.enabled {
		return
	}

	// worldから全HUDデータを一括抽出
	hudData := ExtractHUDData(world)

	// 各ウィジェットにデータを渡して描画
	if sys.gameInfo != nil {
		sys.gameInfo.Draw(screen, hudData.GameInfo)
	}
	if sys.minimap != nil {
		sys.minimap.Draw(screen, hudData.MinimapData)
	}
	if sys.debugOverlay != nil {
		sys.debugOverlay.Draw(screen, hudData.DebugOverlay)
	}
	if sys.messageArea != nil {
		sys.messageArea.Draw(screen, hudData.MessageData)
	}
	if sys.currencyDisplay != nil {
		sys.currencyDisplay.Draw(screen, hudData.CurrencyData)
	}
}
