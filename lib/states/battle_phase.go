package states

import (
	gc "github.com/kijimaD/ruins/lib/components"
	es "github.com/kijimaD/ruins/lib/engine/states"
	"github.com/kijimaD/ruins/lib/gamelog"
	gs "github.com/kijimaD/ruins/lib/systems"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type battlePhase interface {
	// 初期化処理（フェーズ開始時に一度だけ呼ばれる）
	OnInit(st *BattleState, world w.World)
	// 更新処理（毎フレーム呼ばれる）
	OnUpdate(st *BattleState, world w.World) es.Transition
}

// 敵遭遇フェーズ（「敵が現れた」メッセージ表示）
type phaseEnemyEncounter struct{}

func (p *phaseEnemyEncounter) OnInit(st *BattleState, world w.World) {
	// 「敵が現れた」メッセージをログに追加
	gamelog.BattleLog.Append("敵が現れた。")
	// クリック待ち状態にする
	st.isWaitClick = true
}

func (p *phaseEnemyEncounter) OnUpdate(st *BattleState, world w.World) es.Transition {
	// メッセージを表示
	st.reloadMsg(world)
	// エンターキーが押されたら次のフェーズに進む
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.isWaitClick = false
		gamelog.BattleLog.Flush() // メッセージをクリア
		st.phase = &phaseChoosePolicy{}
	}
	return es.Transition{Type: es.TransNone}
}

// 開戦 / 逃走 を選択する
type phaseChoosePolicy struct{}

func (p *phaseChoosePolicy) OnInit(st *BattleState, world w.World) {
	var err error
	st.party, err = worldhelper.NewParty(world, gc.FactionAlly)
	if err != nil {
		// TODO: エラーハンドリング
	}
	st.reloadPolicy(world)
	if err != nil {
		// TODO: エラーハンドリング
	}
}

func (p *phaseChoosePolicy) OnUpdate(st *BattleState, world w.World) es.Transition {
	// 特別な更新処理はない
	return st.ConsumeTransition()
}

// 各キャラの行動選択
type phaseChooseAction struct {
	owner ecs.Entity
}

func (p *phaseChooseAction) OnInit(st *BattleState, world w.World) {
	st.reloadAction(world, p)
}

func (p *phaseChooseAction) OnUpdate(st *BattleState, world w.World) es.Transition {
	// 特別な更新処理はない
	return st.ConsumeTransition()
}

// 行動の対象選択
type phaseChooseTarget struct {
	owner ecs.Entity
	way   ecs.Entity
}

func (p *phaseChooseTarget) OnInit(st *BattleState, world w.World) {
	st.reloadTarget(world, p)
}

func (p *phaseChooseTarget) OnUpdate(st *BattleState, world w.World) es.Transition {
	// 特別な更新処理はない
	return st.ConsumeTransition()
}

type phaseEnemyActionSelect struct{}

func (p *phaseEnemyActionSelect) OnInit(st *BattleState, world w.World) {
	st.handleEnemyActionSelect(world)
}

func (p *phaseEnemyActionSelect) OnUpdate(st *BattleState, world w.World) es.Transition {
	// 特別な更新処理はない
	return st.ConsumeTransition()
}

// 戦闘実行
type phaseExecute struct{}

func (p *phaseExecute) OnInit(st *BattleState, world w.World) {
	// 特別な初期化処理はない
}

func (p *phaseExecute) OnUpdate(st *BattleState, world w.World) es.Transition {
	st.updateEnemyListContainer(world)
	st.reloadMsg(world)
	st.updateMemberContainer(world)

	// 戦闘終了チェック
	result := gs.BattleExtinctionSystem(world)
	switch result {
	case gs.BattleExtinctionNone:
		// 戦闘継続 - コマンド実行処理
		commandCount := st.countBattleCommands(world)
		if commandCount > 0 {
			// 未処理のコマンドがまだ残っている
			st.isWaitClick = true
			if st.keyboardInput.IsEnterJustPressedOnce() {
				gs.BattleCommandSystem(world)
				st.isWaitClick = false
			}
			return es.Transition{Type: es.TransNone}
		}
		// 処理完了
		if st.keyboardInput.IsEnterJustPressedOnce() {
			st.phase = &phaseChoosePolicy{}
			st.isWaitClick = false
			gamelog.BattleLog.Flush()
		}
		return es.Transition{Type: es.TransNone}
	case gs.BattleExtinctionAlly:
		gamelog.BattleLog.Append("全滅した。")
		st.phase = &phaseGameOver{}
		return es.Transition{Type: es.TransNone}
	case gs.BattleExtinctionMonster:
		gamelog.BattleLog.Append("敵を全滅させた。")
		st.phase = &phaseResult{}
		return es.Transition{Type: es.TransNone}
	default:
		return es.Transition{Type: es.TransNone}
	}
}

// リザルト画面
type phaseResult struct {
	actionCount int
}

func (p *phaseResult) OnInit(st *BattleState, world w.World) {
	// 特別な初期化処理はない
}

func (p *phaseResult) OnUpdate(st *BattleState, world w.World) es.Transition {
	st.reloadMsg(world)
	if st.keyboardInput.IsEnterJustPressedOnce() {
		switch p.actionCount {
		case 0:
			// ドロップ処理
			dropResult := gs.BattleDropSystem(world)
			st.resultWindow = st.initResultWindow(world, dropResult)
			st.ui.AddWindow(st.resultWindow)  // ウィンドウをUIに追加
			st.isWaitClick = true
			p.actionCount++
		case 1:
			st.isWaitClick = false
			gamelog.BattleLog.Flush() // メッセージをクリア
			return es.Transition{Type: es.TransPop}
		}
	}
	return es.Transition{Type: es.TransNone}
}

type phaseGameOver struct{}

func (p *phaseGameOver) OnInit(st *BattleState, world w.World) {
	// 特別な初期化処理はない
}

func (p *phaseGameOver) OnUpdate(st *BattleState, world w.World) es.Transition {
	st.reloadMsg(world)
	st.isWaitClick = true
	if st.keyboardInput.IsEnterJustPressedOnce() {
		st.isWaitClick = false
		gamelog.BattleLog.Flush() // メッセージをクリア
		st.phase = &phaseChoosePolicy{}
	}
	return es.Transition{Type: es.TransNone}
}
