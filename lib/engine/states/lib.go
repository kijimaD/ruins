package states

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TransType is a transition type
type TransType int

const (
	// TransNone does nothing
	TransNone TransType = iota
	// TransPop removes the active state and resume the next state
	TransPop
	// TransPush pauses the active state and add new states to the stack
	TransPush
	// TransSwitch removes the active state and replace it by a new one
	TransSwitch
	// TransReplace removes all states and insert a new stack
	TransReplace
	// TransQuit removes all states and quit
	TransQuit
)

// Transition は状態遷移を表す
type Transition[T any] struct {
	Type          TransType
	NewStateFuncs []StateFactory[T]
}

// TransitionFactory はステート遷移を生成するファクトリー関数の型
type TransitionFactory[T any] func() Transition[T]

// State はゲームステートのジェネリックインターフェース
type State[T any] interface {
	// Executed when the state begins
	OnStart(world T)
	// Executed when the state exits
	OnStop(world T)
	// Executed when a new state is pushed over this one
	OnPause(world T)
	// Executed when the state become active again (states pushed over this one have been popped)
	OnResume(world T)
	// Executed on every frame when the state is active
	Update(world T) Transition[T]
	// 描画
	Draw(world T, screen *ebiten.Image)
}

// StateFactory はステートを作成するファクトリー関数の型
type StateFactory[T any] func() State[T]

// StateMachine はジェネリックな状態スタックを管理する
type StateMachine[T any] struct {
	states         []State[T]
	lastTransition Transition[T]
}

// Init は新しいステートマシンを初期化する
func Init[T any](s State[T], world T) StateMachine[T] {
	s.OnStart(world)
	return StateMachine[T]{
		states:         []State[T]{s},
		lastTransition: Transition[T]{Type: TransNone},
	}
}

// Update はステートマシンを更新する
func (sm *StateMachine[T]) Update(world T) {
	// ファクトリー関数からステートを作成
	states := sm.createStatesFromFunc()

	switch sm.lastTransition.Type {
	case TransPop:
		sm.pop(world)
	case TransPush:
		sm.push(world, states)
	case TransSwitch:
		sm.switchState(world, states)
	case TransReplace:
		sm.replace(world, states)
	case TransQuit:
		sm.quit(world)
	}

	if len(sm.states) < 1 {
		return
	}

	// アクティブなステートを更新
	sm.lastTransition = sm.states[len(sm.states)-1].Update(world)
}

// Draw は画面を描画する
func (sm *StateMachine[T]) Draw(world T, screen *ebiten.Image) {
	for _, state := range sm.states {
		state.Draw(world, screen)
	}
}

// createStatesFromFunc はファクトリー関数からステートを作成する
func (sm *StateMachine[T]) createStatesFromFunc() []State[T] {
	if len(sm.lastTransition.NewStateFuncs) == 0 {
		return []State[T]{}
	}

	states := make([]State[T], len(sm.lastTransition.NewStateFuncs))
	for i, factory := range sm.lastTransition.NewStateFuncs {
		states[i] = factory()
	}
	return states
}

// pop はアクティブなステートを削除して次のステートを再開する
func (sm *StateMachine[T]) pop(world T) {
	if len(sm.states) == 0 {
		return
	}

	currentState := sm.states[len(sm.states)-1]
	currentState.OnStop(world)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		resumeState := sm.states[len(sm.states)-1]
		resumeState.OnResume(world)
	}
}

// push はアクティブなステートを一時停止して新しいステートをスタックに追加する
func (sm *StateMachine[T]) push(world T, newStates []State[T]) {
	if len(newStates) == 0 {
		return
	}

	if len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnPause(world)
	}

	for _, state := range newStates[:len(newStates)-1] {
		state.OnStart(world)
		state.OnPause(world)
	}

	activeState := newStates[len(newStates)-1]
	activeState.OnStart(world)

	sm.states = append(sm.states, newStates...)
}

// switchState はアクティブなステートを新しいものに置き換える
func (sm *StateMachine[T]) switchState(world T, newStates []State[T]) {
	if len(newStates) != 1 {
		return
	}

	if len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		newState := newStates[0]

		currentState.OnStop(world)
		newState.OnStart(world)
		sm.states[len(sm.states)-1] = newState
	}
}

// replace はすべてのステートを削除して新しいスタックを挿入する
func (sm *StateMachine[T]) replace(world T, newStates []State[T]) {
	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(world)
			state.OnPause(world)
		}
		activeState := newStates[len(newStates)-1]
		activeState.OnStart(world)
	}
	sm.states = newStates
}

// quit はすべてのステートを削除して終了する
func (sm *StateMachine[T]) quit(world T) {
	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}
}

// GetStates はステートスタックのコピーを返す（テスト用）
func (sm *StateMachine[T]) GetStates() []State[T] {
	states := make([]State[T], len(sm.states))
	copy(states, sm.states)
	return states
}

// GetCurrentState は現在アクティブなステートを返す（テスト用）
func (sm *StateMachine[T]) GetCurrentState() State[T] {
	if len(sm.states) == 0 {
		return nil
	}
	return sm.states[len(sm.states)-1]
}

// GetStateCount はステートスタックの数を返す（テスト用）
func (sm *StateMachine[T]) GetStateCount() int {
	return len(sm.states)
}
