package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
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
	ValueType raw.ValueType
	Amount    int
	Ratio     float64 // 0.0 ~ 1.0
}

func (RecoveryStamina) isEffectType() {}

// ================

type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
