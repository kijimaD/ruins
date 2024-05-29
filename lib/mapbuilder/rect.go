package mapbuilder

type Rect struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

func (r *Rect) Center() (int, int) {
	x := (r.X1 + r.X2) / 2
	y := (r.Y1 + r.Y2) / 2

	return x, y
}
