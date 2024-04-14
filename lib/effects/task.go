package effects

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/utils/mathutil"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func InflictDamage(world w.World, damage EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := damage.EffectType.(Damage)
	if ok {
		pools.HP.Current = mathutil.Max(0, pools.HP.Current-v.Amount)
	}
}

func HealDamage(world w.World, healing EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(Healing)
	if ok {
		switch a := v.Amount.(type) {
		case gc.RatioAmount:
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+a.Calc(pools.HP.Max))
		case gc.NumeralAmount:
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+a.Calc())
		default:
			log.Fatalf("unexpected: %T", a)
		}
	}
}

func RecoverStamina(world w.World, recover EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := recover.EffectType.(RecoveryStamina)
	if ok {
		switch v.ValueType {
		case raw.PercentageType:
			amount := int(float64(pools.SP.Max) * v.Ratio)
			pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+amount)
		case raw.NumeralType:
			pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+v.Amount)
		}
	}
}
