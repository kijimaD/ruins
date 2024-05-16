package loader

import (
	"regexp"

	"github.com/kijimaD/ruins/lib/utils/vutil"
)

// 最大のグリッドサイズ
const MaxGridSize = 50

const (
	exteriorSpriteNumber   = 0
	wallSpriteNumber       = 1
	floorSpriteNumber      = 2
	playerSpriteNumber     = 3
	warpNextSpriteNumber   = 4
	warpEscapeSpriteNumber = 5
)

const (
	// フロア
	charFloor = ' '
	// 壁
	charWall = '#'
	// 操作するプレイヤー
	charPlayer = '@'
	// 壁より外側の埋める部分
	charExterior = '_'
	// 次の階層へ
	charWarpNext = 'O'
	// 脱出
	charWarpEscape = 'X'
)

var regexpValidChars = regexp.MustCompile(`^[ #@+_OX]+$`)

// 1つのパッケージは複数の階層を持つ
type PackageData struct {
	Name   string
	Levels []vutil.Vec2d[byte]
}

// フィールドのタイル
type Tile uint8

const (
	TileEmpty Tile = 0
	TileWall  Tile = 1 << iota
)
