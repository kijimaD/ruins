package loader

import (
	"testing"

	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/resources"
	"github.com/stretchr/testify/assert"
)

func TestFontMetadata(t *testing.T) {
	t.Run("create font metadata", func(t *testing.T) {
		metadata := fontMetadata{
			Fonts: map[string]resources.Font{
				"main": resources.Font{
					// Fontの具体的なフィールドは resources.Font の実装に依存
				},
				"ui": resources.Font{},
			},
		}

		assert.Len(t, metadata.Fonts, 2, "フォントの数が正しくない")
		assert.Contains(t, metadata.Fonts, "main", "mainフォントが存在しない")
		assert.Contains(t, metadata.Fonts, "ui", "uiフォントが存在しない")
	})

	t.Run("empty font metadata", func(t *testing.T) {
		metadata := fontMetadata{
			Fonts: map[string]resources.Font{},
		}

		assert.Empty(t, metadata.Fonts, "空のフォントマップが空でない")
	})
}

func TestSpriteSheetMetadata(t *testing.T) {
	t.Run("create sprite sheet metadata", func(t *testing.T) {
		metadata := spriteSheetMetadata{
			SpriteSheets: map[string]c.SpriteSheet{
				"player": c.SpriteSheet{
					Name:    "player",
					Texture: c.Texture{}, // 実際のテクスチャはテスト環境では空
					Sprites: []c.Sprite{
						{X: 0, Y: 0, Width: 32, Height: 32},
						{X: 32, Y: 0, Width: 32, Height: 32},
					},
				},
				"enemy": c.SpriteSheet{
					Name:    "enemy",
					Texture: c.Texture{},
					Sprites: []c.Sprite{
						{X: 0, Y: 0, Width: 16, Height: 16},
						{X: 16, Y: 0, Width: 16, Height: 16},
					},
				},
			},
		}

		assert.Len(t, metadata.SpriteSheets, 2, "スプライトシートの数が正しくない")

		// playerスプライトシートの検証
		playerSheet, exists := metadata.SpriteSheets["player"]
		assert.True(t, exists, "playerスプライトシートが存在しない")
		assert.Equal(t, "player", playerSheet.Name, "playerの名前が正しくない")
		assert.Len(t, playerSheet.Sprites, 2, "playerのスプライト数が正しくない")
		assert.Equal(t, 32, playerSheet.Sprites[0].Width, "playerのスプライト幅が正しくない")
		assert.Equal(t, 32, playerSheet.Sprites[0].Height, "playerのスプライト高さが正しくない")

		// enemyスプライトシートの検証
		enemySheet, exists := metadata.SpriteSheets["enemy"]
		assert.True(t, exists, "enemyスプライトシートが存在しない")
		assert.Equal(t, "enemy", enemySheet.Name, "enemyの名前が正しくない")
		assert.Len(t, enemySheet.Sprites, 2, "enemyのスプライト数が正しくない")
		assert.Equal(t, 16, enemySheet.Sprites[0].Width, "enemyのスプライト幅が正しくない")
		assert.Equal(t, 16, enemySheet.Sprites[0].Height, "enemyのスプライト高さが正しくない")
	})

	t.Run("sprite sheet name assignment", func(t *testing.T) {
		// LoadSpriteSheets内でName フィールドを設定する処理のテスト
		spriteSheets := map[string]c.SpriteSheet{
			"test1": c.SpriteSheet{
				Texture: c.Texture{},
				Sprites: []c.Sprite{},
			},
			"test2": c.SpriteSheet{
				Texture: c.Texture{},
				Sprites: []c.Sprite{},
			},
		}

		// LoadSpriteSheetsの処理の一部を模倣
		for k, v := range spriteSheets {
			v.Name = k
			spriteSheets[k] = v
		}

		// 名前が正しく設定されていることを確認
		assert.Equal(t, "test1", spriteSheets["test1"].Name, "test1のNameが正しく設定されていない")
		assert.Equal(t, "test2", spriteSheets["test2"].Name, "test2のNameが正しく設定されていない")
	})

	t.Run("empty sprite sheet metadata", func(t *testing.T) {
		metadata := spriteSheetMetadata{
			SpriteSheets: map[string]c.SpriteSheet{},
		}

		assert.Empty(t, metadata.SpriteSheets, "空のスプライトシートマップが空でない")
	})
}

// LoadFonts と LoadSpriteSheets の実際のテストは、
// テスト用のTOMLファイルを準備した統合テストとして実装する必要があります。
// ここでは、TOMLデコード後の処理ロジックのテストに焦点を当てています。
