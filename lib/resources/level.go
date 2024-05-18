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
	// フィールド上で発生したイベント。各stateで補足されて処理される
	StateEvent StateEvent
	// 現在階のフィールド情報
	Level loader.Level
	// 階層数
	Depth int
}

// フィールド上でのイベント
type StateEvent string

const (
	StateEventNone       = StateEvent("NONE")
	StateEventWarpNext   = StateEvent("WARP_NEXT")
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
