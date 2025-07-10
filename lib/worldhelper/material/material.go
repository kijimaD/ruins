package material

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/kijimaD/ruins/lib/worldhelper/simple"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 所持している素材の数を取得する
func GetAmount(name string, world w.World) int {
	result := 0
	gameComponents := world.Components.Game.(*gc.Components)
	simple.OwnedMaterial(func(entity ecs.Entity) {
		n := gameComponents.Name.Get(entity).(*gc.Name)
		if n.Name == name {
			material := gameComponents.Material.Get(entity).(*gc.Material)
			result = material.Amount
		}
	}, world)
	return result
}

func PlusAmount(name string, amount int, world w.World) {
	changeAmount(name, amount, world)
}

func MinusAmount(name string, amount int, world w.World) {
	changeAmount(name, -amount, world)
}

func changeAmount(name string, amount int, world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	simple.OwnedMaterial(func(entity ecs.Entity) {
		n := gameComponents.Name.Get(entity).(*gc.Name)
		if n.Name == name {
			material := gameComponents.Material.Get(entity).(*gc.Material)
			material.Amount = utils.Min(999, utils.Max(0, material.Amount+amount))
		}
	}, world)
}
