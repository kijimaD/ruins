package gamelog

import (
	"testing"
)

func TestLoggerBasicUsage(t *testing.T) {
	// ログをクリア
	FieldLog.Clear()

	// メソッドチェーンでのログ作成をテスト
	New().
		Append("Player").
		Append(" attacks ").
		NPCName("Goblin").
		Append(" for ").
		Damage(15).
		Append(" damage!").
		Log(LogKindField)

	// ログが追加されたかチェック
	if FieldLog.Count() != 1 {
		t.Errorf("Expected 1 log entry, got %d", FieldLog.Count())
	}

	// テキストの内容をチェック
	messages := FieldLog.GetRecent(1)
	expected := "Player attacks Goblin for 15 damage!"
	if messages[0] != expected {
		t.Errorf("Expected '%s', got '%s'", expected, messages[0])
	}

	// 色付きエントリもチェック
	entries := FieldLog.GetRecentEntries(1)
	if len(entries) != 1 {
		t.Errorf("Expected 1 colored entry, got %d", len(entries))
	}

	entry := entries[0]
	if len(entry.Fragments) != 6 {
		t.Errorf("Expected 6 fragments, got %d", len(entry.Fragments))
	}

	// 各フラグメントをチェック
	expectedFragments := []struct {
		text  string
		color string
	}{
		{"Player", "white"},
		{" attacks ", "white"},
		{"Goblin", "yellow"},
		{" for ", "white"},
		{"15", "red"},
		{" damage!", "white"},
	}

	for i, expected := range expectedFragments {
		if entry.Fragments[i].Text != expected.text {
			t.Errorf("Fragment %d: expected text '%s', got '%s'", i, expected.text, entry.Fragments[i].Text)
		}
	}

	// NPCの名前が黄色、ダメージが赤色かチェック
	if entry.Fragments[2].Color != ColorYellow {
		t.Errorf("Expected NPC name to be yellow, got %v", entry.Fragments[2].Color)
	}
	if entry.Fragments[4].Color != ColorRed {
		t.Errorf("Expected damage to be red, got %v", entry.Fragments[4].Color)
	}
}

func TestLoggerColorMethod(t *testing.T) {
	FieldLog.Clear()

	// カスタム色での使用例
	New().
		Color(0, 255, 255). // Cyan
		Append("John").
		ColorRGBA(ColorWhite).
		Append(" considers attacking ").
		ColorRGBA(ColorCyan).
		Append("Orc").
		Log(LogKindField)

	entries := FieldLog.GetRecentEntries(1)
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}

	fragments := entries[0].Fragments
	if len(fragments) != 3 {
		t.Errorf("Expected 3 fragments, got %d", len(fragments))
	}

	// 色のチェック
	if fragments[0].Color != ColorCyan {
		t.Errorf("Expected first fragment to be cyan")
	}
	if fragments[1].Color != ColorWhite {
		t.Errorf("Expected second fragment to be white")
	}
	if fragments[2].Color != ColorCyan {
		t.Errorf("Expected third fragment to be cyan")
	}
}

func TestLoggerItemName(t *testing.T) {
	FieldLog.Clear()

	New().
		Append("You pick up ").
		ItemName("Iron Sword").
		Append(".").
		Log(LogKindField)

	entries := FieldLog.GetRecentEntries(1)
	fragments := entries[0].Fragments

	if fragments[1].Color != ColorCyan {
		t.Errorf("Expected item name to be cyan")
	}
	if fragments[1].Text != "Iron Sword" {
		t.Errorf("Expected item name 'Iron Sword', got '%s'", fragments[1].Text)
	}
}

func TestLoggerPlayerName(t *testing.T) {
	FieldLog.Clear()

	New().
		PlayerName("Hero").
		Append(" enters the dungeon").
		Log(LogKindField)

	entries := FieldLog.GetRecentEntries(1)
	fragments := entries[0].Fragments

	if fragments[0].Color != ColorGreen {
		t.Errorf("Expected player name to be green")
	}
	if fragments[0].Text != "Hero" {
		t.Errorf("Expected player name 'Hero', got '%s'", fragments[0].Text)
	}
}

func TestLoggerMultipleLogs(t *testing.T) {
	FieldLog.Clear()

	// 複数のログを追加
	New().Append("First message").Log(LogKindField)
	New().Append("Second message").Log(LogKindField)
	New().NPCName("Enemy").Append(" appears!").Log(LogKindField)

	if FieldLog.Count() != 3 {
		t.Errorf("Expected 3 log entries, got %d", FieldLog.Count())
	}

	entries := FieldLog.GetRecentEntries(3)
	if len(entries) != 3 {
		t.Errorf("Expected 3 colored entries, got %d", len(entries))
	}

	// 最後のエントリをチェック
	lastEntry := entries[2]
	if len(lastEntry.Fragments) != 2 {
		t.Errorf("Expected 2 fragments in last entry, got %d", len(lastEntry.Fragments))
	}
	if lastEntry.Fragments[0].Color != ColorYellow {
		t.Errorf("Expected enemy name to be yellow")
	}
}

func TestLoggerBattleLog(t *testing.T) {
	BattleLog.Clear()

	New().
		NPCName("Skeleton").
		Append(" attacks you for ").
		Damage(8).
		Append(" damage!").
		Log(LogKindBattle)

	if BattleLog.Count() != 1 {
		t.Errorf("Expected 1 battle log entry, got %d", BattleLog.Count())
	}

	entries := BattleLog.GetRecentEntries(1)
	if len(entries) != 1 {
		t.Errorf("Expected 1 colored battle entry, got %d", len(entries))
	}
}
