package states

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// 2つ以上windowを開くときに同じ位置で開くために使う
	winRect image.Rectangle
)

func setWinRect() image.Rectangle {
	x, y := ebiten.CursorPosition()
	winRect = image.Rect(0, 0, x, y)
	winRect = winRect.Add(image.Point{x + 20, y + 20})

	return winRect
}

func getWinRect() image.Rectangle {
	return winRect
}
