package mapbuilder

import (
	"log"
	"math/rand/v2"

	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils"
	"github.com/kijimaD/ruins/lib/worldhelper/spawner"

	gc "github.com/kijimaD/ruins/lib/components"
)

// 新規に階層を生成する。
// 階層を初期化するので、具体的なコードであり、その分参照を多く含んでいる。循環参照を防ぐためにこの関数はLevel構造体とは同じpackageに属していない。
func NewLevel(world w.World, width gc.Row, height gc.Col) resources.Level {
	gameResources := world.Resources.Game.(*resources.Game)

	chain := SimpleRoomBuilder(width, height)
	chain.Build()

	// 進行ワープホールを生成する
	// FIXME: たまに届かない位置に生成される
	{
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("進行ワープホールの生成に失敗した")
			}
			x := gc.Row(rand.IntN(int(chain.BuildData.Level.TileWidth)))
			y := gc.Col(rand.IntN(int(chain.BuildData.Level.TileHeight)))
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
			x := gc.Row(rand.IntN(int(chain.BuildData.Level.TileWidth)))
			y := gc.Col(rand.IntN(int(chain.BuildData.Level.TileHeight)))
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
			tx := gc.Row(rand.IntN(int(chain.BuildData.Level.TileWidth)))
			ty := gc.Col(rand.IntN(int(chain.BuildData.Level.TileHeight)))
			if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
				failCount++
				continue
			}
			spawner.SpawnOperator(
				world,
				gc.Pixel(int(tx)*int(utils.TileSize)+int(utils.TileSize)/2),
				gc.Pixel(int(ty)*int(utils.TileSize)+int(utils.TileSize)/2),
			)
			break
		}
	}
	// フィールドにNPCを生成する
	{
		failCount := 0
		total := rand.IntN(10 + 10)
		successCount := 0
		for {
			if failCount > 200 {
				log.Fatal("NPCの生成に失敗した")
			}
			tx := gc.Row(rand.IntN(int(chain.BuildData.Level.TileWidth)))
			ty := gc.Col(rand.IntN(int(chain.BuildData.Level.TileHeight)))
			if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
				failCount++
				continue
			}
			spawner.SpawnNPC(
				world,
				gc.Pixel(int(tx)*int(utils.TileSize)+int(utils.TileSize/2)),
				gc.Pixel(int(ty)*int(utils.TileSize)+int(utils.TileSize/2)),
			)
			successCount += 1
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
			chain.BuildData.Level.Entities[i] = spawner.SpawnFloor(world, gc.Row(x), gc.Col(y))
		case TileWall:
			// 近傍4タイルにフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentOrthoAnyFloor(i) {
				chain.BuildData.Level.Entities[i] = spawner.SpawnFieldWall(world, gc.Row(x), gc.Col(y))
			}
		case TileWarpNext:
			chain.BuildData.Level.Entities[i] = spawner.SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		case TileWarpEscape:
			chain.BuildData.Level.Entities[i] = spawner.SpawnFieldWarpEscape(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}
