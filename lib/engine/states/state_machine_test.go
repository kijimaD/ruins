package states

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

// TestContext はテスト用のコンテキスト
type TestContext struct {
	Name string
}

// TestState はテスト用の状態実装
type TestState struct {
	name           string
	onStartCalled  bool
	onStopCalled   bool
	onPauseCalled  bool
	onResumeCalled bool
	updateCalled   bool
	drawCalled     bool
}

func (ts *TestState) String() string {
	return ts.name
}

func (ts *TestState) OnStart(_ TestContext) {
	ts.onStartCalled = true
}

func (ts *TestState) OnStop(_ TestContext) {
	ts.onStopCalled = true
}

func (ts *TestState) OnPause(_ TestContext) {
	ts.onPauseCalled = true
}

func (ts *TestState) OnResume(_ TestContext) {
	ts.onResumeCalled = true
}

func (ts *TestState) Update(_ TestContext) Transition[TestContext] {
	ts.updateCalled = true
	return Transition[TestContext]{Type: TransNone}
}

func (ts *TestState) Draw(_ TestContext, _ *ebiten.Image) {
	ts.drawCalled = true
}

// TestGetStatesMethods はGetStatesメソッド群のテスト
func TestGetStatesMethods(t *testing.T) {
	t.Parallel()
	t.Run("初期状態での動作確認", func(t *testing.T) {
		t.Parallel()
		ctx := TestContext{Name: "TestContext"}
		initialState := &TestState{name: "InitialState"}

		// StateMachineの初期化
		stateMachine := Init(initialState, ctx)

		// GetStatesのテスト
		states := stateMachine.GetStates()
		assert.Len(t, states, 1, "初期状態の数が正しくない")
		assert.Equal(t, "InitialState", states[0].(*TestState).name, "初期状態の名前が正しくない")

		// GetCurrentStateのテスト
		currentState := stateMachine.GetCurrentState()
		assert.NotNil(t, currentState, "現在の状態がnil")
		assert.Equal(t, "InitialState", currentState.(*TestState).name, "現在の状態の名前が正しくない")

		// GetStateCountのテスト
		stateCount := stateMachine.GetStateCount()
		assert.Equal(t, 1, stateCount, "状態数が正しくない")

		// OnStartが呼ばれていることを確認
		assert.True(t, initialState.onStartCalled, "OnStartが呼ばれていない")
	})

	t.Run("状態の不変性確認", func(t *testing.T) {
		t.Parallel()
		ctx := TestContext{Name: "TestContext"}
		initialState := &TestState{name: "InitialState"}
		stateMachine := Init(initialState, ctx)

		// GetStatesで取得したスライスを変更しても元のスタックに影響しないことを確認
		states := stateMachine.GetStates()
		originalLength := len(states)

		// 取得したスライスを変更
		_ = append(states, &TestState{name: "ModifiedState"})

		// 元のスタックは変更されていないことを確認
		newStates := stateMachine.GetStates()
		assert.Len(t, newStates, originalLength, "元の状態スタックが変更されている")
		assert.Equal(t, "InitialState", newStates[0].(*TestState).name, "元の状態が変更されている")
	})

	t.Run("空の状態スタックでの動作", func(t *testing.T) {
		t.Parallel()
		// 空のStateMachineを作成（実際のゲームでは発生しないが、テスト用）
		stateMachine := StateMachine[TestContext]{}

		// GetStatesのテスト
		states := stateMachine.GetStates()
		assert.Len(t, states, 0, "空のスタックの状態数が正しくない")

		// GetCurrentStateのテスト
		currentState := stateMachine.GetCurrentState()
		assert.Nil(t, currentState, "空のスタックの現在状態がnilでない")

		// GetStateCountのテスト
		stateCount := stateMachine.GetStateCount()
		assert.Equal(t, 0, stateCount, "空のスタックの状態数が正しくない")
	})
}

// TestStateMachineTransitions は状態遷移のテスト
func TestStateMachineTransitions(t *testing.T) {
	t.Parallel()
	t.Run("Push遷移のテスト", func(t *testing.T) {
		t.Parallel()
		ctx := TestContext{Name: "TestContext"}
		initialState := &TestState{name: "InitialState"}
		stateMachine := Init(initialState, ctx)

		// Push遷移を実行
		newState := &TestState{name: "PushedState"}
		stateMachine.lastTransition = Transition[TestContext]{
			Type:          TransPush,
			NewStateFuncs: []StateFactory[TestContext]{func() State[TestContext] { return newState }},
		}
		stateMachine.Update(ctx)

		// 状態数の確認
		assert.Equal(t, 2, stateMachine.GetStateCount(), "Push後の状態数が正しくない")

		// 現在の状態の確認
		currentState := stateMachine.GetCurrentState()
		assert.Equal(t, "PushedState", currentState.(*TestState).name, "Push後の現在状態が正しくない")

		// 状態スタックの確認
		states := stateMachine.GetStates()
		assert.Len(t, states, 2, "Push後の状態スタック数が正しくない")
		assert.Equal(t, "InitialState", states[0].(*TestState).name, "Push後の最初の状態が正しくない")
		assert.Equal(t, "PushedState", states[1].(*TestState).name, "Push後の最後の状態が正しくない")

		// 初期状態がPauseされていることを確認
		assert.True(t, initialState.onPauseCalled, "初期状態のOnPauseが呼ばれていない")
		// 新しい状態がStartされていることを確認
		assert.True(t, newState.onStartCalled, "新しい状態のOnStartが呼ばれていない")
	})

	t.Run("Pop遷移のテスト", func(t *testing.T) {
		t.Parallel()
		ctx := TestContext{Name: "TestContext"}
		initialState := &TestState{name: "InitialState"}
		stateMachine := Init(initialState, ctx)

		// まずPushして2つの状態にする
		pushedState := &TestState{name: "PushedState"}
		stateMachine.lastTransition = Transition[TestContext]{
			Type:          TransPush,
			NewStateFuncs: []StateFactory[TestContext]{func() State[TestContext] { return pushedState }},
		}
		stateMachine.Update(ctx)

		// Pop遷移を実行
		stateMachine.lastTransition = Transition[TestContext]{Type: TransPop}
		stateMachine.Update(ctx)

		// 状態数の確認
		assert.Equal(t, 1, stateMachine.GetStateCount(), "Pop後の状態数が正しくない")

		// 現在の状態の確認
		currentState := stateMachine.GetCurrentState()
		assert.Equal(t, "InitialState", currentState.(*TestState).name, "Pop後の現在状態が正しくない")

		// Popされた状態のOnStopが呼ばれていることを確認
		assert.True(t, pushedState.onStopCalled, "Popされた状態のOnStopが呼ばれていない")
		// 再開された状態のOnResumeが呼ばれていることを確認
		assert.True(t, initialState.onResumeCalled, "再開された状態のOnResumeが呼ばれていない")
	})

	t.Run("Switch遷移のテスト", func(t *testing.T) {
		t.Parallel()
		ctx := TestContext{Name: "TestContext"}
		initialState := &TestState{name: "InitialState"}
		stateMachine := Init(initialState, ctx)

		// Switch遷移を実行
		newState := &TestState{name: "SwitchedState"}
		stateMachine.lastTransition = Transition[TestContext]{
			Type:          TransSwitch,
			NewStateFuncs: []StateFactory[TestContext]{func() State[TestContext] { return newState }},
		}
		stateMachine.Update(ctx)

		// 状態数の確認（変わらず1つ）
		assert.Equal(t, 1, stateMachine.GetStateCount(), "Switch後の状態数が正しくない")

		// 現在の状態の確認
		currentState := stateMachine.GetCurrentState()
		assert.Equal(t, "SwitchedState", currentState.(*TestState).name, "Switch後の現在状態が正しくない")

		// 初期状態のOnStopが呼ばれていることを確認
		assert.True(t, initialState.onStopCalled, "初期状態のOnStopが呼ばれていない")
		// 新しい状態のOnStartが呼ばれていることを確認
		assert.True(t, newState.onStartCalled, "新しい状態のOnStartが呼ばれていない")
	})
}
