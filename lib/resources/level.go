package resources

import (
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/loader"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

type Game struct {
	StateEvent StateEvent
	Level      loader.Level
}

// フィールド上でのイベント
type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpEscape = StateEvent("WARP_ESCAPE")
)

// UpdateGameLayoutはゲームウィンドウサイズを更新する
func UpdateGameLayout(world w.World) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
