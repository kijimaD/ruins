package systems

import (
	"math"
	"math/rand/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/worldhelper/material"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 戦闘終了後に経験値や素材を獲得する
// 獲得した素材名を返す
func BattleDropSystem(world w.World) []string {
	rawMaster := world.Resources.RawMaster.(raw.RawMaster)
	gameComponents := world.Components.Game.(*gc.Components)

	// 素材を獲得する
	cands := []string{}
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Attributes,
		gameComponents.DropTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := gameComponents.Name.Get(entity).(*gc.Name)
		dt := rawMaster.GetDropTable(name.Name)
		for i := 0; i < 3; i++ {
			cands = append(cands, dt.SelectByWeight())
		}
	}))
	rand.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	resultCands := cands[0:3]
	for _, cand := range resultCands {
		material.PlusAmount(cand, 1, world)
	}

	// 経験値を獲得する
	world.Manager.Join(
		gameComponents.Name,
		gameComponents.FactionEnemy,
		gameComponents.Pools,
		gameComponents.DropTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		enemyName := gameComponents.Name.Get(entity).(*gc.Name)
		enemyPools := gameComponents.Pools.Get(entity).(*gc.Pools)
		dt := rawMaster.GetDropTable(enemyName.Name)
		world.Manager.Join(
			gameComponents.Name,
			gameComponents.FactionAlly,
			gameComponents.Pools,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			allyPools := gameComponents.Pools.Get(entity).(*gc.Pools)
			levelDiff := enemyPools.Level - allyPools.Level
			multiplier := calcExpMultiplier(levelDiff)
			allyPools.XP += int(dt.XpBase * multiplier)
		}))

	}))

	return resultCands
}

// 倍率を計算する
// diffが正 -> 敵のほうが強い。倍率が高くなる
// diffが負 -> 味方のほうが強い。倍率が低くなる
func calcExpMultiplier(levelDiff int) float64 {
	expMultiplier := 1.0

	if levelDiff > 0 {
		expMultiplier = math.Pow(1.08, float64(levelDiff))

	} else if levelDiff < 0 {
		expMultiplier = math.Pow(0.9, float64(-levelDiff))
	}

	return expMultiplier
}
