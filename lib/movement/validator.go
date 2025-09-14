// Package movement provides movement validation logic shared between player and AI systems.
//
// このパッケージは移動判定の責務を持つ：
//   - エンティティ衝突チェック
//   - マップ境界チェック
//   - 通行可否判定
package movement

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// CanMoveTo は指定位置に移動可能かチェックする
func CanMoveTo(world w.World, tileX, tileY int, movingEntity ecs.Entity) bool {
	// 基本的な境界チェック
	if tileX < 0 || tileY < 0 || tileX >= consts.MapTileWidth || tileY >= consts.MapTileHeight {
		return false
	}

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

	// TODO: タイルの通行可否チェックを追加
	return canMove
}
