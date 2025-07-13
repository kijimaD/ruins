package effects

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestEffectSpawner(t *testing.T) {
	t.Run("create effect spawner", func(t *testing.T) {
		creator := ecs.Entity(456)
		effect := Damage{Amount: 100}
		target := Single{Target: ecs.Entity(789)}

		spawner := EffectSpawner{
			Creator:    &creator,
			EffectType: effect,
			Targets:    target,
		}

		assert.Equal(t, &creator, spawner.Creator, "クリエーターが正しく設定されない")
		assert.Equal(t, effect, spawner.EffectType, "エフェクトタイプが正しく設定されない")
		assert.Equal(t, target, spawner.Targets, "ターゲットが正しく設定されない")
	})

	t.Run("create effect spawner with nil creator", func(t *testing.T) {
		effect := Healing{Amount: gc.NumeralAmount{Numeral: 50}}
		target := Single{Target: ecs.Entity(123)}

		spawner := EffectSpawner{
			Creator:    nil,
			EffectType: effect,
			Targets:    target,
		}

		assert.Nil(t, spawner.Creator, "クリエーターがnilでない")
		assert.Equal(t, effect, spawner.EffectType, "エフェクトタイプが正しく設定されない")
	})
}

func TestAddEffect(t *testing.T) {
	t.Run("add effect to queue", func(t *testing.T) {
		// キューをクリア
		EffectQueue = []EffectSpawner{}

		creator := ecs.Entity(100)
		effect := Damage{Amount: 75}
		target := Single{Target: ecs.Entity(200)}

		// エフェクトを追加
		AddEffect(&creator, effect, target)

		// キューに追加されたことを確認
		assert.Len(t, EffectQueue, 1, "エフェクトがキューに追加されない")
		assert.Equal(t, &creator, EffectQueue[0].Creator, "クリエーターが正しく追加されない")
		assert.Equal(t, effect, EffectQueue[0].EffectType, "エフェクトタイプが正しく追加されない")
		assert.Equal(t, target, EffectQueue[0].Targets, "ターゲットが正しく追加されない")
	})

	t.Run("add multiple effects to queue", func(t *testing.T) {
		// キューをクリア
		EffectQueue = []EffectSpawner{}

		creator1 := ecs.Entity(101)
		effect1 := Damage{Amount: 25}
		target1 := Single{Target: ecs.Entity(201)}

		creator2 := ecs.Entity(102)
		effect2 := Healing{Amount: gc.NumeralAmount{Numeral: 30}}
		target2 := Single{Target: ecs.Entity(202)}

		// 複数のエフェクトを追加
		AddEffect(&creator1, effect1, target1)
		AddEffect(&creator2, effect2, target2)

		// キューに追加されたことを確認
		assert.Len(t, EffectQueue, 2, "複数のエフェクトがキューに追加されない")

		// 最初のエフェクト
		assert.Equal(t, &creator1, EffectQueue[0].Creator, "最初のクリエーターが正しくない")
		assert.Equal(t, effect1, EffectQueue[0].EffectType, "最初のエフェクトタイプが正しくない")
		assert.Equal(t, target1, EffectQueue[0].Targets, "最初のターゲットが正しくない")

		// 二番目のエフェクト
		assert.Equal(t, &creator2, EffectQueue[1].Creator, "二番目のクリエーターが正しくない")
		assert.Equal(t, effect2, EffectQueue[1].EffectType, "二番目のエフェクトタイプが正しくない")
		assert.Equal(t, target2, EffectQueue[1].Targets, "二番目のターゲットが正しくない")
	})

	t.Run("add effect with nil creator", func(t *testing.T) {
		// キューをクリア
		EffectQueue = []EffectSpawner{}

		effect := ConsumptionStamina{Amount: gc.NumeralAmount{Numeral: 10}}
		target := Single{Target: ecs.Entity(300)}

		// nilクリエーターでエフェクトを追加
		AddEffect(nil, effect, target)

		// キューに追加されたことを確認
		assert.Len(t, EffectQueue, 1, "nilクリエーターのエフェクトがキューに追加されない")
		assert.Nil(t, EffectQueue[0].Creator, "クリエーターがnilでない")
		assert.Equal(t, effect, EffectQueue[0].EffectType, "エフェクトタイプが正しく追加されない")
		assert.Equal(t, target, EffectQueue[0].Targets, "ターゲットが正しく追加されない")
	})
}

func TestEffectTypeInterface(t *testing.T) {
	t.Run("all effect types implement interface", func(t *testing.T) {
		// すべてのエフェクトタイプがEffectTypeインターフェースを実装していることを確認
		var effects []EffectType = []EffectType{
			Damage{Amount: 10},
			Healing{Amount: gc.NumeralAmount{Numeral: 20}},
			ConsumptionStamina{Amount: gc.NumeralAmount{Numeral: 5}},
			RecoveryStamina{Amount: gc.RatioAmount{Ratio: 0.3}},
			ItemUse{Item: ecs.Entity(999)},
		}

		assert.Len(t, effects, 5, "すべてのエフェクトタイプが確認されていない")

		for i, effect := range effects {
			assert.NotNil(t, effect, "エフェクト%dがnilである", i)
		}
	})
}
