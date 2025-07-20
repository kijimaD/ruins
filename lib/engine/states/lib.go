package states

import (
	"log"
	"os"
	"reflect"

	"github.com/kijimaD/ruins/lib/logger"
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
	Type          TransType
	NewStateFuncs []StateFactory // ファクトリー関数で動的にステートを作成する。新しいステートは毎回新規作成する
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

// StateFactory はステートを作成するファクトリー関数の型
type StateFactory func() State

// StateMachine contains a stack of states.
// Only the top state is active.
type StateMachine struct {
	states         []State
	lastTransition Transition
	logger         *logger.Logger
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
	l := logger.New(logger.CategoryTransition)
	l.Debug("ステートマシン初期化", "初期ステート", getStateName(s))
	s.OnStart(world)
	return StateMachine{
		states:         []State{s},
		lastTransition: Transition{Type: TransNone},
		logger:         l,
	}
}

// Update updates the state machine
func (sm *StateMachine) Update(world w.World) {
	// ファクトリー関数からステートを作成
	states := sm.createStatesFromFunc()

	switch sm.lastTransition.Type {
	case TransPop:
		sm.logger.Debug("ステート遷移開始", "遷移タイプ", "Pop", "スタック深度", len(sm.states))
		sm._Pop(world)
	case TransPush:
		sm.logger.Debug("ステート遷移開始", "遷移タイプ", "Push", "スタック深度", len(sm.states), "新しいステート数", len(states))
		sm._Push(world, states)
	case TransSwitch:
		sm.logger.Debug("ステート遷移開始", "遷移タイプ", "Switch", "スタック深度", len(sm.states), "新しいステート", getStateName(states[0]))
		sm._Switch(world, states)
	case TransReplace:
		sm.logger.Debug("ステート遷移開始", "遷移タイプ", "Replace", "スタック深度", len(sm.states), "新しいステート数", len(states))
		sm._Replace(world, states)
	case TransQuit:
		sm.logger.Debug("ステート遷移開始", "遷移タイプ", "Quit", "スタック深度", len(sm.states))
		sm._Quit(world)
	}

	if len(sm.states) < 1 {
		os.Exit(0)
	}

	// Run state update function with game systems
	sm.lastTransition = sm.states[len(sm.states)-1].Update(world)

	// Run post-game systems
}

// createStatesFromFunc はStateFactoriesからステートのスライスを作成する
func (sm *StateMachine) createStatesFromFunc() []State {
	if len(sm.lastTransition.NewStateFuncs) == 0 {
		return []State{}
	}

	states := make([]State, len(sm.lastTransition.NewStateFuncs))
	for i, factory := range sm.lastTransition.NewStateFuncs {
		states[i] = factory()
		sm.logger.Debug("ファクトリーからステート作成", "ステート", getStateName(states[i]))
	}
	return states
}

// Draw draws the screen after a state update
func (sm *StateMachine) Draw(world w.World, screen *ebiten.Image) {
	sm.states[len(sm.states)-1].Draw(world, screen)
}

// Remove the active state and resume the next state
func (sm *StateMachine) _Pop(world w.World) {
	currentState := sm.states[len(sm.states)-1]
	sm.logger.Debug("ステート終了", "ステート", getStateName(currentState))
	currentState.OnStop(world)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		resumeState := sm.states[len(sm.states)-1]
		sm.logger.Debug("ステート再開", "ステート", getStateName(resumeState), "新しいスタック深度", len(sm.states))
		// 共通のOnResume処理を実行
		hookOnResume(resumeState, world)
		// 各stateのOnResumeを実行
		resumeState.OnResume(world)
	}
}

// Pause the active state and add new states to the stack
func (sm *StateMachine) _Push(world w.World, newStates []State) {
	if len(newStates) > 0 {
		currentState := sm.states[len(sm.states)-1]
		sm.logger.Debug("ステート一時停止", "ステート", getStateName(currentState))
		currentState.OnPause(world)

		for _, state := range newStates[:len(newStates)-1] {
			sm.logger.Debug("ステート開始（一時停止）", "ステート", getStateName(state))
			state.OnStart(world)
			state.OnPause(world)
		}
		activeState := newStates[len(newStates)-1]
		sm.logger.Debug("ステート開始（アクティブ）", "ステート", getStateName(activeState), "新しいスタック深度", len(sm.states)+len(newStates))

		activeState.OnStart(world)

		sm.states = append(sm.states, newStates...)
	}
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(world w.World, newStates []State) {
	if len(newStates) != 1 {
		sm.logger.Error("Switch遷移でのステート数が不正", "期待値", 1, "実際の値", len(newStates))
		log.Fatal()
	}

	currentState := sm.states[len(sm.states)-1]
	newState := newStates[0]
	sm.logger.Debug("ステート切り替え", "旧ステート", getStateName(currentState), "新ステート", getStateName(newState))

	currentState.OnStop(world)
	newState.OnStart(world)
	sm.states[len(sm.states)-1] = newState
}

// Remove all states and insert a new stack
func (sm *StateMachine) _Replace(world w.World, newStates []State) {
	sm.logger.Debug("全ステート置換開始", "現在のスタック数", len(sm.states), "新しいスタック数", len(newStates))

	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		sm.logger.Debug("ステート終了（置換）", "ステート", getStateName(currentState))
		currentState.OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			sm.logger.Debug("ステート開始（一時停止）", "ステート", getStateName(state))
			state.OnStart(world)
			state.OnPause(world)
		}
		activeState := newStates[len(newStates)-1]
		sm.logger.Debug("ステート開始（アクティブ）", "ステート", getStateName(activeState))
		activeState.OnStart(world)
	}
	sm.states = newStates
}

// Remove all states and quit
func (sm *StateMachine) _Quit(world w.World) {
	sm.logger.Debug("アプリケーション終了開始", "現在のスタック数", len(sm.states))

	for len(sm.states) > 0 {
		currentState := sm.states[len(sm.states)-1]
		sm.logger.Debug("ステート終了（終了）", "ステート", getStateName(currentState))
		currentState.OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	sm.logger.Debug("アプリケーション終了")
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

// getStateName はステートの型名を取得する
func getStateName(state State) string {
	if state == nil {
		return "nil"
	}
	return reflect.TypeOf(state).String()
}
