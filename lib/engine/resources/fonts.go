package resources

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type fonts struct {
	face         text.Face
	titleFace    text.Face
	bigTitleFace text.Face
	toolTipFace  text.Face
}

func loadFonts(tfs *text.GoTextFaceSource) *fonts {
	fontFace := loadFont(tfs, 20)
	titleFontFace := loadFont(tfs, 24)
	bigTitleFontFace := loadFont(tfs, 28)
	toolTipFace := loadFont(tfs, 15)

	return &fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
	}
}

func loadFont(tfs *text.GoTextFaceSource, size float64) text.Face {
	return &text.GoTextFace{
		Source: tfs,
		Size:   size,
	}
}
