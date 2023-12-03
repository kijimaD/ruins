package resources

import (
	"fmt"

	"github.com/kijimaD/sokotwo/lib/engine/math"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	"github.com/kijimaD/sokotwo/lib/utils/vutil"

	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

const (
	TilePlayer   = gloader.TilePlayer
	TileWall     = gloader.TileWall
	TileWarpNext = gloader.TileWarpNext
	TileEmpty    = gloader.TileEmpty
)

type Level struct {
	CurrentNum int
	Grid       vutil.Vec2d[Tile]
	Movements  []MovementType
	Modified   bool
}

// PackageData contains level package data
type PackageData = gloader.PackageData

type Tile = gloader.Tile

// グリッドレイアウト
type GridLayout struct {
	Width  int
	Height int
}

type Game struct {
	Package    PackageData
	Level      Level
	GridLayout GridLayout
}

func InitLevel(world w.World, levelNum int) {
	gameResources := world.Resources.Game.(*Game)

	// Load ui entities
	prefabs := world.Resources.Prefabs.(*Prefabs)
	loader.AddEntities(world, prefabs.Field.PackageInfo)
	levelInfoEntity := loader.AddEntities(world, prefabs.Field.LevelInfo)[0]

	// Load level
	level := gameResources.Package.Levels[levelNum]
	gridLayout := &gameResources.GridLayout
	gridLayout.Width = math.Max(minGridWidth, level.NCols)
	gridLayout.Height = math.Max(minGridHeight, level.NRows)

	UpdateGameLayout(world, gridLayout)

	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]
	grid, levelComponentList := utils.Try2(gloader.LoadLevel(gameResources.Package, levelNum, gridLayout.Width, gridLayout.Height, &gameSpriteSheet))
	loader.AddEntities(world, levelComponentList)
	gameResources.Level = Level{CurrentNum: levelNum, Grid: grid}

	// Set level info text
	world.Components.Engine.Text.Get(levelInfoEntity).(*ec.Text).Text = fmt.Sprintf("B%d", levelNum+1)
}

// UpdateGameLayoutはゲームレイアウトを更新する
func UpdateGameLayout(world w.World, gridLayout *GridLayout) (int, int) {
	gridWidth, gridHeight := minGridWidth, minGridHeight

	if gridLayout != nil {
		gridWidth = gridLayout.Width
		gridHeight = gridLayout.Height
	}

	gameWidth := gridWidth*gridBlockSize + offsetX
	gameHeight := gridHeight*gridBlockSize + offsetY

	fadeOutSprite := &(*world.Resources.SpriteSheets)["intro-bg"].Sprites[0]
	fadeOutSprite.Width = gameWidth
	fadeOutSprite.Height = gameHeight

	world.Resources.ScreenDimensions.Width = gameWidth
	world.Resources.ScreenDimensions.Height = gameHeight

	return gameWidth, gameHeight
}
