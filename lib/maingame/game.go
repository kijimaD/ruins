package maingame

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/config"
	"github.com/kijimaD/ruins/lib/consts"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/loader"
	gr "github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// MainGame はebiten.Game interfaceを満たす
type MainGame struct {
	World        w.World
	StateMachine es.StateMachine[w.World]
}

// Layout はinterface methodのため、シグネチャは変更できない
func (game *MainGame) Layout(_, _ int) (int, int) {
	// TODO: 解像度変更は未実装
	return consts.MinGameWidth, consts.MinGameHeight
}

// Update はゲームの更新処理を行う
func (game *MainGame) Update() error {
	// デバッグ表示をトグルする
	cfg := config.Get()
	if ebiten.IsKeyPressed(ebiten.KeyShift) && inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		// パフォーマンスモニターは攻略に関係ないのでトグルできてよい
		cfg.ShowMonitor = !cfg.ShowMonitor
	}

	if err := game.StateMachine.Update(game.World); err != nil {
		return err
	}

	return nil
}

// Draw はゲームの描画処理を行う
func (game *MainGame) Draw(screen *ebiten.Image) {
	if err := game.StateMachine.Draw(game.World, screen); err != nil {
		log.Fatal(err)
	}

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

	world.Resources.SetScreenDimensions(minGameWidth, minGameHeight)

	// ResourceLoaderを使用してリソースを読み込む
	resourceLoader := loader.NewResourceLoader()

	// Load sprite sheets
	spriteSheets, err := resourceLoader.LoadSpriteSheets()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.SpriteSheets = &spriteSheets

	// load fonts
	fonts, err := resourceLoader.LoadFonts()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.Fonts = &fonts

	font := (*world.Resources.Fonts)["kappa"]
	world.Resources.Faces = &map[string]text.Face{
		"kappa": font.Font,
	}

	// load UI resources
	uir, err := gr.NewUIResources(font.FaceSource)
	if err != nil {
		return w.World{}, err
	}
	world.Resources.UIResources = uir

	// load raws
	rw, err := resourceLoader.LoadRaws()
	if err != nil {
		return w.World{}, err
	}
	world.Resources.RawMaster = rw

	gameResource := &gr.Dungeon{
		ExploredTiles: make(map[gc.GridElement]bool),
		MinimapSettings: gr.MinimapSettings{
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
