package mapplanner

// 壁スプライト番号定数（mapspawnerに移動済み、互換性のため残している）
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

// GetSpriteNumberForWallType は壁タイプに対応するスプライト番号を返す
// 互換性のためmapbuilderに残している（mapspawnerにも同じ関数がある）
func GetSpriteNumberForWallType(wallType WallType) int {
	switch wallType {
	case WallTypeTop:
		return spriteWallTop // 上壁（下に床がある）
	case WallTypeBottom:
		return spriteWallBottom // 下壁（上に床がある）
	case WallTypeLeft:
		return spriteWallLeft // 左壁（右に床がある）
	case WallTypeRight:
		return spriteWallRight // 右壁（左に床がある）
	case WallTypeTopLeft:
		return spriteWallTopLeft // 左上角（右下に床がある）
	case WallTypeTopRight:
		return spriteWallTopRight // 右上角（左下に床がある）
	case WallTypeBottomLeft:
		return spriteWallBottomLeft // 左下角（右上に床がある）
	case WallTypeBottomRight:
		return spriteWallBottomRight // 右下角（左上に床がある）
	case WallTypeGeneric:
		return spriteWallGeneric // 汎用壁（従来の壁）
	default:
		return 1 // デフォルトは従来の壁
	}
}
