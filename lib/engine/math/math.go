package math

// Vector2 type
type Vector2 struct {
	X float64
	Y float64
}

// VectorInt2 type
type VectorInt2 struct {
	X int
	Y int
}

// Min は、2つの整数のうち小さい方を返す
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max は、2つの整数のうち大きい方を返す
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Abs は整数の絶対値を返す
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Mod は数学的な剰余を返す
func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += Abs(b)
	}
	return m
}
