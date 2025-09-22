package states

import (
	"github.com/hajimehoshi/ebiten/v2"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/world"
)

// ExecState はステート管理に乗っかりつつ任意のコマンドを実行するためのダミーステート
// 処理の中身は呼び出し側で注入する
type ExecState struct {
	es.BaseState[w.World]
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

var _ es.State[w.World] = &ExecState{}

// OnPause はステートが一時停止される際に呼ばれる
func (st *ExecState) OnPause(_ w.World) {}

// OnResume はステートが再開される際に呼ばれる
func (st *ExecState) OnResume(world w.World) {
	st.f(world)
}

// OnStart はstate pushされたときはpush時にまとめてインスタンスを作成し、スタックトップに回ってきたときに使う。トップに回ってきたときに初期化しているわけではない。つまり、OnStartにf(world)を書くと複数stateをpushしたときに即実行されてしまう。そうではなく、スタックトップに来たときに実行してほしい。
// ほかのstateではswitchを使っていることが多い。switchではstate stackを再作成して1つだけpushする。なので、stateインスタンス生成とスタックトップに来たときの時差がないため問題が起きない。
func (st *ExecState) OnStart(_ w.World) {}

// OnStop はステートが停止される際に呼ばれる
func (st *ExecState) OnStop(_ w.World) {}

// Update はゲームステートの更新処理を行う
func (st *ExecState) Update(_ w.World) es.Transition[w.World] {
	// BaseStateの共通処理を使用
	if transition := st.ConsumeTransition(); transition.Type != es.TransNone {
		return transition
	}
	return es.Transition[w.World]{Type: es.TransPop}
}

// Draw はゲームステートの描画処理を行う
func (st *ExecState) Draw(_ w.World, _ *ebiten.Image) {
	// 何も表示しないので、ユーザーにはわからない状態
	// デバッグ用に何かを表示したい場合はここに追加
}
