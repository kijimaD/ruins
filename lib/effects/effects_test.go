package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestDamage(t *testing.T) {
	t.Run("create damage effect", func(t *testing.T) {
		damage := Damage{Amount: 50}
		assert.Equal(t, 50, damage.Amount, "ダメージ量が正しく設定されない")

		// isEffectType インターフェースを実装していることを確認
		var effect EffectType = damage
		assert.NotNil(t, effect, "EffectTypeインターフェースを実装していない")
	})

	t.Run("zero damage", func(t *testing.T) {
		damage := Damage{Amount: 0}
		assert.Equal(t, 0, damage.Amount, "ゼロダメージが正しく設定されない")
	})

	t.Run("negative damage", func(t *testing.T) {
		damage := Damage{Amount: -10}
		assert.Equal(t, -10, damage.Amount, "負のダメージが正しく設定されない")
	})
}

func TestHealing(t *testing.T) {
	t.Run("create healing effect with numeral amount", func(t *testing.T) {
		healing := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
		assert.NotNil(t, healing.Amount, "回復量が設定されない")

		// isEffectType インターフェースを実装していることを確認
		var effect EffectType = healing
		assert.NotNil(t, effect, "EffectTypeインターフェースを実装していない")
	})

	t.Run("create healing effect with ratio amount", func(t *testing.T) {
		healing := Healing{Amount: gc.RatioAmount{Ratio: 0.5}}
		assert.NotNil(t, healing.Amount, "回復量が設定されない")
	})
}

func TestConsumptionStamina(t *testing.T) {
	t.Run("create stamina consumption effect", func(t *testing.T) {
		consumption := ConsumptionStamina{Amount: gc.NumeralAmount{Numeral: 20}}
		assert.NotNil(t, consumption.Amount, "消費量が設定されない")

		// isEffectType インターフェースを実装していることを確認
		var effect EffectType = consumption
		assert.NotNil(t, effect, "EffectTypeインターフェースを実装していない")
	})
}

func TestRecoveryStamina(t *testing.T) {
	t.Run("create stamina recovery effect", func(t *testing.T) {
		recovery := RecoveryStamina{Amount: gc.NumeralAmount{Numeral: 15}}
		assert.NotNil(t, recovery.Amount, "回復量が設定されない")

		// isEffectType インターフェースを実装していることを確認
		var effect EffectType = recovery
		assert.NotNil(t, effect, "EffectTypeインターフェースを実装していない")
	})
}

func TestItemUse(t *testing.T) {
	t.Run("create item use effect", func(t *testing.T) {
		item := ecs.Entity(123)
		itemUse := ItemUse{Item: item}
		assert.Equal(t, item, itemUse.Item, "アイテムが正しく設定されない")

		// isEffectType インターフェースを実装していることを確認
		var effect EffectType = itemUse
		assert.NotNil(t, effect, "EffectTypeインターフェースを実装していない")
	})
}
