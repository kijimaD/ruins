package systems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattleExtinctionType(t *testing.T) {
	t.Run("battle extinction type constants", func(t *testing.T) {
		assert.Equal(t, BattleExtinctionType(0), BattleExtinctionNone, "BattleExtinctionNoneの値が正しくない")
		assert.Equal(t, BattleExtinctionType(1), BattleExtinctionAlly, "BattleExtinctionAllyの値が正しくない")
		assert.Equal(t, BattleExtinctionType(2), BattleExtinctionMonster, "BattleExtinctionMonsterの値が正しくない")
	})

	t.Run("battle extinction type comparison", func(t *testing.T) {
		assert.True(t, BattleExtinctionNone < BattleExtinctionAlly, "BattleExtinctionNone < BattleExtinctionAllyではない")
		assert.True(t, BattleExtinctionAlly < BattleExtinctionMonster, "BattleExtinctionAlly < BattleExtinctionMonsterではない")
		assert.True(t, BattleExtinctionNone < BattleExtinctionMonster, "BattleExtinctionNone < BattleExtinctionMonsterではない")
	})
}
