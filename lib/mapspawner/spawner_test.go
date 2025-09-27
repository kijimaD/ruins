package mapspawner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestSpawnLevel_ValidPlan(t *testing.T) {
	t.Parallel()
	// テスト用のワールドを作成
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	// シンプルなMapPlanを作成（エンティティなし）
	plan := mapplanner.NewMapPlan(3, 3)

	// SpawnLevelをテスト（エンティティ生成なしで基本機能をテスト）
	level, err := SpawnLevel(world, plan)
	if err != nil {
		t.Fatalf("SpawnLevel failed: %v", err)
	}

	// Levelの基本プロパティをチェック
	if level.TileWidth != 3 {
		t.Errorf("Expected TileWidth 3, got %d", level.TileWidth)
	}
	if level.TileHeight != 3 {
		t.Errorf("Expected TileHeight 3, got %d", level.TileHeight)
	}
	if len(level.Entities) != 9 { // 3x3 = 9（エンティティ配列のサイズ）
		t.Errorf("Expected 9 entity slots, got %d", len(level.Entities))
	}

	// 空のLevelが正常に作成されることを確認
	allZeroEntities := true
	for _, entity := range level.Entities {
		if entity != 0 {
			allZeroEntities = false
			break
		}
	}
	if !allZeroEntities {
		t.Log("Note: Some entities were generated despite empty plan")
	}
}

func TestSpawnLevel_InvalidPlan(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	// 無効な座標を持つMapPlanを作成
	plan := mapplanner.NewMapPlan(2, 2)
	plan.AddFloor(5, 5) // 範囲外の座標

	// SpawnLevelがエラーを返すことを確認
	_, err := SpawnLevel(world, plan)
	if err == nil {
		t.Error("Expected error for invalid plan, but got none")
	}
}

func TestSpawnLevel_EmptyPlan(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	// 空のMapPlanを作成
	plan := mapplanner.NewMapPlan(2, 2)

	// 空のプランでもエラーなく動作することを確認
	level, err := SpawnLevel(world, plan)
	if err != nil {
		t.Fatalf("SpawnLevel failed with empty plan: %v", err)
	}

	// 空のLevelが作成されることを確認
	if level.TileWidth != 2 {
		t.Errorf("Expected TileWidth 2, got %d", level.TileWidth)
	}
	if level.TileHeight != 2 {
		t.Errorf("Expected TileHeight 2, got %d", level.TileHeight)
	}
	if len(level.Entities) != 4 { // 2x2 = 4
		t.Errorf("Expected 4 entities, got %d", len(level.Entities))
	}

	// 全てのエンティティが空（0）であることを確認
	for i, entity := range level.Entities {
		if entity != 0 {
			t.Errorf("Expected empty entity at index %d, but got %d", i, entity)
		}
	}
}

func TestSpawnEntityFromSpec_WallWithoutSprite(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	spec := mapplanner.EntitySpec{
		X:          0,
		Y:          0,
		EntityType: mapplanner.EntityTypeWall,
		// WallSprite is nil
	}

	_, err := spawnEntityFromSpec(world, spec)
	if err == nil {
		t.Error("Expected error for wall entity without sprite, but got none")
	}
}

func TestSpawnEntityFromSpec_Prop(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	propType := gc.PropTypeTable
	spec := mapplanner.EntitySpec{
		X:          2,
		Y:          2,
		EntityType: mapplanner.EntityTypeProp,
		PropType:   &propType,
	}

	entity, err := spawnEntityFromSpec(world, spec)
	if err != nil {
		t.Fatalf("Failed to spawn prop entity: %v", err)
	}

	if entity == 0 {
		t.Error("Expected non-zero entity, got 0")
	}
}

func TestSpawnEntityFromSpec_PropWithoutType(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	spec := mapplanner.EntitySpec{
		X:          2,
		Y:          2,
		EntityType: mapplanner.EntityTypeProp,
		// PropType is nil
	}

	_, err := spawnEntityFromSpec(world, spec)
	if err == nil {
		t.Error("Expected error for prop entity without type, but got none")
	}
}

func TestSpawnEntityFromSpec_UnknownEntityType(t *testing.T) {
	t.Parallel()
	components := &gc.Components{}
	if err := components.InitializeComponents(&ecs.Manager{}); err != nil {
		t.Fatalf("InitializeComponents failed: %v", err)
	}
	world, _ := world.InitWorld(components)

	spec := mapplanner.EntitySpec{
		X:          1,
		Y:          1,
		EntityType: mapplanner.EntityType(999), // 未知のタイプ
	}

	_, err := spawnEntityFromSpec(world, spec)
	if err == nil {
		t.Error("Expected error for unknown entity type, but got none")
	}
}
