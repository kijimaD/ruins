package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type BattleExtinctionType int

const (
	BattleExtinctionNone BattleExtinctionType = iota
	BattleExtinctionAlly
	BattleExtinctionMonster
)

// 敵や味方の全滅をチェックする
func BattleExtinctionSystem(world w.World) BattleExtinctionType {
	gameComponents := world.Components.Game.(*gc.Components)

	// 味方が全員死んでいたらゲームオーバーにする
	liveAllyCount := 0
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionAlly,
		gameComponents.Attributes,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		if pools.HP.Current == 0 {
			return
		}
		liveAllyCount += 1
	}))
	if liveAllyCount == 0 {
		return BattleExtinctionAlly
	}

	// 敵が全員死んでいたらリザルトフェーズに遷移する
	liveEnemyCount := 0
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Attributes,
		gameComponents.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pools := gameComponents.Pools.Get(entity).(*gc.Pools)
		if pools.HP.Current == 0 {
			return
		}
		liveEnemyCount += 1
	}))
	if liveEnemyCount == 0 {
		return BattleExtinctionMonster
	}

	return BattleExtinctionNone
}
