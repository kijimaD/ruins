package resources

import (
	"log"

	"github.com/golang/freetype/truetype"

	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/engine/utils"
)

// Font structure
type Font struct {
	Font *truetype.Font
}

// UnmarshalTOML fills structure fields from TOML data
func (f *Font) UnmarshalTOML(i interface{}) error {
	data, err := assets.FS.ReadFile(i.(map[string]interface{})["font"].(string))
	if err != nil {
		log.Fatal(err)
	}
	f.Font = utils.Try(truetype.Parse(data))
	return nil
}
