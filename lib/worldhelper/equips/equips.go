package equips

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 装備する
func Equip(world w.World, item ecs.Entity, owner ecs.Entity, slotNumber gc.EquipmentSlotNumber) {
	gameComponents := world.Components.Game.(*gc.Components)
	item.AddComponent(gameComponents.Equipped, &gc.Equipped{Owner: owner, EquipmentSlot: slotNumber})
	item.RemoveComponent(gameComponents.InBackpack)
	item.AddComponent(gameComponents.EquipmentChanged, &gc.EquipmentChanged{})
}

// 装備を外す
func Disarm(world w.World, item ecs.Entity) {
	gameComponents := world.Components.Game.(*gc.Components)
	item.AddComponent(gameComponents.InBackpack, &gc.InBackpack{})
	item.RemoveComponent(gameComponents.Equipped)
	item.AddComponent(gameComponents.EquipmentChanged, &gc.EquipmentChanged{})
}

// 指定キャラクターの装備中の防具一覧を取得する
// 必ず長さ4のスライスを返す
func GetWearEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 4)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Equipped,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.Equipped.Get(entity).(*gc.Equipped)
		if owner == equipped.Owner {
			switch equipped.EquipmentSlot {
			case gc.EquipmentSlotZero:
				entities[0] = &entity
			case gc.EquipmentSlotOne:
				entities[1] = &entity
			case gc.EquipmentSlotTwo:
				entities[2] = &entity
			case gc.EquipmentSlotThree:
				entities[3] = &entity
			default:
				log.Fatal("対応していないEquipmentSlot")
			}
		}
	}))

	return entities
}

// 指定キャラクターの装備中のカード一覧を取得する
// 必ず長さ4のスライスを返す
func GetCardEquipments(world w.World, owner ecs.Entity) []*ecs.Entity {
	entities := make([]*ecs.Entity, 4)

	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Item,
		gameComponents.Equipped,
		gameComponents.Card,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		equipped := gameComponents.Equipped.Get(entity).(*gc.Equipped)
		if owner == equipped.Owner {
			switch equipped.EquipmentSlot {
			case gc.EquipmentSlotZero:
				entities[0] = &entity
			case gc.EquipmentSlotOne:
				entities[1] = &entity
			case gc.EquipmentSlotTwo:
				entities[2] = &entity
			case gc.EquipmentSlotThree:
				entities[3] = &entity
			}
		}
	}))

	return entities
}
