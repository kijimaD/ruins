package game

import (
	"fmt"
	"runtime"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	es "github.com/kijimaD/ruins/lib/engine/states"
	ew "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/resources"
	gr "github.com/kijimaD/ruins/lib/resources"
)

// mainGameはebiten.Game interfaceを満たす
type MainGame struct {
	World        ew.World
	StateMachine es.StateMachine
}

func (game *MainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gr.UpdateGameLayout(game.World)
}

func (game *MainGame) Update() error {
	game.StateMachine.Update(game.World)

	return nil
}

func (game *MainGame) Draw(screen *ebiten.Image) {
	game.StateMachine.Draw(game.World, screen)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	msg := fmt.Sprintf(`FPS: %f
Alloc: %.2fMB
TotalAlloc: %.2fMB
Mallocs: %.2fMB
Frees: %.2fMB
`,

		ebiten.ActualFPS(),
		float64(mem.Alloc/1024/1024),
		float64(mem.TotalAlloc/1024/1024), // 起動後から割り当てられたヒープオブジェクトの数。freeされてもリセットされない
		float64(mem.Mallocs/1024/1024),    // 割り当てられているヒープオブジェクトの数。freeされたら減る
		float64(mem.Frees/1024/1024),      // 解放されたヒープオブジェクトの数
	)
	ebitenutil.DebugPrint(screen, msg)
}

func InitWorld(minGameWidth int, minGameHeight int) ew.World {
	world := ew.InitWorld(&gc.Components{})

	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: minGameWidth, Height: minGameHeight}

	// Load controls
	axes := []string{}
	actions := []string{
		gr.MoveUpAction, gr.MoveDownAction, gr.MoveLeftAction, gr.MoveRightAction,
	}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("metadata/spritesheets/spritesheets.toml")

	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts := loader.LoadFonts("metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts
	world.Resources.DefaultFaces = &map[string]font.Face{
		"kappa": truetype.NewFace((*world.Resources.Fonts)["kappa"].Font, &truetype.Options{
			Size: 24,
			DPI:  72,
		}),
	}

	// load raws
	rw := raw.LoadFromFile("metadata/entities/raw/raw.toml")
	world.Resources.RawMaster = rw

	world.Resources.Game = &resources.Game{}

	return world
}
