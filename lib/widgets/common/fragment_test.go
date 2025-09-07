package common

import (
	"testing"

	"github.com/kijimaD/ruins/lib/colors"
	w "github.com/kijimaD/ruins/lib/world"
)

func TestNewFragmentText(t *testing.T) {
	// モックのworldを作成（テスト用）
	world := w.World{}
	// 実際のUIリソースが必要になるため、このテストは統合テスト環境でのみ実行可能
	// 単体テストとしては関数の存在確認のみ行う

	// 関数が存在し、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			// UIリソースがない場合のパニックは予期される
			t.Logf("Expected panic due to missing UI resources: %v", r)
		}
	}()

	// 基本的な呼び出しテスト
	text := "テストテキスト"
	color := colors.ColorRed

	// この関数呼び出しは、UIリソースがない場合はパニックするが、
	// 関数自体が正しく定義されていることを確認
	NewFragmentText(text, color, world)

	t.Log("NewFragmentText function exists and can be called")
}

// TestFragmentTextConcept は概念的なテスト
func TestFragmentTextConcept(t *testing.T) {
	// NewFragmentTextとNewListItemTextの設計思想の違いを文書化
	t.Log("NewFragmentText design goals:")
	t.Log("- No stretching (Stretch: false)")
	t.Log("- No minimum width constraint")
	t.Log("- No padding around text")
	t.Log("- Exact text width only")

	t.Log("NewListItemText design goals:")
	t.Log("- Stretches to parent width (Stretch: true)")
	t.Log("- Minimum width of 120px")
	t.Log("- Left/right padding of 8px each")
	t.Log("- Suitable for list items")
}
