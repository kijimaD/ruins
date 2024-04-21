package vrt_test

import (
	"errors"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	gs "github.com/kijimaD/ruins/lib/states"
)

func TestAaa(t *testing.T) {
	RunTestGame()
}

// ================

var regularTermination = errors.New("テスト環境における、想定どおりの終了")

type TestGame struct {
	game.MainGame
	gameCount int
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

func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, screen)
	if err != nil {
		panic(err)
	}
}

func RunTestGame() {
	world := game.InitWorld(960, 720)
	g := &TestGame{
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(&gs.MainMenuState{}, world),
		},
		gameCount: 0,
	}

	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		log.Fatal(err)
	}
}
