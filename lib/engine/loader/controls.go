package loader

import (
	"log"

	"github.com/kijimaD/ruins/assets"
	"github.com/kijimaD/ruins/lib/engine/resources"
	"github.com/kijimaD/ruins/lib/engine/utils"

	"github.com/BurntSushi/toml"
)

type controlsConfig struct {
	Controls resources.Controls `toml:"controls"`
}

// LoadControls loads controls from a TOML file
func LoadControls(controlsConfigPath string, axes []string, actions []string) (resources.Controls, resources.InputHandler) {
	var controlsConfig controlsConfig
	bs, err := assets.FS.ReadFile(controlsConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	utils.Try(toml.Decode(string(bs), &controlsConfig))

	var inputHandler resources.InputHandler
	inputHandler.Axes = make(map[string]float64)
	inputHandler.Actions = make(map[string]bool)

	// Check axes
	for _, axis := range axes {
		if _, ok := controlsConfig.Controls.Axes[axis]; !ok {
			utils.LogFatalf("unable to find controls for axis '%s'", axis)
		}
		inputHandler.Axes[axis] = 0
	}

	// Check actions
	for _, action := range actions {
		if _, ok := controlsConfig.Controls.Actions[action]; !ok {
			utils.LogFatalf("unable to find controls for action '%s'", action)
		}
		inputHandler.Actions[action] = false
	}

	return controlsConfig.Controls, inputHandler
}
