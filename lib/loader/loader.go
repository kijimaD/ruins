package loader

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// TODO: 移動する
type GameComponentList struct {
	GridElement      *gc.GridElement
	Player           *gc.Player
	Camera           *gc.Camera
	Warp             *gc.Warp
	Item             *gc.Item
	Name             *gc.Name
	Description      *gc.Description
	InBackpack       *gc.InBackpack
	Equipped         *gc.Equipped
	Consumable       *gc.Consumable
	InParty          *gc.InParty
	Member           *gc.Member
	Pools            *gc.Pools
	ProvidesHealing  *gc.ProvidesHealing
	InflictsDamage   *gc.InflictsDamage
	Attack           *gc.Attack
	Material         *gc.Material
	Recipe           *gc.Recipe
	Wearable         *gc.Wearable
	Attributes       *gc.Attributes
	EquipmentChanged *gc.EquipmentChanged
	Card             *gc.Card

	Position     *gc.Position
	SpriteRender *ec.SpriteRender
	BlockView    *gc.BlockView
	BlockPass    *gc.BlockPass
}

type Entity struct {
	Components GameComponentList
}

// 現在の階層
type Level struct {
	// 横のタイル数
	TileWidth gc.Row
	// 縦のタイル数
	TileHeight gc.Col
	// 1タイルあたりのピクセル数。タイルは正方形のため、縦横で同じピクセル数になる
	TileSize int
	// タイルエンティティ群
	Entities []ecs.Entity
}

// タイル座標から、タイルスライスのインデックスを求める
func (l *Level) XYTileIndex(tx int, ty int) int {
	return ty*int(l.TileWidth) + tx
}

// タイルスライスのインデックスからタイル座標を求める
func (l *Level) XYTileCoord(idx int) (int, int) {
	x := idx % int(l.TileWidth)
	y := idx / int(l.TileWidth)
	return x, y
}

// xy座標から、該当するエンティティを求める
func (l *Level) AtEntity(x int, y int) ecs.Entity {
	tx := x / l.TileSize
	ty := y / l.TileSize
	idx := l.XYTileIndex(tx, ty)

	return l.Entities[idx]
}

// ステージ幅。横の全体ピクセル数
func (l *Level) Width() int {
	return int(l.TileWidth) * l.TileSize
}

// ステージ縦。縦の全体ピクセル数
func (l *Level) Height() int {
	return int(l.TileHeight) * l.TileSize
}
