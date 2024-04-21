package vrt_test

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	gs "github.com/kijimaD/ruins/lib/states"
	"github.com/stretchr/testify/assert"
)

// 1プロセスで、複数RunGameを呼べないため、ここで複数ケースのテストができない
// https://github.com/hajimehoshi/ebiten/blob/be771268ede283303836afc5823389429b87fddd/run.go#L290
// Don't call RunGame or RunGameWithOptions twice or more in one process.
func TestRunTestGame(t *testing.T) {
	RunTestGame(t, &gs.MainMenuState{}, "MainMenu")
}

// ================

var regularTermination = errors.New("テスト環境における、想定どおりの終了")

type TestGame struct {
	game.MainGame
	gameCount  int
	T          *testing.T
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

func (g *TestGame) Draw(screen *ebiten.Image) {
	g.StateMachine.Draw(g.World, screen)

	file, err := os.Create(fmt.Sprintf("%s.png", g.outputPath))
	assert.NoError(g.T, err)
	defer file.Close()

	err = png.Encode(file, screen)
	assert.NoError(g.T, err)
}

func RunTestGame(t *testing.T, state states.State, outputPath string) {
	world := game.InitWorld(960, 720)

	g := &TestGame{
		MainGame: game.MainGame{
			World:        world,
			StateMachine: es.Init(state, world),
		},
		gameCount:  0,
		T:          t,
		outputPath: outputPath,
	}

	if err := ebiten.RunGame(g); err != nil && err != regularTermination {
		assert.NoError(t, err)
	}
}
