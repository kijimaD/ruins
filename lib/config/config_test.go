package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("無効な値が修正される", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			WindowWidth:  100,       // 最小値以下
			WindowHeight: 50,        // 最小値以下
			TargetFPS:    0,         // 無効
			PProfPort:    80,        // 範囲外
			LogLevel:     "invalid", // 無効なログレベル
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 320, cfg.WindowWidth)
		assert.Equal(t, 240, cfg.WindowHeight)
		assert.Equal(t, 60, cfg.TargetFPS)
		assert.Equal(t, 6060, cfg.PProfPort)
		assert.Equal(t, "info", cfg.LogLevel) // 無効な値はinfoに修正
	})

	t.Run("有効な値は変更されない", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			WindowWidth:  1920,
			WindowHeight: 1080,
			TargetFPS:    144,
			PProfPort:    8080,
			LogLevel:     "debug",
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 1920, cfg.WindowWidth)
		assert.Equal(t, 1080, cfg.WindowHeight)
		assert.Equal(t, 144, cfg.TargetFPS)
		assert.Equal(t, 8080, cfg.PProfPort)
		assert.Equal(t, "debug", cfg.LogLevel)
	})
}

//nolint:paralleltest // シングルトンテストのため並列実行不可
func TestManager(t *testing.T) {
	t.Run("Get()がシングルトンとして動作する", func(t *testing.T) {
		// 注意: これはパラレルテストできない（シングルトンのため）
		Reset() // テスト用にリセット

		cfg1 := Get()
		cfg2 := Get()

		assert.Same(t, cfg1, cfg2, "Get()は同じインスタンスを返すべき")
	})

	t.Run("MustGet()が設定を返す", func(t *testing.T) {
		Reset() // テスト用にリセット

		cfg := MustGet()
		assert.NotNil(t, cfg)
		assert.Equal(t, ProfileProduction, cfg.Profile) // デフォルトは本番
		assert.Equal(t, 960, cfg.WindowWidth)           // 本番プロファイルのデフォルト値
	})
}

func TestString(t *testing.T) {
	t.Parallel()
	cfg := &Config{
		Profile:       ProfileDevelopment,
		WindowWidth:   1280,
		WindowHeight:  720,
		Debug:         true,
		LogLevel:      "debug",
		LogCategories: "battle=debug,render=warn",
		StartingState: "debug_menu",
	}

	str := cfg.String()
	assert.Contains(t, str, "Profile: development")
	assert.Contains(t, str, "WindowWidth: 1280")
	assert.Contains(t, str, "WindowHeight: 720")
	assert.Contains(t, str, "Debug: true")
	assert.Contains(t, str, "LogLevel: debug")
	assert.Contains(t, str, "LogCategories: battle=debug,render=warn")
	assert.Contains(t, str, "StartingState: debug_menu")
}
