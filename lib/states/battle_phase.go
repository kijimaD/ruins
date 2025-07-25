package states

import (
	ecs "github.com/x-hgg-x/goecs/v2"
)

type battlePhase interface {
	isBattlePhase()
}

// 敵遭遇フェーズ（「敵が現れた」メッセージ表示）
type phaseEnemyEncounter struct{}

func (p *phaseEnemyEncounter) isBattlePhase() {}

// 開戦 / 逃走 を選択する
type phaseChoosePolicy struct{}

func (p *phaseChoosePolicy) isBattlePhase() {}

// 各キャラの行動選択
type phaseChooseAction struct {
	owner ecs.Entity
}

func (p *phaseChooseAction) isBattlePhase() {}

// 行動の対象選択
type phaseChooseTarget struct {
	owner ecs.Entity
	way   ecs.Entity
}

func (p *phaseChooseTarget) isBattlePhase() {}

type phaseEnemyActionSelect struct{}

func (p *phaseEnemyActionSelect) isBattlePhase() {}

// 戦闘実行
type phaseExecute struct{}

func (p *phaseExecute) isBattlePhase() {}

// リザルト画面
type phaseResult struct {
	actionCount int
}

func (p *phaseResult) isBattlePhase() {}

type phaseGameOver struct{}

func (p *phaseGameOver) isBattlePhase() {}
