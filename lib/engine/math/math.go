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

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Mod(a, b int) int {
	m := a % b
	if m < 0 {
		m += Abs(b)
	}
	return m
}
