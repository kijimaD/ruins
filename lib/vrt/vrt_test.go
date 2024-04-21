package vrt_test

import (
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	gs "github.com/kijimaD/ruins/lib/states"
)

func TestMain(m *testing.M) {
	RunTestGame(m)
}

func TestAaa(t *testing.T) {
	fmt.Println("===============")
}

// ================

var regularTermination = errors.New("regular termination")

type TestGame struct {
	game.MainGame
	m         *testing.M
	gameCount int
}

func (g *TestGame) Update() error {
	// テストの前に実行される
	g.StateMachine.Update(g.World)

	// 1フレームだけ実行する。描画するため
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

func RunTestGame(m *testing.M) {
	ebiten.SetWindowSize(960, 720)
	ebiten.SetInitFocused(false)
	ebiten.SetWindowTitle("Testing...")

	world := game.InitWorld(960, 720)
	g := &TestGame{
		m: m,
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(&gs.MainMenuState{}, world),
		},
	}

	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		log.Fatal(err)
	}
}
