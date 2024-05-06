package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

// カメラを考慮した画像配置オプションをセットする
func SetTranslate(world w.World, op *ebiten.DrawImageOptions, cameraX int, cameraY int) {
	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置
	op.GeoM.Translate(float64(cameraX), float64(cameraY))
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
	op.GeoM.Scale(1, 1)
}
