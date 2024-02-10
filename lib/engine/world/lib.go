package world

import (
	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// World is the main ECS structure
type World struct {
	Manager    *ecs.Manager
	Components *c.Components
	Resources  *resources.Resources
}

// InitWorld initializes the world
func InitWorld(gameComponents interface{}) World {
	manager := ecs.NewManager()
	components := c.InitComponents(manager, gameComponents)
	resources := resources.InitResources()

	return World{
		Manager:    manager,
		Components: components,
		Resources:  resources,
	}
}
