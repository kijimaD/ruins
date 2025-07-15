package effects

import ecs "github.com/x-hgg-x/goecs/v2"

// Targets はターゲットのインターフェース
type Targets interface {
	isTarget()
}

// Party はパーティターゲット
type Party struct{}

func (Party) isTarget() {}

// Single は単体ターゲット
type Single struct {
	Target ecs.Entity
}

func (Single) isTarget() {}

// None はターゲットなし
type None struct{}

func (None) isTarget() {}
