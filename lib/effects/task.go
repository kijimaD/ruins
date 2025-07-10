package effects

import (
	"fmt"
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func InflictDamage(world w.World, damage EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := damage.EffectType.(Damage)
	if ok {
		pools.HP.Current = utils.Max(0, pools.HP.Current-v.Amount)

		name := gameComponents.Name.Get(target).(*gc.Name)
		entry := fmt.Sprintf("%sに%dのダメージ。", name.Name, v.Amount)
		gamelog.BattleLog.Append(entry)

		if pools.HP.Current == 0 {
			gamelog.BattleLog.Append(fmt.Sprintf("%sは倒れた。", name.Name))
		}
	}
}

func HealDamage(world w.World, healing EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(Healing)
	if !ok {
		log.Print("Healingがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.HP.Current = utils.Min(pools.HP.Max, pools.HP.Current+at.Calc(pools.HP.Max))
	case gc.NumeralAmount:
		pools.HP.Current = utils.Min(pools.HP.Max, pools.HP.Current+at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

func ConsumeStamina(world w.World, consume EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := consume.EffectType.(ConsumptionStamina)
	if !ok {
		log.Print("ConsumeStaminaがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.SP.Current = utils.Max(0, pools.SP.Current-at.Calc(pools.SP.Max))
	case gc.NumeralAmount:
		pools.SP.Current = utils.Max(0, pools.SP.Current-at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

func RecoverStamina(world w.World, recover EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := recover.EffectType.(RecoveryStamina)
	if !ok {
		log.Print("RecoverStaminaがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.SP.Current = utils.Min(pools.SP.Max, pools.SP.Current+at.Calc(pools.SP.Max))
	case gc.NumeralAmount:
		pools.SP.Current = utils.Min(pools.SP.Max, pools.SP.Current+at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

func WarpNextTask(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext
}

func WarpEscapeTask(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpEscape
}
