package game

import (
	"fmt"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	er "github.com/kijimaD/ruins/lib/engine/resources"
	es "github.com/kijimaD/ruins/lib/engine/states"
	gr "github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// MainGame はebiten.Game interfaceを満たす
type MainGame struct {
	World        w.World
	StateMachine es.StateMachine
}

// Layout はinterface methodのため、シグネチャは変更できない
func (game *MainGame) Layout(_, _ int) (int, int) {
	x, y := gr.UpdateGameLayout(game.World)

	return int(x), int(y)
}

// Update はゲームの更新処理を行う
func (game *MainGame) Update() error {
	// デバッグ表示をトグルする
	cfg := config.Get()
	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		// パフォーマンスモニターは攻略に関係ないのでトグルできてよい
		cfg.ShowMonitor = !cfg.ShowMonitor
	}
	if cfg.Debug && inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		cfg.ShowAIDebug = !cfg.ShowAIDebug
		cfg.NoEncounter = !cfg.NoEncounter
	}

	game.StateMachine.Update(game.World)

	return nil
}

// Draw はゲームの描画処理を行う
func (game *MainGame) Draw(screen *ebiten.Image) {
	game.StateMachine.Draw(game.World, screen)

	cfg := config.Get()
	if cfg.ShowMonitor {
		msg := getPerformanceInfo()
		ebitenutil.DebugPrint(screen, msg)
	}
}

// getPerformanceInfo はパフォーマンス情報を文字列として返す
func getPerformanceInfo() string {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	// 最後のGCからの経過時間を計算
	var lastGCTime string
	if mem.LastGC > 0 {
		lastGC := time.Unix(0, int64(mem.LastGC))
		elapsed := time.Since(lastGC)
		lastGCTime = fmt.Sprintf("%.2fs", elapsed.Seconds())
	} else {
		lastGCTime = "N/A"
	}

	return fmt.Sprintf(`FPS: %.1f
Alloc: %.2fMB
HeapInuse: %.2fMB
StackInuse: %.2fMB
Sys: %.2fMB
NextGC: %.2fMB
TotalAlloc: %.2fMB
Mallocs: %d
Frees: %d
GC: %d
LastGC: %s
PauseTotalNs: %.2fms
Goroutines: %d
`,
		ebiten.ActualFPS(),
		float64(mem.Alloc/1024/1024),      // 現在割り当てられているメモリ
		float64(mem.HeapInuse/1024/1024),  // ヒープで実際に使用中のメモリ
		float64(mem.StackInuse/1024/1024), // スタックで使用中のメモリ
		float64(mem.Sys/1024/1024),        // OSから取得した総メモリ
		float64(mem.NextGC/1024/1024),     // 次回GC実行予定サイズ
		float64(mem.TotalAlloc/1024/1024), // 起動後から割り当てられたヒープオブジェクトの累計バイト数
		mem.Mallocs,                       // 割り当てられたヒープオブジェクトの回数
		mem.Frees,                         // 解放されたヒープオブジェクトの回数
		mem.NumGC,                         // GC実行回数
		lastGCTime,                        // 最後のGC実行からの経過時間
		float64(mem.PauseTotalNs)/1000000, // GC停止時間の累計（ミリ秒）
		runtime.NumGoroutine(),            // 実行中のGoroutine数
	)
}

// InitWorld はゲームワールドを初期化する
func InitWorld(minGameWidth int, minGameHeight int) (w.World, error) {
	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		return w.World{}, err
	}

	world.Resources.ScreenDimensions = &er.ScreenDimensions{Width: minGameWidth, Height: minGameHeight}

	// ResourceManagerを使用してリソースを読み込む
	resourceManager := gr.NewDefaultResourceManager()

	// Load sprite sheets
	spriteSheets, err := resourceManager.LoadSpriteSheets()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts, err := resourceManager.LoadFonts()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.Fonts = &fonts

	defaultFont := (*world.Resources.Fonts)["kappa"]
	world.Resources.DefaultFaces = &map[string]text.Face{
		"kappa": defaultFont.Font,
	}

	// load UI resources
	uir, err := er.NewUIResources(defaultFont.FaceSource)
	if err != nil {
		return w.World{}, err
	}
	world.Resources.UIResources = uir

	// load raws
	rw, err := resourceManager.LoadRaws()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.RawMaster = rw

	gameResource := &gr.Dungeon{
		ExploredTiles: make(map[string]bool),
		Minimap: gr.MinimapSettings{
			Width:   150,
			Height:  150,
			OffsetX: 10,
			OffsetY: 10,
			Scale:   3,
		},
	}
	gameResource.SetStateEvent(gr.StateEventNone)
	world.Resources.Dungeon = gameResource

	return world, nil
}
