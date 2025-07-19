package cmd

import (
	"log"
	"net/http"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/pkg/profile"
	"github.com/urfave/cli/v2"

	_ "net/http/pprof" // pprofのHTTPエンドポイントを登録するためのインポート

	es "github.com/kijimaD/ruins/lib/engine/states"
	gs "github.com/kijimaD/ruins/lib/states"
)

// CmdPlay はゲームをプレイするコマンド
var CmdPlay = &cli.Command{
	Name:        "play",
	Usage:       "play",
	Description: "play game",
	Action:      runPlay,
	Flags:       []cli.Flag{},
}

func runPlay(_ *cli.Context) error {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(consts.MinGameWidth, consts.MinGameHeight)
	ebiten.SetWindowTitle("ruins")

	// プロファイラ。WASMは除外する
	if runtime.GOOS != "js" {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
		go func() {
			log.Fatal(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	world, err := game.InitWorld(consts.MinGameWidth, consts.MinGameHeight)
	if err != nil {
		return err
	}

	err = ebiten.RunGame(&game.MainGame{
		World:        world,
		StateMachine: es.Init(&gs.MainMenuState{}, world),
	})
	if err != nil {
		return err
	}

	return nil
}
