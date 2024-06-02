package effects

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ================

// ダメージを与える
type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

// ================

// 体力を回復する
type Healing struct {
	Amount gc.Amounter
}

func (Healing) isEffectType() {}

// ================

// スタミナを回復する
type RecoveryStamina struct {
	Amount gc.Amounter
}

func (RecoveryStamina) isEffectType() {}

// ================

// アイテムを使用する
type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}

// ================

// 次階層に移動する
type WarpNext struct{}

func (WarpNext) isEffectType() {}

// ================

// 脱出
type WarpEscape struct{}

func (WarpEscape) isEffectType() {}
