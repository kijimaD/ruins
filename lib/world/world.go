// Package world はゲームワールドの実装を提供する。
package world

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/resources"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// System は全てのゲームシステムが実装すべきインターフェース
// Wはworld型を表すジェネリック型パラメータ
type System[W any] interface {
	// String はシステム名を返す
	// map[string]Systemのキーとして使用される
	String() string

	// Draw は描画処理を行う
	Draw(world W, screen *ebiten.Image) error

	// Update は更新処理を行う
	Update(world W) error
}

// World はゲーム全体に必要な情報を保持する
type World struct {
	Manager    *ecs.Manager
	Components *gc.Components
	Resources  *resources.Resources
	Systems    map[string]System[World]
}

// InitWorld は初期化する
func InitWorld(c *gc.Components) (World, error) {
	manager := ecs.NewManager()
	err := c.InitializeComponents(manager)
	if err != nil {
		return World{}, err
	}
	return World{
		Manager:    manager,
		Components: c,
		Resources:  resources.InitGameResources(),
		Systems:    make(map[string]System[World]),
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
