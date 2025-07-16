package effects

import (
	"fmt"
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/mathutil"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// InflictDamage はダメージを与える
func InflictDamage(world w.World, damage EffectSpawner, target ecs.Entity) {
	pools := world.Components.Pools.Get(target).(*gc.Pools)
	v, ok := damage.EffectType.(Damage)
	if ok {
		pools.HP.Current = mathutil.Max(0, pools.HP.Current-v.Amount)

		name := world.Components.Name.Get(target).(*gc.Name)
		entry := fmt.Sprintf("%sに%dのダメージ。", name.Name, v.Amount)
		gamelog.BattleLog.Append(entry)

		if pools.HP.Current == 0 {
			gamelog.BattleLog.Append(fmt.Sprintf("%sは倒れた。", name.Name))
		}
	}
}

// HealDamage はダメージを回復する
func HealDamage(world w.World, healing EffectSpawner, target ecs.Entity) {
	pools := world.Components.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(Healing)
	if !ok {
		log.Print("Healingがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+at.Calc(pools.HP.Max))
	case gc.NumeralAmount:
		pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

// ConsumeStamina はスタミナを消費する
func ConsumeStamina(world w.World, consume EffectSpawner, target ecs.Entity) {
	pools := world.Components.Pools.Get(target).(*gc.Pools)
	v, ok := consume.EffectType.(ConsumptionStamina)
	if !ok {
		log.Print("ConsumeStaminaがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.SP.Current = mathutil.Max(0, pools.SP.Current-at.Calc(pools.SP.Max))
	case gc.NumeralAmount:
		pools.SP.Current = mathutil.Max(0, pools.SP.Current-at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

// RecoverStamina はスタミナを回復する
func RecoverStamina(world w.World, recoveryEffect EffectSpawner, target ecs.Entity) {
	pools := world.Components.Pools.Get(target).(*gc.Pools)
	v, ok := recoveryEffect.EffectType.(RecoveryStamina)
	if !ok {
		log.Print("RecoverStaminaがついてない")
	}
	switch at := v.Amount.(type) {
	case gc.RatioAmount:
		pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+at.Calc(pools.SP.Max))
	case gc.NumeralAmount:
		pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+at.Calc())
	default:
		log.Fatalf("unexpected: %T", at)
	}
}

// WarpNextTask は次のフロアにワープする
func WarpNextTask(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpNext
}

// WarpEscapeTask はゲームから脱出する
func WarpEscapeTask(world w.World) {
	gameResources := world.Resources.Game.(*resources.Game)
	gameResources.StateEvent = resources.StateEventWarpEscape
}
