package mapbuilder

import (
	"log"
	"math/rand"

	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/utils/consts"

	gc "github.com/kijimaD/ruins/lib/components"
)

// 新規に階層を生成する。
// 階層を初期化するので、具体的なコードであり、その分参照を多く含んでいる。循環参照を防ぐためにこの関数はLevel構造体とは同じpackageに属していない。
func NewLevel(world w.World, width gc.Row, height gc.Col) resources.Level {
	gameResources := world.Resources.Game.(*resources.Game)

	chain := SimpleRoomBuilder(width, height)
	chain.Build()

	// ワープホールを生成する
	// FIXME: たまに届かない位置に生成される
	{
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("進行ワープホールの生成に失敗した")
			}
			x := rand.Intn(int(chain.BuildData.Level.TileWidth))
			y := rand.Intn(int(chain.BuildData.Level.TileHeight))
			tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
			if chain.BuildData.IsSpawnableTile(x, y) {
				chain.BuildData.Tiles[tileIdx] = TileWarpNext

				break
			}
			failCount++
		}
	}

	if gameResources.Depth%5 == 0 {
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("帰還ワープホールの生成に失敗した")
			}
			x := rand.Intn(int(chain.BuildData.Level.TileWidth))
			y := rand.Intn(int(chain.BuildData.Level.TileHeight))
			if chain.BuildData.IsSpawnableTile(x, y) {
				tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
				chain.BuildData.Tiles[tileIdx] = TileWarpEscape

				break
			}
			failCount++
		}
	}
	// フィールドにプレイヤーを配置する
	{
		failCount := 0
		for {
			if failCount > 200 {
				log.Fatal("プレイヤーの生成に失敗した")
			}
			x := rand.Intn(int(chain.BuildData.Level.TileWidth))
			y := rand.Intn(int(chain.BuildData.Level.TileHeight))
			if chain.BuildData.IsSpawnableTile(x, y) {
				resources.SpawnPlayer(world, x*consts.TileSize+consts.TileSize/2, y*consts.TileSize+consts.TileSize/2)
				break
			}
			failCount++
		}
	}
	{
		failCount := 0
		NPCCount := 0
		for {
			if failCount > 200 {
				log.Fatal("NPCの生成に失敗した")
			}
			x := rand.Intn(int(chain.BuildData.Level.TileWidth))
			y := rand.Intn(int(chain.BuildData.Level.TileHeight))
			// TODO: キャラがかぶって生成されるのを直す
			if chain.BuildData.IsSpawnableTile(x, y) {
				resources.SpawnNPC(world, x*consts.TileSize+consts.TileSize/2, y*consts.TileSize+consts.TileSize/2)
				NPCCount += 1
				if NPCCount > 10 {
					break
				}
			}
			failCount++
		}
	}

	// tilesを元にエンティティを生成する
	for i, t := range chain.BuildData.Tiles {
		x, y := chain.BuildData.Level.XYTileCoord(i)
		switch t {
		case TileFloor:
			chain.BuildData.Level.Entities[i] = resources.SpawnFloor(world, gc.Row(x), gc.Col(y))
		case TileWall:
			// 近傍4タイルにフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentOrthoAnyFloor(i) {
				chain.BuildData.Level.Entities[i] = resources.SpawnFieldWall(world, gc.Row(x), gc.Col(y))
			}
		case TileWarpNext:
			chain.BuildData.Level.Entities[i] = resources.SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		case TileWarpEscape:
			chain.BuildData.Level.Entities[i] = resources.SpawnFieldWarpEscape(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}
