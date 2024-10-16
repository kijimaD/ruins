package states

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	w "github.com/kijimaD/ruins/lib/engine/world"
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

// ================

type ItemCategoryType string

var (
	// 道具
	ItemCategoryTypeItem ItemCategoryType = "ITEM"
	// 手札
	ItemCategoryTypeCard ItemCategoryType = "CARD"
	// 装備
	ItemCategoryTypeWearable ItemCategoryType = "WEARABLE"
	// 素材
	ItemCategoryTypeMaterial ItemCategoryType = "MATERIAL"
)

// ================

// 単に実装形式を合わせるためのintarface
type haveCategory interface {
	setCategory(world w.World, category ItemCategoryType)
	categoryReload(world w.World)
}
