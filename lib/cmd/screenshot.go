package cmd

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/states"
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
	case gs.CraftMenuState{}.String():
		st := &gs.CraftMenuState{}
		st.SetCategory(states.ItemCategoryTypeItem)
		vrt.RunTestGame(st, mode)
	case gs.DebugMenuState{}.String():
		vrt.RunTestGame(&gs.DebugMenuState{}, mode)
	case gs.DungeonSelectState{}.String():
		vrt.RunTestGame(&gs.DungeonSelectState{}, mode)
	case gs.EquipMenuState{}.String():
		vrt.RunTestGame(&gs.EquipMenuState{}, mode)
	case gs.FieldMenuState{}.String():
		vrt.RunTestGame(&gs.FieldMenuState{}, mode)
	case gs.HomeMenuState{}.String():
		vrt.RunTestGame(&gs.HomeMenuState{}, mode)
	case gs.IntroState{}.String():
		vrt.RunTestGame(&gs.IntroState{}, mode)
	case gs.InventoryMenuState{}.String():
		st := &gs.InventoryMenuState{}
		st.SetCategory(states.ItemCategoryTypeCard)
		vrt.RunTestGame(st, mode)
	case gs.MainMenuState{}.String():
		vrt.RunTestGame(&gs.MainMenuState{}, mode)
	case gs.RayFieldState{}.String():
		vrt.RunTestGame(&gs.RayFieldState{}, mode)
	default:
		return fmt.Errorf("対応してないステートが指定された: %s", mode)
	}

	return nil
}
