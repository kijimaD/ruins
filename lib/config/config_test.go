package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("本番プロファイルのデフォルト設定が正しく読み込まれる", func(t *testing.T) {
		t.Parallel()
		
		// 環境変数をクリア
		os.Unsetenv("RUINS_PROFILE")
		os.Unsetenv("RUINS_DEBUG")
		os.Unsetenv("RUINS_DEBUG_PPROF")
		os.Unsetenv("RUINS_PROFILE_MEMORY")
		defer func() {
			os.Unsetenv("RUINS_PROFILE")
			os.Unsetenv("RUINS_DEBUG")
			os.Unsetenv("RUINS_DEBUG_PPROF")
			os.Unsetenv("RUINS_PROFILE_MEMORY")
		}()
		
		cfg, err := Load()
		require.NoError(t, err)
		
		assert.Equal(t, ProfileProduction, cfg.Profile)
		assert.Equal(t, 960, cfg.WindowWidth)
		assert.Equal(t, 720, cfg.WindowHeight)
		assert.Equal(t, false, cfg.Fullscreen)
		assert.Equal(t, false, cfg.Debug) // 本番ではfalse
		assert.Equal(t, false, cfg.DebugPProf) // 本番ではfalse
		assert.Equal(t, 6060, cfg.PProfPort)
		assert.Equal(t, "main_menu", cfg.StartingState)
		assert.Equal(t, false, cfg.SkipIntro)
		assert.Equal(t, 60, cfg.TargetFPS)
		assert.Equal(t, false, cfg.ProfileMemory) // 本番ではfalse
		assert.Equal(t, false, cfg.ProfileCPU)
		assert.Equal(t, false, cfg.ProfileMutex)
		assert.Equal(t, false, cfg.ProfileTrace)
		assert.Equal(t, ".", cfg.ProfilePath)
	})
	
	t.Run("開発プロファイルのデフォルト設定が正しく読み込まれる", func(t *testing.T) {
		t.Parallel()
		
		// 全ての環境変数をクリアしてから開発プロファイルを設定
		os.Unsetenv("RUINS_WINDOW_WIDTH")
		os.Unsetenv("RUINS_WINDOW_HEIGHT")
		os.Unsetenv("RUINS_DEBUG")
		os.Unsetenv("RUINS_STARTING_STATE")
		os.Unsetenv("RUINS_TARGET_FPS")
		os.Setenv("RUINS_PROFILE", "development")
		defer func() {
			os.Unsetenv("RUINS_PROFILE")
			os.Unsetenv("RUINS_WINDOW_WIDTH")
			os.Unsetenv("RUINS_WINDOW_HEIGHT")
			os.Unsetenv("RUINS_DEBUG")
			os.Unsetenv("RUINS_STARTING_STATE")
			os.Unsetenv("RUINS_TARGET_FPS")
		}()
		
		cfg, err := Load()
		require.NoError(t, err)
		
		assert.Equal(t, ProfileDevelopment, cfg.Profile)
		assert.Equal(t, 800, cfg.WindowWidth) // 開発では小さめ
		assert.Equal(t, 600, cfg.WindowHeight) // 開発では小さめ
		assert.Equal(t, false, cfg.Fullscreen)
		assert.Equal(t, true, cfg.Debug) // 開発ではtrue
		assert.Equal(t, true, cfg.DebugPProf) // 開発ではtrue
		assert.Equal(t, 6060, cfg.PProfPort)
		assert.Equal(t, "debug_menu", cfg.StartingState) // 開発ではdebug_menu
		assert.Equal(t, true, cfg.SkipIntro) // 開発ではtrue
		assert.Equal(t, 60, cfg.TargetFPS)
		assert.Equal(t, true, cfg.ProfileMemory) // 開発ではtrue
		assert.Equal(t, false, cfg.ProfileCPU)
		assert.Equal(t, false, cfg.ProfileMutex)
		assert.Equal(t, false, cfg.ProfileTrace)
		assert.Equal(t, "./profiles", cfg.ProfilePath) // 開発では専用フォルダ
	})

	t.Run("環境変数がプロファイルデフォルトを上書きする", func(t *testing.T) {
		t.Parallel()
		
		// 全ての関連環境変数をクリア
		os.Unsetenv("RUINS_PROFILE")
		os.Unsetenv("RUINS_WINDOW_WIDTH")
		os.Unsetenv("RUINS_WINDOW_HEIGHT") 
		os.Unsetenv("RUINS_DEBUG")
		os.Unsetenv("RUINS_DEBUG_PPROF")
		os.Unsetenv("RUINS_STARTING_STATE")
		os.Unsetenv("RUINS_TARGET_FPS")
		os.Unsetenv("RUINS_PROFILE_MEMORY")
		
		// 本番プロファイルを設定してから個別の環境変数で上書き
		os.Setenv("RUINS_PROFILE", "production")
		os.Setenv("RUINS_WINDOW_WIDTH", "1280")
		os.Setenv("RUINS_WINDOW_HEIGHT", "1024")
		os.Setenv("RUINS_DEBUG", "true") // 本番でもデバッグを有効に
		os.Setenv("RUINS_STARTING_STATE", "debug_menu")
		os.Setenv("RUINS_TARGET_FPS", "120")
		defer func() {
			os.Unsetenv("RUINS_PROFILE")
			os.Unsetenv("RUINS_WINDOW_WIDTH")
			os.Unsetenv("RUINS_WINDOW_HEIGHT")
			os.Unsetenv("RUINS_DEBUG")
			os.Unsetenv("RUINS_DEBUG_PPROF")
			os.Unsetenv("RUINS_STARTING_STATE")
			os.Unsetenv("RUINS_TARGET_FPS")
			os.Unsetenv("RUINS_PROFILE_MEMORY")
		}()

		cfg, err := Load()
		require.NoError(t, err)

		assert.Equal(t, ProfileProduction, cfg.Profile)
		assert.Equal(t, 1280, cfg.WindowWidth) // 環境変数で上書き
		assert.Equal(t, 1024, cfg.WindowHeight) // 環境変数で上書き
		assert.Equal(t, true, cfg.Debug) // 環境変数で上書き
		assert.Equal(t, "debug_menu", cfg.StartingState) // 環境変数で上書き
		assert.Equal(t, 120, cfg.TargetFPS) // 環境変数で上書き
		// 他の設定は本番プロファイルのデフォルト
		assert.Equal(t, false, cfg.DebugPProf)
		assert.Equal(t, false, cfg.ProfileMemory)
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("無効な値が修正される", func(t *testing.T) {
		t.Parallel()
		
		cfg := &Config{
			WindowWidth:  100, // 最小値以下
			WindowHeight: 50,  // 最小値以下
			TargetFPS:    0,   // 無効
			PProfPort:    80,  // 範囲外
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 320, cfg.WindowWidth)
		assert.Equal(t, 240, cfg.WindowHeight)
		assert.Equal(t, 60, cfg.TargetFPS)
		assert.Equal(t, 6060, cfg.PProfPort)
	})

	t.Run("有効な値は変更されない", func(t *testing.T) {
		t.Parallel()
		
		cfg := &Config{
			WindowWidth:  1920,
			WindowHeight: 1080,
			TargetFPS:    144,
			PProfPort:    8080,
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 1920, cfg.WindowWidth)
		assert.Equal(t, 1080, cfg.WindowHeight)
		assert.Equal(t, 144, cfg.TargetFPS)
		assert.Equal(t, 8080, cfg.PProfPort)
	})
}

func TestManager(t *testing.T) {
	t.Parallel()

	t.Run("Get()がシングルトンとして動作する", func(t *testing.T) {
		// 注意: これはパラレルテストできない（シングルトンのため）
		Reset() // テスト用にリセット

		cfg1 := Get()
		cfg2 := Get()

		assert.Same(t, cfg1, cfg2, "Get()は同じインスタンスを返すべき")
	})

	t.Run("MustGet()が設定を返す", func(t *testing.T) {
		// 環境変数をクリア
		os.Unsetenv("RUINS_PROFILE")
		os.Unsetenv("RUINS_WINDOW_WIDTH")
		os.Unsetenv("RUINS_WINDOW_HEIGHT")
		os.Unsetenv("RUINS_DEBUG")
		os.Unsetenv("RUINS_STARTING_STATE")
		os.Unsetenv("RUINS_TARGET_FPS")
		defer func() {
			os.Unsetenv("RUINS_PROFILE")
			os.Unsetenv("RUINS_WINDOW_WIDTH")
			os.Unsetenv("RUINS_WINDOW_HEIGHT")
			os.Unsetenv("RUINS_DEBUG")
			os.Unsetenv("RUINS_STARTING_STATE")
			os.Unsetenv("RUINS_TARGET_FPS")
		}()
		
		Reset() // テスト用にリセット

		cfg := MustGet()
		assert.NotNil(t, cfg)
		assert.Equal(t, ProfileProduction, cfg.Profile) // デフォルトは本番
		assert.Equal(t, 960, cfg.WindowWidth) // 本番プロファイルのデフォルト値
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Profile:       ProfileDevelopment,
		WindowWidth:   1280,
		WindowHeight:  720,
		Debug:         true,
		StartingState: "debug_menu",
	}

	str := cfg.String()
	assert.Contains(t, str, "Profile: development")
	assert.Contains(t, str, "WindowWidth: 1280")
	assert.Contains(t, str, "WindowHeight: 720")
	assert.Contains(t, str, "Debug: true")
	assert.Contains(t, str, "StartingState: debug_menu")
}