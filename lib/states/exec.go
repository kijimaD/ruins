package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
)

// ステート管理に乗っかりつつ任意のコマンドを実行するためのダミーステート
// 処理の中身は呼び出し側で注入する
type ExecState struct {
	es.BaseState
	f func(w.World)
}

func (st ExecState) String() string {
	return "Exec"
}

// NewExecState は新しいExecStateを作成する
func NewExecState(f func(w.World)) *ExecState {
	return &ExecState{
		f: f,
	}
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

func (st *ExecState) Update(world w.World) es.Transition {
	// BaseStateの共通処理を使用
	if transition := st.ConsumeTransition(); transition.Type != es.TransNone {
		return transition
	}
	return es.Transition{Type: es.TransPop}
}

func (st *ExecState) Draw(world w.World, screen *ebiten.Image) {
	// 何も表示しないので、ユーザーにはわからない状態
	// デバッグ用に何かを表示したい場合はここに追加
}
