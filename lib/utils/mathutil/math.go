package mathutil

// xとyの小さい方を返す
func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// xとyの大きい方を返す
func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// valueを[min, max]の範囲に制限する
func Clamp[T int | float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// xの絶対値を返す
func Abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
