package mapspawner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestWarpPortalNoDuplication_StringMapPlanner(t *testing.T) {
	t.Parallel()

	// StringMapPlannerを使用してワープポータルが含まれるマップを作成
	tileMap := []string{
		"#####",
		"#fff#",
		"#fff#",
		"#fff#",
		"#####",
	}

	entityMap := []string{
		".....",
		".....",
		"..w..",
		"..e..",
		".....",
	}

	// StringMapPlannerでチェーンを作成
	chain := mapplanner.NewStringMapPlannerWithMaps(tileMap, entityMap, 12345)

	// BuildPlanでEntityPlanを生成
	plan, err := mapplanner.BuildPlan(chain)
	if err == nil {
		CompleteWallSprites(plan)
	}
	if err != nil {
		t.Fatalf("BuildPlan failed: %v", err)
	}

	// ワープポータルをカウント
	warpNextCount := 0
	warpEscapeCount := 0

	for _, entity := range plan.Entities {
		switch entity.EntityType {
		case mapplanner.EntityTypeWarpNext:
			warpNextCount++
		case mapplanner.EntityTypeWarpEscape:
			warpEscapeCount++
		}
	}

	// それぞれのワープポータルが1つずつだけ生成されていることを確認
	if warpNextCount != 1 {
		t.Errorf("進行ワープポータル数が期待値と異なる: 期待値=1, 実際=%d", warpNextCount)
	}

	if warpEscapeCount != 1 {
		t.Errorf("帰還ワープポータル数が期待値と異なる: 期待値=1, 実際=%d", warpEscapeCount)
	}

	// 総数も確認
	totalWarpPortals := warpNextCount + warpEscapeCount
	if totalWarpPortals != 2 {
		t.Errorf("ワープポータル総数が期待値と異なる: 期待値=2, 実際=%d", totalWarpPortals)
	}
}

func TestWarpPortalNoDuplication_NonStringMapPlanner(t *testing.T) {
	t.Parallel()

	// 通常のプランナー（SmallRoomPlannerなど）を使用した場合
	width, height := gc.Tile(10), gc.Tile(10)
	chain := mapplanner.NewSmallRoomPlanner(width, height, 12345)

	// ワールドを初期化
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	// BuildPlanAndSpawnでLevelを生成（ワープポータル重複が起きるかの確認）
	_, _, _, err := BuildPlanAndSpawn(world, chain, mapplanner.PlannerTypeTown)
	if err != nil {
		t.Fatalf("BuildPlanAndSpawn failed: %v", err)
	}

	// Dungeonリソースを初期化（帰還ワープホール配置のため、5の倍数の階層）
	gameResource := &resources.Dungeon{}
	gameResource.SetStateEvent(resources.StateEventNone)
	gameResource.Depth = 5 // 5の倍数にして帰還ワープホールを配置
	world.Resources.Dungeon = gameResource

	// ワープポータルプランナーを追加してワープポータル生成をテスト
	warpPlanner := mapplanner.NewWarpPortalPlanner(world, mapplanner.PlannerTypeTown)
	chain.With(warpPlanner)

	// EntityPlanを生成してワープポータル追加処理を実行（テスト用）
	plan, err := mapplanner.BuildPlan(chain)
	if err == nil {
		CompleteWallSprites(plan)
	}
	if err != nil {
		t.Fatalf("BuildPlan failed: %v", err)
	}

	// ワープポータルをカウント
	warpNextCount := 0
	warpEscapeCount := 0

	for _, entity := range plan.Entities {
		switch entity.EntityType {
		case mapplanner.EntityTypeWarpNext:
			warpNextCount++
		case mapplanner.EntityTypeWarpEscape:
			warpEscapeCount++
		}
	}

	// 通常のプランナーでもワープポータルが適切に生成されていることを確認
	if warpNextCount != 1 {
		t.Errorf("進行ワープポータル数が期待値と異なる: 期待値=1, 実際=%d", warpNextCount)
	}

	if warpEscapeCount != 1 {
		t.Errorf("帰還ワープポータル数が期待値と異なる: 期待値=1, 実際=%d", warpEscapeCount)
	}
}
