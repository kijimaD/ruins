package mathutil

// Min returns the smaller of x or y.
func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger of x or y.
func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Clamp returns value clamped to the range [min, max].
func Clamp[T int | float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Abs returns the absolute value of x.
func Abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
