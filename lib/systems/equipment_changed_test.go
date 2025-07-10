package systems

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestMaxHP(t *testing.T) {
	t.Run("calculate max HP with base stats", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     10,
				Modifier: 0,
				Total:    10,
			},
			Strength: gc.Attribute{
				Base:     5,
				Modifier: 0,
				Total:    5,
			},
			Sensation: gc.Attribute{
				Base:     3,
				Modifier: 0,
				Total:    3,
			},
		}
		pools := &gc.Pools{Level: 1}

		result := maxHP(attrs, pools)
		// 30 + (10*8 + 5 + 3) * (1 + (1-1)*0.03) = 30 + 88 * 1 = 118
		expected := 118
		assert.Equal(t, expected, result, "maxHPの計算が正しくない")
	})

	t.Run("calculate max HP with level bonus", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     10,
				Modifier: 0,
				Total:    10,
			},
			Strength: gc.Attribute{
				Base:     5,
				Modifier: 0,
				Total:    5,
			},
			Sensation: gc.Attribute{
				Base:     3,
				Modifier: 0,
				Total:    3,
			},
		}
		pools := &gc.Pools{Level: 5}

		result := maxHP(attrs, pools)
		// 30 + (10*8 + 5 + 3) * (1 + (5-1)*0.03) = 30 + 88 * 1.12 = 30 + 98.56 = 128.56 -> 128
		expected := 128
		assert.Equal(t, expected, result, "レベルボーナス込みのmaxHPの計算が正しくない")
	})

	t.Run("calculate max HP with high stats", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     20,
				Modifier: 5,
				Total:    25,
			},
			Strength: gc.Attribute{
				Base:     15,
				Modifier: 3,
				Total:    18,
			},
			Sensation: gc.Attribute{
				Base:     10,
				Modifier: 2,
				Total:    12,
			},
		}
		pools := &gc.Pools{Level: 10}

		result := maxHP(attrs, pools)
		// 30 + (25*8 + 18 + 12) * (1 + (10-1)*0.03) = 30 + 230 * 1.27 = 30 + 292.1 = 322.1 -> 322
		expected := 322
		assert.Equal(t, expected, result, "高ステータスでのmaxHPの計算が正しくない")
	})
}

func TestMaxSP(t *testing.T) {
	t.Run("calculate max SP with base stats", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     10,
				Modifier: 0,
				Total:    10,
			},
			Dexterity: gc.Attribute{
				Base:     8,
				Modifier: 0,
				Total:    8,
			},
			Agility: gc.Attribute{
				Base:     6,
				Modifier: 0,
				Total:    6,
			},
		}
		pools := &gc.Pools{Level: 1}

		result := maxSP(attrs, pools)
		// (10*2 + 8 + 6) * (1 + (1-1)*0.02) = 34 * 1 = 34
		expected := 34
		assert.Equal(t, expected, result, "maxSPの計算が正しくない")
	})

	t.Run("calculate max SP with level bonus", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     10,
				Modifier: 0,
				Total:    10,
			},
			Dexterity: gc.Attribute{
				Base:     8,
				Modifier: 0,
				Total:    8,
			},
			Agility: gc.Attribute{
				Base:     6,
				Modifier: 0,
				Total:    6,
			},
		}
		pools := &gc.Pools{Level: 5}

		result := maxSP(attrs, pools)
		// (10*2 + 8 + 6) * (1 + (5-1)*0.02) = 34 * 1.08 = 36.72 -> 36
		expected := 36
		assert.Equal(t, expected, result, "レベルボーナス込みのmaxSPの計算が正しくない")
	})

	t.Run("calculate max SP with high stats", func(t *testing.T) {
		attrs := &gc.Attributes{
			Vitality: gc.Attribute{
				Base:     20,
				Modifier: 5,
				Total:    25,
			},
			Dexterity: gc.Attribute{
				Base:     15,
				Modifier: 3,
				Total:    18,
			},
			Agility: gc.Attribute{
				Base:     12,
				Modifier: 2,
				Total:    14,
			},
		}
		pools := &gc.Pools{Level: 10}

		result := maxSP(attrs, pools)
		// (25*2 + 18 + 14) * (1 + (10-1)*0.02) = 82 * 1.18 = 96.76 -> 96
		expected := 96
		assert.Equal(t, expected, result, "高ステータスでのmaxSPの計算が正しくない")
	})
}
