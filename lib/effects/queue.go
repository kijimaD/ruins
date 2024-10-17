package effects

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"
)

var EffectQueue []EffectSpawner

type EffectType interface {
	isEffectType()
}

// queueの中身
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

// 単数or複数Targetを処理する。最終的にAffectEntityが呼ばれるのは同じ
func TargetApplicator(world w.World, es EffectSpawner) {
	switch e := es.EffectType.(type) {
	case Damage, Healing, ConsumptionStamina, RecoveryStamina:
		v, ok := es.Targets.(Single)
		if ok {
			AffectEntity(world, es, utils.GetPtr(v.Target))
		}
		_, ok = es.Targets.(Party)
		if ok {
			gameComponents := world.Components.Game.(*gc.Components)
			world.Manager.Join(
				gameComponents.FactionAlly,
				gameComponents.InParty,
			).Visit(ecs.Visit(func(entity ecs.Entity) {
				AffectEntity(world, es, utils.GetPtr(entity))
			}))
		}
	case WarpNext, WarpEscape:
		_, ok := es.Targets.(None)
		if !ok {
			log.Fatal("Warp EffectのTargetはNoneである必要がある")
		}
		AffectEntity(world, es, nil)
	case ItemUse:
		_, ok := es.Targets.(Single)
		if ok {
			// アイテムは複数のComponents->Effectに分解されてキューに追加される
			ItemTrigger(nil, e.Item, es.Targets, world)
		}
	default:
		log.Fatalf("TargetApplicator, 対応してないEffectType: %T", e)
	}
}

func AffectEntity(world w.World, es EffectSpawner, target *ecs.Entity) {
	switch e := es.EffectType.(type) {
	case Damage:
		InflictDamage(world, es, *target)
	case Healing:
		HealDamage(world, es, *target)
	case ConsumptionStamina:
		ConsumeStamina(world, es, *target)
	case RecoveryStamina:
		RecoverStamina(world, es, *target)
	case WarpNext:
		WarpNextTask(world)
	case WarpEscape:
		WarpEscapeTask(world)
	default:
		log.Fatalf("AffectEntity, 対応してないEffectType: %T", e)
	}
}
