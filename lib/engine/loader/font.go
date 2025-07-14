package loader

import (
	"log"

	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/engine/resources"
	"github.com/kijimaD/ruins/lib/engine/utils"

	"github.com/BurntSushi/toml"
)

type fontMetadata struct {
	// resources.Fontは UnmarshalTOML を実装しており、そのロジックでフォントに変換して構造体に入れられる
	Fonts map[string]resources.Font `toml:"font"`
}

// LoadFonts はフォントの設定ymlを読み込む
func LoadFonts(fontPath string) map[string]resources.Font {
	var fontMetadata fontMetadata
	bs, err := assets.FS.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	utils.Try(toml.Decode(string(bs), &fontMetadata))
	return fontMetadata.Fonts
}
