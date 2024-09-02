package resources

import (
	"log"
	"math/rand"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
	"github.com/kijimaD/ruins/lib/loader"
	"github.com/kijimaD/ruins/lib/mapbuilder"
	"github.com/kijimaD/ruins/lib/utils/consts"
)

// 冒険出発から帰還までを1セットとした情報を保持する。
// 冒険出発から帰還までは複数階層が存在し、複数階層を通しての情報を保持する必要がある。
type Game struct {
	// フィールド上で発生したイベント。各stateで補足されて処理される
	StateEvent StateEvent
	// 現在階のフィールド情報
	Level loader.Level
	// 階層数
	Depth int
}

// 新規に階層を生成する。
func NewLevel(world w.World, width gc.Row, height gc.Col) loader.Level {
	gameResources := world.Resources.Game.(*Game)

	chain := mapbuilder.SimpleRoomBuilder(width, height)
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
				chain.BuildData.Tiles[tileIdx] = mapbuilder.TileWarpNext

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
				chain.BuildData.Tiles[tileIdx] = mapbuilder.TileWarpEscape

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
				SpawnPlayer(world, x*consts.TileSize+consts.TileSize/2, y*consts.TileSize+consts.TileSize/2)
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
				SpawnNPC(world, x*consts.TileSize+consts.TileSize/2, y*consts.TileSize+consts.TileSize/2)
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
		case mapbuilder.TileFloor:
			chain.BuildData.Level.Entities[i] = SpawnFloor(world, gc.Row(x), gc.Col(y))
		case mapbuilder.TileWall:
			// 近傍4タイルにフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentOrthoAnyFloor(i) {
				chain.BuildData.Level.Entities[i] = SpawnFieldWall(world, gc.Row(x), gc.Col(y))
			}
		case mapbuilder.TileWarpNext:
			chain.BuildData.Level.Entities[i] = SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		case mapbuilder.TileWarpEscape:
			chain.BuildData.Level.Entities[i] = SpawnFieldWarpEscape(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}
