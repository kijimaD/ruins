package mapspawner

import (
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
)

// getSpriteKeyForWallType は壁タイプからスプライトキーを取得する
func getSpriteKeyForWallType(wallType mapplanner.WallType) string {
	switch wallType {
	case mapplanner.WallTypeTop:
		return "wall_top" // 上壁（下に床がある）
	case mapplanner.WallTypeBottom:
		return "wall_bottom" // 下壁（上に床がある）
	case mapplanner.WallTypeLeft:
		return "wall_left" // 左壁（右に床がある）
	case mapplanner.WallTypeRight:
		return "wall_right" // 右壁（左に床がある）
	case mapplanner.WallTypeTopLeft:
		return "wall_corner_tl" // 左上角（右下に床がある）
	case mapplanner.WallTypeTopRight:
		return "wall_corner_tr" // 右上角（左下に床がある）
	case mapplanner.WallTypeBottomLeft:
		return "wall_corner_bl" // 左下角（右上に床がある）
	case mapplanner.WallTypeBottomRight:
		return "wall_corner_br" // 右下角（左上に床がある）
	case mapplanner.WallTypeGeneric:
		return "wall_generic" // 汎用壁
	default:
		return "wall_generic" // 不明な場合は汎用壁
	}
}
