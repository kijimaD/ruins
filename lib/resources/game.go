package resources

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kijimaD/sokotwo/lib/engine/math"
	"github.com/kijimaD/sokotwo/lib/engine/utils"
	"github.com/kijimaD/sokotwo/lib/utils/vutil"

	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
)

type StateEvent string

const (
	StateEventNone           = StateEvent("NONE")
	StateEventMainTransition = StateEvent("MAIN_TRANSITION")
)

const (
	offsetX       = 0
	offsetY       = 80
	gridBlockSize = 32
	minGridWidth  = 30
	minGridHeight = 20
)

// Tileはsystemなどでも使う。systemから直接gloaderを扱わせたくないので、ここでエクスポートする
const (
	TilePlayer     = gloader.TilePlayer
	TileWall       = gloader.TileWall
	TileWarpNext   = gloader.TileWarpNext
	TileWarpEscape = gloader.TileWarpEscape
	TileEmpty      = gloader.TileEmpty
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
	StateEvent StateEvent
	Package    PackageData
	Level      Level
	GridLayout GridLayout
}

// levelNum: 今いる階数
func InitLevel(world w.World, levelNum int) {
	gameResources := world.Resources.Game.(*Game)

	// Load ui entities
	prefabs := world.Resources.Prefabs.(*Prefabs)
	loader.AddEntities(world, prefabs.Field.PackageInfo)
	levelInfoEntity := loader.AddEntities(world, prefabs.Field.LevelInfo)[0]

	rand.Seed(time.Now().UnixNano())
	randLevelNum := rand.Intn(len(gameResources.Package.Levels))

	level := gameResources.Package.Levels[randLevelNum]
	gridLayout := &gameResources.GridLayout
	gridLayout.Width = math.Max(minGridWidth, level.NCols)
	gridLayout.Height = math.Max(minGridHeight, level.NRows)

	UpdateGameLayout(world, gridLayout)

	gameSpriteSheet := (*world.Resources.SpriteSheets)["game"]
	grid, levelComponentList := utils.Try2(gloader.LoadLevel(gameResources.Package, randLevelNum, levelNum, gridLayout.Width, gridLayout.Height, &gameSpriteSheet))
	loader.AddEntities(world, levelComponentList)
	gameResources.Level = Level{CurrentNum: levelNum, Grid: grid}

	// Set level info text
	world.Components.Engine.Text.Get(levelInfoEntity).(*ec.Text).Text = fmt.Sprintf("%dF", levelNum)
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
