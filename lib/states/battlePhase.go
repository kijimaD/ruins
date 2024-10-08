package states

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type battlePhase interface {
	isBattlePhase()
}

// 開戦 / 逃走
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
	way   gc.Card
}

func (p *phaseChooseTarget) isBattlePhase() {}

// 戦闘実行
type phaseExecute struct{}

func (p *phaseExecute) isBattlePhase() {}

// リザルト画面
type phaseResult struct{}

func (p *phaseResult) isBattlePhase() {}
