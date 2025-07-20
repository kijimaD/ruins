package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultResourceManager(t *testing.T) {
	t.Parallel()
	rm := NewDefaultResourceManager()
	assert.NotNil(t, rm)

	// デフォルトマネージャーの型確認
	drm, ok := rm.(*DefaultResourceManager)
	require.True(t, ok)

	// デフォルトパスの確認
	assert.Equal(t, "metadata/fonts/fonts.toml", drm.config.FontsPath)
	assert.Equal(t, "metadata/spritesheets/spritesheets.toml", drm.config.SpriteSheetsPath)
	assert.Equal(t, "metadata/entities/raw/raw.toml", drm.config.RawsPath)
}

func TestNewResourceManager(t *testing.T) {
	t.Parallel()
	config := ResourceConfig{
		FontsPath:        "custom/fonts.toml",
		SpriteSheetsPath: "custom/sprites.toml",
		RawsPath:         "custom/raw.toml",
	}

	rm := NewResourceManager(config)
	assert.NotNil(t, rm)

	drm, ok := rm.(*DefaultResourceManager)
	require.True(t, ok)

	// カスタムパスの確認
	assert.Equal(t, config.FontsPath, drm.config.FontsPath)
	assert.Equal(t, config.SpriteSheetsPath, drm.config.SpriteSheetsPath)
	assert.Equal(t, config.RawsPath, drm.config.RawsPath)
}

func TestLoadFonts(t *testing.T) {
	t.Parallel()
	t.Run("正常にフォントを読み込める", func(t *testing.T) {
		t.Parallel()
		rm := NewDefaultResourceManager()
		fonts, err := rm.LoadFonts()

		assert.NoError(t, err)
		assert.NotNil(t, fonts)
		assert.Greater(t, len(fonts), 0)

		// キャッシュされていることの確認
		drm := rm.(*DefaultResourceManager)
		assert.Equal(t, fonts, drm.cache.Fonts)
	})

	t.Run("キャッシュから読み込む", func(t *testing.T) {
		t.Parallel()
		rm := NewDefaultResourceManager()

		// 1回目の読み込み
		fonts1, err1 := rm.LoadFonts()
		require.NoError(t, err1)

		// 2回目の読み込み（キャッシュから）
		fonts2, err2 := rm.LoadFonts()
		require.NoError(t, err2)

		// 同じオブジェクトを参照していることを確認
		assert.Equal(t, fonts1, fonts2)
	})

	t.Run("存在しないファイルパスの場合", func(t *testing.T) {
		t.Parallel()
		config := ResourceConfig{
			FontsPath: "invalid/path/fonts.toml",
		}
		rm := NewResourceManager(config)

		fonts, err := rm.LoadFonts()
		assert.Error(t, err)
		assert.Nil(t, fonts)
		assert.Contains(t, err.Error(), "フォントファイルの読み込みに失敗")
	})
}

func TestLoadSpriteSheets(t *testing.T) {
	t.Parallel()
	t.Run("正常にスプライトシートを読み込める", func(t *testing.T) {
		t.Parallel()
		rm := NewDefaultResourceManager()
		sprites, err := rm.LoadSpriteSheets()

		assert.NoError(t, err)
		assert.NotNil(t, sprites)
		assert.Greater(t, len(sprites), 0)

		// 各スプライトシートに名前が設定されていることを確認
		for name, sprite := range sprites {
			assert.Equal(t, name, sprite.Name)
		}

		// キャッシュされていることの確認
		drm := rm.(*DefaultResourceManager)
		assert.Equal(t, sprites, drm.cache.SpriteSheets)
	})
}

func TestLoadRaws(t *testing.T) {
	t.Parallel()
	t.Run("正常にRawデータを読み込める", func(t *testing.T) {
		t.Parallel()
		rm := NewDefaultResourceManager()
		rawMaster, err := rm.LoadRaws()

		assert.NoError(t, err)
		assert.NotNil(t, rawMaster)

		// キャッシュされていることの確認
		drm := rm.(*DefaultResourceManager)
		assert.Equal(t, rawMaster, drm.cache.RawMaster)
	})
}

func TestLoadAll(t *testing.T) {
	t.Parallel()
	t.Run("すべてのリソースを一括で読み込める", func(t *testing.T) {
		t.Parallel()
		rm := NewDefaultResourceManager()
		axes := []string{}
		actions := []string{}

		err := rm.LoadAll(axes, actions)
		assert.NoError(t, err)

		// すべてのリソースがキャッシュされていることを確認
		drm := rm.(*DefaultResourceManager)
		assert.NotNil(t, drm.cache.Fonts)
		assert.NotNil(t, drm.cache.SpriteSheets)
		assert.NotNil(t, drm.cache.RawMaster)
	})

	t.Run("一部のリソース読み込みに失敗した場合", func(t *testing.T) {
		t.Parallel()
		config := ResourceConfig{
			FontsPath:        "invalid/fonts.toml",
			SpriteSheetsPath: "metadata/spritesheets/spritesheets.toml",
			RawsPath:         "metadata/entities/raw/raw.toml",
		}
		rm := NewResourceManager(config)

		err := rm.LoadAll([]string{}, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "フォントの読み込みに失敗")
	})
}

func TestClearCache(t *testing.T) {
	t.Parallel()
	rm := NewDefaultResourceManager()
	drm := rm.(*DefaultResourceManager)

	// リソースを読み込んでキャッシュを作成
	err := rm.LoadAll([]string{}, []string{})
	require.NoError(t, err)

	// キャッシュが存在することを確認
	assert.NotNil(t, drm.cache.Fonts)

	// キャッシュをクリア
	drm.ClearCache()

	// キャッシュがクリアされていることを確認
	assert.Nil(t, drm.cache.Fonts)
	assert.Nil(t, drm.cache.SpriteSheets)
	assert.Nil(t, drm.cache.RawMaster)
}

func TestGetResourcePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		basePath string
		filename string
		expected string
	}{
		{
			name:     "通常のパス結合",
			basePath: "metadata/fonts",
			filename: "fonts.toml",
			expected: "metadata/fonts/fonts.toml",
		},
		{
			name:     "basePathが空の場合",
			basePath: "",
			filename: "fonts.toml",
			expected: "fonts.toml",
		},
		{
			name:     "filenameが空の場合",
			basePath: "metadata/fonts",
			filename: "",
			expected: "metadata/fonts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GetResourcePath(tt.basePath, tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}
