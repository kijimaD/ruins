// Package mapplanner の敵NPC配置プランナー
package mapplanner

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// 敵NPC配置用の定数
const (
	// 敵NPC生成関連
	baseHostileNPCCount    = 5   // 敵NPC生成の基本数
	randomHostileNPCCount  = 5   // 敵NPC生成のランダム追加数（0-4の範囲）
	maxHostileNPCFailCount = 200 // 敵NPC生成の最大失敗回数
)

// NPCSpec はNPC配置仕様を表す
type NPCSpec struct {
	X       int    // X座標
	Y       int    // Y座標
	NPCType string // NPCタイプ
}

// HostileNPCPlanner は敵NPC配置を担当するプランナー
// 会話可能なNPCはConversationNPCPlannerで配置する
type HostileNPCPlanner struct {
	world       w.World
	plannerType PlannerType
}

// NewHostileNPCPlanner は敵NPCプランナーを作成する
func NewHostileNPCPlanner(world w.World, plannerType PlannerType) *HostileNPCPlanner {
	return &HostileNPCPlanner{
		world:       world,
		plannerType: plannerType,
	}
}

// PlanMeta は敵NPC配置情報をMetaPlanに追加する
func (n *HostileNPCPlanner) PlanMeta(planData *MetaPlan) {
	// NPCsフィールドが存在しない場合は初期化
	if planData.NPCs == nil {
		planData.NPCs = []NPCSpec{}
	}

	// 敵NPCの配置
	if !n.plannerType.SpawnEnemies {
		return // 敵をスポーンしない設定の場合は何もしない
	}

	failCount := 0
	total := baseHostileNPCCount + planData.RandomSource.Intn(randomHostileNPCCount)
	successCount := 0

	for successCount < total && failCount <= maxHostileNPCFailCount {
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

	if failCount > maxHostileNPCFailCount {
		// エラーは記録するが、エラーを返さずに部分的な配置で続行
		fmt.Printf("敵NPC配置の試行回数が上限に達しました。配置数: %d/%d\n", successCount, total)
	}
}
