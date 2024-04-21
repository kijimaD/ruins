package vrt

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
)

var regularTermination = errors.New("テスト環境における、想定どおりの終了")

type TestGame struct {
	game.MainGame
	gameCount  int
	outputPath string
}

func (g *TestGame) Update() error {
	// テストの前に実行される
	g.StateMachine.Update(g.World)

	// 1フレームだけ実行する。更新→描画の順なので、1度は更新しないと描画されない
	if g.gameCount < 1 {
		g.gameCount += 1
		return nil
	}

	// エラーを返さないと、実行終了しない
	return regularTermination
}

const outputDirName = "vrtimages"
const dirPerm = 0o755

func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	_ = os.Mkdir(outputDirName, dirPerm)
	file, err := os.Create(path.Join(outputDirName, fmt.Sprintf("%s.png", g.outputPath)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, screen)
	if err != nil {
		log.Fatal(err)
	}
}

func RunTestGame(state states.State, outputPath string) {
	world := game.InitWorld(960, 720)

	g := &TestGame{
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(state, world),
		},
		gameCount:  0,
		outputPath: outputPath,
	}

	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		log.Fatal(err)
	}
}
