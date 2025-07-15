package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ================

// Damage はダメージを与える
type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

// ================

// Healing は体力を回復する
type Healing struct {
	Amount gc.Amounter
}

func (Healing) isEffectType() {}

// ================

// ConsumptionStamina はスタミナを消費する
type ConsumptionStamina struct {
	Amount gc.Amounter
}

func (ConsumptionStamina) isEffectType() {}

// ================

// RecoveryStamina はスタミナを回復する
type RecoveryStamina struct {
	Amount gc.Amounter
}

func (RecoveryStamina) isEffectType() {}

// ================

// ItemUse はアイテムを使用する
type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}

// ================

// WarpNext は次階層に移動する
type WarpNext struct{}

func (WarpNext) isEffectType() {}

// ================

// WarpEscape は脱出
type WarpEscape struct{}

func (WarpEscape) isEffectType() {}
