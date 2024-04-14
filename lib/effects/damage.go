package effects

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
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
		pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+v.Amount)
	}
}

func HealDamageByRatio(world w.World, healing EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(HealingByRatio)
	if ok {
		amount := int(float64(pools.HP.Max) * v.Amount)
		pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+amount)
	} else {
		log.Fatalf("expect: HealingByRatio, actual: %T", v)
	}
}

func RecoverStaminaByRatio(world w.World, healing EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(RecoveryStaminaByRatio)
	if ok {
		amount := int(float64(pools.SP.Max) * v.Amount)
		pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+amount)
	} else {
		log.Fatalf("expect: RecoverStaminaByRatio, actual: %T", v)
	}
}
