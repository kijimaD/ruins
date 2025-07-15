package systems

import (
	"math"
	"math/rand/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

const (
	// LevelUpThreshold はレベルアップに必要な経験値の闾値
	LevelUpThreshold = 100
)

// DropResult は UI用に渡す実行結果
type DropResult struct {
	// 獲得した素材名
	MaterialNames []string
	// 獲得前の経験値
	XPBefore map[ecs.Entity]int
	// 獲得後の経験値
	XPAfter map[ecs.Entity]int
	// レベルアップしたかどうか
	IsLevelUp map[ecs.Entity]bool
}

// BattleDropSystem は戦闘終了後に経験値や素材を獲得する
// 獲得した素材名を返す
func BattleDropSystem(world w.World) DropResult {
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	result := DropResult{
		MaterialNames: []string{},
		XPBefore:      map[ecs.Entity]int{},
		XPAfter:       map[ecs.Entity]int{},
		IsLevelUp:     map[ecs.Entity]bool{},
	}

	// 素材を獲得する
	cands := []string{}
	world.Manager.Join(
		world.Components.Game.Name,
		world.Components.Game.FactionEnemy,
		world.Components.Game.Attributes,
		world.Components.Game.DropTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		name := world.Components.Game.Name.Get(entity).(*gc.Name)
		dt := rawMaster.GetDropTable(name.Name)
		for i := 0; i < 3; i++ {
			cands = append(cands, dt.SelectByWeight())
		}
	}))
	rand.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	result.MaterialNames = cands[0:3]
	for _, cand := range result.MaterialNames {
		worldhelper.PlusAmount(cand, 1, world)
	}
	result.XPBefore = getMemberXP(world)

	// 経験値を獲得する
	world.Manager.Join(
		world.Components.Game.Name,
		world.Components.Game.FactionEnemy,
		world.Components.Game.Pools,
		world.Components.Game.DropTable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		enemyName := world.Components.Game.Name.Get(entity).(*gc.Name)
		enemyPools := world.Components.Game.Pools.Get(entity).(*gc.Pools)
		dt := rawMaster.GetDropTable(enemyName.Name)
		world.Manager.Join(
			world.Components.Game.Name,
			world.Components.Game.FactionAlly,
			world.Components.Game.Pools,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			allyPools := world.Components.Game.Pools.Get(entity).(*gc.Pools)
			levelDiff := enemyPools.Level - allyPools.Level
			multiplier := calcExpMultiplier(levelDiff)
			allyPools.XP += int(dt.XpBase * multiplier)
		}))
	}))
	result.XPAfter = getMemberXP(world)

	// 経験値を見てレベルを上げる
	world.Manager.Join(
		world.Components.Game.Name,
		world.Components.Game.FactionAlly,
		world.Components.Game.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pools := world.Components.Game.Pools.Get(entity).(*gc.Pools)
		if pools.XP >= LevelUpThreshold {
			result.IsLevelUp[entity] = true

			pools.Level++
			pools.XP = 0
		}
	}))

	return result
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

// メンバーごとの経験値を取得する
func getMemberXP(world w.World) map[ecs.Entity]int {
	xpMap := map[ecs.Entity]int{}

	world.Manager.Join(
		world.Components.Game.Name,
		world.Components.Game.FactionAlly,
		world.Components.Game.Pools,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pools := world.Components.Game.Pools.Get(entity).(*gc.Pools)
		xpMap[entity] = pools.XP
	}))

	return xpMap
}
