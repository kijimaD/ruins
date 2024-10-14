package resources

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	fontFaceRegular = "assets/fonts/NotoSans-Regular.ttf"
	fontFaceBold    = "assets/fonts/NotoSans-Bold.ttf"
)

type fonts struct {
	face         text.Face
	titleFace    text.Face
	bigTitleFace text.Face
	toolTipFace  text.Face
}

func loadFonts(tfs *text.GoTextFaceSource) (*fonts, error) {
	fontFace, err := loadFont(tfs, 20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := loadFont(tfs, 24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := loadFont(tfs, 28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := loadFont(tfs, 15)
	if err != nil {
		return nil, err
	}

	return &fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
	}, nil
}

func loadFont(tfs *text.GoTextFaceSource, size float64) (text.Face, error) {
	return &text.GoTextFace{
		Source: tfs,
		Size:   size,
	}, nil
}
