package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// EquipmentChangedSystem は装備変更のダーティフラグが立ったら、ステータス補正まわりを再計算する
// TODO: 最大HP/SPの更新はここでやったほうがよさそう
// TODO: マイナスにならないようにする
func EquipmentChangedSystem(world w.World) bool {
	running := false
	world.Manager.Join(
		world.Components.EquipmentChanged,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		running = true
		entity.RemoveComponent(world.Components.EquipmentChanged)
	}))

	if !running {
		return false
	}

	// 初期化
	world.Manager.Join(
		world.Components.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)

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
		world.Components.ItemLocationEquipped,
		world.Components.Wearable,
	).Visit(ecs.Visit(func(item ecs.Entity) {
		equipped := world.Components.ItemLocationEquipped.Get(item).(*gc.LocationEquipped)
		wearable := world.Components.Wearable.Get(item).(*gc.Wearable)

		owner := equipped.Owner
		attrs := world.Components.Attributes.Get(owner).(*gc.Attributes)

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

	world.Manager.Join(
		world.Components.Pools,
		world.Components.Attributes,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		pools := world.Components.Pools.Get(entity).(*gc.Pools)
		attrs := world.Components.Attributes.Get(entity).(*gc.Attributes)

		pools.HP.Max = maxHP(attrs, pools)
		pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current)
		pools.SP.Max = maxSP(attrs, pools)
		pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current)
	}))

	return true
}

// 30+(体力*8+力+感覚)
func maxHP(attrs *gc.Attributes, pools *gc.Pools) int {
	return 30 + attrs.Vitality.Total*8 + attrs.Strength.Total + attrs.Sensation.Total
}

// 体力*2+器用さ+素早さ
func maxSP(attrs *gc.Attributes, pools *gc.Pools) int {
	return attrs.Vitality.Total*2 + attrs.Dexterity.Total + attrs.Agility.Total
}
