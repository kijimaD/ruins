package effects

import ecs "github.com/x-hgg-x/goecs/v2"

type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

type Healing struct {
	Amount int
}

func (Healing) isEffectType() {}

type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
