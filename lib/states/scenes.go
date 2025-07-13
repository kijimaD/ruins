package states

import (
	"strings"

	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// 汎用アイテム入手イベント
func ItemGetEvent1() []es.State {
	ss := []es.State{}
	ss = push(ss, &MessageState{text: "「倉庫だな。役立ちそうなものはもらっていこう」"})
	ss = push(ss, NewExecState(func(world w.World) {
		// TODO: アイテム入手テーブルから獲得するようにする
		worldhelper.PlusAmount("鉄", 1, world)
		gamelog.SceneLog.Append("鉄を1個手に入れた")
		worldhelper.PlusAmount("木の棒", 1, world)
		gamelog.SceneLog.Append("木の棒を1個手に入れた")
		worldhelper.PlusAmount("フェライトコア", 1, world)
		gamelog.SceneLog.Append("フェライトコアを2個手に入れた")
	}))
	ss = push(ss, &MessageState{
		textFunc: utils.GetPtr(func() string {
			return strings.Join(gamelog.SceneLog.Pop(), "\n")
		})})

	return ss
}

// 汎用戦闘イベント開始
func RaidEvent1() []es.State {
	ss := []es.State{}
	ss = push(ss, &MessageState{text: "「何か動いた」\n「...敵だ!」"})
	ss = push(ss, &BattleState{})
	ss = push(ss, &MessageState{text: "「びっくりしたな」\n「おや、何か落ちてるぞ」"})
	ss = push(ss, NewExecState(func(world w.World) {
		worldhelper.PlusAmount("鉄", 1, world)
		// FIXME: どうして複数回実行したときSceneLogが蓄積してないのかわからない
		// Flush()してないのに...
		// アイテム入手→戦闘イベントとすると前の画面が表示される
		gamelog.SceneLog.Append("鉄を1個手に入れた")
	}))
	ss = push(ss, &MessageState{
		textFunc: utils.GetPtr(func() string {
			return strings.Join(gamelog.SceneLog.Pop(), "\n")
		})})

	return ss
}

func push(stack []es.State, value es.State) []es.State {
	return append([]es.State{value}, stack...)
}
