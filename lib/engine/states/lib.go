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

// State はゲームステートのジェネリックインターフェース
type State[T any] interface {
	// Executed when the state begins
	OnStart(ctx T)
	// Executed when the state exits
	OnStop(ctx T)
	// Executed when a new state is pushed over this one
	OnPause(ctx T)
	// Executed when the state become active again (states pushed over this one have been popped)
	OnResume(ctx T)
	// Executed on every frame when the state is active
	Update(ctx T) Transition[T]
	// 描画
	Draw(ctx T, screen *ebiten.Image)
}

// StateFactory はステートを作成するファクトリー関数の型
type StateFactory[T any] func() State[T]

// StateMachine はジェネリックな状態スタックを管理する
type StateMachine[T any] struct {
	states         []State[T]
	lastTransition Transition[T]
}

// Init は新しいステートマシンを初期化する
func Init[T any](s State[T], ctx T) StateMachine[T] {
	s.OnStart(ctx)
	return StateMachine[T]{
		states:         []State[T]{s},
		lastTransition: Transition[T]{Type: TransNone},
	}
}

// Update はステートマシンを更新する
func (sm *StateMachine[T]) Update(ctx T) {
	// ファクトリー関数からステートを作成
	states := sm.createStatesFromFunc()

	switch sm.lastTransition.Type {
	case TransPop:
		sm.pop(ctx)
	case TransPush:
		sm.push(ctx, states)
	case TransSwitch:
		sm.switchState(ctx, states)
	case TransReplace:
		sm.replace(ctx, states)
	case TransQuit:
		sm.quit(ctx)
	}

	if len(sm.states) < 1 {
		return
	}

	// アクティブなステートを更新
	sm.lastTransition = sm.states[len(sm.states)-1].Update(ctx)
}

// Draw は画面を描画する
func (sm *StateMachine[T]) Draw(ctx T, screen *ebiten.Image) {
	for _, state := range sm.states {
		state.Draw(ctx, screen)
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
func (sm *StateMachine[T]) pop(ctx T) {
	if len(sm.states) == 0 {
		return
	}

	currentState := sm.states[len(sm.states)-1]
	currentState.OnStop(ctx)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		resumeState := sm.states[len(sm.states)-1]
		resumeState.OnResume(ctx)
	}
}

// push はアクティブなステートを一時停止して新しいステートをスタックに追加する
func (sm *StateMachine[T]) push(ctx T, newStates []State[T]) {
	if len(newStates) == 0 {
		return
	}

	if len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnPause(ctx)
	}

	for _, state := range newStates[:len(newStates)-1] {
		state.OnStart(ctx)
		state.OnPause(ctx)
	}

	activeState := newStates[len(newStates)-1]
	activeState.OnStart(ctx)

	sm.states = append(sm.states, newStates...)
}

// switchState はアクティブなステートを新しいものに置き換える
func (sm *StateMachine[T]) switchState(ctx T, newStates []State[T]) {
	if len(newStates) != 1 {
		return
	}

	if len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		newState := newStates[0]

		currentState.OnStop(ctx)
		newState.OnStart(ctx)
		sm.states[len(sm.states)-1] = newState
	}
}

// replace はすべてのステートを削除して新しいスタックを挿入する
func (sm *StateMachine[T]) replace(ctx T, newStates []State[T]) {
	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnStop(ctx)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(ctx)
			state.OnPause(ctx)
		}
		activeState := newStates[len(newStates)-1]
		activeState.OnStart(ctx)
	}
	sm.states = newStates
}

// quit はすべてのステートを削除して終了する
func (sm *StateMachine[T]) quit(ctx T) {
	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		currentState.OnStop(ctx)
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
