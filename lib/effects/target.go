package effects

import ecs "github.com/x-hgg-x/goecs/v2"

type Targets interface {
	isTarget()
}

type Party struct{}

func (Party) isTarget() {}

type Single struct {
	Target ecs.Entity
}

func (Single) isTarget() {}

type None struct{}

func (None) isTarget() {}
