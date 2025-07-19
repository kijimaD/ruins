package effects

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TargetSelector はターゲット選択戦略のインターフェース
type TargetSelector interface {
	// SelectTargets は指定された条件に基づいてターゲットを選択する
	SelectTargets(world w.World) ([]ecs.Entity, error)

	// String はセレクタの説明を返す
	String() string
}

// TargetSingle は単体ターゲットセレクタ
type TargetSingle struct {
	Entity ecs.Entity
}

// SelectTargets は単体ターゲットを選択する
func (s TargetSingle) SelectTargets(_ w.World) ([]ecs.Entity, error) {
	return []ecs.Entity{s.Entity}, nil
}

func (s TargetSingle) String() string {
	return "TargetSingle"
}

// TargetParty はパーティ全体ターゲットセレクタ
type TargetParty struct{}

// SelectTargets はパーティ全体をターゲットとして選択する
func (p TargetParty) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		targets = append(targets, entity)
	}))
	return targets, nil
}

func (p TargetParty) String() string {
	return "TargetParty"
}

// TargetAllEnemies はすべての敵ターゲットセレクタ
type TargetAllEnemies struct{}

// SelectTargets はすべての敵をターゲットとして選択する
func (a TargetAllEnemies) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionEnemy,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		targets = append(targets, entity)
	}))
	return targets, nil
}

func (a TargetAllEnemies) String() string {
	return "TargetAllEnemies"
}

// TargetAliveParty は生きているパーティメンバーのみをターゲットとする
type TargetAliveParty struct{}

// SelectTargets は生存しているパーティメンバーをターゲットとして選択する
func (a TargetAliveParty) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		poolsComponent := world.Components.Pools.Get(entity)
		if poolsComponent == nil {
			return // Poolsコンポーネントがない場合はスキップ
		}
		pools := poolsComponent.(*gc.Pools)
		if pools.HP.Current > 0 {
			targets = append(targets, entity)
		}
	}))
	return targets, nil
}

func (a TargetAliveParty) String() string {
	return "TargetAliveParty"
}

// TargetDeadParty は死亡しているパーティメンバーのみをターゲットとする
type TargetDeadParty struct{}

// SelectTargets は死亡しているパーティメンバーをターゲットとして選択する
func (d TargetDeadParty) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
		world.Components.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		poolsComponent := world.Components.Pools.Get(entity)
		if poolsComponent == nil {
			return // Poolsコンポーネントがない場合はスキップ
		}
		pools := poolsComponent.(*gc.Pools)
		if pools.HP.Current == 0 {
			targets = append(targets, entity)
		}
	}))
	return targets, nil
}

func (d TargetDeadParty) String() string {
	return "TargetDeadParty"
}

// TargetNone はターゲット不要のエフェクト用セレクタ
type TargetNone struct{}

// SelectTargets はターゲット不要のエフェクト用の空のターゲットリストを返す
func (n TargetNone) SelectTargets(_ w.World) ([]ecs.Entity, error) {
	return []ecs.Entity{}, nil
}

func (n TargetNone) String() string {
	return "TargetNone"
}

// AddTargetedEffect はターゲットセレクタを使用してエフェクトをキューに追加する便利関数
func (p *Processor) AddTargetedEffect(effect Effect, creator *ecs.Entity, selector TargetSelector, world w.World) error {
	return p.AddTargetedEffectWithLogger(effect, creator, selector, nil, world)
}

// AddTargetedEffectWithLogger はターゲットセレクタとLoggerを使用してエフェクトをキューに追加する便利関数
func (p *Processor) AddTargetedEffectWithLogger(effect Effect, creator *ecs.Entity, selector TargetSelector, logger GameLogAppender, world w.World) error {
	targets, err := selector.SelectTargets(world)
	if err != nil {
		return fmt.Errorf("ターゲット選択失敗 %s: %w", selector, err)
	}

	if len(targets) == 0 {
		p.logger.Debug("ターゲットが見つかりませんでした: %s", selector)
	}

	p.AddEffectWithLogger(effect, creator, logger, targets...)
	return nil
}
