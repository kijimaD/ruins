package game

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	gc "github.com/kijimaD/ruins/lib/components"
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
	game.StateMachine.Update(game.World)

	return nil
}

// Draw はゲームの描画処理を行う
func (game *MainGame) Draw(screen *ebiten.Image) {
	game.StateMachine.Draw(game.World, screen)

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	msg := fmt.Sprintf(`FPS: %f
Alloc: %.2fMB
TotalAlloc: %.2fMB
Mallocs: %.2fMB
Frees: %.2fMB
`,

		ebiten.ActualFPS(),
		float64(mem.Alloc/1024/1024),
		float64(mem.TotalAlloc/1024/1024), // 起動後から割り当てられたヒープオブジェクトの数。freeされてもリセットされない
		float64(mem.Mallocs/1024/1024),    // 割り当てられているヒープオブジェクトの数。freeされたら減る
		float64(mem.Frees/1024/1024),      // 解放されたヒープオブジェクトの数
	)
	ebitenutil.DebugPrint(screen, msg)
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

	world.Resources.Game = &gr.Game{
		StateEvent: gr.StateEventNone,
	}

	return world, nil
}
