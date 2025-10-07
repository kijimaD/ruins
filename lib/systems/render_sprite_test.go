package systems

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

func TestSpriteImageCache(t *testing.T) {
	t.Parallel()
	t.Run("sprite image cache initialization", func(t *testing.T) {
		t.Parallel()
		assert.NotNil(t, spriteImageCache, "spriteImageCacheがnilになっている")
		assert.Empty(t, spriteImageCache, "spriteImageCacheが空でない")
	})

	t.Run("sprite image cache is map", func(t *testing.T) {
		t.Parallel()
		// キャッシュがmap型であることを確認
		cache := spriteImageCache
		expectedType := make(map[spriteImageCacheKey]*ebiten.Image)
		assert.IsType(t, expectedType, cache, "spriteImageCacheの型が正しくない")
	})
}

// spriteImageCacheの操作テスト（実際の画像なしでテスト）
func TestSpriteImageCacheOperations(t *testing.T) {
	t.Parallel()
	t.Run("cache operations", func(t *testing.T) {
		t.Parallel()
		// 初期状態の確認
		initialLen := len(spriteImageCache)

		testKey := spriteImageCacheKey{
			SpriteSheetName: "test_sheet",
			SpriteKey:       "test_sprite",
		}

		// キーが存在しないことを確認
		_, exists := spriteImageCache[testKey]
		assert.False(t, exists, "存在しないキーがtrueを返している")

		// キャッシュに値を設定（nilでテスト）
		spriteImageCache[testKey] = nil

		// キーが存在することを確認
		_, exists = spriteImageCache[testKey]
		assert.True(t, exists, "設定したキーが存在しない")

		// サイズが増えたことを確認
		assert.Equal(t, initialLen+1, len(spriteImageCache), "キャッシュサイズが正しくない")

		// キャッシュをクリア（テスト後の処理）
		delete(spriteImageCache, testKey)

		// 元の状態に戻ったことを確認
		assert.Equal(t, initialLen, len(spriteImageCache), "キャッシュクリア後のサイズが正しくない")
	})
}
