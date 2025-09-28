package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// Plan はPlannerChainを初期化して追加プランナーを統合し、EntityPlanを構築する
func Plan(world w.World, width, height int, seed uint64, plannerType PlannerType) (*EntityPlan, error) {
	// PlannerChainを初期化
	var chain *PlannerChain
	if plannerType.Name == PlannerTypeRandom.Name {
		chain = NewRandomPlanner(gc.Tile(width), gc.Tile(height), seed)
	} else {
		chain = plannerType.PlannerFunc(gc.Tile(width), gc.Tile(height), seed)
	}
	// ワープポータルプランナーを追加する
	// プランナータイプによってはもうすでに計画されているので判定する
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
