package loader

import (
	c "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/utils"

	"github.com/BurntSushi/toml"
)

type spriteSheetMetadata struct {
	SpriteSheets map[string]c.SpriteSheet `toml:"sprite_sheet"`
}

// LoadSpriteSheets loads sprite sheets from a TOML file
func LoadSpriteSheets(spriteSheetMetadataPath string) map[string]c.SpriteSheet {
	var spriteSheetMetadata spriteSheetMetadata
	utils.Try(toml.DecodeFile(spriteSheetMetadataPath, &spriteSheetMetadata))
	return spriteSheetMetadata.SpriteSheets
}
