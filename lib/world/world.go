// Package world はゲームワールドの実装を提供する。
package world

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// World はゲーム全体に必要な情報を保持する
type World struct {
	Manager    *ecs.Manager
	Components *gc.GameComponents
	Resources  *resources.Resources
}

// InitWorld は初期化する
func InitWorld(gameComponents *gc.Components) World {
	manager := ecs.NewManager()
	err := gameComponents.InitializeComponents(manager)
	if err != nil {
		log.Fatal(err)
	}
	resources := resources.InitResources()

	return World{
		Manager: manager,
		Components: &gc.GameComponents{
			Game: gameComponents,
		},
		Resources: resources,
	}
}
