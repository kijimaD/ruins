package effects

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestParty(t *testing.T) {
	t.Run("create party target", func(t *testing.T) {
		party := Party{}

		// isTarget インターフェースを実装していることを確認
		var target Targets = party
		assert.NotNil(t, target, "Targetsインターフェースを実装していない")
	})
}

func TestSingle(t *testing.T) {
	t.Run("create single target", func(t *testing.T) {
		entity := ecs.Entity(456)
		single := Single{Target: entity}

		assert.Equal(t, entity, single.Target, "ターゲットエンティティが正しく設定されない")

		// isTarget インターフェースを実装していることを確認
		var target Targets = single
		assert.NotNil(t, target, "Targetsインターフェースを実装していない")
	})

	t.Run("create single target with zero entity", func(t *testing.T) {
		entity := ecs.Entity(0)
		single := Single{Target: entity}

		assert.Equal(t, entity, single.Target, "ゼロエンティティが正しく設定されない")
	})
}

func TestNone(t *testing.T) {
	t.Run("create none target", func(t *testing.T) {
		none := None{}

		// isTarget インターフェースを実装していることを確認
		var target Targets = none
		assert.NotNil(t, target, "Targetsインターフェースを実装していない")
	})
}

func TestTargetsInterface(t *testing.T) {
	t.Run("all target types implement interface", func(t *testing.T) {
		// すべてのターゲットタイプがTargetsインターフェースを実装していることを確認
		var targets []Targets = []Targets{
			Party{},
			Single{Target: ecs.Entity(123)},
			None{},
		}

		assert.Len(t, targets, 3, "すべてのターゲットタイプが確認されていない")

		for i, target := range targets {
			assert.NotNil(t, target, "ターゲット%dがnilである", i)
		}
	})

	t.Run("targets can be used in effect spawner", func(t *testing.T) {
		// 各ターゲットタイプがEffectSpawnerで使用できることを確認
		damage := Damage{Amount: 50}

		targets := []Targets{
			Party{},
			Single{Target: ecs.Entity(789)},
			None{},
		}

		for _, target := range targets {
			spawner := EffectSpawner{
				Creator:    nil,
				EffectType: damage,
				Targets:    target,
			}
			assert.NotNil(t, spawner.Targets, "ターゲットがEffectSpawnerで使用できない")
		}
	})
}
