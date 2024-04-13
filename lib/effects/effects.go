package effects

import ecs "github.com/x-hgg-x/goecs/v2"

// ================
type Damage struct {
	Amount int
}

func (Damage) isEffectType() {}

// ================
type Healing struct {
	Amount int
}

func (Healing) isEffectType() {}

// ================
// 全体から割合分を加算して回復する
// 例: 最大HPが100で0.5指定すると、回復量は50
type HealingByRatio struct {
	Amount float64 // 0.0 ~ 1.0
}

func (HealingByRatio) isEffectType() {}

// ================
type ItemUse struct {
	Item ecs.Entity
}

func (ItemUse) isEffectType() {}
