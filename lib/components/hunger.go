package components

// HungerLevel は空腹度の段階を表す
type HungerLevel int

const (
	// HungerSatiated は満腹状態
	HungerSatiated HungerLevel = iota
	// HungerNormal は普通状態
	HungerNormal
	// HungerHungry は空腹状態
	HungerHungry
	// HungerStarving は飢餓状態
	HungerStarving
)

// String はHungerLevelの文字列表現を返す
func (h HungerLevel) String() string {
	switch h {
	case HungerSatiated:
		return "Full"
	case HungerNormal:
		return "Normal"
	case HungerHungry:
		return "Hungry"
	case HungerStarving:
		return "Starving"
	default:
		return "Unknown"
	}
}

// Hunger はプレイヤー専用の空腹度システム
type Hunger struct {
	Value int // 空腹度（0が満腹、値が大きいほど空腹）
}

// GetLevel は現在の空腹度レベルを取得する
func (h *Hunger) GetLevel() HungerLevel {
	switch {
	case h.Value <= 100:
		return HungerSatiated
	case h.Value <= 300:
		return HungerNormal
	case h.Value <= 500:
		return HungerHungry
	default:
		return HungerStarving
	}
}

// Increase は空腹度を増加させる（行動によって腹が減る）
func (h *Hunger) Increase(amount int) {
	h.Value += amount
	if h.Value < 0 {
		h.Value = 0
	}
}

// Decrease は空腹度を減少させる（食事によって満腹になる）
func (h *Hunger) Decrease(amount int) {
	h.Value -= amount
	if h.Value < 0 {
		h.Value = 0
	}
}

// GetStatusPenalty は空腹度によるペナルティを取得する
func (h *Hunger) GetStatusPenalty() int {
	level := h.GetLevel()
	switch level {
	case HungerStarving:
		return -20 // 飢餓状態では大きなペナルティ
	case HungerHungry:
		return -10 // 空腹状態では中程度のペナルティ
	default:
		return 0 // 満腹・普通では影響なし
	}
}

// NewHunger は新しいHungerを作成する
func NewHunger() *Hunger {
	return &Hunger{Value: 0} // 初期状態は満腹
}
