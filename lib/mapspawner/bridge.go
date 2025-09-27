package mapspawner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
)

// BuildPlanFromTiles は既存のタイル配列からMapPlanを構築する
// 既存のPlannerChainとの橋渡し用の関数
func BuildPlanFromTiles(planData *mapplaner.PlannerMap) (*mapplaner.MapPlan, error) {
	plan := mapplaner.NewMapPlan(int(planData.Level.TileWidth), int(planData.Level.TileHeight))

	// プレイヤー開始位置を設定（タイル配列ベースの場合は中央付近）
	width := int(planData.Level.TileWidth)
	height := int(planData.Level.TileHeight)
	centerX := width / 2
	centerY := height / 2

	// スポーン可能な位置を探す
	playerX, playerY := centerX, centerY
	found := false

	// 複数の候補位置を試す
	attempts := []struct{ x, y int }{
		{width / 2, height / 2},         // 中央
		{width / 4, height / 4},         // 左上寄り
		{3 * width / 4, height / 4},     // 右上寄り
		{width / 4, 3 * height / 4},     // 左下寄り
		{3 * width / 4, 3 * height / 4}, // 右下寄り
	}

	// 最適な位置を探す
	for _, pos := range attempts {
		tileIdx := planData.Level.XYTileIndex(gc.Tile(pos.x), gc.Tile(pos.y))
		if int(tileIdx) < len(planData.Tiles) && planData.Tiles[tileIdx] == mapplaner.TileFloor {
			playerX, playerY = pos.x, pos.y
			found = true
			break
		}
	}

	// 見つからない場合は全体をスキャン
	if !found {
		for _i, tile := range planData.Tiles {
			if tile == mapplaner.TileFloor {
				i := resources.TileIdx(_i)
				x, y := planData.Level.XYTileCoord(i)
				playerX, playerY = int(x), int(y)
				found = true
				break
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("プレイヤー配置可能な床タイルが見つかりません")
	}

	// プレイヤー位置を設定
	plan.SetPlayerStartPosition(playerX, playerY)

	// タイルを走査してMapPlanを構築
	for _i, tile := range planData.Tiles {
		i := resources.TileIdx(_i)
		x, y := planData.Level.XYTileCoord(i)

		switch tile {
		case mapplaner.TileFloor:
			plan.AddFloor(int(x), int(y))

		case mapplaner.TileWall:
			// 近傍8タイル（直交・斜め）にフロアがあるときだけ壁にする
			if planData.AdjacentAnyFloor(i) {
				// 壁タイプを判定してスプライト番号を決定
				wallType := planData.GetWallType(i)
				spriteNumber := getSpriteNumberForWallType(wallType)
				plan.AddWall(int(x), int(y), spriteNumber)
			}

		case mapplaner.TileWarpNext:
			plan.AddWarpNext(int(x), int(y))

		case mapplaner.TileWarpEscape:
			plan.AddWarpEscape(int(x), int(y))

		case mapplaner.TileEmpty:
			// 空のタイルはエンティティを生成しない
			continue

		default:
			return nil, fmt.Errorf("未知のタイルタイプ: %d", tile)
		}
	}

	return plan, nil
}

// 壁スプライト番号定数（mapplaner/level.goから移動）
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
func getSpriteNumberForWallType(wallType mapplaner.WallType) int {
	switch wallType {
	case mapplaner.WallTypeTop:
		return spriteWallTop // 上壁（下に床がある）
	case mapplaner.WallTypeBottom:
		return spriteWallBottom // 下壁（上に床がある）
	case mapplaner.WallTypeLeft:
		return spriteWallLeft // 左壁（右に床がある）
	case mapplaner.WallTypeRight:
		return spriteWallRight // 右壁（左に床がある）
	case mapplaner.WallTypeTopLeft:
		return spriteWallTopLeft // 左上角（右下に床がある）
	case mapplaner.WallTypeTopRight:
		return spriteWallTopRight // 右上角（左下に床がある）
	case mapplaner.WallTypeBottomLeft:
		return spriteWallBottomLeft // 左下角（右上に床がある）
	case mapplaner.WallTypeBottomRight:
		return spriteWallBottomRight // 右下角（左上に床がある）
	case mapplaner.WallTypeGeneric:
		return spriteWallGeneric // 汎用壁
	default:
		return spriteWallGeneric // 不明な場合は汎用壁
	}
}
