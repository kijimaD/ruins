package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVector2(t *testing.T) {
	t.Run("create vector2", func(t *testing.T) {
		v := Vector2{X: 3.5, Y: 4.2}
		assert.Equal(t, 3.5, v.X, "Xの値が正しく設定されない")
		assert.Equal(t, 4.2, v.Y, "Yの値が正しく設定されない")
	})

	t.Run("zero vector2", func(t *testing.T) {
		v := Vector2{X: 0, Y: 0}
		assert.Equal(t, 0.0, v.X, "ゼロベクトルのXが正しくない")
		assert.Equal(t, 0.0, v.Y, "ゼロベクトルのYが正しくない")
	})

	t.Run("negative vector2", func(t *testing.T) {
		v := Vector2{X: -1.5, Y: -2.3}
		assert.Equal(t, -1.5, v.X, "負のXの値が正しく設定されない")
		assert.Equal(t, -2.3, v.Y, "負のYの値が正しく設定されない")
	})
}

func TestVectorInt2(t *testing.T) {
	t.Run("create vector int2", func(t *testing.T) {
		v := VectorInt2{X: 10, Y: 20}
		assert.Equal(t, 10, v.X, "Xの値が正しく設定されない")
		assert.Equal(t, 20, v.Y, "Yの値が正しく設定されない")
	})

	t.Run("zero vector int2", func(t *testing.T) {
		v := VectorInt2{X: 0, Y: 0}
		assert.Equal(t, 0, v.X, "ゼロベクトルのXが正しくない")
		assert.Equal(t, 0, v.Y, "ゼロベクトルのYが正しくない")
	})

	t.Run("negative vector int2", func(t *testing.T) {
		v := VectorInt2{X: -100, Y: -200}
		assert.Equal(t, -100, v.X, "負のXの値が正しく設定されない")
		assert.Equal(t, -200, v.Y, "負のYの値が正しく設定されない")
	})
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a is smaller", 3, 5, 3},
		{"b is smaller", 7, 2, 2},
		{"equal values", 4, 4, 4},
		{"negative values", -5, -3, -5},
		{"positive and negative", -5, 3, -5},
		{"zero and positive", 0, 5, 0},
		{"zero and negative", 0, -5, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "Min(%d, %d)の結果が正しくない", tt.a, tt.b)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a is larger", 5, 3, 5},
		{"b is larger", 2, 7, 7},
		{"equal values", 4, 4, 4},
		{"negative values", -3, -5, -3},
		{"positive and negative", -5, 3, 3},
		{"zero and positive", 0, 5, 5},
		{"zero and negative", 0, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "Max(%d, %d)の結果が正しくない", tt.a, tt.b)
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		expected int
	}{
		{"positive value", 5, 5},
		{"negative value", -5, 5},
		{"zero", 0, 0},
		{"large positive", 1000000, 1000000},
		{"large negative", -1000000, 1000000},
		{"min int boundary", -2147483647, 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.x)
			assert.Equal(t, tt.expected, result, "Abs(%d)の結果が正しくない", tt.x)
		})
	}
}

func TestMod(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"positive mod positive", 7, 3, 1},
		{"positive mod positive exact", 9, 3, 0},
		{"negative mod positive", -7, 3, 2},
		{"negative mod positive exact", -9, 3, 0},
		{"positive mod negative", 7, -3, 1},
		{"negative mod negative", -7, -3, 2},
		{"zero mod positive", 0, 5, 0},
		{"large numbers", 1000000, 7, 1},
		{"negative large numbers", -1000000, 7, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Mod(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "Mod(%d, %d)の結果が正しくない", tt.a, tt.b)
		})
	}
}

func TestModBehavior(t *testing.T) {
	t.Run("mod always returns non-negative", func(t *testing.T) {
		// Modは常に非負の値を返すことを確認
		testCases := []struct {
			a int
			b int
		}{
			{-10, 3},
			{-5, 2},
			{-100, 7},
			{10, -3},
			{-10, -3},
		}

		for _, tc := range testCases {
			result := Mod(tc.a, tc.b)
			assert.GreaterOrEqual(t, result, 0, "Mod(%d, %d)は非負の値を返すべき", tc.a, tc.b)
			assert.Less(t, result, Abs(tc.b), "Mod(%d, %d)の結果は|b|より小さいべき", tc.a, tc.b)
		}
	})

	t.Run("mod mathematical properties", func(t *testing.T) {
		// 数学的性質のテスト
		// (a + b) mod m = ((a mod m) + (b mod m)) mod m
		a, b, m := 13, 17, 5
		left := Mod(a+b, m)
		right := Mod(Mod(a, m)+Mod(b, m), m)
		assert.Equal(t, left, right, "Modの加法性が満たされない")
	})
}
