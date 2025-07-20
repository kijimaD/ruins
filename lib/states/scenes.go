package states

import (
	"strings"

	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/helpers"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
)

// GetItemGetEvent1Factories は汎用アイテム入手イベントのファクトリー関数配列を返す
func GetItemGetEvent1Factories() []es.StateFactory {
	factories := []es.StateFactory{}
	factories = append(factories, NewMessageStateWithText("「倉庫だな。役立ちそうなものはもらっていこう」"))
	factories = append(factories, NewExecStateWithFunc(func(world w.World) {
		// TODO: アイテム入手テーブルから獲得するようにする
		worldhelper.PlusAmount("鉄", 1, world)
		gamelog.SceneLog.Append("鉄を1個手に入れた")
		worldhelper.PlusAmount("木の棒", 1, world)
		gamelog.SceneLog.Append("木の棒を1個手に入れた")
		worldhelper.PlusAmount("フェライトコア", 1, world)
		gamelog.SceneLog.Append("フェライトコアを2個手に入れた")
	}))
	factories = append(factories, func() es.State {
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
	factories = append(factories, NewMessageStateWithText("「何か動いた」\n「...敵だ!」"))
	factories = append(factories, NewBattleState)
	factories = append(factories, NewMessageStateWithText("「びっくりしたな」\n「おや、何か落ちてるぞ」"))
	factories = append(factories, NewExecStateWithFunc(func(world w.World) {
		worldhelper.PlusAmount("鉄", 1, world)
		// FIXME: どうして複数回実行したときSceneLogが蓄積してないのかわからない
		// Flush()してないのに...
		// アイテム入手→戦闘イベントとすると前の画面が表示される
		gamelog.SceneLog.Append("鉄を1個手に入れた")
	}))
	factories = append(factories, func() es.State {
		return &MessageState{
			textFunc: helpers.GetPtr(func() string {
				return strings.Join(gamelog.SceneLog.Pop(), "\n")
			}),
		}
	})
	return factories
}
