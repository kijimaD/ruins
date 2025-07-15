package world

import (
	gc "github.com/kijimaD/ruins/lib/components"
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

// World は後方互換性のためのデフォルト型
// TODO: 具体に移動する
type World struct {
	Manager    *ecs.Manager
	Components *gc.GameComponents
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

// InitWorld は後方互換性のためのラッパー関数
func InitWorld(gameComponents *gc.Components) World {
	manager := ecs.NewManager()
	err := gameComponents.InitializeComponents(manager)
	if err != nil {
		// 既存コードとの互換性のため、panicさせる
		panic(err)
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
