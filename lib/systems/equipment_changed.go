package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 装備変更のダーティフラグが立ったら、ステータス補正まわりを再計算する
// TODO: 最大HP/SPの更新はここでやったほうがよさそう
// TODO: マイナスにならないようにする
func EquipmentChangedSystem(world w.World) bool {
	running := false
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.EquipmentChanged,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		running = true
		entity.RemoveComponent(gameComponents.EquipmentChanged)
	}))

	if !running {
		return false
	}

	// 初期化
	world.Manager.Join(
		gameComponents.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		attrs := gameComponents.Attributes.Get(entity).(*gc.Attributes)

		attrs.Vitality.Modifier = 0
		attrs.Vitality.Total = attrs.Vitality.Base
		attrs.Strength.Modifier = 0
		attrs.Strength.Total = attrs.Strength.Base
		attrs.Sensation.Modifier = 0
		attrs.Sensation.Total = attrs.Sensation.Base
		attrs.Dexterity.Modifier = 0
		attrs.Dexterity.Total = attrs.Dexterity.Base
		attrs.Agility.Modifier = 0
		attrs.Agility.Total = attrs.Agility.Base
		attrs.Defense.Modifier = 0
		attrs.Defense.Total = attrs.Defense.Base
	}))

	world.Manager.Join(
		gameComponents.Equipped,
		gameComponents.Weapon,
	).Visit(ecs.Visit(func(item ecs.Entity) {
		equipped := gameComponents.Equipped.Get(item).(*gc.Equipped)
		weapon := gameComponents.Weapon.Get(item).(*gc.Weapon)

		owner := equipped.Owner
		attrs := gameComponents.Attributes.Get(owner).(*gc.Attributes)

		attrs.Vitality.Modifier += weapon.EquipBonus.Vitality
		attrs.Vitality.Total = attrs.Vitality.Base + attrs.Vitality.Modifier
		attrs.Strength.Modifier += weapon.EquipBonus.Strength
		attrs.Strength.Total = attrs.Strength.Base + attrs.Strength.Modifier
		attrs.Sensation.Modifier += weapon.EquipBonus.Sensation
		attrs.Sensation.Total = attrs.Sensation.Base + attrs.Sensation.Modifier
		attrs.Dexterity.Modifier += weapon.EquipBonus.Dexterity
		attrs.Dexterity.Total = attrs.Dexterity.Base + attrs.Dexterity.Modifier
		attrs.Agility.Modifier += weapon.EquipBonus.Agility
		attrs.Agility.Total = attrs.Agility.Base + attrs.Agility.Modifier
	}))

	world.Manager.Join(
		gameComponents.Equipped,
		gameComponents.Wearable,
	).Visit(ecs.Visit(func(item ecs.Entity) {
		equipped := gameComponents.Equipped.Get(item).(*gc.Equipped)
		wearable := gameComponents.Wearable.Get(item).(*gc.Wearable)

		owner := equipped.Owner
		attrs := gameComponents.Attributes.Get(owner).(*gc.Attributes)

		attrs.Defense.Modifier += wearable.Defense
		attrs.Defense.Total = attrs.Defense.Base + attrs.Defense.Modifier

		attrs.Vitality.Modifier += wearable.EquipBonus.Vitality
		attrs.Vitality.Total = attrs.Vitality.Base + attrs.Vitality.Modifier
		attrs.Strength.Modifier += wearable.EquipBonus.Strength
		attrs.Strength.Total = attrs.Strength.Base + attrs.Strength.Modifier
		attrs.Sensation.Modifier += wearable.EquipBonus.Sensation
		attrs.Sensation.Total = attrs.Sensation.Base + attrs.Sensation.Modifier
		attrs.Dexterity.Modifier += wearable.EquipBonus.Dexterity
		attrs.Dexterity.Total = attrs.Dexterity.Base + attrs.Dexterity.Modifier
		attrs.Agility.Modifier += wearable.EquipBonus.Agility
		attrs.Agility.Total = attrs.Agility.Base + attrs.Agility.Modifier
	}))

	return true
}
