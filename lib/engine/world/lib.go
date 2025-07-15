package world

import (
	c "github.com/kijimaD/ruins/lib/engine/components"
	"github.com/kijimaD/ruins/lib/engine/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// Generic は型安全なワールド型
type Generic[T c.ComponentInitializer] struct {
	Manager    *ecs.Manager
	Components *c.Components[T]
	Resources  *resources.Resources
}

// InitGeneric は型安全なワールド初期化
func InitGeneric[T c.ComponentInitializer](gameComponents T) (Generic[T], error) {
	manager := ecs.NewManager()
	components, err := c.InitComponents(manager, gameComponents)
	if err != nil {
		return Generic[T]{}, err
	}
	resources := resources.InitResources()

	return Generic[T]{
		Manager:    manager,
		Components: components,
		Resources:  resources,
	}, nil
}
