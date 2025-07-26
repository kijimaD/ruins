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
// stateは先頭から実行されていく。複数stateの定義が直感的に見えるようにする
// 1 2 3 ...
func then(stack []es.StateFactory, value es.StateFactory) []es.StateFactory {
	return append([]es.StateFactory{value}, stack...)
}

// GetItemGetEvent1Factories は汎用アイテム入手イベントのファクトリー関数配列を返す
func GetItemGetEvent1Factories() []es.StateFactory {
	factories := []es.StateFactory{}

	factories = then(factories, NewMessageStateWithText("「倉庫だな。役立ちそうなものはもらっていこう」"))
	factories = then(factories, NewExecStateWithFunc(func(world w.World) {
		// TODO: アイテム入手テーブルから獲得するようにする
		worldhelper.PlusAmount("鉄", 1, world)
		gamelog.SceneLog.Append("鉄を1個手に入れた。")
		worldhelper.PlusAmount("木の棒", 1, world)
		gamelog.SceneLog.Append("木の棒を1個手に入れた。")
		worldhelper.PlusAmount("フェライトコア", 1, world)
		gamelog.SceneLog.Append("フェライトコアを2個手に入れた。")
	}))
	factories = then(factories, func() es.State {
		return &MessageState{
			textFunc: helpers.GetPtr(func() string {
				return strings.Join(gamelog.SceneLog.Pop(), "\n")
			}),
		}
	})

	return factories
}

// GetRaidEvent1Factories は汎用戦闘イベントのファクトリー関数配列を返す
func GetRaidEvent1Factories() []es.StateFactory {
	factories := []es.StateFactory{}

	factories = then(factories, NewMessageStateWithText("「何か動いた」\n「...敵だ!」"))
	factories = then(factories, NewBattleState)
	factories = then(factories, NewMessageStateWithText("「びっくりしたな」\n「おや、何か落ちてるぞ」"))
	factories = then(factories, NewExecStateWithFunc(func(world w.World) {
		worldhelper.PlusAmount("鉄", 1, world)
		gamelog.SceneLog.Append("鉄を1個手に入れた")
	}))
	factories = then(factories, func() es.State {
		return &MessageState{
			textFunc: helpers.GetPtr(func() string {
				return strings.Join(gamelog.SceneLog.Pop(), "\n")
			}),
		}
	})

	return factories
}
