// Package world はゲームワールドの実装を提供する。
package world

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// World はゲーム全体に必要な情報を保持する
type World struct {
	Manager    *ecs.Manager
	Components *gc.Components
	Resources  *resources.Resources
}

// InitWorld は初期化する
func InitWorld(c *gc.Components) (World, error) {
	manager := ecs.NewManager()
	err := c.InitializeComponents(manager)
	if err != nil {
		return World{}, err
	}
	gameResources := resources.InitGameResources()

	return World{
		Manager:    manager,
		Components: c,
		Resources:  gameResources,
	}, nil
}

// GetManager は World interfaceを満たすためのメソッド
func (w World) GetManager() *ecs.Manager {
	return w.Manager
}

// GetComponents は World interfaceを満たすためのメソッド
func (w World) GetComponents() interface{} {
	return w.Components
}
