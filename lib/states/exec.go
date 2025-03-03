package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kijimaD/ruins/lib/engine/states"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

// ステート管理に乗っかりつつ任意のコマンドを実行するためのダミーステート
type ExecState struct {
	f func(w.World)
}

func (st ExecState) String() string {
	return "Exec"
}

// State interface ================

var _ es.State = &ExecState{}

func (st *ExecState) OnPause(world w.World) {}

func (st *ExecState) OnResume(world w.World) {
	st.f(world)
}

// state pushされたときはpush時にまとめてインスタンスを作成し、スタックトップに回ってきたときに使う。トップに回ってきたときに初期化しているわけではない。つまり、OnStartにf(world)を書くと複数stateをpushしたときに即実行されてしまう。そうではなく、スタックトップに来たときに実行してほしい。
// ほかのstateではswitchを使っていることが多い。switchではstate stackを再作成して1つだけpushする。なので、stateインスタンス生成とスタックトップに来たときの時差がないため問題が起きない。
func (st *ExecState) OnStart(world w.World) {}

func (st *ExecState) OnStop(world w.World) {}

func (st *ExecState) Update(world w.World) states.Transition {
	return states.Transition{Type: states.TransPop}
}

func (st *ExecState) Draw(world w.World, screen *ebiten.Image) {}
