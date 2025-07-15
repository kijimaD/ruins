package mapbuilder

import gc "github.com/kijimaD/ruins/lib/components"

// Rect は矩形を表す構造体
type Rect struct {
	X1 gc.Row
	X2 gc.Row
	Y1 gc.Col
	Y2 gc.Col
}

// Center は矩形の中心座標を返す
func (r *Rect) Center() (gc.Row, gc.Col) {
	x := (r.X1 + r.X2) / 2
	y := (r.Y1 + r.Y2) / 2

	return x, y
}
