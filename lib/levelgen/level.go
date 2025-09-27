package levelgen

import (
	"errors"
	"fmt"
	"log"

	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/mapspawner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"

	gc "github.com/kijimaD/ruins/lib/components"
)

// エラーメッセージ
var (
	ErrPlayerStartNotFound = errors.New("プレイヤーのスタート位置が見つかりません")
	ErrMapGenerationFailed = errors.New("マップ生成に失敗しました")
)

// NewLevel は新規に階層を生成する。
// mapplanerとmapspawnerを統合して完全なゲームレベルを生成する
func NewLevel(world w.World, width gc.Tile, height gc.Tile, seed uint64, plannerType mapplanner.PlannerType) (resources.Level, error) {
	log.Printf("NewLevel開始: PlannerType=%s UseFixedPortalPos=%v", plannerType.Name, plannerType.UseFixedPortalPos)

	var chain *mapplanner.PlannerChain
	var playerX, playerY int

	// mapspawner方式でEntityPlanを生成し、Levelをスポーンする
	chain = createPlannerChain(plannerType, width, height, seed)

	// BuildPlanAndSpawnでEntityPlanからLevelまで一括生成（プレイヤー位置も取得）
	level, px, py, err := mapspawner.BuildPlanAndSpawn(world, chain, plannerType)
	if err != nil {
		return resources.Level{}, fmt.Errorf("EntityPlanからのLevel生成に失敗: %w", err)
	}
	playerX, playerY = px, py

	// プレイヤーを移動する
	if err := worldhelper.MovePlayerToPosition(world, playerX, playerY); err != nil {
		return resources.Level{}, fmt.Errorf("プレイヤー移動エラー: %w", err)
	}

	return level, nil
}

// createPlannerChain は指定されたプランナータイプに応じてプランナーチェーンを作成する
func createPlannerChain(plannerType mapplanner.PlannerType, width gc.Tile, height gc.Tile, seed uint64) *mapplanner.PlannerChain {
	// ランダム選択の場合は特別処理
	if plannerType.Name == mapplanner.PlannerTypeRandom.Name {
		return mapplanner.NewRandomPlanner(width, height, seed)
	}

	return plannerType.PlannerFunc(width, height, seed)
}
