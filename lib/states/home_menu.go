// 拠点でのコマンド選択画面
package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/loader"
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	"github.com/kijimaD/sokotwo/lib/resources"
)

type HomeMenuState struct{}

// State interface ================

func (st *HomeMenuState) OnPause(world w.World) {}

func (st *HomeMenuState) OnResume(world w.World) {}

func (st *HomeMenuState) OnStart(world w.World) {
	prefabs := world.Resources.Prefabs.(*resources.Prefabs)
	loader.AddEntities(world, prefabs.Menu.HomeMenu)
}

func (st *HomeMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *HomeMenuState) Update(world w.World) states.Transition {
	return states.Transition{}
}
