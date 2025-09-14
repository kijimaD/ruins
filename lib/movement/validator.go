// Package movement provides movement validation logic shared between player and AI systems.
//
// このパッケージは移動判定の責務を持つ：
//   - エンティティ衝突チェック
//   - マップ境界チェック
//   - 通行可否判定
//
// 循環importを避けるため、systemsパッケージとai_inputパッケージの両方から使用される
package movement

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// CanMoveTo は指定位置に移動可能かチェックする
func CanMoveTo(world w.World, tileX, tileY int, movingEntity ecs.Entity) bool {
	// 他のエンティティとの衝突チェック
	canMove := true
	world.Manager.Join(
		world.Components.GridElement,
		world.Components.BlockPass,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 自分自身は除外
		if entity == movingEntity {
			return
		}

		gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
		if int(gridElement.X) == tileX && int(gridElement.Y) == tileY {
			canMove = false
		}
	}))

	// TODO: マップの境界チェックやタイルの通行可否チェックを追加
	return canMove
}
