package mapspawner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// PlanAndSpawn はEntityPlanを構築してLevelをスポーンする
func PlanAndSpawn(world w.World, width, height int, seed uint64, plannerType mapplanner.PlannerType) (resources.Level, int, int, error) {
	// PlannerChainを初期化
	var chain *mapplanner.PlannerChain
	if plannerType.Name == mapplanner.PlannerTypeRandom.Name {
		chain = mapplanner.NewRandomPlanner(gc.Tile(width), gc.Tile(height), seed)
	} else {
		chain = plannerType.PlannerFunc(gc.Tile(width), gc.Tile(height), seed)
	}

	// EntityPlan構築
	plan, err := mapplanner.Plan(world, chain, plannerType)
	if err != nil {
		return resources.Level{}, 0, 0, fmt.Errorf("EntityPlan構築エラー: %w", err)
	}

	// 壁スプライト番号を補完
	CompleteWallSprites(plan)

	// プレイヤー位置を取得
	playerX, playerY, hasPlayerPos := plan.GetPlayerStartPosition()
	if !hasPlayerPos {
		return resources.Level{}, 0, 0, fmt.Errorf("EntityPlanにプレイヤー開始位置が設定されていません")
	}

	// EntityPlanからLevelをスポーン
	level, err := SpawnLevel(world, plan)
	if err != nil {
		return resources.Level{}, 0, 0, fmt.Errorf("level生成エラー: %w", err)
	}

	return level, playerX, playerY, nil
}
