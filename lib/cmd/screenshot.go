package cmd

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/messagedata"
	gs "github.com/kijimaD/ruins/lib/states"
	"github.com/kijimaD/ruins/lib/vrt"
	w "github.com/kijimaD/ruins/lib/world"
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
	case "DebugMenu":
		vrt.RunTestGame(gs.NewDebugMenuState(), mode)
	case gs.DungeonState{}.String():
		// 固定seed値を使用する
		const seedVal = 1
		vrt.RunTestGame(&gs.DungeonState{
			Depth:       1,
			Seed:        seedVal,
			BuilderType: mapplanner.PlannerTypeSmallRoom,
		}, mode)
	case gs.EquipMenuState{}.String():
		vrt.RunTestGame(&gs.EquipMenuState{}, mode)
	case "GameOver":
		vrt.RunTestGame(gs.NewGameOverMessageState(), mode)
	case "Town":
		// 固定seed値を使用する
		const townSeedVal = 1
		stateFactory := gs.NewDungeonState(1, gs.WithSeed(townSeedVal), gs.WithBuilderType(mapplanner.PlannerTypeTown))
		vrt.RunTestGame(stateFactory(), mode)
	case gs.InventoryMenuState{}.String():
		vrt.RunTestGame(&gs.InventoryMenuState{}, mode)
	case "LoadMenu":
		vrt.RunTestGame(gs.NewLoadMenuState(), mode)
	case gs.MainMenuState{}.String():
		vrt.RunTestGame(&gs.MainMenuState{}, mode)
	case gs.MessageState{}.String():
		messageData := messagedata.NewDialogMessage(
			"これはメッセージウィンドウのVRTテストです。\n\n表示状態の確認用メッセージです。",
			"VRTテスト",
		).WithChoice(
			"選択肢1", func(_ w.World) {},
		).WithChoice(
			"選択肢2", func(_ w.World) {},
		)
		vrt.RunTestGame(gs.NewMessageState(messageData), mode)
	case "SaveMenu":
		vrt.RunTestGame(gs.NewSaveMenuState(), mode)
	case gs.ShopMenuState{}.String():
		vrt.RunTestGame(&gs.ShopMenuState{}, mode)
	default:
		return fmt.Errorf("スクリーンショット実行時に対応してないステートが指定された: %s", mode)
	}

	return nil
}
