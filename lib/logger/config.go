package logger

import (
	"os"
	"strings"
)

// Config はログ設定
type Config struct {
	// デフォルトログレベル
	DefaultLevel Level

	// コンテキスト別ログレベル
	ContextLevels map[Context]Level

	// タイムスタンプ形式
	TimeFormat string
}

// デフォルト設定
var defaultConfig = Config{
	DefaultLevel:  LevelInfo,
	ContextLevels: make(map[Context]Level),
	TimeFormat:    "2006-01-02T15:04:05.000Z",
}

// グローバル設定（初期化時に環境変数から読み込み）
var globalConfig = loadConfig()

// loadConfig は環境変数から設定を読み込む
func loadConfig() Config {
	config := defaultConfig

	// LOG_LEVEL環境変数
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.DefaultLevel = parseLevel(level)
	}

	// LOG_CONTEXTS環境変数 (例: "battle=debug,render=warn")
	if contexts := os.Getenv("LOG_CONTEXTS"); contexts != "" {
		config.ContextLevels = parseContextLevels(contexts)
	}

	return config
}

// parseContextLevels はコンテキスト別レベル設定を解析する
func parseContextLevels(s string) map[Context]Level {
	result := make(map[Context]Level)
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), "=")
		if len(parts) == 2 {
			context := Context(parts[0])
			level := parseLevel(parts[1])
			result[context] = level
		}
	}
	return result
}

// SetConfig はグローバル設定を更新する（主にテスト用）
func SetConfig(config Config) {
	globalConfig = config
}

// ResetConfig は設定をデフォルトに戻す（主にテスト用）
func ResetConfig() {
	globalConfig = loadConfig()
}