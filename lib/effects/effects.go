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

// ValueTypeで使うフィールドが分岐する
// 数値タイプだと、その数値がそのまま回復量となる
// 例: 50指定すると、回復量は50
// 割合タイプだと、全体からの数値割合分が回復量となる
// 例: 最大HPが100で0.5指定すると、回復量は50
type Healing struct {
	ValueType gc.ValueType
	Amount    int
	Ratio     float64 // 0.0 ~ 1.0
}

func (Healing) isEffectType() {}

// ================

// スタミナ
type RecoveryStamina struct {
	ValueType gc.ValueType
	Amount    int
	Ratio     float64 // 0.0 ~ 1.0
}

func (RecoveryStamina) isEffectType() {}

// ================

type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
