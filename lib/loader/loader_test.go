package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultResourceLoader(t *testing.T) {
	t.Parallel()
	rl := NewResourceLoader()
	assert.NotNil(t, rl)

	// デフォルトローダーの型確認
	drl, ok := rl.(*DefaultResourceLoader)
	require.True(t, ok)

	// デフォルトパスの確認
	assert.Equal(t, "metadata/fonts/fonts.toml", drl.config.FontsPath)
	assert.Equal(t, "metadata/spritesheets/spritesheets.toml", drl.config.SpriteSheetsPath)
	assert.Equal(t, "metadata/entities/raw/raw.toml", drl.config.RawsPath)
}

func TestLoadFonts(t *testing.T) {
	t.Parallel()
	t.Run("正常にフォントを読み込める", func(t *testing.T) {
		t.Parallel()
		rl := NewResourceLoader()
		fonts, err := rl.LoadFonts()

		assert.NoError(t, err)
		assert.NotNil(t, fonts)
		assert.Greater(t, len(fonts), 0)

		// キャッシュされていることの確認
		drl := rl.(*DefaultResourceLoader)
		assert.Equal(t, fonts, drl.cache.Fonts)
	})

	t.Run("キャッシュから読み込む", func(t *testing.T) {
		t.Parallel()
		rl := NewResourceLoader()

		// 1回目の読み込み
		fonts1, err1 := rl.LoadFonts()
		require.NoError(t, err1)

		// 2回目の読み込み（キャッシュから）
		fonts2, err2 := rl.LoadFonts()
		require.NoError(t, err2)

		// 同じオブジェクトを参照していることを確認
		assert.Equal(t, fonts1, fonts2)
	})

}

func TestLoadSpriteSheets(t *testing.T) {
	t.Parallel()
	t.Run("正常にスプライトシートを読み込める", func(t *testing.T) {
		t.Parallel()
		rl := NewResourceLoader()
		sprites, err := rl.LoadSpriteSheets()

		assert.NoError(t, err)
		assert.NotNil(t, sprites)
		assert.Greater(t, len(sprites), 0)

		// 各スプライトシートに名前が設定されていることを確認
		for name, sprite := range sprites {
			assert.Equal(t, name, sprite.Name)
		}

		// キャッシュされていることの確認
		drl := rl.(*DefaultResourceLoader)
		assert.Equal(t, sprites, drl.cache.SpriteSheets)
	})
}

func TestLoadRaws(t *testing.T) {
	t.Parallel()
	t.Run("正常にRawデータを読み込める", func(t *testing.T) {
		t.Parallel()
		rl := NewResourceLoader()
		rawMaster, err := rl.LoadRaws()

		assert.NoError(t, err)
		assert.NotNil(t, rawMaster)

		// キャッシュされていることの確認
		drl := rl.(*DefaultResourceLoader)
		assert.Equal(t, rawMaster, drl.cache.RawMaster)
	})
}
