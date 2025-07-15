package worldhelper

import (
	"github.com/kijimaD/ruins/lib/utils"
	ecs "github.com/x-hgg-x/goecs/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// GetAmount は所持している素材の数を取得する
func GetAmount(name string, world w.World) int {
	result := 0
	QueryOwnedMaterial(func(entity ecs.Entity) {
		n := world.Components.Game.Name.Get(entity).(*gc.Name)
		if n.Name == name {
			material := world.Components.Game.Material.Get(entity).(*gc.Material)
			result = material.Amount
		}
	}, world)
	return result
}

// PlusAmount は素材の数を増やす
func PlusAmount(name string, amount int, world w.World) {
	changeAmount(name, amount, world)
}

// MinusAmount は素材の数を減らす
func MinusAmount(name string, amount int, world w.World) {
	changeAmount(name, -amount, world)
}

func changeAmount(name string, amount int, world w.World) {
	QueryOwnedMaterial(func(entity ecs.Entity) {
		n := world.Components.Game.Name.Get(entity).(*gc.Name)
		if n.Name == name {
			material := world.Components.Game.Material.Get(entity).(*gc.Material)
			material.Amount = utils.Min(999, utils.Max(0, material.Amount+amount))
		}
	}, world)
}
