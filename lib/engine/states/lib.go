package states

import (
	"log"
	"os"

	w "github.com/kijimaD/ruins/lib/world"

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

// Transition is a state transition
type Transition struct {
	Type      TransType
	NewStates []State
}

// State is a game state
type State interface {
	// Executed when the state begins
	OnStart(world w.World)
	// Executed when the state exits
	OnStop(world w.World)
	// Executed when a new state is pushed over this one
	OnPause(world w.World)
	// Executed when the state become active again (states pushed over this one have been popped)
	OnResume(world w.World)
	// Executed on every frame when the state is active
	Update(world w.World) Transition
	// 描画
	Draw(world w.World, screen *ebiten.Image)
}

// StateMachine contains a stack of states.
// Only the top state is active.
type StateMachine struct {
	states         []State
	lastTransition Transition
}

// 全てのstateのOnResume時に実行される共通処理
// メソッドにする
func hookOnResume(state State, _ w.World) {
	// StateWithTransitionインターフェースを実装している場合は遷移をクリアする
	if stateWithTrans, ok := state.(StateWithTransition); ok {
		stateWithTrans.ClearTransition()
	}
}

// Init creates a new state machine with an initial state
func Init(s State, world w.World) StateMachine {
	s.OnStart(world)
	return StateMachine{[]State{s}, Transition{TransNone, []State{}}}
}

// Update updates the state machine
func (sm *StateMachine) Update(world w.World) {
	switch sm.lastTransition.Type {
	case TransPop:
		sm._Pop(world)
	case TransPush:
		sm._Push(world, sm.lastTransition.NewStates)
	case TransSwitch:
		sm._Switch(world, sm.lastTransition.NewStates)
	case TransReplace:
		sm._Replace(world, sm.lastTransition.NewStates)
	case TransQuit:
		sm._Quit(world)
	}

	if len(sm.states) < 1 {
		os.Exit(0)
	}

	// Run state update function with game systems
	sm.lastTransition = sm.states[len(sm.states)-1].Update(world)

	// Run post-game systems
}

// Draw draws the screen after a state update
func (sm *StateMachine) Draw(world w.World, screen *ebiten.Image) {
	sm.states[len(sm.states)-1].Draw(world, screen)
}

// Remove the active state and resume the next state
func (sm *StateMachine) _Pop(world w.World) {
	sm.states[len(sm.states)-1].OnStop(world)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		// 共通のOnResume処理を実行
		hookOnResume(sm.states[len(sm.states)-1], world)
		// 各stateのOnResumeを実行
		sm.states[len(sm.states)-1].OnResume(world)
	}
}

// Pause the active state and add new states to the stack
func (sm *StateMachine) _Push(world w.World, newStates []State) {
	if len(newStates) > 0 {
		sm.states[len(sm.states)-1].OnPause(world)

		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(world)
			state.OnPause(world)
		}
		newStates[len(newStates)-1].OnStart(world)

		sm.states = append(sm.states, newStates...)
	}
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(world w.World, newStates []State) {
	if len(newStates) != 1 {
		log.Fatal()
	}

	sm.states[len(sm.states)-1].OnStop(world)
	newStates[0].OnStart(world)
	sm.states[len(sm.states)-1] = newStates[0]
}

// Remove all states and insert a new stack
func (sm *StateMachine) _Replace(world w.World, newStates []State) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(world)
			state.OnPause(world)
		}
		newStates[len(newStates)-1].OnStart(world)
	}
	sm.states = newStates
}

// Remove all states and quit
func (sm *StateMachine) _Quit(world w.World) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}
	os.Exit(0)
}

// GetStates は現在の状態スタックを返す（テスト用）
func (sm *StateMachine) GetStates() []State {
	// スライスのコピーを返して、外部からの変更を防ぐ
	result := make([]State, len(sm.states))
	copy(result, sm.states)
	return result
}

// GetCurrentState は現在アクティブな状態を返す（テスト用）
// 状態が存在しない場合はnilを返す
func (sm *StateMachine) GetCurrentState() State {
	if len(sm.states) == 0 {
		return nil
	}
	return sm.states[len(sm.states)-1]
}

// GetStateCount は状態スタックの要素数を返す（テスト用）
func (sm *StateMachine) GetStateCount() int {
	return len(sm.states)
}
