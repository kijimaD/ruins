package systems

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcExpMultiplier(t *testing.T) {
	tests := []struct {
		input  int
		expect float64
	}{
		{
			input:  -5,
			expect: 0.59,
		},
		{
			input:  -4,
			expect: 0.66,
		},
		{
			input:  -3,
			expect: 0.73,
		},
		{
			input:  -2,
			expect: 0.81,
		},
		{
			input:  -1,
			expect: 0.9,
		},
		{
			input:  0,
			expect: 1.0,
		},
		{
			input:  1,
			expect: 1.08,
		},
		{
			input:  2,
			expect: 1.17,
		},
		{
			input:  3,
			expect: 1.26,
		},
		{
			input:  4,
			expect: 1.36,
		},
		{
			input:  5,
			expect: 1.47,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := roundUnder2(t, calcExpMultiplier(tt.input))
			assert.Equal(t, tt.expect, got)
		})
	}
}

func roundUnder2(t *testing.T, v float64) float64 {
	t.Helper()

	return math.Round(v*100) / 100
}
