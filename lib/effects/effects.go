package effects

type Damage struct {
	Amount int
}

func (d Damage) isEffectType() {}

type Healing struct {
	Amount int
}

func (d Healing) isEffectType() {}
