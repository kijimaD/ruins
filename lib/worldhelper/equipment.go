package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Equip は装備する
func Equip(world w.World, item ecs.Entity, owner ecs.Entity, slotNumber gc.EquipmentSlotNumber) {
	item.AddComponent(world.Components.ItemLocationEquipped, &gc.LocationEquipped{Owner: owner, EquipmentSlot: slotNumber})
	item.RemoveComponent(world.Components.ItemLocationInBackpack)
	item.AddComponent(world.Components.EquipmentChanged, &gc.EquipmentChanged{})
}

// Disarm は装備を外す
func Disarm(world w.World, item ecs.Entity) {
	item.AddComponent(world.Components.ItemLocationInBackpack, &gc.ItemLocationInBackpack)
	item.RemoveComponent(world.Components.ItemLocationEquipped)
	item.AddComponent(world.Components.EquipmentChanged, &gc.EquipmentChanged{})
}

// GetWearEquipments は指定キャラクターの装備中の防具一覧を取得する
// 必ず長さ4のスライスを返す
func GetWearEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 4)

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationEquipped,
		world.Components.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
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

// GetWeaponEquipments は指定キャラクターの装備中の武器一覧を取得する
// 必ず長さ8のスライスを返す
func GetWeaponEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 8)

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationEquipped,
		world.Components.Weapon,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
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
