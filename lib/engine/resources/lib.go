package resources

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kijimaD/ruins/lib/engine/components"
)

// Resources contains references to data not related to any entity
type Resources struct {
	ScreenDimensions *ScreenDimensions
	Controls         *Controls
	InputHandler     *InputHandler
	SpriteSheets     *map[string]components.SpriteSheet
	Fonts            *map[string]Font
	DefaultFaces     *map[string]text.Face
	AudioContext     *audio.Context
	AudioPlayers     *map[string]*audio.Player
	Prefabs          interface{}
	Game             interface{}
	RawMaster        interface{}
}

// InitResources initializes resources
func InitResources() *Resources {
	return &Resources{Controls: &Controls{}, InputHandler: &InputHandler{}}
}
