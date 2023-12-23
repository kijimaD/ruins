package effects

import (
	w "github.com/kijimaD/sokotwo/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var EffectQueue []EffectSpawner

type EffectType interface {
	isEffectType()
}

type EffectSpawner struct {
	Creator    *ecs.Entity
	EffectType EffectType
	Targets    Targets
}

func AddEffect(creator *ecs.Entity, effectType EffectType, targets Targets) {
	EffectQueue = append(EffectQueue, EffectSpawner{
		Creator:    creator,
		EffectType: effectType,
		Targets:    targets,
	})
}

// キューに貯められたEffectSpawnerを処理する
func RunEffectQueue(world w.World) {
	for {
		if len(EffectQueue) > 0 {
			effect := EffectQueue[0]
			EffectQueue = EffectQueue[1:]
			TargetApplicator(world, effect)
		} else {
			break
		}
	}
}

// EffectSpawnerからEffectを生成する
func TargetApplicator(world w.World, es EffectSpawner) {
	switch e := es.EffectType.(type) {
	case Damage:
		v, ok := es.Targets.(Single)
		if ok {
			AffectEntity(world, es, v.Target)
		}
	case Healing:
		v, ok := es.Targets.(Single)
		if ok {
			AffectEntity(world, es, v.Target)
		}
	case ItemUse:
		_, ok := es.Targets.(Single)
		if ok {
			ItemTrigger(nil, e.Item, es.Targets, world)
		}

	}
}

func AffectEntity(world w.World, es EffectSpawner, target ecs.Entity) {
	switch es.EffectType.(type) {
	case Damage:
		InflictDamage(world, es, target)
	case Healing:
		HealDamage(world, es, target)
	}
}
