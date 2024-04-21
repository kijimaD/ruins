package main

import (
	"log"
	"os"

	_ "net/http/pprof"

	"github.com/kijimaD/ruins/lib/cmd"
)

const (
	minGameWidth  = 960
	minGameHeight = 720
)

func main() {
	app := cmd.NewMainApp()
	err := cmd.RunMainApp(app, os.Args...)
	if err != nil {
		log.Fatal(err)
	}
}
