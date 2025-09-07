package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	w "github.com/kijimaD/ruins/lib/world"
)

var (
	hudRenderingSystem *HUDRenderingSystem // メモ
)

// HUDSystem はゲームの HUD 情報を描画する
func HUDSystem(world w.World, screen *ebiten.Image) {
	// HUDRenderingSystemの初期化（初回のみ）
	if hudRenderingSystem == nil {
		hudRenderingSystem = NewHUDRenderingSystem(world)
	}

	// HUDRenderingSystemを更新・描画
	hudRenderingSystem.Update(world)
	hudRenderingSystem.Run(world, screen)
}
