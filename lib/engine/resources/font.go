package resources

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/kijimaD/ruins/assets"
)

// Font structure
type Font struct {
	Font text.Face
}

// UnmarshalTOML fills structure fields from TOML data
func (f *Font) UnmarshalTOML(i interface{}) error {
	fontFile, err := assets.FS.Open(i.(map[string]interface{})["font"].(string))
	if err != nil {
		return err
	}
	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		return err
	}
	font := &text.GoTextFace{
		Source: s,
		Size:   24,
	}
	f.Font = font

	return nil
}
