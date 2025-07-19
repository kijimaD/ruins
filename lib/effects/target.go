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

// SingleTarget は単体ターゲットセレクタ
type SingleTarget struct {
	Entity ecs.Entity
}

func (s SingleTarget) SelectTargets(world w.World) ([]ecs.Entity, error) {
	return []ecs.Entity{s.Entity}, nil
}

func (s SingleTarget) String() string {
	return "SingleTarget"
}

// PartyTargets はパーティ全体ターゲットセレクタ
type PartyTargets struct{}

func (p PartyTargets) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		targets = append(targets, entity)
	}))
	return targets, nil
}

func (p PartyTargets) String() string {
	return "PartyTargets"
}

// AllEnemies はすべての敵ターゲットセレクタ
type AllEnemies struct{}

func (a AllEnemies) SelectTargets(world w.World) ([]ecs.Entity, error) {
	var targets []ecs.Entity
	world.Manager.Join(
		world.Components.FactionEnemy,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		targets = append(targets, entity)
	}))
	return targets, nil
}

func (a AllEnemies) String() string {
	return "AllEnemies"
}

// AlivePartyMembers は生きているパーティメンバーのみをターゲットとする
type AlivePartyMembers struct{}

func (a AlivePartyMembers) SelectTargets(world w.World) ([]ecs.Entity, error) {
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

func (a AlivePartyMembers) String() string {
	return "AlivePartyMembers"
}

// DeadPartyMembers は死亡しているパーティメンバーのみをターゲットとする
type DeadPartyMembers struct{}

func (d DeadPartyMembers) SelectTargets(world w.World) ([]ecs.Entity, error) {
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

func (d DeadPartyMembers) String() string {
	return "DeadPartyMembers"
}

// NoTarget はターゲット不要のエフェクト用セレクタ
type NoTarget struct{}

func (n NoTarget) SelectTargets(world w.World) ([]ecs.Entity, error) {
	return []ecs.Entity{}, nil
}

func (n NoTarget) String() string {
	return "NoTarget"
}

// AddTargetedEffect はターゲットセレクタを使用してエフェクトをキューに追加する便利関数
func (p *Processor) AddTargetedEffect(effect Effect, creator *ecs.Entity, selector TargetSelector, world w.World) error {
	targets, err := selector.SelectTargets(world)
	if err != nil {
		return fmt.Errorf("ターゲット選択失敗 %s: %w", selector, err)
	}

	if len(targets) == 0 {
		p.logger.Debug("ターゲットが見つかりませんでした: %s", selector)
	}

	return p.AddEffect(effect, creator, targets...)
}
