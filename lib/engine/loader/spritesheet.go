package loader

import (
	"log"

	"github.com/kijimaD/ruins/assets"
	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/utils"

	"github.com/BurntSushi/toml"
)

type spriteSheetMetadata struct {
	SpriteSheets map[string]c.SpriteSheet `toml:"sprite_sheet"`
}

// LoadSpriteSheets loads sprite sheets from a TOML file
func LoadSpriteSheets(spriteSheetMetadataPath string) map[string]c.SpriteSheet {
	var spriteSheetMetadata spriteSheetMetadata
	bs, err := assets.FS.ReadFile(spriteSheetMetadataPath)
	if err != nil {
		log.Fatal(err)
	}
	utils.Try(toml.Decode(string(bs), &spriteSheetMetadata))
	for k, v := range spriteSheetMetadata.SpriteSheets {
		v.Name = k
		spriteSheetMetadata.SpriteSheets[k] = v
	}
	return spriteSheetMetadata.SpriteSheets
}
