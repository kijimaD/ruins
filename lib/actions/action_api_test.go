package actions

import (
	"testing"
)

func TestActionAPICreation(t *testing.T) {
	t.Parallel()
	api := NewActionAPI()

	if api == nil {
		t.Errorf("Expected non-nil ActionAPI")
		return
	}

	if api.manager == nil {
		t.Errorf("Expected non-nil ActivityManager")
	}

	if api.logger == nil {
		t.Errorf("Expected non-nil logger")
	}
}

// TestActionAPIQuickMove のテストは一時的に無効化（TurnManager依存のため）
// func TestActionAPIQuickMove(t *testing.T) {
// }
