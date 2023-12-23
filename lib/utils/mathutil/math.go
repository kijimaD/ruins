package mathutil

func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}
