package states

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStateNoReuse はステートが再利用されないことを確認するテスト
func TestStateNoReuse(t *testing.T) {
	t.Parallel()

	t.Run("StateFactoryから毎回異なるインスタンスが作成される", func(t *testing.T) {
		t.Parallel()

		// ファクトリー関数を定義
		factory := func() State {
			return &TestState{name: "TestState"}
		}

		// 複数回ファクトリー関数を実行
		state1 := factory()
		state2 := factory()
		state3 := factory()

		// 各インスタンスが異なることを確認（ポインタアドレスが異なる）
		assert.NotSame(t, state1, state2, "state1とstate2は異なるインスタンスである必要があります")
		assert.NotSame(t, state2, state3, "state2とstate3は異なるインスタンスである必要があります")
		assert.NotSame(t, state1, state3, "state1とstate3は異なるインスタンスである必要があります")
	})

	t.Run("TransitionのStateFactoriesが実行時に新しいインスタンスを作成する", func(t *testing.T) {
		t.Parallel()
		world := createTestWorld(t)

		// カウンターを使用して各インスタンスを追跡
		instanceCount := 0
		factory := func() State {
			instanceCount++
			return &TestStateWithID{
				TestState: TestState{name: "TestState"},
				ID:        instanceCount,
			}
		}

		// Transitionを作成
		transition := Transition{
			Type:          TransPush,
			NewStateFuncs: []StateFactory{factory, factory},
		}

		// StateMachineを初期化
		initialState := &TestState{name: "Initial"}
		sm := Init(initialState, world)
		sm.lastTransition = transition

		// 最初の実行
		sm.Update(world)
		assert.Equal(t, 2, instanceCount, "2つのステートが作成されるべき")

		// 同じTransitionで再実行
		sm.lastTransition = transition
		sm.Update(world)
		assert.Equal(t, 4, instanceCount, "さらに2つの新しいステートが作成されるべき")
	})

	t.Run("複数のPush操作で毎回新しいインスタンスが作成される", func(t *testing.T) {
		t.Parallel()
		world := createTestWorld(t)

		// 作成されたインスタンスを追跡
		createdStates := []*TestStateWithID{}
		idCounter := 0

		factory := func() State {
			idCounter++
			state := &TestStateWithID{
				TestState: TestState{name: "BattleState"},
				ID:        idCounter,
			}
			createdStates = append(createdStates, state)
			return state
		}

		// StateMachineを初期化
		sm := Init(&TestState{name: "Initial"}, world)

		// 複数回同じStateFactoryでPush
		for i := 0; i < 3; i++ {
			sm.lastTransition = Transition{
				Type:          TransPush,
				NewStateFuncs: []StateFactory{factory},
			}
			sm.Update(world)
		}

		// 3つの異なるインスタンスが作成されたことを確認
		assert.Equal(t, 3, len(createdStates), "3つのステートが作成されるべき")

		// 各インスタンスのIDが異なることを確認
		assert.Equal(t, 1, createdStates[0].ID)
		assert.Equal(t, 2, createdStates[1].ID)
		assert.Equal(t, 3, createdStates[2].ID)

		// ポインタが異なることを確認
		for i := 0; i < len(createdStates)-1; i++ {
			for j := i + 1; j < len(createdStates); j++ {
				assert.NotSame(t, createdStates[i], createdStates[j],
					"インスタンス%dとインスタンス%dは異なるべき", i, j)
			}
		}
	})
}

// TestStateWithID はIDを持つテスト用ステート
type TestStateWithID struct {
	TestState
	ID int
}
