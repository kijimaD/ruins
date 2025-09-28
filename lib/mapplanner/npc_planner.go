// Package mapplanner のNPC配置プランナー - 責務分離によりmapspawnerから移動
package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// NPC配置用の定数
const (
	// NPC生成関連
	baseNPCCount    = 5   // NPC生成の基本数
	randomNPCCount  = 5   // NPC生成のランダム追加数（0-4の範囲）
	maxNPCFailCount = 200 // NPC生成の最大失敗回数
)

// NPCSpec はNPC配置仕様を表す
type NPCSpec struct {
	X       int    // X座標
	Y       int    // Y座標
	NPCType string // NPCタイプ
}

// NPCPlanner はNPC配置を担当するプランナー
type NPCPlanner struct {
	world       w.World
	plannerType PlannerType
}

// NewNPCPlanner はNPCプランナーを作成する
func NewNPCPlanner(world w.World, plannerType PlannerType) *NPCPlanner {
	return &NPCPlanner{
		world:       world,
		plannerType: plannerType,
	}
}

// PlanMeta はNPC配置情報をMetaPlanに追加する
func (n *NPCPlanner) PlanMeta(planData *MetaPlan) {
	if !n.plannerType.SpawnEnemies {
		return // 敵をスポーンしない設定の場合は何もしない
	}

	failCount := 0
	total := baseNPCCount + planData.RandomSource.Intn(randomNPCCount)
	successCount := 0

	// NPCsフィールドが存在しない場合は初期化
	if planData.NPCs == nil {
		planData.NPCs = []NPCSpec{}
	}

	for successCount < total && failCount <= maxNPCFailCount {
		tx := gc.Tile(planData.RandomSource.Intn(int(planData.Level.TileWidth)))
		ty := gc.Tile(planData.RandomSource.Intn(int(planData.Level.TileHeight)))

		if !planData.IsSpawnableTile(n.world, tx, ty) {
			failCount++
			continue
		}

		// NPCタイプを選択（現在は固定、将来的にはテーブル化）
		npcType := "火の玉" // TODO: テーブルで選ぶ

		planData.NPCs = append(planData.NPCs, NPCSpec{
			X:       int(tx),
			Y:       int(ty),
			NPCType: npcType,
		})

		successCount++
		failCount = 0
	}

	if failCount > maxNPCFailCount {
		// エラーは記録するが、エラーを返さずに部分的な配置で続行
		fmt.Printf("NPC配置の試行回数が上限に達しました。配置数: %d/%d\n", successCount, total)
	}
}
