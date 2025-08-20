package cmd

import (
	"fmt"

	gs "github.com/kijimaD/ruins/lib/states"
	"github.com/kijimaD/ruins/lib/vrt"
	"github.com/urfave/cli/v2"
)

// CmdScreenshot はスクリーンショットを撮影するコマンド
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
	case gs.BattleState{}.String():
		vrt.RunTestGame(&gs.BattleState{}, mode)
	case gs.CraftMenuState{}.String():
		vrt.RunTestGame(&gs.CraftMenuState{}, mode)
	case gs.DebugMenuState{}.String():
		vrt.RunTestGame(&gs.DebugMenuState{}, mode)
	case gs.DungeonMenuState{}.String():
		vrt.RunTestGame(&gs.DungeonMenuState{}, mode)
	case gs.DungeonSelectState{}.String():
		vrt.RunTestGame(&gs.DungeonSelectState{}, mode)
	case gs.EquipMenuState{}.String():
		vrt.RunTestGame(&gs.EquipMenuState{}, mode)
	case gs.ExecState{}.String():
		vrt.RunTestGame(&gs.ExecState{}, mode)
	case gs.HomeMenuState{}.String():
		vrt.RunTestGame(&gs.HomeMenuState{}, mode)
	case gs.IntroState{}.String():
		vrt.RunTestGame(&gs.IntroState{}, mode)
	case gs.InventoryMenuState{}.String():
		vrt.RunTestGame(&gs.InventoryMenuState{}, mode)
	case gs.LoadMenuState{}.String():
		vrt.RunTestGame(&gs.LoadMenuState{}, mode)
	case gs.MainMenuState{}.String():
		vrt.RunTestGame(&gs.MainMenuState{}, mode)
	case gs.MessageState{}.String():
		vrt.RunTestGame(&gs.MessageState{}, mode)
	case gs.SaveMenuState{}.String():
		vrt.RunTestGame(&gs.SaveMenuState{}, mode)
	case gs.PartySetupState{}.String():
		vrt.RunTestGame(&gs.PartySetupState{}, mode)
	case gs.DungeonState{}.String():
		vrt.RunTestGame(&gs.DungeonState{}, mode)
	case gs.GameOverState{}.String():
		vrt.RunTestGame(&gs.GameOverState{}, mode)
	default:
		return fmt.Errorf("スクリーンショット実行時に対応してないステートが指定された: %s", mode)
	}

	return nil
}
