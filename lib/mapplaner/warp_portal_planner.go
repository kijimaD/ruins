// Package mapplanner のワープポータル配置プランナー - 責務分離によりmapspawnerから移動
package mapplanner

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// ワープポータル配置用の定数
const (
	// 配置処理関連
	maxPortalPlacementAttempts = 200 // ワープポータル配置処理の最大試行回数
	// ワープホール関連
	escapePortalInterval = 5 // 帰還ワープホール配置間隔（n階層ごと）
)

// WarpPortalPlanner はワープポータル配置を担当するプランナー
type WarpPortalPlanner struct {
	world       w.World
	plannerType PlannerType
}

// NewWarpPortalPlanner はワープポータルプランナーを作成する
func NewWarpPortalPlanner(world w.World, plannerType PlannerType) *WarpPortalPlanner {
	return &WarpPortalPlanner{
		world:       world,
		plannerType: plannerType,
	}
}

// BuildMeta はワープポータルをMetaPlanに追加する
func (w *WarpPortalPlanner) BuildMeta(buildData *MetaPlan) {
	// プランナーが既にワープポータルを配置済みかどうかを確認
	existingWarpCount := len(buildData.WarpPortals)

	// 進行ワープホールを配置
	if w.plannerType.UseFixedPortalPos {
		if existingWarpCount == 0 {
			// 街の公民館（下部の部屋）の中央にワープポータルを配置
			centerX := int(buildData.Level.TileWidth) / 2
			centerY := int(buildData.Level.TileHeight) / 2
			// 公民館の中央: Y1=centerY+10, Y2=centerY+22 の中央 = centerY+16
			warpX := centerX
			warpY := centerY + 16

			// 小さなマップの場合は範囲内に調整
			maxY := int(buildData.Level.TileHeight) - 1
			if warpY >= maxY {
				warpY = maxY - 1
			}

			buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
				X:    warpX,
				Y:    warpY,
				Type: WarpPortalNext,
			})
		}
	} else {
		// ダンジョンの場合は通常のランダム配置
		for attempt := 0; attempt < maxPortalPlacementAttempts; attempt++ {
			x := buildData.RandomSource.Intn(int(buildData.Level.TileWidth))
			y := buildData.RandomSource.Intn(int(buildData.Level.TileHeight))

			if buildData.IsSpawnableTile(w.world, gc.Tile(x), gc.Tile(y)) {
				buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
					X:    x,
					Y:    y,
					Type: WarpPortalNext,
				})
				break
			}
		}
	}

	// 帰還ワープホール配置（5階層ごと、またはデバッグ用）
	// 既に帰還ワープポータルが存在するかチェック
	hasEscapePortal := false
	for _, portal := range buildData.WarpPortals {
		if portal.Type == WarpPortalEscape {
			hasEscapePortal = true
			break
		}
	}

	if !hasEscapePortal && w.world.Resources.Dungeon != nil && w.world.Resources.Dungeon.Depth%escapePortalInterval == 0 {
		if w.plannerType.UseFixedPortalPos {
			centerX := int(buildData.Level.TileWidth) / 2
			centerY := int(buildData.Level.TileHeight) / 2
			// 図書館（知識が集まる場所）に帰還ポータルを配置
			escapeX := centerX - 3
			escapeY := centerY - 15 // 図書館の中心

			buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
				X:    escapeX,
				Y:    escapeY,
				Type: WarpPortalEscape,
			})
		} else {
			// ダンジョンの場合は通常のランダム配置
			for attempt := 0; attempt < maxPortalPlacementAttempts; attempt++ {
				x := buildData.RandomSource.Intn(int(buildData.Level.TileWidth))
				y := buildData.RandomSource.Intn(int(buildData.Level.TileHeight))

				if buildData.IsSpawnableTile(w.world, gc.Tile(x), gc.Tile(y)) {
					buildData.WarpPortals = append(buildData.WarpPortals, WarpPortal{
						X:    x,
						Y:    y,
						Type: WarpPortalEscape,
					})
					break
				}
			}
		}
	}
}
