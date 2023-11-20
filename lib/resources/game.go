package resources

import (
	w "github.com/x-hgg-x/goecsengine/world"
)

// グリッドレイアウト
type GridLayout struct {
	Width  int
	Height int
}

type Game struct {
	GridLayout GridLayout
}

func UpdateGameLayout(world w.World, gridLayout *GridLayout) (int, int) {
	return 100, 100
}
