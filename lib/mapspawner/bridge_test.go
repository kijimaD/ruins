package mapspawner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSpawnFromMetaPlan_SimpleFloorAndWall(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)
	world.Resources.RawMaster = createMapspawnerTestRawMaster()

	// 3x3のシンプルなマップを作成
	width, height := 3, 3
	seed := uint64(42)
	plannerType := mapplanner.PlannerType{
		Name:         "SmallRoom",
		SpawnEnemies: false,
		SpawnItems:   false,
		PlannerFunc:  mapplanner.PlannerTypeSmallRoom.PlannerFunc,
	}

	// MetaPlanを生成
	metaPlan, err := mapplanner.Plan(world, width, height, seed, plannerType)
	require.NoError(t, err, "Plan failed")

	// SpawnFromMetaPlanをテスト
	level, err := Spawn(world, metaPlan)
	require.NoError(t, err, "SpawnFromMetaPlan failed")

	// Levelの基本プロパティをチェック
	if level.TileWidth != gc.Tile(width) {
		t.Errorf("Expected TileWidth %d, got %d", width, level.TileWidth)
	}
	if level.TileHeight != gc.Tile(height) {
		t.Errorf("Expected TileHeight %d, got %d", height, level.TileHeight)
	}

	// エンティティが生成されていることを確認
	if len(level.Entities) == 0 {
		t.Error("Expected some entities, but got none")
	}

	// 非ゼロエンティティ（実際に生成されたエンティティ）が存在することを確認
	hasNonZeroEntity := false
	for _, entity := range level.Entities {
		if entity != 0 {
			hasNonZeroEntity = true
			break
		}
	}
	if !hasNonZeroEntity {
		t.Error("Expected at least one non-zero entity")
	}
}

func TestGetSpriteKeyForWallType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		wallType mapplanner.WallType
		expected string
	}{
		{mapplanner.WallTypeTop, "wall_top"},
		{mapplanner.WallTypeBottom, "wall_bottom"},
		{mapplanner.WallTypeLeft, "wall_left"},
		{mapplanner.WallTypeRight, "wall_right"},
		{mapplanner.WallTypeTopLeft, "wall_corner_tl"},
		{mapplanner.WallTypeTopRight, "wall_corner_tr"},
		{mapplanner.WallTypeBottomLeft, "wall_corner_bl"},
		{mapplanner.WallTypeBottomRight, "wall_corner_br"},
		{mapplanner.WallTypeGeneric, "wall_generic"},
	}

	for _, tc := range testCases {
		actual := getSpriteKeyForWallType(tc.wallType)
		if actual != tc.expected {
			t.Errorf("壁タイプ %s のスプライトキーが間違っています。期待値: %s, 実際: %s",
				tc.wallType.String(), tc.expected, actual)
		}
	}
}

// TestSpawnFromMetaPlan_Integration は実際のPlannerChainとの統合テスト
func TestSpawnFromMetaPlan_Integration(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)
	world.Resources.RawMaster = createMapspawnerTestRawMaster()

	// 実際のSmallRoomPlannerを使用
	width, height := 10, 10
	seed := uint64(12345)
	plannerType := mapplanner.PlannerType{
		Name:         "SmallRoom",
		SpawnEnemies: false,
		SpawnItems:   false,
		PlannerFunc:  mapplanner.PlannerTypeSmallRoom.PlannerFunc,
	}

	// MetaPlanを生成
	metaPlan, err := mapplanner.Plan(world, width, height, seed, plannerType)
	require.NoError(t, err, "Plan failed")

	// SpawnFromMetaPlanをテスト
	level, err := Spawn(world, metaPlan)
	require.NoError(t, err, "SpawnFromMetaPlan integration test failed")

	// 基本的な整合性をチェック
	if level.TileWidth != gc.Tile(width) {
		t.Errorf("Expected TileWidth %d, got %d", width, level.TileWidth)
	}
	if level.TileHeight != gc.Tile(height) {
		t.Errorf("Expected TileHeight %d, got %d", height, level.TileHeight)
	}

	// エンティティが生成されていることを確認
	if len(level.Entities) == 0 {
		t.Error("Expected some entities from SmallRoomPlanner, but got none")
	}

	// 非ゼロエンティティが存在することを確認
	hasNonZeroEntity := false
	for _, entity := range level.Entities {
		if entity != 0 {
			hasNonZeroEntity = true
			break
		}
	}
	if !hasNonZeroEntity {
		t.Error("Expected at least one non-zero entity from SmallRoomPlanner")
	}
}
