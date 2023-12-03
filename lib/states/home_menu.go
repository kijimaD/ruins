// 拠点でのコマンド選択画面
package states

import (
	"github.com/kijimaD/sokotwo/lib/engine/states"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
)

type HomeMenuState struct{}

// State interface ================

func (st *HomeMenuState) OnPause(world w.World) {}

func (st *HomeMenuState) OnResume(world w.World) {}

func (st *HomeMenuState) OnStart(world w.World) {}

func (st *HomeMenuState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

func (st *HomeMenuState) Update(world w.World) states.Transition {
	return states.Transition{}
}
