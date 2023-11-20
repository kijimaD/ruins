package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	gr "github.com/kijimaD/sokotwo/lib/resources"
	gs "github.com/kijimaD/sokotwo/lib/states"
	ec "github.com/x-hgg-x/goecsengine/components"
	es "github.com/x-hgg-x/goecsengine/states"
	ew "github.com/x-hgg-x/goecsengine/world"
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
	ebiten.RunGame(&mainGame{
		world:        world,
		stateMachine: es.Init(&gs.MainMenuState{}, world),
	})
}
