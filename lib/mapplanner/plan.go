package mapplanner

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	w "github.com/kijimaD/ruins/lib/world"
)

const (
	// MaxPlanRetries はプランナーチェーンの最大再試行回数
	MaxPlanRetries = 10
)

var (
	// ErrConnectivity は接続性エラーを表す
	ErrConnectivity = errors.New("マップ接続性エラー: プレイヤーからワープホールに到達できません")
	// ErrPlayerPlacement はプレイヤー配置エラーを表す
	ErrPlayerPlacement = errors.New("プレイヤー配置可能な床タイルが見つかりません")
)

// Plan はPlannerChainを初期化して追加プランナーを統合し、EntityPlanを構築する
// 接続性検証に失敗した場合は規定回数まで再試行する
func Plan(world w.World, width, height int, seed uint64, plannerType PlannerType) (*EntityPlan, error) {
	var lastErr error

	// 最大再試行回数まで繰り返す
	for attempt := 0; attempt < MaxPlanRetries; attempt++ {
		// 再試行時は異なるシードを使用
		currentSeed := seed + uint64(attempt*1000)

		plan, err := attemptPlan(world, width, height, currentSeed, plannerType)
		if err == nil {
			return plan, nil
		}

		lastErr = err
		// 接続性エラー以外は即座に失敗
		if !isConnectivityError(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("プラン生成に%d回失敗しました。最後のエラー: %w", MaxPlanRetries, lastErr)
}

// attemptPlan は単一回のプラン生成を試行する
func attemptPlan(world w.World, width, height int, seed uint64, plannerType PlannerType) (*EntityPlan, error) {
	// PlannerChainを初期化
	var chain *PlannerChain
	if plannerType.Name == PlannerTypeRandom.Name {
		chain = NewRandomPlanner(gc.Tile(width), gc.Tile(height), seed)
	} else {
		chain = plannerType.PlannerFunc(gc.Tile(width), gc.Tile(height), seed)
	}

	// RawMasterを設定
	if world.Resources != nil && world.Resources.RawMaster != nil {
		if rawMaster, ok := world.Resources.RawMaster.(*raw.Master); ok {
			chain.PlanData.RawMaster = rawMaster
		}
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
	plan, err := chain.PlanData.BuildPlan()
	if err != nil {
		return nil, fmt.Errorf("EntityPlan構築エラー: %w", err)
	}

	// 計画の妥当性と接続性をチェック
	if err := plan.Validate(); err != nil {
		return nil, fmt.Errorf("計画検証エラー: %w", err)
	}

	return plan, nil
}

// isConnectivityError は接続性エラーかどうかを判定する
func isConnectivityError(err error) bool {
	return errors.Is(err, ErrConnectivity) || errors.Is(err, ErrPlayerPlacement)
}
