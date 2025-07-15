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

// 共通の文字列定数
const (
	// UI表示用の定数
	TextNoDescription = "説明なし" // アイテムの説明がない場合の表示文字列
	TextClose         = "閉じる"  // メニューやウィンドウを閉じる際の表示文字列
)

// ItemCategoryType はアイテムのカテゴリーを表す
type ItemCategoryType string

var (
	// ItemCategoryTypeItem は道具を表す
	ItemCategoryTypeItem ItemCategoryType = "ITEM"
	// ItemCategoryTypeCard は手札を表す
	ItemCategoryTypeCard ItemCategoryType = "CARD"
	// ItemCategoryTypeWearable は装備を表す
	ItemCategoryTypeWearable ItemCategoryType = "WEARABLE"
	// ItemCategoryTypeMaterial は素材を表す
	ItemCategoryTypeMaterial ItemCategoryType = "MATERIAL"
)
