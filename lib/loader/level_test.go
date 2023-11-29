package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTileContains(t *testing.T) {
	ptile := TilePlayer
	wtile := TileWall
	etile := TileEmpty

	assert.False(t, ptile.Contains(wtile))
	assert.True(t, ptile.Contains(ptile))
	assert.True(t, ptile.Contains(etile))

	assert.False(t, wtile.Contains(ptile))
	assert.True(t, wtile.Contains(wtile))
	assert.True(t, wtile.Contains(etile))

	assert.False(t, etile.Contains(wtile))
	assert.False(t, etile.Contains(ptile))
	assert.True(t, etile.Contains(etile))
}

func TestTileContainsAny(t *testing.T) {
	ptile := TilePlayer
	wtile := TileWall
	etile := TileEmpty

	assert.False(t, ptile.ContainsAny(wtile))
	assert.True(t, ptile.ContainsAny(ptile))
	assert.False(t, ptile.ContainsAny(etile))

	assert.False(t, wtile.ContainsAny(ptile))
	assert.True(t, wtile.ContainsAny(wtile))
	assert.False(t, wtile.ContainsAny(etile))

	assert.False(t, etile.ContainsAny(wtile))
	assert.False(t, etile.ContainsAny(ptile))
	assert.False(t, etile.ContainsAny(etile))
}

// よくわからない...
func TestTileSet(t *testing.T) {
	ptile := TilePlayer
	wtile := TileWall
	etile := TileEmpty

	ptile.Set(etile)
	assert.Equal(t, TilePlayer, ptile)
	etile.Set(ptile)
	assert.Equal(t, TilePlayer, ptile)

	wtile.Set(ptile)
	assert.Equal(t, TilePlayer, ptile)
	ptile.Set(wtile)
	assert.Equal(t, Tile(3), ptile) // 定義から飛び出す
}

// よくわからない...
func TestTileRemove(t *testing.T) {
	ptile := TilePlayer
	etile := TileEmpty

	ptile.Remove(etile)
	assert.Equal(t, TilePlayer, ptile)
	etile.Remove(ptile)
	assert.Equal(t, TilePlayer, ptile)
}
