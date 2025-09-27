package mapplanner

import (
	"fmt"

	w "github.com/kijimaD/ruins/lib/world"
)

// Plan はPlannerChainに追加プランナーを統合してEntityPlanを構築する
func Plan(world w.World, chain *PlannerChain, plannerType PlannerType) (*EntityPlan, error) {
	// ワープポータルプランナーを追加（StringMapPlannerは既に独自にワープポータルを処理しているため条件判定）
	if len(chain.PlanData.WarpPortals) == 0 {
		warpPlanner := NewWarpPortalPlanner(world, plannerType)
		chain.With(warpPlanner)
	}

	// NPCプランナーを追加
	if plannerType.SpawnEnemies {
		npcPlanner := NewNPCPlanner(world, plannerType)
		chain.With(npcPlanner)
	}

	// アイテムプランナーを追加
	if plannerType.SpawnItems {
		itemPlanner := NewItemPlanner(world, plannerType)
		chain.With(itemPlanner)
	}

	// Propsプランナーを追加（町タイプで固定Props配置）
	propsPlanner := NewPropsPlanner(world, plannerType)
	chain.With(propsPlanner)

	// プランナーチェーンを実行
	chain.Plan()

	// PlanDataからEntityPlanを構築
	plan, err := chain.PlanData.BuildPlanFromTiles()
	if err != nil {
		return nil, fmt.Errorf("EntityPlan構築エラー: %w", err)
	}

	return plan, nil
}
