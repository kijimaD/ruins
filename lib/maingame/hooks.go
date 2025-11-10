package maingame

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	w "github.com/kijimaD/ruins/lib/world"
)

// afterDrawHook は各stateのDraw完了後に呼ばれるフック関数
func afterDrawHook(stateIndex, stateCount int, _ w.World, screen *ebiten.Image) error {
	// stateが複数あるとき初回のstate描画後に1回だけオーバーレイを描画する
	if stateCount > 1 && stateIndex == 0 {
		bounds := screen.Bounds()
		width := float32(bounds.Dx())
		height := float32(bounds.Dy())

		// 現在の画面をコピー
		src := ebiten.NewImage(bounds.Dx(), bounds.Dy())
		src.DrawImage(screen, nil)

		// ブラー効果を2パス（水平→垂直）で適用
		// 水平ブラー
		const blurRadius = 4
		tmp := ebiten.NewImage(bounds.Dx(), bounds.Dy())
		for x := -blurRadius; x <= blurRadius; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), 0)
			op.ColorScale.ScaleAlpha(1.0 / float32(blurRadius*2+1))
			op.Blend = ebiten.BlendSourceOver
			tmp.DrawImage(src, op)
		}

		// 垂直ブラー
		screen.Clear()
		for y := -blurRadius; y <= blurRadius; y++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, float64(y))
			op.ColorScale.ScaleAlpha(1.0 / float32(blurRadius*2+1))
			op.Blend = ebiten.BlendSourceOver
			screen.DrawImage(tmp, op)
		}

		// 半透明の黒をオーバーレイ
		overlayColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
		vector.FillRect(screen, 0, 0, width, height, overlayColor, false)
	}

	return nil
}
