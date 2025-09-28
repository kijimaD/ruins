package mapspawner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/stretchr/testify/require"
)

func TestBuildPlan_SimpleFloorAndWall(t *testing.T) {
	t.Parallel()
	// 3x3のシンプルなマップを作成
	width, height := 3, 3
	chain := mapplanner.NewPlannerChain(gc.Tile(width), gc.Tile(height), 42)

	// タイル配列を手動で設定
	chain.PlanData.Tiles = []mapplanner.Tile{
		mapplanner.TileWall, mapplanner.TileWall, mapplanner.TileWall, // Row 0
		mapplanner.TileWall, mapplanner.TileFloor, mapplanner.TileWall, // Row 1
		mapplanner.TileWall, mapplanner.TileWall, mapplanner.TileWall, // Row 2
	}

	// BuildPlanをテスト
	plan, err := chain.PlanData.BuildPlan()
	require.NoError(t, err, "BuildPlan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// EntityPlanの基本プロパティをチェック
	if plan.Width != width {
		t.Errorf("Expected Width %d, got %d", width, plan.Width)
	}
	if plan.Height != height {
		t.Errorf("Expected Height %d, got %d", height, plan.Height)
	}

	// エンティティ数をチェック（床1個 + 隣接する壁8個）
	// 実際の数は隣接判定ロジックに依存するため、最低限のチェックのみ
	if len(plan.Entities) == 0 {
		t.Error("Expected some entities, but got none")
	}

	// 床エンティティが含まれていることを確認
	hasFloor := false
	for _, entity := range plan.Entities {
		if entity.EntityType == mapplanner.EntityTypeFloor && entity.X == 1 && entity.Y == 1 {
			hasFloor = true
			break
		}
	}
	if !hasFloor {
		t.Error("Expected floor entity at (1,1), but not found")
	}
}

func TestBuildPlan_EmptyMap(t *testing.T) {
	t.Parallel()
	// 空のマップを作成
	width, height := 2, 2
	chain := mapplanner.NewPlannerChain(gc.Tile(width), gc.Tile(height), 42)

	// 全て空のタイル
	chain.PlanData.Tiles = []mapplanner.Tile{
		mapplanner.TileEmpty, mapplanner.TileEmpty,
		mapplanner.TileEmpty, mapplanner.TileEmpty,
	}

	// 空のマップではプレイヤー位置が見つからずエラーになることを期待
	_, err := chain.PlanData.BuildPlan()
	if err == nil {
		t.Fatalf("Expected error for empty map, but got nil")
	}

	expectedMsg := "プレイヤー配置可能な床タイルが見つかりません"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestBuildPlan_WarpTiles(t *testing.T) {
	t.Parallel()
	// ワープタイルを含むマップを作成
	width, height := 2, 2
	chain := mapplanner.NewPlannerChain(gc.Tile(width), gc.Tile(height), 42)

	chain.PlanData.Tiles = []mapplanner.Tile{
		mapplanner.TileFloor, mapplanner.TileFloor,
		mapplanner.TileFloor, mapplanner.TileFloor,
	}

	// ワープポータルエンティティを追加
	chain.PlanData.WarpPortals = append(chain.PlanData.WarpPortals, mapplanner.WarpPortal{
		X:    0,
		Y:    0,
		Type: mapplanner.WarpPortalNext,
	})
	chain.PlanData.WarpPortals = append(chain.PlanData.WarpPortals, mapplanner.WarpPortal{
		X:    1,
		Y:    1,
		Type: mapplanner.WarpPortalEscape,
	})

	plan, err := chain.PlanData.BuildPlan()
	require.NoError(t, err, "BuildPlan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// ワープエンティティが含まれていることを確認
	hasWarpNext := false
	hasWarpEscape := false
	hasFloors := 0

	for _, entity := range plan.Entities {
		switch entity.EntityType {
		case mapplanner.EntityTypeWarpNext:
			if entity.X == 0 && entity.Y == 0 {
				hasWarpNext = true
			}
		case mapplanner.EntityTypeWarpEscape:
			if entity.X == 1 && entity.Y == 1 {
				hasWarpEscape = true
			}
		case mapplanner.EntityTypeFloor:
			hasFloors++
		}
	}

	if !hasWarpNext {
		t.Error("Expected WarpNext entity at (0,0), but not found")
	}
	if !hasWarpEscape {
		t.Error("Expected WarpEscape entity at (1,1), but not found")
	}
	if hasFloors != 4 {
		t.Errorf("Expected 4 floor entities, got %d", hasFloors)
	}
}

func TestGetSpriteNumberForWallType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		wallType mapplanner.WallType
		expected int
	}{
		{mapplanner.WallTypeTop, 10},
		{mapplanner.WallTypeBottom, 11},
		{mapplanner.WallTypeLeft, 12},
		{mapplanner.WallTypeRight, 13},
		{mapplanner.WallTypeGeneric, 1},
	}

	for _, tc := range testCases {
		actual := getSpriteNumberForWallType(tc.wallType)
		if actual != tc.expected {
			t.Errorf("壁タイプ %s のスプライト番号が間違っています。期待値: %d, 実際: %d",
				tc.wallType.String(), tc.expected, actual)
		}
	}
}

// TestBuildPlan_Integration は実際のBuilderChainとの統合テスト
func TestBuildPlan_Integration(t *testing.T) {
	t.Parallel()
	// 実際のSmallRoomBuilderを使用
	width, height := 10, 10
	chain, err := mapplanner.NewSmallRoomPlanner(gc.Tile(width), gc.Tile(height), 12345)
	require.NoError(t, err, "NewSmallRoomPlanner failed")

	// マップを生成
	require.NoError(t, chain.Plan(), "Plan failed")

	// BuildPlanをテスト
	plan, err := chain.PlanData.BuildPlan()
	require.NoError(t, err, "BuildPlan integration test failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// 基本的な整合性をチェック
	if plan.Width != width {
		t.Errorf("Expected Width %d, got %d", width, plan.Width)
	}
	if plan.Height != height {
		t.Errorf("Expected Height %d, got %d", height, plan.Height)
	}

	// エンティティが生成されていることを確認
	if len(plan.Entities) == 0 {
		t.Error("Expected some entities from SmallRoomBuilder, but got none")
	}

	// 床エンティティが含まれていることを確認
	hasFloor := false
	for _, entity := range plan.Entities {
		if entity.EntityType == mapplanner.EntityTypeFloor {
			hasFloor = true
			break
		}
	}
	if !hasFloor {
		t.Error("Expected at least one floor entity from SmallRoomBuilder")
	}
}
