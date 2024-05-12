package loader

import (
	gc "github.com/kijimaD/ruins/lib/components"
	ec "github.com/kijimaD/ruins/lib/engine/components"
)

type GameComponentList struct {
	GridElement      *gc.GridElement
	Player           *gc.Player
	Camera           *gc.Camera
	Wall             *gc.Wall
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
