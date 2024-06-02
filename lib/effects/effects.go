package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ================

// ダメージ
type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

// ================

// 体力回復
type Healing struct {
	Amount gc.Amounter
}

func (Healing) isEffectType() {}

// ================

// スタミナ回復
type RecoveryStamina struct {
	Amount gc.Amounter
}

func (RecoveryStamina) isEffectType() {}

// ================

// アイテム使用
type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
