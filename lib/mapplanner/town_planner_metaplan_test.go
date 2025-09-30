package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/assert"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestNewTownPlannerMetaPlan(t *testing.T) {
	t.Parallel()

	// テスト用のワールドを作成
	components := &gc.Components{}
	assert.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)
	world.Resources.RawMaster = createTownTestRawMaster()

	// TownPlannerで街マップを生成
	width, height := 50, 50
	chain := NewTownPlanner(gc.Tile(width), gc.Tile(height), 123)
	chain.PlanData.RawMaster = createTownTestRawMaster()

	// マップを構築
	chain.Plan()

	// MetaPlanが正しく生成されているかテスト
	metaPlan := &chain.PlanData

	// サイズ確認
	assert.Equal(t, gc.Tile(50), metaPlan.Level.TileWidth, "TileWidth should be 50")
	assert.Equal(t, gc.Tile(50), metaPlan.Level.TileHeight, "TileHeight should be 50")

	// タイル数確認
	expectedTiles := 50 * 50
	assert.Equal(t, expectedTiles, len(metaPlan.Tiles), "Should have correct number of tiles")

	// NPCは現在スキップ中
	assert.Equal(t, 0, len(metaPlan.NPCs), "Town NPCs are currently disabled")

	// ワープポータルが配置されているか確認
	assert.Equal(t, 1, len(metaPlan.WarpPortals), "Should have exactly one warp portal")

	// Propsが配置されているか確認
	assert.Greater(t, len(metaPlan.Props), 0, "Should have Props")

	// 壁と床タイルが適切に設定されているか確認
	wallCount := 0
	floorCount := 0
	for _, tile := range metaPlan.Tiles {
		if !tile.Walkable {
			wallCount++
		} else {
			floorCount++
		}
	}

	assert.Greater(t, wallCount, 0, "Should have wall tiles")
	assert.Greater(t, floorCount, 0, "Should have floor tiles")

	t.Logf("街マップ生成成功: 壁=%d, 床=%d, NPC=%d(スキップ中), Props=%d, ワープ=%d",
		wallCount, floorCount, len(metaPlan.NPCs), len(metaPlan.Props), len(metaPlan.WarpPortals))
}

func TestTownLayoutValidation(t *testing.T) {
	t.Parallel()

	// 正常なレイアウトのテスト
	tileMap, entityMap := getTownLayout()
	err := validateTownLayout(tileMap, entityMap)
	assert.NoError(t, err, "Valid town layout should not have errors")

	// ワープホール数の確認
	warpCount := 0
	for _, row := range entityMap {
		for _, char := range row {
			if char == 'w' {
				warpCount++
			}
		}
	}
	assert.Equal(t, 1, warpCount, "Should have exactly one warp portal")

	// プレイヤー開始位置の確認
	playerCount := 0
	for _, row := range entityMap {
		for _, char := range row {
			if char == '@' {
				playerCount++
			}
		}
	}
	assert.Equal(t, 1, playerCount, "Should have exactly one player start position")
}

// createTownTestRawMaster はテスト用のRawMasterを作成する
func createTownTestRawMaster() *raw.Master {
	rawMaster, _ := raw.Load(`
[[Tiles]]
Name = "Floor"
Description = "床タイル"
Walkable = true

[[Tiles]]
Name = "Wall"
Description = "壁タイル"
Walkable = false

[[Tiles]]
Name = "Empty"
Description = "空のタイル"
Walkable = false

[[Tiles]]
Name = "Dirt"
Description = "土タイル"
Walkable = true
`)
	return &rawMaster
}
