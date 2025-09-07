package hud

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/config"
	w "github.com/kijimaD/ruins/lib/world"
)

// Container はすべてのHUDウィジェットを統合管理する
type Container struct {
	gameInfo     *GameInfo
	debugOverlay *DebugOverlay
	minimap      *Minimap
	messageArea  *MessageArea
	enabled      bool
}

// NewContainer は新しいContainerを作成する
func NewContainer(world w.World) *Container {
	return &Container{
		gameInfo:     NewGameInfo(),
		debugOverlay: NewDebugOverlay(),
		minimap:      NewMinimap(),
		messageArea:  NewMessageArea(world),
		enabled:      true,
	}
}

// SetEnabled は有効/無効を設定する
func (container *Container) SetEnabled(enabled bool) {
	container.enabled = enabled
}

// SetGameInfoEnabled はゲーム情報表示の有効/無効を設定する
func (container *Container) SetGameInfoEnabled(enabled bool) {
	if container.gameInfo != nil {
		container.gameInfo.SetEnabled(enabled)
	}
}

// SetDebugOverlayEnabled はデバッグオーバーレイの有効/無効を設定する
func (container *Container) SetDebugOverlayEnabled(enabled bool) {
	if container.debugOverlay != nil {
		container.debugOverlay.SetEnabled(enabled)
	}
}

// SetMinimapEnabled はミニマップの有効/無効を設定する
func (container *Container) SetMinimapEnabled(enabled bool) {
	if container.minimap != nil {
		container.minimap.SetEnabled(enabled)
	}
}

// SetMessageAreaEnabled はメッセージエリアの有効/無効を設定する
func (container *Container) SetMessageAreaEnabled(enabled bool) {
	if container.messageArea != nil {
		container.messageArea.SetEnabled(enabled)
	}
}

// Update はHUDコンテナを更新する
func (container *Container) Update(world w.World) {
	if !container.enabled {
		return
	}

	// 各ウィジェットを更新
	if container.gameInfo != nil {
		container.gameInfo.Update(world)
	}
	if container.debugOverlay != nil {
		container.debugOverlay.Update(world)
	}
	if container.minimap != nil {
		container.minimap.Update(world)
	}
	if container.messageArea != nil {
		container.messageArea.Update(world)
	}
}

// Draw はHUDコンテナを描画する
func (container *Container) Draw(world w.World, screen *ebiten.Image) {
	if !container.enabled {
		return
	}

	// 基本ゲーム情報を描画
	if container.gameInfo != nil {
		container.gameInfo.Draw(world, screen)
	}

	// AI デバッグ表示フラグが有効な時のみAI情報表示
	cfg := config.Get()
	if cfg.ShowAIDebug && container.debugOverlay != nil {
		container.debugOverlay.Draw(world, screen)
	}

	// ミニマップを描画
	if container.minimap != nil {
		container.minimap.Draw(world, screen)
	}

	// ログメッセージを描画
	if container.messageArea != nil {
		container.messageArea.Draw(world, screen)
	}
}
