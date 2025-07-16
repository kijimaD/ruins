package mathutil

// Min はxとyの小さい方を返す
func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max はxとyの大きい方を返す
func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Clamp はvalueを[min, max]の範囲に制限する
func Clamp[T int | float64](value, minVal, maxVal T) T {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Abs はxの絶対値を返す
func Abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
