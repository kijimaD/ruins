package states

import (
	"testing"

	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	t.Run("先頭に追加される", func(t *testing.T) {
		result := Push(
			[]es.State{
				&MessageState{},
				&BattleState{},
			},
			&ExecState{},
		)
		expect := []es.State{
			&ExecState{},
			&MessageState{},
			&BattleState{},
		}
		assert.Equal(t, expect, result)
	})
}
