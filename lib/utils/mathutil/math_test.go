package mathutil

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{"x is smaller", 3, 5, 3},
		{"y is smaller", 7, 2, 2},
		{"equal values", 4, 4, 4},
		{"negative values", -5, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.x, tt.y); got != tt.want {
				t.Errorf("Min(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestMinFloat(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		y    float64
		want float64
	}{
		{"x is smaller", 3.5, 5.2, 3.5},
		{"y is smaller", 7.8, 2.1, 2.1},
		{"equal values", 4.0, 4.0, 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.x, tt.y); got != tt.want {
				t.Errorf("Min(%f, %f) = %f, want %f", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{"x is larger", 5, 3, 5},
		{"y is larger", 2, 7, 7},
		{"equal values", 4, 4, 4},
		{"negative values", -3, -5, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.x, tt.y); got != tt.want {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		min   int
		max   int
		want  int
	}{
		{"value within range", 5, 0, 10, 5},
		{"value below min", -5, 0, 10, 0},
		{"value above max", 15, 0, 10, 10},
		{"value equals min", 0, 0, 10, 0},
		{"value equals max", 10, 0, 10, 10},
		{"negative range", -5, -10, -1, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.value, tt.min, tt.max); got != tt.want {
				t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestClampFloat(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		min   float64
		max   float64
		want  float64
	}{
		{"value within range", 5.5, 0.0, 10.0, 5.5},
		{"value below min", -5.5, 0.0, 10.0, 0.0},
		{"value above max", 15.5, 0.0, 10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.value, tt.min, tt.max); got != tt.want {
				t.Errorf("Clamp(%f, %f, %f) = %f, want %f", tt.value, tt.min, tt.max, got, tt.want)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name string
		x    int
		want int
	}{
		{"positive value", 5, 5},
		{"negative value", -5, 5},
		{"zero", 0, 0},
		{"large negative", -1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%d) = %d, want %d", tt.x, got, tt.want)
			}
		})
	}
}

func TestAbsFloat(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"positive value", 5.5, 5.5},
		{"negative value", -5.5, 5.5},
		{"zero", 0.0, 0.0},
		{"small negative", -0.001, 0.001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.x); got != tt.want {
				t.Errorf("Abs(%f) = %f, want %f", tt.x, got, tt.want)
			}
		})
	}
}
