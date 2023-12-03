package loader

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/kijimaD/sokotwo/assets"
	gc "github.com/kijimaD/sokotwo/lib/components"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/math"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	"github.com/kijimaD/sokotwo/lib/utils/vutil"
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
	// 壁より外側の埋め合わせる部分
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
	TilePlayer Tile = 1 << iota
	TileWall
	TileWarpNext
	TileWarpEscape
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

// 設定ファイルをタイルとして読み出す。コンポーネント生成はしない
// あとで階を生成するときに、タイルを元にコンポーネントを生成する
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
	warpNextCount := 0
	warpEscapeCount := 0
	for _, line := range lines {
		gridWidth = math.Max(gridWidth, len(line))
		playerCount += strings.Count(line, string(charPlayer))
		warpNextCount += strings.Count(line, string(charWarpNext))
		warpEscapeCount += strings.Count(line, string(charWarpEscape))
	}

	if gridWidth > MaxGridSize || gridHeight > MaxGridSize {
		return nil, fmt.Errorf("level size must be less than %dx%d", MaxGridSize, MaxGridSize)
	}
	if playerCount != 1 {
		return nil, fmt.Errorf("invalid level: level must have one player")
	}
	if warpNextCount != 1 {
		return nil, fmt.Errorf("invalid level: level must have one next warp hole")
	}
	if warpEscapeCount != 1 {
		return nil, fmt.Errorf("invalid level: level must have one escape warp hole")
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

// フロアに外壁を置く
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

// 階層データからエンティティ(コンポーネント群)を生成する
// selectLevel: 選択したフロアデータのインデックス
// levelNum: 今いる階数
func LoadLevel(packageData PackageData, selectLevel, levelNum, layoutWidth, layoutHeight int, gameSpriteSheet *ec.SpriteSheet) (vutil.Vec2d[Tile], loader.EntityComponentList, error) {
	componentList := loader.EntityComponentList{}

	grid := packageData.Levels[selectLevel]
	gridWidth := grid.NCols
	gridHeight := grid.NRows

	horizontalPadding := layoutWidth - gridWidth
	horizontalPaddingBefore := horizontalPadding / 2
	horizontalPaddingAfter := horizontalPadding - horizontalPaddingBefore

	verticalPadding := layoutHeight - gridHeight
	verticalPaddingBefore := verticalPadding / 2
	verticalPaddingAfter := verticalPadding - verticalPaddingBefore

	tiles := make([]Tile, 0, gridWidth*gridHeight)

	for iLine := 0; iLine < verticalPaddingBefore; iLine++ {
		for iCol := 0; iCol < layoutWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	for iGridLine := 0; iGridLine < gridHeight; iGridLine++ {
		iLine := iGridLine + verticalPaddingBefore

		for iCol := 0; iCol < horizontalPaddingBefore; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}

		for iGridCol := 0; iGridCol < gridWidth; iGridCol++ {
			char := *grid.Get(iGridLine, iGridCol)
			iCol := iGridCol + horizontalPaddingBefore

			switch char {
			case charFloor:
				tiles = append(tiles, TileEmpty)
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charExterior:
				tiles = append(tiles, TileEmpty)
				createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charWall:
				tiles = append(tiles, TileWall)
				createWallEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charPlayer:
				tiles = append(tiles, TilePlayer)
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createPlayerEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charWarpNext:
				tiles = append(tiles, TileWarpNext)
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)
				createWarpNextEntity(&componentList, gameSpriteSheet, iLine, iCol)
			case charWarpEscape:
				tiles = append(tiles, TileWarpEscape)
				createFloorEntity(&componentList, gameSpriteSheet, iLine, iCol)

				const EscapeFloorCycle = 5 // 5階ごとに脱出フロア
				if levelNum%EscapeFloorCycle == 0 {
					createWarpEscapeEntity(&componentList, gameSpriteSheet, iLine, iCol)
				}
			default:
				return vutil.Vec2d[Tile]{}, loader.EntityComponentList{}, fmt.Errorf("invalid level: invalid char '%c'", char)
			}
		}

		for iCol := layoutWidth - horizontalPaddingAfter; iCol < layoutWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	for iLine := layoutHeight - verticalPaddingAfter; iLine < layoutHeight; iLine++ {
		for iCol := 0; iCol < layoutWidth; iCol++ {
			createExteriorEntity(&componentList, gameSpriteSheet, iLine, iCol)
		}
	}

	gameGrid := utils.Try(vutil.NewVec2d(gridHeight, gridWidth, tiles))
	return gameGrid, componentList, nil
}

func createFloorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: floorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createExteriorEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: exteriorSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createWallEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: wallSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Wall:        &gc.Wall{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createPlayerEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: playerSpriteNumber},
		Transform:    &ec.Transform{Depth: 1},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Player:      &gc.Player{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createWarpNextEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: warpNextSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Warp:        &gc.Warp{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}

func createWarpEscapeEntity(componentList *loader.EntityComponentList, gameSpriteSheet *ec.SpriteSheet, line, col int) {
	componentList.Engine = append(componentList.Engine, loader.EngineComponentList{
		SpriteRender: &ec.SpriteRender{SpriteSheet: gameSpriteSheet, SpriteNumber: warpEscapeSpriteNumber},
		Transform:    &ec.Transform{},
	})
	componentList.Game = append(componentList.Game, gameComponentList{
		Warp:        &gc.Warp{},
		GridElement: &gc.GridElement{Line: line, Col: col},
	})
}
