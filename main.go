package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/sokotwo/lib/components"
	ec "github.com/kijimaD/sokotwo/lib/engine/components"
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	er "github.com/kijimaD/sokotwo/lib/engine/resources"
	es "github.com/kijimaD/sokotwo/lib/engine/states"
	ew "github.com/kijimaD/sokotwo/lib/engine/world"
	gloader "github.com/kijimaD/sokotwo/lib/loader"
	gr "github.com/kijimaD/sokotwo/lib/resources"
	gs "github.com/kijimaD/sokotwo/lib/states"
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
	world := ew.InitWorld(&gc.Components{})

	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: minGameWidth, Height: minGameHeight}

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("metadata/spritesheets/spritesheets.toml")

	textureImage := ebiten.NewImage(minGameWidth, minGameHeight)
	textureImage.Fill(color.RGBA{B: 255})
	spriteSheets["intro-bg"] = ec.SpriteSheet{Texture: ec.Texture{Image: textureImage}, Sprites: []ec.Sprite{{Width: minGameWidth, Height: minGameHeight}}}
	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts := loader.LoadFonts("metadata/fonts/fonts.toml")
	world.Resources.Fonts = &fonts

	// load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			MainMenu: gloader.PreloadEntities("metadata/entities/ui/main_menu.toml", world),
		},
		Intro: gloader.PreloadEntities("metadata/entities/ui/intro.toml", world),
		Field: gr.FieldPrefabs{
			LevelInfo:   gloader.PreloadEntities("metadata/entities/ui/level.toml", world),
			PackageInfo: gloader.PreloadEntities("metadata/entities/ui/package.toml", world),
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
