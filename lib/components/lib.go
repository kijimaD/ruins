package components

import ecs "github.com/x-hgg-x/goecs/v2"

type Components struct {
	GridElement *ecs.SliceComponent
	Player      *ecs.NullComponent
	Wall        *ecs.NullComponent
	Warp        *ecs.NullComponent
	Item        *ecs.NullComponent
	Name        *ecs.SliceComponent
	Description *ecs.SliceComponent
	InBackpack  *ecs.NullComponent
	Consumable  *ecs.NullComponent
}

type GridElement struct {
	Line int
	Col  int
}

type Player struct{}

type Wall struct{}

type Warp struct {
	Mode warpMode
}

// アイテム枠に入るもの
type Item struct{}

// 消耗品
type Consumable struct{}

// 表示名
type Name struct {
	Name string
}

// 説明
type Description struct {
	Description string
}

// 所持品
type InBackpack struct{}
