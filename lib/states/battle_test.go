package states

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattleStateFixedEnemies(t *testing.T) {
	t.Parallel()

	t.Run("固定敵が設定される", func(t *testing.T) {
		t.Parallel()

		// 固定敵でBattleStateを作成
		state := &BattleState{
			FixedEnemies: []string{"軽戦車", "火の玉"},
		}

		assert.Equal(t, []string{"軽戦車", "火の玉"}, state.FixedEnemies, "固定敵が正しく設定されるべき")
		assert.Len(t, state.FixedEnemies, 2, "固定敵が2つ設定されるべき")
	})

	t.Run("空の固定敵リスト", func(t *testing.T) {
		t.Parallel()

		// 空の固定敵でBattleStateを作成
		state := &BattleState{
			FixedEnemies: []string{},
		}

		assert.Empty(t, state.FixedEnemies, "固定敵リストが空であるべき")
	})

	t.Run("固定敵なし", func(t *testing.T) {
		t.Parallel()

		// 固定敵なしでBattleStateを作成
		state := &BattleState{}

		assert.Nil(t, state.FixedEnemies, "固定敵がnilであるべき")
	})
}

func TestNewBattleStateWithEnemies(t *testing.T) {
	t.Parallel()

	t.Run("NewBattleStateWithEnemiesファクトリーが正しく動作する", func(t *testing.T) {
		t.Parallel()

		enemies := []string{"軽戦車", "火の玉"}
		factory := NewBattleStateWithEnemies(enemies)
		state := factory().(*BattleState)

		assert.Equal(t, enemies, state.FixedEnemies, "指定した敵リストが設定されるべき")
	})

	t.Run("空リストでファクトリーを作成", func(t *testing.T) {
		t.Parallel()

		enemies := []string{}
		factory := NewBattleStateWithEnemies(enemies)
		state := factory().(*BattleState)

		assert.Empty(t, state.FixedEnemies, "空の敵リストが設定されるべき")
	})

	t.Run("単一の敵でファクトリーを作成", func(t *testing.T) {
		t.Parallel()

		enemies := []string{"灰の偶像"}
		factory := NewBattleStateWithEnemies(enemies)
		state := factory().(*BattleState)

		assert.Equal(t, enemies, state.FixedEnemies, "単一の敵が設定されるべき")
		assert.Len(t, state.FixedEnemies, 1, "敵が1つ設定されるべき")
	})
}
