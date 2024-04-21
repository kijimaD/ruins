package vrt_test

import (
	"testing"

	gs "github.com/kijimaD/ruins/lib/states"
	"github.com/kijimaD/ruins/lib/vrt"
)

// 1プロセスで、複数RunGameを呼べないため、ここで複数ケースのテストができない
// https://github.com/hajimehoshi/ebiten/blob/be771268ede283303836afc5823389429b87fddd/run.go#L290
// Don't call RunGame or RunGameWithOptions twice or more in one process.
func TestRunTestGame(t *testing.T) {
	vrt.RunTestGame(t, &gs.MainMenuState{}, "MainMenu")
}
