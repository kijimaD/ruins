package resources

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/ruins/lib/components"
)

// Resources contains references to data not related to any entity
type Resources struct {
	ScreenDimensions *ScreenDimensions
	SpriteSheets     *map[string]components.SpriteSheet
	Fonts            *map[string]Font
	Faces            *map[string]text.Face
	Dungeon          interface{}
	RawMaster        interface{}
	UIResources      *UIResources
	TurnManager      interface{}
}

// ScreenDimensions contains current screen dimensions
type ScreenDimensions struct {
	Width  int
	Height int
}

// InitResources initializes resources
func InitResources() *Resources {
	return &Resources{}
}
