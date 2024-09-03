package resources

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils/consts"
)

// UpdateGameLayoutはゲームウィンドウサイズを更新する。
func UpdateGameLayout(world w.World) (gc.Pixel, gc.Pixel) {
	const (
		offsetX       gc.Pixel = 0
		offsetY       gc.Pixel = 80
		minGridWidth  gc.Pixel = 30
		minGridHeight gc.Pixel = 20
	)

	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*consts.TileSize + offsetX
	gameHeight := gridHeight*consts.TileSize + offsetY

	world.Resources.ScreenDimensions.Width = int(gameWidth)
	world.Resources.ScreenDimensions.Height = int(gameHeight)

	return gameWidth, gameHeight
}
