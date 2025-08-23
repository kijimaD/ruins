package mapbuilder

import (
	"log"

	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"

	gc "github.com/kijimaD/ruins/lib/components"
)

// NewLevel は新規に階層を生成する。
// 階層を初期化するので、具体的なコードであり、その分参照を多く含んでいる。循環参照を防ぐためにこの関数はLevel構造体とは同じpackageに属していない。
func NewLevel(world w.World, width gc.Row, height gc.Col, seed uint64) resources.Level {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	chain := NewRandomBuilder(width, height, seed)
	chain.Build()

	// 進行ワープホールを生成する
	// FIXME: たまに届かない位置に生成される
	{
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("進行ワープホールの生成に失敗した")
			}
			x := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			y := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
			tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
			if chain.BuildData.IsSpawnableTile(world, x, y) {
				chain.BuildData.Tiles[tileIdx] = TileWarpNext

				break
			}
			failCount++
		}
	}
	// 帰還ワープホールを生成する
	if gameResources.Depth%5 == 0 {
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("帰還ワープホールの生成に失敗した")
			}
			x := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			y := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
			if chain.BuildData.IsSpawnableTile(world, x, y) {
				tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
				chain.BuildData.Tiles[tileIdx] = TileWarpEscape

				break
			}
			failCount++
		}
	}
	// フィールドに操作対象キャラを配置する
	{
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("操作対象キャラの生成に失敗した")
			}
			tx := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			ty := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
			if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
				failCount++
				continue
			}
			worldhelper.SpawnOperator(
				world,
				gc.Pixel(int(tx)*int(consts.TileSize)+int(consts.TileSize)/2),
				gc.Pixel(int(ty)*int(consts.TileSize)+int(consts.TileSize)/2),
			)
			break
		}
	}
	// フィールドにNPCを生成する
	{
		failCount := 0
		total := 5 + chain.BuildData.RandomSource.Intn(5)
		successCount := 0
		for {
			if failCount > 200 {
				log.Fatal("NPCの生成に失敗した")
			}
			tx := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			ty := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
			if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
				failCount++
				continue
			}
			worldhelper.SpawnNPC(
				world,
				gc.Pixel(int(tx)*int(consts.TileSize)+int(consts.TileSize/2)),
				gc.Pixel(int(ty)*int(consts.TileSize)+int(consts.TileSize/2)),
			)
			successCount++
			failCount = 0
			if successCount > total {
				break
			}
		}
	}

	// tilesを元にタイルエンティティを生成する
	for _i, t := range chain.BuildData.Tiles {
		i := resources.TileIdx(_i)
		x, y := chain.BuildData.Level.XYTileCoord(i)
		switch t {
		case TileFloor:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFloor(world, gc.Row(x), gc.Col(y))
		case TileWall:
			// 近傍8タイル（直交・斜め）にフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentAnyFloor(i) {
				// 壁タイプを判定してスプライト番号を決定
				wallType := chain.BuildData.GetWallType(i)
				spriteNumber := getSpriteNumberForWallType(wallType)
				chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWallWithSprite(world, gc.Row(x), gc.Col(y), spriteNumber)
			}
		case TileWarpNext:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		case TileWarpEscape:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWarpEscape(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}

// getSpriteNumberForWallType は壁タイプに対応するスプライト番号を返す
func getSpriteNumberForWallType(wallType WallType) int {
	switch wallType {
	case WallTypeTop:
		return 10 // 上壁（下に床がある）
	case WallTypeBottom:
		return 11 // 下壁（上に床がある）
	case WallTypeLeft:
		return 12 // 左壁（右に床がある）
	case WallTypeRight:
		return 13 // 右壁（左に床がある）
	case WallTypeTopLeft:
		return 14 // 左上角（右下に床がある）
	case WallTypeTopRight:
		return 15 // 右上角（左下に床がある）
	case WallTypeBottomLeft:
		return 16 // 左下角（右上に床がある）
	case WallTypeBottomRight:
		return 17 // 右下角（左上に床がある）
	case WallTypeGeneric:
		return 1 // 汎用壁（従来の壁）
	default:
		return 1 // デフォルトは従来の壁
	}
}
