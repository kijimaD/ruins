package vutil

import (
	"testing"
)

func TestNewVec2d(t *testing.T) {
	t.Run("valid dimensions", func(t *testing.T) {
		data := []int{1, 2, 3, 4, 5, 6}
		v, err := NewVec2d(2, 3, data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if v.NRows != 2 || v.NCols != 3 {
			t.Errorf("Expected dimensions 2x3, got %dx%d", v.NRows, v.NCols)
		}
	})

	t.Run("invalid dimensions - negative", func(t *testing.T) {
		data := []int{1, 2, 3, 4}
		_, err := NewVec2d(-1, 2, data)
		if err == nil {
			t.Error("Expected error for negative dimensions")
		}
	})

	t.Run("invalid dimensions - zero", func(t *testing.T) {
		data := []int{1, 2, 3, 4}
		_, err := NewVec2d(0, 2, data)
		if err == nil {
			t.Error("Expected error for zero dimensions")
		}
	})

	t.Run("mismatched data length", func(t *testing.T) {
		data := []int{1, 2, 3}
		_, err := NewVec2d(2, 3, data)
		if err == nil {
			t.Error("Expected error for mismatched data length")
		}
	})
}

func TestVec2d_Get(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	t.Run("valid index", func(t *testing.T) {
		result := v.Get(0, 1)
		if result == nil {
			t.Error("Expected non-nil result for valid index")
		}
		if *result != 2 {
			t.Errorf("Expected value 2, got %d", *result)
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		result := v.Get(5, 5)
		if result != nil {
			t.Error("Expected nil result for out of bounds index")
		}
	})

	t.Run("negative index", func(t *testing.T) {
		result := v.Get(-1, 0)
		if result != nil {
			t.Error("Expected nil result for negative index")
		}
	})
}

func TestVec2d_GetSafe(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	t.Run("valid index", func(t *testing.T) {
		result, err := v.GetSafe(1, 2)
		if err != nil {
			t.Errorf("Expected no error for valid index, got %v", err)
		}
		if result != 6 {
			t.Errorf("Expected value 6, got %d", result)
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		_, err := v.GetSafe(5, 5)
		if err == nil {
			t.Error("Expected error for out of bounds index")
		}
	})
}

func TestVec2d_Set(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	t.Run("valid index", func(t *testing.T) {
		err := v.Set(1, 1, 99)
		if err != nil {
			t.Errorf("Expected no error for valid index, got %v", err)
		}
		if v.Data[4] != 99 {
			t.Errorf("Expected value 99, got %d", v.Data[4])
		}
	})

	t.Run("out of bounds", func(t *testing.T) {
		err := v.Set(5, 5, 99)
		if err == nil {
			t.Error("Expected error for out of bounds index")
		}
	})
}

func TestVec2d_IsValidIndex(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	tests := []struct {
		name string
		row  int
		col  int
		want bool
	}{
		{"valid index (0,0)", 0, 0, true},
		{"valid index (1,2)", 1, 2, true},
		{"out of bounds row", 2, 1, false},
		{"out of bounds col", 1, 3, false},
		{"negative row", -1, 1, false},
		{"negative col", 1, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.IsValidIndex(tt.row, tt.col); got != tt.want {
				t.Errorf("IsValidIndex(%d, %d) = %v, want %v", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

func TestVec2d_Size(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	if got := v.Size(); got != 6 {
		t.Errorf("Size() = %d, want 6", got)
	}
}

func TestVec2d_Dimensions(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6}
	v, _ := NewVec2d(2, 3, data)

	rows, cols := v.Dimensions()
	if rows != 2 || cols != 3 {
		t.Errorf("Dimensions() = (%d, %d), want (2, 3)", rows, cols)
	}
}
