package loader

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/kijimaD/sokotwo/assets"
	"github.com/kijimaD/sokotwo/lib/engine/math"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	"github.com/kijimaD/sokotwo/lib/utils/vutil"
)

// 最大のグリッドサイズ
const MaxGridSize = 50

const (
	charFloor  = ' '
	charWall   = '#'
	charPlayer = '@'
	// 壁より外側の部分
	charExterior = '_'
)

var regexpValidChars = regexp.MustCompile(`^[ #@+]+$`)

// 1つのパッケージは複数の階層を持つ
type PackageData struct {
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

func LoadPackage(packageName string) (packageData PackageData, packageErr error) {
	packageData.Name = packageName

	// Load file
	file := utils.Try(assets.FS.ReadFile(fmt.Sprintf("levels/%s.xsb", packageName)))
	lines := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	lines = append(lines, "")

	// Split levels
	levels := [][]string{}
	currentLevel := []string{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 && len(currentLevel) > 0 {
			levels = append(levels, currentLevel)
			currentLevel = []string{}
		} else if regexpValidChars.MatchString(line) {
			currentLevel = append(currentLevel, line)
		}
	}

	for iLevel, level := range levels {
		if grid, err := normalizeLevel(level); err == nil {
			gridWidth := len(grid[0])
			gridHeight := len(grid)

			blocks := make([]byte, 0, gridWidth*gridHeight)
			for _, line := range grid {
				blocks = append(blocks, line...)
			}
			packageData.Levels = append(packageData.Levels, utils.Try(vutil.NewVec2d(gridHeight, gridWidth, blocks)))
		} else {
			packageErr = fmt.Errorf("error when loading level %d: %s", iLevel+1, err.Error())
			break
		}
	}

	if len(packageData.Levels) == 0 {
		if packageErr != nil {
			log.Println(packageErr)
		}
		log.Fatal("invalid package: no valid levels in package")
	}
	return
}

// レベルの正規化。テキストで作成されたフロアのデータを変換する。論理的に正しいかのチェックもする
func normalizeLevel(lines []string) ([][]byte, error) {
	gridWidth := 0
	gridHeight := len(lines)
	playerCount := 0
	for _, line := range lines {
		gridWidth = math.Max(gridWidth, len(line))
		playerCount += strings.Count(line, string(charPlayer))
	}

	if gridWidth > MaxGridSize || gridHeight > MaxGridSize {
		return nil, fmt.Errorf("level size must be less than %dx%d", MaxGridSize, MaxGridSize)
	}
	if playerCount != 1 {
		return nil, fmt.Errorf("invalid level: level must have one player")
	}

	grid := make([][]byte, len(lines))

	for iLine := range lines {
		chars := []byte(lines[iLine])

		deltaLen := gridWidth - len(chars)
		for iSlice := 0; iSlice < deltaLen; iSlice++ {
			chars = append(chars, charFloor)
		}

		grid[iLine] = chars
	}

	// 横
	for iLine := 0; iLine < gridHeight; iLine++ {
		fillExterior(grid, iLine, 0, gridWidth, gridHeight)
		fillExterior(grid, iLine, gridWidth-1, gridWidth, gridHeight)
	}

	// 縦
	for iCol := 0; iCol < gridWidth; iCol++ {
		fillExterior(grid, 0, iCol, gridWidth, gridHeight)
		fillExterior(grid, gridHeight-1, iCol, gridWidth, gridHeight)
	}

	return grid, nil
}

// フロアに物を置く
func fillExterior(grid [][]byte, line, col, gridWidth, gridHeight int) {
	if grid[line][col] != charFloor {
		return
	}

	fillQueue := &[]struct{ line, col int }{{line, col}}

	for len(*fillQueue) > 0 {
		elem := (*fillQueue)[0]
		*fillQueue = (*fillQueue)[1:]

		colLeft := elem.col
		for colLeft > 0 && grid[elem.line][colLeft-1] == charFloor {
			colLeft--
		}

		colRight := elem.col
		for colRight < gridWidth-1 && grid[elem.line][colRight+1] == charFloor {
			colRight++
		}

		for iCol := colLeft; iCol <= colRight; iCol++ {
			grid[elem.line][iCol] = charExterior

			if elem.line > 0 && grid[elem.line-1][iCol] == charFloor {
				*fillQueue = append(*fillQueue, struct{ line, col int }{elem.line - 1, iCol})
			}

			if elem.line < gridHeight-1 && grid[elem.line+1][iCol] == charFloor {
				*fillQueue = append(*fillQueue, struct{ line, col int }{elem.line + 1, iCol})
			}
		}
	}
}
