package world

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestInitWorld(t *testing.T) {
	t.Run("InitWorldが動作する", func(t *testing.T) {
		gameComponents := &gc.Components{}

		world := InitWorld(gameComponents)

		assert.NotNil(t, world.Manager)
		assert.NotNil(t, world.Components)
		assert.NotNil(t, world.Resources)
		assert.NotNil(t, world.Components)
	})
}
