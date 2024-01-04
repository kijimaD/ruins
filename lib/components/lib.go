package components

import ecs "github.com/x-hgg-x/goecs/v2"

type Components struct {
	GridElement     *ecs.SliceComponent
	Player          *ecs.NullComponent
	Wall            *ecs.NullComponent
	Warp            *ecs.NullComponent
	Item            *ecs.NullComponent
	Consumable      *ecs.SliceComponent
	Name            *ecs.SliceComponent
	Description     *ecs.SliceComponent
	InBackpack      *ecs.NullComponent
	InParty         *ecs.NullComponent
	Member          *ecs.NullComponent
	Pools           *ecs.SliceComponent
	ProvidesHealing *ecs.SliceComponent
	InflictsDamage  *ecs.SliceComponent
}

type GridElement struct {
	Line int
	Col  int
}

// フィールドでの移動体
type Player struct{}

// 壁
type Wall struct{}

// ワープパッド
type Warp struct {
	Mode warpMode
}

// アイテム枠に入るもの
type Item struct{}

// 消耗品
type Consumable struct {
	UsableScene UsableSceneType
	Target      Target
}

// 対象
type Target struct {
	TargetFaction TargetFactionType // 対象派閥
	TargetWhole   bool              // 全体対象
}

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

// パーティに参加している
type InParty struct{}

// 冒険に参加できるメンバー
type Member struct{}

// 最大値と現在値を持つようなパラメータ
type Pool struct {
	Max     int
	Current int
}

// メンバーに関連するパラメータ群
type Pools struct {
	HP    Pool
	SP    Pool
	Level int
}

type ProvidesHealing struct {
	Amount int
}

type InflictsDamage struct {
	Amount int
}
