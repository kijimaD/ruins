package gamelog

import "testing"

func TestLogKind_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		logKind  LogKind
		expected string
	}{
		{LogKindField, "Field"},
		{LogKindBattle, "Battle"},
		{LogKindScene, "Scene"},
		{LogKind(-1), "Unknown"},  // 無効な値
		{LogKind(999), "Unknown"}, // 定義されていない値
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			t.Parallel()
			actual := tc.logKind.String()
			if actual != tc.expected {
				t.Errorf("LogKind(%d).String() = %s, expected %s", tc.logKind, actual, tc.expected)
			}
		})
	}
}

func TestLogKind_Constants(t *testing.T) {
	t.Parallel()

	t.Run("定数値の確認", func(t *testing.T) {
		t.Parallel()
		if LogKindField != 0 {
			t.Errorf("Expected LogKindField to be 0, got %d", LogKindField)
		}
		if LogKindBattle != 1 {
			t.Errorf("Expected LogKindBattle to be 1, got %d", LogKindBattle)
		}
		if LogKindScene != 2 {
			t.Errorf("Expected LogKindScene to be 2, got %d", LogKindScene)
		}
	})
}
