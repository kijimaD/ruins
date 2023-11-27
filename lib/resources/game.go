package resources

import (
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

// グリッドレイアウト
type GridLayout struct {
	Width  int
	Height int
}

type Game struct {
	GridLayout GridLayout
}

// UpdateGameLayoutはゲームレイアウトを更新する
func UpdateGameLayout(world w.World, gridLayout *GridLayout) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	if gridLayout != nil {
		gridWidth = gridLayout.Width
		gridHeight = gridLayout.Height
	}

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	fadeOutSprite := &(*world.Resources.SpriteSheets)["intro-bg"].Sprites[0]
	fadeOutSprite.Width = gameWidth
	fadeOutSprite.Height = gameHeight

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
