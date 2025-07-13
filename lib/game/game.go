package game

import (
	"fmt"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	es "github.com/kijimaD/ruins/lib/engine/states"
	ew "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	gr "github.com/kijimaD/ruins/lib/resources"
)

// mainGameはebiten.Game interfaceを満たす
type MainGame struct {
	World        ew.World
	StateMachine es.StateMachine
}

// interface methodのため、シグネチャは変更できない
func (game *MainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	x, y := gr.UpdateGameLayout(game.World)

	return int(x), int(y)
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

	// ResourceManagerを使用してリソースを読み込む
	resourceManager := gr.NewDefaultResourceManager()

	// Load controls
	axes := []string{}
	actions := []string{
		gr.MoveUpAction, gr.MoveDownAction, gr.MoveLeftAction, gr.MoveRightAction,
	}
	controls, inputHandler, err := resourceManager.LoadControls(axes, actions)
	if err != nil {
		log.Fatal(err)
	}
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets, err := resourceManager.LoadSpriteSheets()
	if err != nil {
		log.Fatal(err)
	}
	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts, err := resourceManager.LoadFonts()
	if err != nil {
		log.Fatal(err)
	}
	world.Resources.Fonts = &fonts

	defaultFont := (*world.Resources.Fonts)["kappa"]
	world.Resources.DefaultFaces = &map[string]text.Face{
		"kappa": defaultFont.Font,
	}

	// load UI resources
	uir, err := er.NewUIResources(defaultFont.FaceSource)
	if err != nil {
		log.Fatal(err)
	}
	world.Resources.UIResources = uir

	// load raws
	rw, err := resourceManager.LoadRaws()
	if err != nil {
		log.Fatal(err)
	}
	world.Resources.RawMaster = rw

	world.Resources.Game = &resources.Game{}

	return world
}
