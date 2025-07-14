package states

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func getCenterWinRect() image.Rectangle {
	screenWidth, screenHeight := ebiten.WindowSize()
	windowWidth, windowHeight := 400, 400 // パーティ選択ウィンドウに合わせてサイズを拡大

	// ウィンドウの中心が画面の中心に来るように左上角の座標を計算
	x := screenWidth/2 - windowWidth/2
	y := screenHeight/2 - windowHeight/2

	rect := image.Rect(x, y, x+windowWidth, y+windowHeight)
	return rect
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
