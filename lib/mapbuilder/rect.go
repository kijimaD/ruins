package mapbuilder

import gc "github.com/kijimaD/ruins/lib/components"

type Rect struct {
	X1 gc.Row
	X2 gc.Row
	Y1 gc.Col
	Y2 gc.Col
}

func (r *Rect) Center() (gc.Row, gc.Col) {
	x := (r.X1 + r.X2) / 2
	y := (r.Y1 + r.Y2) / 2

	return x, y
}
