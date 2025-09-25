package cmd

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/messagedata"
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
	case gs.HomeMenuState{}.String():
		vrt.RunTestGame(&gs.HomeMenuState{}, mode)
	case gs.InventoryMenuState{}.String():
		vrt.RunTestGame(&gs.InventoryMenuState{}, mode)
	case gs.LoadMenuState{}.String():
		vrt.RunTestGame(&gs.LoadMenuState{}, mode)
	case gs.MainMenuState{}.String():
		vrt.RunTestGame(&gs.MainMenuState{}, mode)
	case gs.MessageState{}.String():
		messageData := messagedata.NewDialogMessage(
			"これはメッセージウィンドウのVRTテストです。\n\n表示状態の確認用メッセージです。",
			"VRTテスト",
		).WithChoice(
			"選択肢1", func() {},
		).WithChoice(
			"選択肢2", func() {},
		)
		vrt.RunTestGame(gs.NewMessageState(messageData), mode)
	case gs.SaveMenuState{}.String():
		vrt.RunTestGame(&gs.SaveMenuState{}, mode)
	case gs.DungeonState{}.String():
		// いい感じのseed値。画面内に敵がいると動いて差分が出てしまうので、いないものを選んだ
		const seedVal = 4012
		vrt.RunTestGame(&gs.DungeonState{Depth: 1, Seed: seedVal}, mode)
	case "GameOver":
		vrt.RunTestGame(gs.NewGameOverMessageState(), mode)
	default:
		return fmt.Errorf("スクリーンショット実行時に対応してないステートが指定された: %s", mode)
	}

	return nil
}
