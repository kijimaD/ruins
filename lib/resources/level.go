package resources

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 冒険出発から帰還までを1セットとした情報を保持する。
// 冒険出発から帰還までは複数階層が存在し、複数階層を通しての情報を保持する必要がある。
type Game struct {
	// フィールド上で発生したイベント。各stateで補足されて処理される
	StateEvent StateEvent
	// 現在階のフィールド情報
	Level Level
	// 階層数
	Depth int
}

// 現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数。タイルは正方形のため、縦横で同じピクセル数になる
	TileSize int
	// タイルエンティティ群
	Entities []ecs.Entity
}

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx int, ty int) int {
	return ty*int(l.TileWidth) + tx
}

// タイルスライスのインデックスからタイル座標を求める
func (l *Level) XYTileCoord(idx int) (int, int) {
	x := idx % int(l.TileWidth)
	y := idx / int(l.TileWidth)
	return x, y
}

// xy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x int, y int) ecs.Entity {
	tx := x / l.TileSize
	ty := y / l.TileSize
	idx := l.XYTileIndex(tx, ty)

	return l.Entities[idx]
}

// ステージ幅。横の全体ピクセル数
func (l *Level) Width() int {
	return int(l.TileWidth) * l.TileSize
}

// ステージ縦。縦の全体ピクセル数
func (l *Level) Height() int {
	return int(l.TileHeight) * l.TileSize
}
