// Package mapplanner の会話NPC配置プランナー
// 敵NPCとは別に、会話可能なNPCを配置する
package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// ConversationNPCPlanner は会話可能NPCの配置を担当するプランナー
type ConversationNPCPlanner struct {
	world w.World
}

// NewConversationNPCPlanner は会話NPCプランナーを作成する
func NewConversationNPCPlanner(world w.World) *ConversationNPCPlanner {
	return &ConversationNPCPlanner{
		world: world,
	}
}

// PlanMeta は会話可能NPCをMetaPlanに追加する
func (c *ConversationNPCPlanner) PlanMeta(planData *MetaPlan) {
	// NPCsフィールドが存在しない場合は初期化
	if planData.NPCs == nil {
		planData.NPCs = []NPCSpec{}
	}

	// プレイヤーの近くに会話NPCを配置
	c.placeDialogueNPCs(planData)
}

// placeDialogueNPCs は会話可能NPCをプレイヤーの近くに配置する
func (c *ConversationNPCPlanner) placeDialogueNPCs(planData *MetaPlan) {
	// プレイヤーの開始位置を取得
	playerX, playerY, hasPlayer := planData.GetPlayerStartPosition()
	if !hasPlayer {
		return
	}

	// プレイヤーの近くに配置する候補位置（プレイヤーから1-5マス離れた位置）
	candidateOffsets := []struct{ dx, dy int }{
		// 1マス離れた位置
		{1, 0}, {0, 1}, {-1, 0}, {0, -1},
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1},
		// 2マス離れた位置
		{2, 0}, {0, 2}, {-2, 0}, {0, -2},
		{2, 2}, {-2, 2}, {2, -2}, {-2, -2},
		{2, 1}, {1, 2}, {-2, 1}, {-1, 2},
		// 3マス離れた位置
		{3, 0}, {0, 3}, {-3, 0}, {0, -3},
		{3, 3}, {-3, 3}, {3, -3}, {-3, -3},
		// 4-5マス離れた位置
		{4, 0}, {0, 4}, {5, 0}, {0, 5},
	}

	// 配置する会話NPCのタイプ一覧
	npcTypes := []string{"老兵"}

	for _, npcType := range npcTypes {
		for _, offset := range candidateOffsets {
			x := playerX + offset.dx
			y := playerY + offset.dy

			// マップ範囲内かチェック
			if x < 0 || y < 0 || x >= int(planData.Level.TileWidth) || y >= int(planData.Level.TileHeight) {
				continue
			}

			if planData.IsSpawnableTile(c.world, gc.Tile(x), gc.Tile(y)) {
				planData.NPCs = append(planData.NPCs, NPCSpec{
					X:       x,
					Y:       y,
					NPCType: npcType,
				})
				break
			}
		}
	}
}
