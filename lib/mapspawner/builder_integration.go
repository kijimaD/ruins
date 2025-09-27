package mapspawner

import (
	"fmt"

	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// PlanAndSpawn はPlannerChainを実行してEntityPlanを生成し、Levelをスポーンする
func PlanAndSpawn(world w.World, chain *mapplanner.PlannerChain, plannerType mapplanner.PlannerType) (resources.Level, int, int, error) {
	// ワープポータルプランナーを追加（StringMapPlannerは既に独自にワープポータルを処理しているため条件判定）
	if len(chain.PlanData.GetWarpPortals()) == 0 {
		warpPlanner := mapplanner.NewWarpPortalPlanner(world, plannerType)
		chain.With(warpPlanner)
	}

	// NPCプランナーを追加
	if plannerType.SpawnEnemies {
		npcPlanner := mapplanner.NewNPCPlanner(world, plannerType)
		chain.With(npcPlanner)
	}

	// アイテムプランナーを追加
	if plannerType.SpawnItems {
		itemPlanner := mapplanner.NewItemPlanner(world, plannerType)
		chain.With(itemPlanner)
	}

	// Propsプランナーを追加（町タイプで固定Props配置）
	propsPlanner := mapplanner.NewPropsPlanner(world, plannerType)
	chain.With(propsPlanner)

	// プランナーチェーンを実行
	chain.Plan()

	// PlanDataからEntityPlanを構築
	plan, err := chain.PlanData.BuildPlanFromTiles()
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
