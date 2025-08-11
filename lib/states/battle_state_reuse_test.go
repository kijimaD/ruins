package states

import (
	"testing"

	es "github.com/kijimaD/ruins/lib/engine/states"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
)

// TestBattleStateNoReuse はBattleStateが再利用されないことを確認するテスト
func TestBattleStateNoReuse(t *testing.T) {
	t.Parallel()

	t.Run("NewBattleStateファクトリーが毎回新しいインスタンスを作成する", func(t *testing.T) {
		t.Parallel()

		// ファクトリー関数から複数回インスタンスを作成
		state1 := NewBattleState()
		state2 := NewBattleState()
		state3 := NewBattleState()

		// 各インスタンスが異なることを確認
		assert.NotSame(t, state1, state2, "state1とstate2は異なるインスタンスである必要があります")
		assert.NotSame(t, state2, state3, "state2とstate3は異なるインスタンスである必要があります")
		assert.NotSame(t, state1, state3, "state1とstate3は異なるインスタンスである必要があります")
	})

	t.Run("debugMenuTransから戦闘開始時に毎回新しいBattleStateが作成される", func(t *testing.T) {
		t.Parallel()

		// 戦闘開始のTransitionを取得（「戦闘開始(複数)」に変更）
		var battleTransition es.Transition
		for _, item := range debugMenuTrans {
			if item.label == "戦闘開始(複数)" {
				battleTransition = item.getTransFunc()
				break
			}
		}

		// Transitionが見つかったことを確認
		assert.Equal(t, es.TransPush, battleTransition.Type)
		assert.Len(t, battleTransition.NewStateFuncs, 1)

		// ファクトリーから複数回ステートを作成
		factory := battleTransition.NewStateFuncs[0]
		state1 := factory()
		state2 := factory()
		state3 := factory()

		// 各インスタンスが異なることを確認
		assert.NotSame(t, state1, state2, "毎回異なるBattleStateインスタンスが作成されるべき")
		assert.NotSame(t, state2, state3, "毎回異なるBattleStateインスタンスが作成されるべき")
		assert.NotSame(t, state1, state3, "毎回異なるBattleStateインスタンスが作成されるべき")

		// 型がBattleStateであることを確認
		_, ok1 := state1.(*BattleState)
		_, ok2 := state2.(*BattleState)
		_, ok3 := state3.(*BattleState)
		assert.True(t, ok1, "state1はBattleState型であるべき")
		assert.True(t, ok2, "state2はBattleState型であるべき")
		assert.True(t, ok3, "state3はBattleState型であるべき")
	})

	t.Run("debugMenuTransから複数回戦闘開始で異なるインスタンスが作成される", func(t *testing.T) {
		t.Parallel()

		// 戦闘開始のTransitionを直接テスト
		var battleEntry *struct {
			label        string
			f            func(world w.World)
			getTransFunc func() es.Transition
		}

		for i := range debugMenuTrans {
			if debugMenuTrans[i].label == "戦闘開始(複数)" {
				battleEntry = &debugMenuTrans[i]
				break
			}
		}
		assert.NotNil(t, battleEntry, "戦闘開始メニュー項目が見つかるべき")

		// 複数回Transitionを取得してファクトリーからステートを作成
		createdStates := []es.State{}

		for i := 0; i < 3; i++ {
			transition := battleEntry.getTransFunc()
			assert.Equal(t, es.TransPush, transition.Type)
			assert.Len(t, transition.NewStateFuncs, 1)

			// ファクトリーからステートを作成
			newState := transition.NewStateFuncs[0]()
			createdStates = append(createdStates, newState)
		}

		// 作成された全てのステートが異なるインスタンスであることを確認
		for i := 0; i < len(createdStates)-1; i++ {
			for j := i + 1; j < len(createdStates); j++ {
				assert.NotSame(t, createdStates[i], createdStates[j],
					"戦闘ステート%dと%dは異なるインスタンスであるべき", i, j)
			}
		}
	})
}

// TestStateReusePrevention は以前の問題（ステート再利用による即座終了）が解決されていることを確認
func TestStateReusePrevention(t *testing.T) {
	t.Parallel()

	t.Run("BattleStateのtransフィールドが新しいインスタンスで初期化される", func(t *testing.T) {
		t.Parallel()

		// 最初のBattleStateを作成してTransPopを設定（OnStartは呼ばない）
		battle1 := NewBattleState().(*BattleState)
		battle1.SetTransition(es.Transition{Type: es.TransPop})

		// GetTransitionでTransPopが設定されていることを確認
		trans1 := battle1.GetTransition()
		assert.NotNil(t, trans1)
		assert.Equal(t, es.TransPop, trans1.Type)

		// 新しいBattleStateを作成
		battle2 := NewBattleState().(*BattleState)

		// 新しいインスタンスではtransがnilであることを確認
		trans2 := battle2.GetTransition()
		assert.Nil(t, trans2, "新しいBattleStateインスタンスのtransフィールドはnilであるべき")

		// ConsumeTransitionがTransNoneを返すことを確認
		consumed := battle2.ConsumeTransition()
		assert.Equal(t, es.TransNone, consumed.Type, "初期状態ではTransNoneが返されるべき")
	})
}
