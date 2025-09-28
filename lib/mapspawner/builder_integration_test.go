package mapspawner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplanner"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/world"
	"github.com/stretchr/testify/require"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestMapPlannerBuildPlan(t *testing.T) {
	t.Parallel()
	// SmallRoomBuilderチェーンを作成
	width, height := 8, 8
	chain := mapplanner.NewSmallRoomPlanner(gc.Tile(width), gc.Tile(height), 42)

	// BuildPlanをテスト
	plan, err := mapplanner.BuildPlan(chain)
	if err == nil {
		completeWallSprites(plan)
	}
	require.NoError(t, err, "BuildPlan failed")

	// EntityPlanの基本プロパティをチェック
	if plan.Width != width {
		t.Errorf("Expected Width %d, got %d", width, plan.Width)
	}
	if plan.Height != height {
		t.Errorf("Expected Height %d, got %d", height, plan.Height)
	}

	// エンティティが生成されていることを確認
	if len(plan.Entities) == 0 {
		t.Error("Expected some entities, but got none")
	}
}

func TestBuildPlanAndSpawn(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)

	// マップサイズとシード
	width, height := 6, 6
	seed := uint64(123)

	// EntityPlan構築とSpawnLevelを個別にテスト（NPCとアイテム生成を無効化）
	plannerType := mapplanner.PlannerType{
		Name:         "SmallRoom",
		SpawnEnemies: false, // テストではNPC生成を無効化
		SpawnItems:   false, // テストではアイテム生成を無効化
		PlannerFunc:  mapplanner.PlannerTypeSmallRoom.PlannerFunc,
	}

	// EntityPlan構築
	plan, err := mapplanner.Plan(world, width, height, seed, plannerType)
	require.NoError(t, err, "Plan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// SpawnLevel
	level, err := Spawn(world, plan)
	require.NoError(t, err, "SpawnLevel failed")

	// Levelの基本プロパティをチェック
	if level.TileWidth != gc.Tile(width) {
		t.Errorf("Expected TileWidth %d, got %d", width, level.TileWidth)
	}
	if level.TileHeight != gc.Tile(height) {
		t.Errorf("Expected TileHeight %d, got %d", height, level.TileHeight)
	}
	if len(level.Entities) != width*height {
		t.Errorf("Expected %d entities, got %d", width*height, len(level.Entities))
	}

	// 少なくとも1つの非ゼロエンティティが存在することを確認
	hasNonZeroEntity := false
	for _, entity := range level.Entities {
		if entity != 0 {
			hasNonZeroEntity = true
			break
		}
	}
	if !hasNonZeroEntity {
		t.Error("Expected at least one non-zero entity, but all entities are zero")
	}
}

func TestBuildPlanAndSpawn_TownBuilder(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)

	// マップサイズとシード
	width, height := 15, 15
	seed := uint64(456)

	// EntityPlan構築とSpawnLevelを個別にテスト（NPCとアイテム生成を無効化）
	plannerType := mapplanner.PlannerType{
		Name:              "Town",
		SpawnEnemies:      false, // テストではNPC生成を無効化
		SpawnItems:        false, // テストではアイテム生成を無効化
		UseFixedPortalPos: true,  // ポータル位置を固定
		PlannerFunc:       mapplanner.PlannerTypeTown.PlannerFunc,
	}

	// EntityPlan構築
	plan, err := mapplanner.Plan(world, width, height, seed, plannerType)
	require.NoError(t, err, "Plan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// SpawnLevel
	level, err := Spawn(world, plan)
	require.NoError(t, err, "SpawnLevel with TownBuilder failed")

	// Levelの基本プロパティをチェック
	if level.TileWidth != gc.Tile(width) {
		t.Errorf("Expected TileWidth %d, got %d", width, level.TileWidth)
	}
	if level.TileHeight != gc.Tile(height) {
		t.Errorf("Expected TileHeight %d, got %d", height, level.TileHeight)
	}
	if len(level.Entities) != width*height {
		t.Errorf("Expected %d entities, got %d", width*height, len(level.Entities))
	}
}

func TestTownBuilderWithPortals(t *testing.T) {
	t.Parallel()
	// BuilderChainを作成してタイル配置をテスト
	// TownPlannerは固定の50x50マップを生成する
	width, height := 50, 50
	chain := mapplanner.NewTownPlanner(gc.Tile(width), gc.Tile(height), 123)

	// マップを構築
	chain.Plan()

	// 中央座標
	centerX := width / 2
	centerY := height / 2

	// 中央のタイルが床タイルかを確認
	centerIdx := chain.PlanData.Level.XYTileIndex(gc.Tile(centerX), gc.Tile(centerY))
	centerTile := chain.PlanData.Tiles[centerIdx]
	t.Logf("Center tile at (%d,%d): %v", centerX, centerY, centerTile)

	if centerTile != (raw.TileRaw{Walkable: true}) {
		t.Errorf("Expected center tile to be floor, got %v", centerTile)
	}

	// 公民館の中心位置にワープポータルタイルを手動で配置（50x50マップに合わせて調整）
	communityHallX := centerX
	communityHallY := centerY + 16 // 50x50マップで公民館の中央位置
	if communityHallY >= height {
		communityHallY = height - 1
	}
	portalIdx := chain.PlanData.Level.XYTileIndex(gc.Tile(communityHallX), gc.Tile(communityHallY))
	// 直接タイルアクセスが必要な場合は専用メソッドを追加検討
	chain.PlanData.Tiles[portalIdx] = (raw.TileRaw{Walkable: true})
	// ワープポータルエンティティを追加
	chain.PlanData.WarpPortals = append(chain.PlanData.WarpPortals, mapplanner.WarpPortal{
		X:    communityHallX,
		Y:    communityHallY,
		Type: mapplanner.WarpPortalNext,
	})

	// BuildPlanを使用してEntityPlanを生成
	plan, err := chain.PlanData.BuildPlan()
	require.NoError(t, err, "BuildPlan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// ワープポータルエンティティが含まれているかチェック（公民館の中央）
	hasWarpPortal := false
	for _, entity := range plan.Entities {
		if entity.EntityType == mapplanner.EntityTypeWarpNext &&
			entity.X == communityHallX && entity.Y == communityHallY {
			hasWarpPortal = true
			break
		}
	}

	if !hasWarpPortal {
		t.Errorf("Expected warp portal at community hall (%d,%d) but found none", communityHallX, communityHallY)
	}

	// 実際にSpawnLevelでエンティティが生成されるかテスト
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)

	level, err := Spawn(world, plan)
	require.NoError(t, err, "SpawnLevel failed")

	// 公民館の中央にエンティティが生成されているかチェック
	portalEntityIdx := level.XYTileIndex(gc.Tile(communityHallX), gc.Tile(communityHallY))
	portalEntity := level.Entities[portalEntityIdx]
	t.Logf("Portal entity at (%d,%d): %v", communityHallX, communityHallY, portalEntity)

	if portalEntity == 0 {
		t.Errorf("Expected warp portal entity at community hall (%d,%d) but found none", communityHallX, communityHallY)
	}
}

// 街のBuildPlanAndSpawnでポータルが配置されるかの統合テスト
func TestTownBuildPlanAndSpawnFullFlow(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)

	// TownPlannerは固定の50x50マップを生成する
	width, height := 50, 50
	seed := uint64(123)

	// EntityPlan構築とSpawnLevelを個別に実行（街の設定で）
	plan, err := mapplanner.Plan(world, width, height, seed, mapplanner.PlannerTypeTown)
	require.NoError(t, err, "Plan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// SpawnLevel
	level, err := Spawn(world, plan)
	require.NoError(t, err, "SpawnLevel failed")

	// 中央座標
	centerX := width / 2
	centerY := height / 2
	// 公民館の中央座標（50x50マップに合わせて調整）
	communityHallX := centerX
	communityHallY := centerY + 16 // 50x50マップで公民館の中央位置
	if communityHallY >= height {
		communityHallY = height - 1
	}

	t.Logf("Testing BuildPlanAndSpawn with town center at (%d,%d)", centerX, centerY)
	t.Logf("Expected warp portal at community hall (%d,%d)", communityHallX, communityHallY)
	t.Logf("BuilderTypeTown.UseFixedPortalPos: %v", mapplanner.PlannerTypeTown.UseFixedPortalPos)

	// 公民館の中央にエンティティが生成されているかチェック
	portalEntityIdx := level.XYTileIndex(gc.Tile(communityHallX), gc.Tile(communityHallY))
	portalEntity := level.Entities[portalEntityIdx]
	t.Logf("Portal entity at (%d,%d): %v", communityHallX, communityHallY, portalEntity)

	if portalEntity == 0 {
		t.Errorf("Expected portal entity at community hall (%d,%d) but found none", communityHallX, communityHallY)
	}
}

func TestBuildPlanAndSpawn_BigRoomBuilder(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	require.NoError(t, components.InitializeComponents(&ecs.Manager{}), "InitializeComponents failed")
	world, _ := world.InitWorld(components)

	// マップサイズとシード
	width, height := 12, 12
	seed := uint64(789)

	// EntityPlan構築とSpawnLevelを個別にテスト（NPCとアイテム生成を無効化）
	plannerType := mapplanner.PlannerType{
		Name:         "BigRoom",
		SpawnEnemies: false, // テストではNPC生成を無効化
		SpawnItems:   false, // テストではアイテム生成を無効化
		PlannerFunc:  mapplanner.PlannerTypeBigRoom.PlannerFunc,
	}

	// EntityPlan構築
	plan, err := mapplanner.Plan(world, width, height, seed, plannerType)
	require.NoError(t, err, "Plan failed")

	// 壁スプライト番号を補完
	completeWallSprites(plan)

	// SpawnLevel
	level, err := Spawn(world, plan)
	require.NoError(t, err, "SpawnLevel with BigRoomBuilder failed")

	// Levelの基本プロパティをチェック
	if level.TileWidth != gc.Tile(width) {
		t.Errorf("Expected TileWidth %d, got %d", width, level.TileWidth)
	}
	if level.TileHeight != gc.Tile(height) {
		t.Errorf("Expected TileHeight %d, got %d", height, level.TileHeight)
	}
	if len(level.Entities) != width*height {
		t.Errorf("Expected %d entities, got %d", width*height, len(level.Entities))
	}
}

// TestBuildPlan_Reproducible は同じシードで同じ結果が得られることを確認
func TestBuildPlan_Reproducible(t *testing.T) {
	t.Parallel()
	width, height := 7, 7
	seed := uint64(999)

	// 同じパラメータで2回実行
	chain1 := mapplanner.NewSmallRoomPlanner(gc.Tile(width), gc.Tile(height), seed)
	plan1, err1 := mapplanner.BuildPlan(chain1)
	if err1 != nil {
		t.Fatalf("First BuildPlan failed: %v", err1)
	}
	completeWallSprites(plan1)

	chain2 := mapplanner.NewSmallRoomPlanner(gc.Tile(width), gc.Tile(height), seed)
	plan2, err2 := mapplanner.BuildPlan(chain2)
	if err2 != nil {
		t.Fatalf("Second BuildPlan failed: %v", err2)
	}
	completeWallSprites(plan2)

	// 結果が同じであることを確認
	if plan1.Width != plan2.Width {
		t.Errorf("Width mismatch: %d vs %d", plan1.Width, plan2.Width)
	}
	if plan1.Height != plan2.Height {
		t.Errorf("Height mismatch: %d vs %d", plan1.Height, plan2.Height)
	}
	if len(plan1.Entities) != len(plan2.Entities) {
		t.Errorf("Entity count mismatch: %d vs %d", len(plan1.Entities), len(plan2.Entities))
	}

	// エンティティの配置が同じであることを確認
	for i, entity1 := range plan1.Entities {
		if i >= len(plan2.Entities) {
			break
		}
		entity2 := plan2.Entities[i]
		if entity1.X != entity2.X || entity1.Y != entity2.Y || entity1.EntityType != entity2.EntityType {
			t.Errorf("Entity %d mismatch: (%d,%d,%v) vs (%d,%d,%v)",
				i, entity1.X, entity1.Y, entity1.EntityType,
				entity2.X, entity2.Y, entity2.EntityType)
		}
	}
}
