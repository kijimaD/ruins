package cmd

import (
	"log"
	"net/http"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/pkg/profile"
	"github.com/urfave/cli/v2"

	_ "net/http/pprof"

	es "github.com/kijimaD/ruins/lib/engine/states"
	gs "github.com/kijimaD/ruins/lib/states"
)

var CmdPlay = &cli.Command{
	Name:        "play",
	Usage:       "play",
	Description: "play game",
	Action:      runPlay,
	Flags:       []cli.Flag{},
}

func runPlay(_ *cli.Context) error {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(utils.MinGameWidth, utils.MinGameHeight)
	ebiten.SetWindowTitle("ruins")

	// プロファイラ。WASMは除外する
	if runtime.GOOS != "js" {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
		go func() {
			log.Fatal(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	world := game.InitWorld(utils.MinGameWidth, utils.MinGameHeight)
	ebiten.RunGame(&game.MainGame{
		World:        world,
		StateMachine: es.Init(&gs.MainMenuState{}, world),
	})

	return nil
}
