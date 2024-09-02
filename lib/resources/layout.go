package resources

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils/consts"
)

// UpdateGameLayoutはゲームウィンドウサイズを更新する。
func UpdateGameLayout(world w.World) (int, int) {
	const (
		offsetX       = 0
		offsetY       = 80
		minGridWidth  = 30
		minGridHeight = 20
	)

	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*consts.TileSize + offsetX
	gameHeight := gridHeight*consts.TileSize + offsetY

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
