package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
	gr "github.com/kijimaD/sokotwo/lib/resources"
	gs "github.com/kijimaD/sokotwo/lib/states"
	ec "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/loader"
	er "github.com/x-hgg-x/goecsengine/resources"
	es "github.com/x-hgg-x/goecsengine/states"
	ew "github.com/x-hgg-x/goecsengine/world"
)

const (
	minGameWidth  = 960
	minGameHeight = 720
)

// mainGameはebiten.Game interfaceを満たす
type mainGame struct {
	world        ew.World
	stateMachine es.StateMachine
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	var gridLayout *gr.GridLayout

	if game.world.Resources.Game != nil {
		gridLayout = &game.world.Resources.Game.(*gr.Game).GridLayout
	}

	return gr.UpdateGameLayout(game.world, gridLayout)
}

func (game *mainGame) Update() error {
	game.stateMachine.Update(game.world)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.stateMachine.Draw(game.world, screen)
}

func main() {
	world := ew.InitWorld(&ec.Components{})

	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: minGameWidth, Height: minGameHeight}

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("assets/metadata/spritesheets/spritesheets.toml")
	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts := loader.LoadFonts("assets/metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts

	// load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			MainMenu: gloader.PreloadEntities("assets/metadata/entities/ui/main_menu.toml", world),
		},
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(minGameWidth, minGameHeight)
	ebiten.SetWindowTitle("sokotwo")

	ebiten.RunGame(&mainGame{
		world:        world,
		stateMachine: es.Init(&gs.MainMenuState{}, world),
	})
}
