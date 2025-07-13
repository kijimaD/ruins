package states

import (
	"testing"

	"github.com/kijimaD/ruins/lib/engine/resources"
	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	t.Run("先頭に追加される", func(t *testing.T) {
		execState := NewExecState(func(w.World) {})
		result := push(
			[]es.State{
				&MessageState{},
				&BattleState{},
			},
			execState,
		)
		expect := []es.State{
			execState,
			&MessageState{},
			&BattleState{},
		}
		assert.Equal(t, expect, result)
	})
}

func TestExecStateAutoExecution(t *testing.T) {
	t.Run("OnResumeで関数自動実行", func(t *testing.T) {
		executed := false
		execState := NewExecState(func(world w.World) {
			executed = true
		})

		// テスト用のワールドを作成
		world := w.World{
			Resources: &resources.Resources{},
		}

		// 初期状態では関数が実行されていない
		assert.False(t, executed)

		// OnResumeで関数が実行される
		execState.OnResume(world)
		assert.True(t, executed)

		// Updateで即座にTransPopが返される
		trans := execState.Update(world)
		assert.Equal(t, es.TransPop, trans.Type)
	})
}
