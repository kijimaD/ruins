package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ================

type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

// ================

type Healing struct {
	Amount gc.Amounter
}

func (Healing) isEffectType() {}

// ================

// スタミナ
type RecoveryStamina struct {
	Amount gc.Amounter
}

func (RecoveryStamina) isEffectType() {}

// ================

type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
