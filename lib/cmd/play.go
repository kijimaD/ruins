package cmd

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/maingame"
	"github.com/kijimaD/ruins/lib/worldhelper"
	"github.com/pkg/profile"
	"github.com/urfave/cli/v2"

	_ "net/http/pprof" // pprofのHTTPエンドポイントを登録するためのインポート

	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/mapplanner"
	gs "github.com/kijimaD/ruins/lib/states"
	w "github.com/kijimaD/ruins/lib/world"
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
	// 設定を読み込み
	cfg := config.Get()

	// ログ設定を読み込み
	logger.LoadFromConfig(cfg.LogLevel, cfg.LogCategories)

	// デバッグモードの場合は設定を表示
	if cfg.Debug {
		log.Printf("Configuration loaded:\n%s", cfg.String())
	}

	// ウィンドウ設定
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(cfg.WindowWidth, cfg.WindowHeight)
	ebiten.SetWindowTitle("ruins")

	// フルスクリーン設定
	if cfg.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	// FPS設定
	if cfg.TargetFPS != 60 {
		ebiten.SetTPS(cfg.TargetFPS)
	}

	// プロファイラー設定（WASMは除外）
	if runtime.GOOS != "js" && cfg.DebugPProf {
		var profileOptions []func(*profile.Profile)

		if cfg.ProfileMemory {
			profileOptions = append(profileOptions, profile.MemProfile)
		}
		if cfg.ProfileCPU {
			profileOptions = append(profileOptions, profile.CPUProfile)
		}
		if cfg.ProfileMutex {
			profileOptions = append(profileOptions, profile.MutexProfile)
		}
		if cfg.ProfileTrace {
			profileOptions = append(profileOptions, profile.TraceProfile)
		}

		// デフォルトでメモリプロファイルを有効化
		if len(profileOptions) == 0 {
			profileOptions = append(profileOptions, profile.MemProfile)
		}

		profileOptions = append(profileOptions, profile.ProfilePath(cfg.ProfilePath))
		defer profile.Start(profileOptions...).Stop()

		// pprofサーバー起動
		pprofAddr := fmt.Sprintf("localhost:%d", cfg.PProfPort)
		go func() {
			log.Printf("pprof server starting on http://%s", pprofAddr)
			log.Fatal(http.ListenAndServe(pprofAddr, nil))
		}()
	}

	world, err := maingame.InitWorld(cfg.WindowWidth, cfg.WindowHeight)
	if err != nil {
		return err
	}

	// デバッグ用データ初期化
	worldhelper.InitDebugData(world)

	// 開始ステートの決定
	var initialState es.State[w.World]
	switch cfg.StartingState {
	case "town":
		stateFactory := gs.NewDungeonStateWithBuilder(1, mapplanner.PlannerTypeTown)
		initialState = stateFactory()
	case "main_menu":
		initialState = &gs.MainMenuState{}
	default:
		log.Fatalf("無効なstate: %s", cfg.StartingState)
	}

	err = ebiten.RunGame(&maingame.MainGame{
		World:        world,
		StateMachine: es.Init(initialState, world),
	})
	if err != nil {
		return err
	}

	return nil
}
