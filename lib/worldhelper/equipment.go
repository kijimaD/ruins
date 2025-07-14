package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 装備する
func Equip(world w.World, item ecs.Entity, owner ecs.Entity, slotNumber gc.EquipmentSlotNumber) {
	gameComponents := world.Components.Game.(*gc.Components)
	item.AddComponent(gameComponents.ItemLocationEquipped, &gc.LocationEquipped{Owner: owner, EquipmentSlot: slotNumber})
	item.RemoveComponent(gameComponents.ItemLocationInBackpack)
	item.AddComponent(gameComponents.EquipmentChanged, &gc.EquipmentChanged{})
}

// 装備を外す
func Disarm(world w.World, item ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	item.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	item.RemoveComponent(gameComponents.ItemLocationEquipped)
	item.AddComponent(gameComponents.EquipmentChanged, &gc.EquipmentChanged{})
}

// 指定キャラクターの装備中の防具一覧を取得する
// 必ず長さ4のスライスを返す
func GetWearEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 4)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationEquipped,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
		if owner == equipped.Owner {
			for i := range entities {
				if equipped.EquipmentSlot != gc.EquipmentSlotNumber(i) {
					continue
				}
				entities[i] = &entity
			}
		}
	}))

	return entities
}

// 指定キャラクターの装備中のカード一覧を取得する
// 必ず長さ8のスライスを返す
func GetCardEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 8)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.ItemLocationEquipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
		if owner == equipped.Owner {
			for i := range entities {
				if equipped.EquipmentSlot != gc.EquipmentSlotNumber(i) {
					continue
				}
				entities[i] = &entity
			}
		}
	}))

	return entities
}
