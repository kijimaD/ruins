package gamelog

import (
	"testing"

	"github.com/kijimaD/ruins/lib/colors"
)

func TestLogEntry_Text(t *testing.T) {
	t.Parallel()

	t.Run("複数フラグメントのテキスト結合", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{
				{Color: colors.ColorWhite, Text: "Hello "},
				{Color: colors.ColorRed, Text: "World"},
				{Color: colors.ColorWhite, Text: "!"},
			},
		}

		expected := "Hello World!"
		actual := entry.Text()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("単一フラグメントのテキスト", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{
				{Color: colors.ColorGreen, Text: "Success"},
			},
		}

		expected := "Success"
		actual := entry.Text()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("空のフラグメントリスト", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{},
		}

		expected := ""
		actual := entry.Text()

		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestLogEntry_IsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("空のエントリ", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{},
		}

		if !entry.IsEmpty() {
			t.Error("Expected entry to be empty")
		}
	})

	t.Run("nilフラグメント", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: nil,
		}

		if !entry.IsEmpty() {
			t.Error("Expected entry with nil fragments to be empty")
		}
	})

	t.Run("非空のエントリ", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{
				{Color: colors.ColorWhite, Text: "Not empty"},
			},
		}

		if entry.IsEmpty() {
			t.Error("Expected entry to not be empty")
		}
	})

	t.Run("空テキストのフラグメントがある場合", func(t *testing.T) {
		t.Parallel()
		entry := LogEntry{
			Fragments: []LogFragment{
				{Color: colors.ColorWhite, Text: ""},
			},
		}

		// フラグメントが存在するので空ではない
		if entry.IsEmpty() {
			t.Error("Expected entry with empty text fragment to not be empty")
		}
	})
}
