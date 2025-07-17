package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput はos.Stdoutの出力をキャプチャする
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestLoggerNew(t *testing.T) {
	logger := New(ContextBattle)
	if logger.context != ContextBattle {
		t.Errorf("期待値: %s, 実際: %s", ContextBattle, logger.context)
	}
	if len(logger.fields) != 0 {
		t.Errorf("fieldsは空であるべき")
	}
}

func TestLoggerWithField(t *testing.T) {
	logger := New(ContextBattle)
	newLogger := logger.WithField("key", "value")

	if len(logger.fields) != 0 {
		t.Errorf("元のロガーは変更されないべき")
	}
	if newLogger.fields["key"] != "value" {
		t.Errorf("フィールドが追加されていない")
	}
}

func TestLoggerWithFields(t *testing.T) {
	logger := New(ContextBattle)
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	newLogger := logger.WithFields(fields)

	if len(logger.fields) != 0 {
		t.Errorf("元のロガーは変更されないべき")
	}
	if newLogger.fields["key1"] != "value1" || newLogger.fields["key2"] != 42 {
		t.Errorf("フィールドが正しく追加されていない")
	}
}

func TestLogLevelFiltering(t *testing.T) {
	// テスト用の設定
	SetConfig(Config{
		DefaultLevel:  LevelInfo,
		ContextLevels: make(map[Context]Level),
		TimeFormat:    "2006-01-02T15:04:05.000Z",
	})
	defer ResetConfig()

	logger := New(ContextBattle)

	// Debugログは出力されない
	output := captureOutput(func() {
		logger.Debug("デバッグメッセージ")
	})
	if output != "" {
		t.Errorf("DEBUGレベルのログは出力されないべき")
	}

	// Infoログは出力される
	output = captureOutput(func() {
		logger.Info("情報メッセージ")
	})
	if output == "" {
		t.Errorf("INFOレベルのログは出力されるべき")
	}
}

func TestContextLevelFiltering(t *testing.T) {
	// コンテキスト別設定
	SetConfig(Config{
		DefaultLevel: LevelWarn,
		ContextLevels: map[Context]Level{
			ContextBattle: LevelDebug,
		},
		TimeFormat: "2006-01-02T15:04:05.000Z",
	})
	defer ResetConfig()

	// Battleコンテキストはデバッグレベルが有効
	battleLogger := New(ContextBattle)
	output := captureOutput(func() {
		battleLogger.Debug("戦闘デバッグ")
	})
	if output == "" {
		t.Errorf("Battleコンテキストのデバッグログは出力されるべき")
	}

	// Moveコンテキストはデフォルト（Warn）レベル
	moveLogger := New(ContextMove)
	output = captureOutput(func() {
		moveLogger.Info("移動情報")
	})
	if output != "" {
		t.Errorf("Moveコンテキストの情報ログは出力されないべき")
	}
}

func TestJSONOutput(t *testing.T) {
	SetConfig(Config{
		DefaultLevel:  LevelDebug,
		ContextLevels: make(map[Context]Level),
		TimeFormat:    "2006-01-02T15:04:05.000Z",
	})
	defer ResetConfig()

	logger := New(ContextBattle)
	output := captureOutput(func() {
		logger.Info("テストメッセージ", "key1", "value1", "key2", 42)
	})

	// JSON解析
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("JSON解析エラー: %v", err)
	}

	// 必須フィールドの確認
	if entry["level"] != "INFO" {
		t.Errorf("levelが正しくない: %v", entry["level"])
	}
	if entry["context"] != "battle" {
		t.Errorf("contextが正しくない: %v", entry["context"])
	}
	if entry["message"] != "テストメッセージ" {
		t.Errorf("messageが正しくない: %v", entry["message"])
	}
	if entry["key1"] != "value1" {
		t.Errorf("key1が正しくない: %v", entry["key1"])
	}
	if entry["key2"] != float64(42) { // JSONでは数値はfloat64になる
		t.Errorf("key2が正しくない: %v", entry["key2"])
	}
}

func TestIsDebugEnabled(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		context       Context
		expectEnabled bool
	}{
		{
			name: "デフォルトレベルがDebug",
			config: Config{
				DefaultLevel:  LevelDebug,
				ContextLevels: make(map[Context]Level),
			},
			context:       ContextBattle,
			expectEnabled: true,
		},
		{
			name: "デフォルトレベルがInfo",
			config: Config{
				DefaultLevel:  LevelInfo,
				ContextLevels: make(map[Context]Level),
			},
			context:       ContextBattle,
			expectEnabled: false,
		},
		{
			name: "コンテキスト別設定でDebug有効",
			config: Config{
				DefaultLevel: LevelInfo,
				ContextLevels: map[Context]Level{
					ContextBattle: LevelDebug,
				},
			},
			context:       ContextBattle,
			expectEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetConfig(tt.config)
			defer ResetConfig()

			logger := New(tt.context)
			if logger.IsDebugEnabled() != tt.expectEnabled {
				t.Errorf("期待値: %v, 実際: %v", tt.expectEnabled, logger.IsDebugEnabled())
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"fatal", LevelFatal},
		{"FATAL", LevelFatal},
		{"unknown", LevelInfo}, // デフォルト
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%q) = %v, 期待値: %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseContextLevels(t *testing.T) {
	input := "battle=debug,render=warn,invalid"
	result := parseContextLevels(input)

	if result[ContextBattle] != LevelDebug {
		t.Errorf("battleコンテキストのレベルが正しくない")
	}
	if result[ContextRender] != LevelWarn {
		t.Errorf("renderコンテキストのレベルが正しくない")
	}
	if _, exists := result["invalid"]; exists {
		t.Errorf("無効な形式は無視されるべき")
	}
}

func TestLoggerOutput(t *testing.T) {
	SetConfig(Config{
		DefaultLevel:  LevelDebug,
		ContextLevels: make(map[Context]Level),
		TimeFormat:    "2006-01-02T15:04:05.000Z",
	})
	defer ResetConfig()

	logger := New(ContextBattle).WithField("session", "test123")

	// 各レベルのテスト
	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		level    string
		contains []string
	}{
		{
			name:     "Debug",
			logFunc:  logger.Debug,
			level:    "DEBUG",
			contains: []string{"デバッグメッセージ", "DEBUG", "battle", "session", "test123"},
		},
		{
			name:     "Info",
			logFunc:  logger.Info,
			level:    "INFO",
			contains: []string{"情報メッセージ", "INFO", "battle"},
		},
		{
			name:     "Warn",
			logFunc:  logger.Warn,
			level:    "WARN",
			contains: []string{"警告メッセージ", "WARN", "battle"},
		},
		{
			name:     "Error",
			logFunc:  logger.Error,
			level:    "ERROR",
			contains: []string{"エラーメッセージ", "ERROR", "battle"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.logFunc(tt.contains[0])
			})

			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("出力に %q が含まれていない: %s", expected, output)
				}
			}
		})
	}
}