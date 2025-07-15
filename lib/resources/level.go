package resources

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Game は冒険出発から帰還までを1セットとした情報を保持する。
// 冒険出発から帰還までは複数階層が存在し、複数階層を通しての情報を保持する必要がある。
type Game struct {
	// フィールド上で発生したイベント。各stateで補足されて処理される
	StateEvent StateEvent
	// 現在階のフィールド情報
	Level Level
	// 階層数
	Depth int
}

// Level は現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数。タイルは正方形のため、縦横で同じピクセル数になる
	TileSize gc.Pixel
	// タイルエンティティ群
	Entities []ecs.Entity
	// 視界を表現する黒背景
	// 階層移動でリセットされる
	VisionImage *ebiten.Image
}

// XYTileIndex はタイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx gc.Row, ty gc.Col) TileIdx {
	return TileIdx(int(ty)*int(l.TileWidth) + int(tx))
}

// XYTileCoord はタイルスライスのインデックスからタイル座標を求める
func (l *Level) XYTileCoord(idx TileIdx) (gc.Pixel, gc.Pixel) {
	x := int(idx) % int(l.TileWidth)
	y := int(idx) / int(l.TileWidth)

	return gc.Pixel(x), gc.Pixel(y)
}

// AtEntity はxy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x gc.Pixel, y gc.Pixel) ecs.Entity {
	tx := gc.Row(int(x) / int(l.TileSize))
	ty := gc.Col(int(y) / int(l.TileSize))
	idx := l.XYTileIndex(tx, ty)

	return l.Entities[idx]
}

// Width はステージ幅。横の全体ピクセル数
func (l *Level) Width() gc.Pixel {
	return gc.Pixel(int(l.TileWidth) * int(l.TileSize))
}

// Height はステージ縦。縦の全体ピクセル数
func (l *Level) Height() gc.Pixel {
	return gc.Pixel(int(l.TileHeight) * int(l.TileSize))
}
