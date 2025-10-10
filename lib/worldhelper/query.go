package worldhelper

import (
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// QueryOwnedStackable は所持中のスタック可能アイテムを取得する
func QueryOwnedStackable(world w.World, f func(entity ecs.Entity)) {
	world.Manager.Join(
		world.Components.Stackable,
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
