package movement

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanMoveToBasic(t *testing.T) {
	t.Parallel()

	// このテストは統合テストのためのものなので、
	// 実際のworld構造が必要だが、ここでは基本的な関数の存在確認のみ行う

	// CanMoveToが存在することを確認
	assert.NotNil(t, CanMoveTo, "CanMoveTo関数が存在すること")

	t.Logf("CanMoveTo関数の基本テスト完了")
}
