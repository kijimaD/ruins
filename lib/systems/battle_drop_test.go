package systems

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestCalcExpMultiplier(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got := roundUnder2(t, calcExpMultiplier(tt.input))
			assert.Equal(t, tt.expect, got)
		})
	}
}

func roundUnder2(t *testing.T, v float64) float64 {
	t.Helper()

	return math.Round(v*100) / 100
}

func TestDropResult(t *testing.T) {
	t.Parallel()
	t.Run("create drop result", func(t *testing.T) {
		t.Parallel()
		result := DropResult{
			MaterialNames: []string{"鉄の鉱石", "銅の鉱石"},
			XPBefore:      map[ecs.Entity]int{},
			XPAfter:       map[ecs.Entity]int{},
			IsLevelUp:     map[ecs.Entity]bool{},
		}

		assert.Len(t, result.MaterialNames, 2, "素材名の数が正しくない")
		assert.Contains(t, result.MaterialNames, "鉄の鉱石", "鉄の鉱石が含まれていない")
		assert.Contains(t, result.MaterialNames, "銅の鉱石", "銅の鉱石が含まれていない")
		assert.NotNil(t, result.XPBefore, "XPBeforeがnilになっている")
		assert.NotNil(t, result.XPAfter, "XPAfterがnilになっている")
		assert.NotNil(t, result.IsLevelUp, "IsLevelUpがnilになっている")
	})

	t.Run("empty drop result", func(t *testing.T) {
		t.Parallel()
		result := DropResult{
			MaterialNames: []string{},
			XPBefore:      map[ecs.Entity]int{},
			XPAfter:       map[ecs.Entity]int{},
			IsLevelUp:     map[ecs.Entity]bool{},
		}

		assert.Empty(t, result.MaterialNames, "素材名が空でない")
		assert.Empty(t, result.XPBefore, "XPBeforeが空でない")
		assert.Empty(t, result.XPAfter, "XPAfterが空でない")
		assert.Empty(t, result.IsLevelUp, "IsLevelUpが空でない")
	})
}

func TestLevelUpThreshold(t *testing.T) {
	t.Parallel()
	t.Run("level up threshold constant", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 100, LevelUpThreshold, "レベルアップの閾値が正しくない")
	})
}
