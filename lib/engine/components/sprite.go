package components

import (
	"bytes"
	"log"

	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/engine/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const DungeonTileSize = 32

// スプライトは1つ1つの意味をなす画像の位置を示す情報
// 1ファイルに対して複数のスプライトが定義されている
type Sprite struct {
	// Horizontal position of the sprite in the sprite sheet
	X int
	// Vertical position of the sprite in the sprite sheet
	Y int
	// Width of the sprite
	Width int
	// Height of the sprite
	Height int
}

// 複数のスプライトが格納された画像ファイル
type Texture struct {
	// Texture image
	Image *ebiten.Image
}

// UnmarshalText fills structure fields from text data
func (t *Texture) UnmarshalText(text []byte) error {
	bs, err := assets.FS.ReadFile(string(text))
	if err != nil {
		log.Fatal(err)
	}
	textureImage, _ := utils.Try2(ebitenutil.NewImageFromReader(bytes.NewReader(bs)))
	t.Image = textureImage
	return nil
}

// 画像ファイルであるテクスチャと、その位置ごとの解釈であるスプライトのマッピング
type SpriteSheet struct {
	// Texture image
	Texture Texture `toml:"texture_image"`
	// List of sprites
	Sprites []Sprite
}

// SpriteRender component
type SpriteRender struct {
	// Reference sprite sheet
	SpriteSheet *SpriteSheet
	// Index of the sprite on the sprite sheet
	SpriteNumber int
	// Draw options
	Options ebiten.DrawImageOptions
}
