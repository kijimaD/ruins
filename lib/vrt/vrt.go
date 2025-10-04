package vrt

import (
	"errors"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path"

	"github.com/kijimaD/ruins/lib/config"
	gs "github.com/kijimaD/ruins/lib/systems"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/maingame"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// エラーを返さないと実行終了しないため
var errRegularTermination = errors.New("テスト環境における、想定どおりの終了")

// TestGame はビジュアルリグレッションテスト用のゲーム構造体
type TestGame struct {
	maingame.MainGame
	gameCount  int
	outputPath string
}

// Update はゲームの更新処理を行う
func (g *TestGame) Update() error {
	// テストの前に実行される
	g.StateMachine.Update(g.World)

	// 10フレームだけ実行する。更新→描画の順なので、1度は更新しないと描画されない
	if g.gameCount < 10 {
		g.gameCount++
		return nil
	}

	// エラーを返さないと、実行終了しない
	return errRegularTermination
}

const outputDirName = "vrtimages"
const dirPerm = 0o755

// Draw はゲームの描画処理を行う
func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	// テストでは保存しない
	if flag.Lookup("test.v") != nil {
		return
	}

	if err := os.Mkdir(outputDirName, dirPerm); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	file, err := os.Create(path.Join(outputDirName, fmt.Sprintf("%s.png", g.outputPath)))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	err = png.Encode(file, screen)
	if err != nil {
		log.Fatal(err)
	}
}

// RunTestGame はテストゲームを実行してスクリーンショットを保存する
func RunTestGame(state es.State[w.World], outputPath string) {
	// VRT用にアニメーションを無効化（シングルトンインスタンスを直接変更）
	cfg := config.Get()
	originalConfig := *cfg
	cfg.DisableAnimation = true
	// テスト終了後に設定を復元
	defer func() {
		*cfg = originalConfig
	}()

	world, err := maingame.InitWorld(960, 720)
	if err != nil {
		panic(fmt.Sprintf("InitWorld failed: %v", err))
	}

	// デバッグデータを初期化
	worldhelper.InitDebugData(world)

	// 装備変更後にステータスを更新
	if changed := gs.EquipmentChangedSystem(world); !changed {
		log.Println("Equipment change was not detected")
	}

	g := &TestGame{
		MainGame: maingame.MainGame{
			World:        world,
			StateMachine: es.Init(state, world),
		},
		gameCount:  0,
		outputPath: outputPath,
	}

	if err := ebiten.RunGame(g); err != nil && err != errRegularTermination {
		log.Fatal(err)
	}
}
