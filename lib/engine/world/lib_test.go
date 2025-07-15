package world

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestInitGeneric(t *testing.T) {
	t.Run("型安全なInitGenericが動作する", func(t *testing.T) {
		gameComponents := &gc.Components{}

		world, err := InitGeneric(gameComponents)

		assert.NoError(t, err)
		assert.NotNil(t, world.Manager)
		assert.NotNil(t, world.Components)
		assert.NotNil(t, world.Resources)
		assert.NotNil(t, world.Components)

		// 型安全性の確認
		assert.IsType(t, &gc.Components{}, world.Components.Game)
	})

	t.Run("型安全性が保たれている", func(t *testing.T) {
		gameComponents := &gc.Components{}

		world, err := InitGeneric(gameComponents)

		assert.NoError(t, err)
		// 型アサーションが不要で、直接アクセスできる
		assert.NotNil(t, world.Components.Game.Position)
		assert.NotNil(t, world.Components.Game.Velocity)
	})
}
