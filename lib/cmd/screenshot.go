package cmd

import (
	"fmt"

	gs "github.com/kijimaD/ruins/lib/states"
	"github.com/kijimaD/ruins/lib/vrt"
	"github.com/urfave/cli/v2"
)

var CmdScreenshot = &cli.Command{
	Name:        "screenshot",
	Usage:       "screenshot",
	Description: "screenshot game",
	Action:      runScreenshot,
	Flags:       []cli.Flag{},
}

func runScreenshot(ctx *cli.Context) error {
	mode := ctx.Args().Get(0)
	if mode == "" {
		return fmt.Errorf("引数が不足している。ステート名が必要")
	}

	switch mode {
	case "MainMenu":
		vrt.RunTestGame(&gs.MainMenuState{}, mode)
	case "Intro":
		vrt.RunTestGame(&gs.IntroState{}, mode)
	}

	return nil
}
