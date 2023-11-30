package components

import ecs "github.com/x-hgg-x/goecs/v2"

type Components struct {
	GridElement *ecs.SliceComponent
	Player      *ecs.NullComponent
	Wall        *ecs.NullComponent
}

type GridElement struct {
	Line int
	Col  int
}

type Player struct{}

type Wall struct{}
