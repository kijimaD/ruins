package states

import (
	"strings"

	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/helpers"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// then は次に実行するステートを追加する
// ここで設定したstate factory要素は順に先頭にpushされる。そのとき時系列が逆転する
// なのでstate factoryは、あとに表示するstateが先にくる
// state factory 3, 2, 1
// state stack 1, 2, 3
func then(stack []es.StateFactory, value es.StateFactory) []es.StateFactory {
	return append([]es.StateFactory{value}, stack...)
}

// GetItemGetEvent1Factories は汎用アイテム入手イベントのファクトリー関数配列を返す
func GetItemGetEvent1Factories() []es.StateFactory {
	factories := []es.StateFactory{}

	factories = then(factories, NewMessageStateWithText("「倉庫だな。役立ちそうなものはもらっていこう。」"))
	factories = then(factories, NewExecStateWithFunc(func(world w.World) {
		// TODO: アイテム入手テーブルから獲得するようにする
		worldhelper.PlusAmount("鉄", 1, world)
		gamelog.New(gamelog.FieldLog).
			ItemName("鉄").
			Append("を1個手に入れた。").
			Log()
		worldhelper.PlusAmount("木の棒", 1, world)
		gamelog.New(gamelog.FieldLog).
			ItemName("木の棒").
			Append("を1個手に入れた。").
			Log()
		worldhelper.PlusAmount("フェライトコア", 1, world)
		gamelog.New(gamelog.FieldLog).
			ItemName("フェライトコア").
			Append("を2個手に入れた。").
			Log()
	}))
	factories = then(factories, func() es.State {
		return &MessageState{
			textFunc: helpers.GetPtr(func() string {
				history := gamelog.SceneLog.GetHistory()
				gamelog.SceneLog.Clear()
				return strings.Join(history, "\n")
			}),
		}
	})

	return factories
}
