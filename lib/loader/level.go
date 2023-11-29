package loader

import (
	"github.com/kijimaD/sokotwo/lib/utils/vutil"
)

const (
	charFloor  = ' '
	charWall   = '#'
	charPlayer = '@'
)

// 1つのダンジョンは複数の階層を持つ
type DungeonData struct {
	Name   string
	Levels []vutil.Vec2d[byte]
}

// フィールドのタイル
type Tile uint8

const (
	TilePlayer Tile = 1 << iota
	TileWall
	TileEmpty Tile = 0
)

// レシーバのゲームタイルが、引数のタイルを含んでいるかチェックする
// 同じならばTrue、引数のタイルが空白ならばTrue
func (t *Tile) Contains(other Tile) bool {
	return (*t & other) == other
}

func (t *Tile) ContainsAny(other Tile) bool {
	return (*t & other) != 0
}

// タイルをセットする。TileEmptyは上書きされる
func (t *Tile) Set(other Tile) {
	*t |= other
}

// タイルを削除
func (t *Tile) Remove(other Tile) {
	*t &= 0xFF ^ other
}
