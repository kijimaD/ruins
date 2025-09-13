package worldhelper

import (
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// QueryOwnedMaterial は所持中の素材
// TODO: worldを先に置く
// Join対象の組み合わせが重要であるから、そこだけ関数にすればよいのではないか
func QueryOwnedMaterial(f func(entity ecs.Entity), world w.World) {
	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(f))
}

// QueryPlayer はプレイヤー
func QueryPlayer(world w.World, f func(entity ecs.Entity)) {
	world.Manager.Join(
		world.Components.Player,
		world.Components.FactionAlly,
	).Visit(ecs.Visit(f))
}
