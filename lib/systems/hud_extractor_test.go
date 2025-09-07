package systems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTileColorInfo はTileColorInfoの型エイリアスをテスト
func TestTileColorInfo(t *testing.T) {
	t.Parallel()
	colorInfo := TileColorInfo{
		R: 255,
		G: 128,
		B: 64,
		A: 200,
	}

	// hud.TileColorInfoと同じ構造であることを確認
	var hudColorInfo = colorInfo

	assert.Equal(t, uint8(255), hudColorInfo.R)
	assert.Equal(t, uint8(128), hudColorInfo.G)
	assert.Equal(t, uint8(64), hudColorInfo.B)
	assert.Equal(t, uint8(200), hudColorInfo.A)
}
