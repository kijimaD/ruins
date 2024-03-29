package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/loader"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	es "github.com/kijimaD/ruins/lib/engine/states"
	ew "github.com/kijimaD/ruins/lib/engine/world"
	gloader "github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/raw"
	gr "github.com/kijimaD/ruins/lib/resources"
	gs "github.com/kijimaD/ruins/lib/states"
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
	// stateによっては背景に隠れて見えない
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.CurrentFPS()))
	game.stateMachine.Draw(game.world, screen)
}

func main() {
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

	// load prefabs
	world.Resources.Prefabs = &gr.Prefabs{
		Menu: gr.MenuPrefabs{
			MainMenu:      gloader.PreloadEntities("metadata/entities/ui/main_menu.toml", world),
			HomeMenu:      gloader.PreloadEntities("metadata/entities/ui/home_menu.toml", world),
			DungeonSelect: gloader.PreloadEntities("metadata/entities/ui/dungeon_select.toml", world),
			FieldMenu:     gloader.PreloadEntities("metadata/entities/ui/field_menu.toml", world),
			DebugMenu:     gloader.PreloadEntities("metadata/entities/ui/debug_menu.toml", world),
			InventoryMenu: gloader.PreloadEntities("metadata/entities/ui/inventory_menu.toml", world),
			CraftMenu:     gloader.PreloadEntities("metadata/entities/ui/craft_menu.toml", world),
			CampMenu:      gloader.PreloadEntities("metadata/entities/ui/camp_menu.toml", world),
			EquipMenu:     gloader.PreloadEntities("metadata/entities/ui/equip_menu.toml", world),
		},
		Intro: gloader.PreloadEntities("metadata/entities/ui/intro.toml", world),
		Field: gr.FieldPrefabs{
			LevelInfo:   gloader.PreloadEntities("metadata/entities/ui/level.toml", world),
			PackageInfo: gloader.PreloadEntities("metadata/entities/ui/package.toml", world),
		},
	}

	// load raws
	rw := raw.LoadFromFile("metadata/entities/raw/raw.toml")
	world.Resources.RawMaster = rw

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(minGameWidth, minGameHeight)
	ebiten.SetWindowTitle("ruins")

	ebiten.RunGame(&mainGame{
		world:        world,
		stateMachine: es.Init(&gs.MainMenuState{}, world),
	})
}
