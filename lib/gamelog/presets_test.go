package gamelog

import (
	"testing"
)

func TestPresetFunctions(t *testing.T) {
	FieldLog.Clear()

	// プリセット関数のテスト
	New().
		Success("勝利しました！").
		Log(LogKindField)

	New().
		Warning("注意が必要です").
		Log(LogKindField)

	New().
		Error("エラーが発生しました").
		Log(LogKindField)

	New().
		PlayerName("Hero").
		Append("が").
		Location("洞窟").
		Append("で").
		ItemName("宝箱").
		Append("を発見した。").
		Log(LogKindField)

	New().
		Action("攻撃").
		Append("で").
		NPCName("ゴブリン").
		Append("に").
		Damage(25).
		Append("ダメージ！").
		Log(LogKindField)

	New().
		Money("500").
		Append("G獲得した。").
		Log(LogKindField)

	// ログ数の確認
	if FieldLog.Count() != 6 {
		t.Errorf("Expected 6 log entries, got %d", FieldLog.Count())
	}

	// 色付きエントリの確認
	entries := FieldLog.GetRecentEntries(6)
	if len(entries) != 6 {
		t.Errorf("Expected 6 colored entries, got %d", len(entries))
	}

	// 各エントリの色の確認
	testCases := []struct {
		entryIndex    int
		fragmentIndex int
		expectedColor string
		expectedText  string
	}{
		{0, 0, "green", "勝利しました！"},
		{1, 0, "yellow", "注意が必要です"},
		{2, 0, "red", "エラーが発生しました"},
		{3, 0, "green", "Hero"},  // PlayerName
		{3, 2, "orange", "洞窟"},   // Location
		{3, 4, "cyan", "宝箱"},     // ItemName
		{4, 0, "purple", "攻撃"},   // Action
		{4, 2, "yellow", "ゴブリン"}, // NPCName
		{4, 4, "red", "25"},      // Damage
		{5, 0, "yellow", "500"},  // Money
	}

	for _, tc := range testCases {
		if tc.entryIndex >= len(entries) {
			continue
		}
		entry := entries[tc.entryIndex]
		if tc.fragmentIndex >= len(entry.Fragments) {
			continue
		}
		fragment := entry.Fragments[tc.fragmentIndex]

		if fragment.Text != tc.expectedText {
			t.Errorf("Entry %d, Fragment %d: expected text '%s', got '%s'",
				tc.entryIndex, tc.fragmentIndex, tc.expectedText, fragment.Text)
		}

		// 色の確認（簡単なチェック）
		switch tc.expectedColor {
		case "green":
			if fragment.Color != ColorGreen {
				t.Errorf("Expected green color for '%s'", tc.expectedText)
			}
		case "yellow":
			if fragment.Color != ColorYellow {
				t.Errorf("Expected yellow color for '%s'", tc.expectedText)
			}
		case "red":
			if fragment.Color != ColorRed {
				t.Errorf("Expected red color for '%s'", tc.expectedText)
			}
		case "orange":
			if fragment.Color != ColorOrange {
				t.Errorf("Expected orange color for '%s'", tc.expectedText)
			}
		case "cyan":
			if fragment.Color != ColorCyan {
				t.Errorf("Expected cyan color for '%s'", tc.expectedText)
			}
		case "purple":
			if fragment.Color != ColorPurple {
				t.Errorf("Expected purple color for '%s'", tc.expectedText)
			}
		}
	}
}

func TestBattlePresets(t *testing.T) {
	BattleLog.Clear()

	// 戦闘専用プリセット
	New().
		Encounter("強敵が現れた！").
		Log(LogKindBattle)

	New().
		Victory("勝利した！").
		Log(LogKindBattle)

	New().
		Defeat("敗北した...").
		Log(LogKindBattle)

	New().
		Magic("ファイアボール").
		Append("を唱えた！").
		Log(LogKindBattle)

	if BattleLog.Count() != 4 {
		t.Errorf("Expected 4 battle log entries, got %d", BattleLog.Count())
	}

	entries := BattleLog.GetRecentEntries(4)

	// Encounter は赤色
	if entries[0].Fragments[0].Color != ColorRed {
		t.Errorf("Expected red color for Encounter")
	}

	// Victory は緑色
	if entries[1].Fragments[0].Color != ColorGreen {
		t.Errorf("Expected green color for Victory")
	}

	// Defeat は赤色
	if entries[2].Fragments[0].Color != ColorRed {
		t.Errorf("Expected red color for Defeat")
	}

	// Magic は紫色
	if entries[3].Fragments[0].Color != ColorMagenta {
		t.Errorf("Expected magenta color for Magic")
	}
}

func TestSystemPresets(t *testing.T) {
	FieldLog.Clear()

	// システム関連のプリセット
	New().
		System("システムが初期化されました").
		Log(LogKindField)

	entries := FieldLog.GetRecentEntries(1)
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	// System は水色（シアン）
	if entries[0].Fragments[0].Color != ColorCyan {
		t.Errorf("Expected cyan color for System")
	}
}
