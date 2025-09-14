package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetActionCost(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		actionID ActionID
		expected int
	}{
		{"移動アクション", ActionMove, 100},
		{"待機アクション", ActionWait, 100},
		{"攻撃アクション", ActionAttack, 100},
		{"アイテム拾得", ActionPickupItem, 50},
		{"ワープアクション", ActionWarp, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cost := GetActionInfo(tt.actionID).MoveCost
			assert.Equal(t, tt.expected, cost, "アクションコストが期待値と一致")
		})
	}
}

func TestGetActionCostUnknown(t *testing.T) {
	t.Parallel()
	// 存在しないアクションIDに対してはActionNullInfo（コスト0）を返す
	unknownActionID := ActionID(999)
	cost := GetActionInfo(unknownActionID).MoveCost
	assert.Equal(t, 0, cost, "未知のアクションはActionNull（コスト0）")
}

func TestSetActionCost(t *testing.T) {
	t.Parallel()
	// テスト用にコストを変更
	originalCost := GetActionInfo(ActionMove).MoveCost
	newCost := 75

	SetActionCost(ActionMove, newCost)
	assert.Equal(t, newCost, GetActionInfo(ActionMove).MoveCost, "アクションコストが正しく設定される")

	// 元に戻す
	SetActionCost(ActionMove, originalCost)
	assert.Equal(t, originalCost, GetActionInfo(ActionMove).MoveCost, "アクションコストが元に戻る")
}
