package maingame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	w "github.com/kijimaD/ruins/lib/world"
)

// afterDrawHook は複数stateがスタックされている場合にオーバーレイを描画するフック関数
func afterDrawHook(stateIndex, stateCount int, _ w.World, screen *ebiten.Image) error {
	// stateが複数あるとき初回のstate描画後に1回だけオーバーレイを描画する
	if stateCount > 1 && stateIndex == 0 {
		bounds := screen.Bounds()
		width := float32(bounds.Dx())
		height := float32(bounds.Dy())

		// 半透明の黒
		overlayColor := color.NRGBA{R: 0, G: 0, B: 0, A: 140}

		vector.FillRect(screen, 0, 0, width, height, overlayColor, false)
	}

	return nil
}
