package main

import (
	"log"
	"runtime"

	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/profile"

	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/game"
	gs "github.com/kijimaD/ruins/lib/states"
)

const (
	minGameWidth  = 960
	minGameHeight = 720
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(minGameWidth, minGameHeight)
	ebiten.SetWindowTitle("ruins")

	// プロファイラ。WASMは除外する
	if runtime.GOOS != "js" {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	}
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	world := game.InitWorld(minGameWidth, minGameHeight)
	ebiten.RunGame(&game.MainGame{
		World:        world,
		StateMachine: es.Init(&gs.MainMenuState{}, world),
	})
}
