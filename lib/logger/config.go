package logger

import (
	"os"
	"strings"
)

// Config はログ設定
type Config struct {
	// デフォルトログレベル
	DefaultLevel Level

	// カテゴリ別ログレベル
	CategoryLevels map[Category]Level

	// タイムスタンプ形式
	TimeFormat string
}

// グローバル設定（初期化時は空の設定）
var globalConfig = Config{
	CategoryLevels: make(map[Category]Level),
	TimeFormat:     "2006-01-02T15:04:05.000Z",
}


// parseCategoryLevels はカテゴリ別レベル設定を解析する
func parseCategoryLevels(s string) map[Category]Level {
	result := make(map[Category]Level)
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.Split(strings.TrimSpace(pair), "=")
		if len(parts) == 2 {
			category := Category(parts[0])
			level := parseLevel(parts[1])
			result[category] = level
		}
	}
	return result
}

// SetConfig はグローバル設定を更新する（主にテスト用）
func SetConfig(config Config) {
	globalConfig = config
}

// ResetConfig は設定をリセットする（主にテスト用）
func ResetConfig() {
	globalConfig = Config{
		CategoryLevels: make(map[Category]Level),
		TimeFormat:     "2006-01-02T15:04:05.000Z",
	}
}

// SetLogLevelFromConfig はconfigパッケージの設定からログレベルを設定する
func SetLogLevelFromConfig(logLevel string) {
	config := globalConfig
	config.DefaultLevel = parseLevel(logLevel)
	globalConfig = config
}

// LoadFromConfig はconfigパッケージの設定を読み込む
func LoadFromConfig(logLevel string) {
	config := Config{
		DefaultLevel:   parseLevel(logLevel),
		CategoryLevels: make(map[Category]Level),
		TimeFormat:     "2006-01-02T15:04:05.000Z",
	}
	
	// LOG_CATEGORIES環境変数 (例: "battle=debug,render=warn") 
	if categories := os.Getenv("LOG_CATEGORIES"); categories != "" {
		config.CategoryLevels = parseCategoryLevels(categories)
	}

	globalConfig = config
}
