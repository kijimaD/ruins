package mapspawner

import (
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
)

// completeWallSprites はEntityPlan内の壁エンティティのスプライト番号を補完する
func completeWallSprites(plan *mapplanner.EntityPlan) {
	for i := range plan.Entities {
		entity := &plan.Entities[i]
		if entity.EntityType == mapplanner.EntityTypeWall && entity.WallType != nil {
			// WallTypeからスプライト番号を決定
			spriteNumber := getSpriteNumberForWallType(*entity.WallType)
			entity.WallSprite = &spriteNumber
			entity.WallType = nil // スプライト番号が決定されたのでWallTypeをクリア
		}
	}
}

// 壁スプライト番号定数
const (
	spriteWallTop         = 10 // 上壁
	spriteWallBottom      = 11 // 下壁
	spriteWallLeft        = 12 // 左壁
	spriteWallRight       = 13 // 右壁
	spriteWallTopLeft     = 14 // 左上角
	spriteWallTopRight    = 15 // 右上角
	spriteWallBottomLeft  = 16 // 左下角
	spriteWallBottomRight = 17 // 右下角
	spriteWallGeneric     = 1  // 汎用壁
)

// getSpriteNumberForWallType は壁タイプからスプライト番号を取得する
func getSpriteNumberForWallType(wallType mapplanner.WallType) int {
	switch wallType {
	case mapplanner.WallTypeTop:
		return spriteWallTop // 上壁（下に床がある）
	case mapplanner.WallTypeBottom:
		return spriteWallBottom // 下壁（上に床がある）
	case mapplanner.WallTypeLeft:
		return spriteWallLeft // 左壁（右に床がある）
	case mapplanner.WallTypeRight:
		return spriteWallRight // 右壁（左に床がある）
	case mapplanner.WallTypeTopLeft:
		return spriteWallTopLeft // 左上角（右下に床がある）
	case mapplanner.WallTypeTopRight:
		return spriteWallTopRight // 右上角（左下に床がある）
	case mapplanner.WallTypeBottomLeft:
		return spriteWallBottomLeft // 左下角（右上に床がある）
	case mapplanner.WallTypeBottomRight:
		return spriteWallBottomRight // 右下角（左上に床がある）
	case mapplanner.WallTypeGeneric:
		return spriteWallGeneric // 汎用壁
	default:
		return spriteWallGeneric // 不明な場合は汎用壁
	}
}
