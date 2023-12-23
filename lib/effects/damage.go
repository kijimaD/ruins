package effects

import (
	gc "github.com/kijimaD/sokotwo/lib/components"
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func InflictDamage(world w.World, damage EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := damage.EffectType.(Damage)
	if ok {
		pools.HP.Current = max(0, pools.HP.Current-v.Amount)
	}
}

func HealDamage(world w.World, healing EffectSpawner, target ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	pools := gameComponents.Pools.Get(target).(*gc.Pools)
	v, ok := healing.EffectType.(Healing)
	if ok {
		pools.HP.Current = min(pools.HP.Max, pools.HP.Current+v.Amount)
	}
}