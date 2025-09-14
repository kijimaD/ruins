package turns

import (
	"testing"

	"github.com/kijimaD/ruins/lib/actions"
	"github.com/stretchr/testify/assert"
)

func TestGetActionCost(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		actionID actions.ActionID
		expected int
	}{
		{"移動アクション", actions.ActionMove, 100},
		{"待機アクション", actions.ActionWait, 100},
		{"攻撃アクション", actions.ActionAttack, 100},
		{"アイテム拾得", actions.ActionPickupItem, 50},
		{"ワープアクション", actions.ActionWarp, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cost := actions.GetActionInfo(tt.actionID).MoveCost
			assert.Equal(t, tt.expected, cost, "アクションコストが期待値と一致")
		})
	}
}

func TestGetActionCostUnknown(t *testing.T) {
	t.Parallel()
	// 存在しないアクションIDに対してはActionNullInfo（コスト0）を返す
	unknownActionID := actions.ActionID(999)
	cost := actions.GetActionInfo(unknownActionID).MoveCost
	assert.Equal(t, 0, cost, "未知のアクションはActionNull（コスト0）")
}

func TestSetActionCost(t *testing.T) {
	t.Parallel()
	// テスト用にコストを変更
	originalCost := actions.GetActionInfo(actions.ActionMove).MoveCost
	newCost := 75

	actions.SetActionCost(actions.ActionMove, newCost)
	assert.Equal(t, newCost, actions.GetActionInfo(actions.ActionMove).MoveCost, "アクションコストが正しく設定される")

	// 元に戻す
	actions.SetActionCost(actions.ActionMove, originalCost)
	assert.Equal(t, originalCost, actions.GetActionInfo(actions.ActionMove).MoveCost, "アクションコストが元に戻る")
}
