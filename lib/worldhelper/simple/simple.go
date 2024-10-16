package simple

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 所持中の素材
// TODO: worldを先に置く
// Join対象の組み合わせが重要であるから、そこだけ関数にすればよいのではないか
func OwnedMaterial(f func(entity ecs.Entity), world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.Material,
		gameComponents.ItemLocationInBackpack,
	).Visit(ecs.Visit(f))
}

// パーティメンバー
func InPartyMember(world w.World, f func(entity ecs.Entity)) {
	gameComponents := world.Components.Game.(*gc.Components)
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(f))
}
