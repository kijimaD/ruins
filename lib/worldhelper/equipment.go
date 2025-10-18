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
// 必ず長さ2のスライスを返す（0: 近接武器, 1: 遠距離武器）
func GetWeaponEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 2)

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationEquipped,
		world.Components.Weapon,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
		if owner == equipped.Owner {
			// 武器スロット番号は4番から開始（0-3番は防具用）
			// 4番: 近接武器, 5番: 遠距離武器
			slotIndex := int(equipped.EquipmentSlot) - 4
			if slotIndex >= 0 && slotIndex < len(entities) {
				entities[slotIndex] = &entity
			}
		}
	}))

	return entities
}

// GetMeleeWeaponSlot は近接武器スロット番号を返す
func GetMeleeWeaponSlot() gc.EquipmentSlotNumber {
	return gc.EquipmentSlotNumber(4)
}

// GetRangedWeaponSlot は遠距離武器スロット番号を返す
func GetRangedWeaponSlot() gc.EquipmentSlotNumber {
	return gc.EquipmentSlotNumber(5)
}
