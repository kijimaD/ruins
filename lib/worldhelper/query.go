package worldhelper

import (
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// QueryOwnedMaterial は所持中の素材
// TODO: worldを先に置く
// Join対象の組み合わせが重要であるから、そこだけ関数にすればよいのではないか
func QueryOwnedMaterial(f func(entity ecs.Entity), world w.World) {
	gameComponents := world.Components.Game
	world.Manager.Join(
		gameComponents.Material,
		gameComponents.ItemLocationInBackpack,
	).Visit(ecs.Visit(f))
}

// QueryInPartyMember はパーティメンバー
func QueryInPartyMember(world w.World, f func(entity ecs.Entity)) {
	gameComponents := world.Components.Game
	world.Manager.Join(
		gameComponents.FactionAlly,
		gameComponents.InParty,
	).Visit(ecs.Visit(f))
}
