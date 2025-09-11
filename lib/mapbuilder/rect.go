package mapbuilder

import gc "github.com/kijimaD/ruins/lib/components"

// Rect は矩形を表す構造体
type Rect struct {
	X1 gc.Tile
	X2 gc.Tile
	Y1 gc.Tile
	Y2 gc.Tile
}

// Center は矩形の中心座標を返す
func (r *Rect) Center() (gc.Tile, gc.Tile) {
	x := (r.X1 + r.X2) / 2
	y := (r.Y1 + r.Y2) / 2

	return x, y
}
