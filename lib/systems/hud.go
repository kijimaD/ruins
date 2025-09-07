package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/widgets/hud"
	w "github.com/kijimaD/ruins/lib/world"
)

var (
	hudContainer *hud.Container // HUDコンテナ
)

// HUDSystem はゲームの HUD 情報を描画する
func HUDSystem(world w.World, screen *ebiten.Image) {
	// HUDコンテナの初期化（初回のみ）
	if hudContainer == nil {
		hudContainer = hud.NewContainer(world)
	}

	// HUDコンテナを更新・描画
	hudContainer.Update(world)
	hudContainer.Draw(world, screen)
}
