package levelgen

import (
	"errors"
	"fmt"
	"log"

	"github.com/kijimaD/ruins/lib/mapplaner"
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
func NewLevel(world w.World, width gc.Tile, height gc.Tile, seed uint64, plannerType mapplaner.PlannerType) (resources.Level, error) {
	log.Printf("NewLevel開始: PlannerType=%s UseFixedPortalPos=%v", plannerType.Name, plannerType.UseFixedPortalPos)

	var chain *mapplaner.PlannerChain
	var playerX, playerY int

	// mapspawner方式でMapPlanを生成し、Levelをスポーンする
	chain = createPlannerChain(plannerType, width, height, seed)

	// BuildPlanAndSpawnでMapPlanからLevelまで一括生成
	level, err := mapspawner.BuildPlanAndSpawn(world, chain, plannerType)
	if err != nil {
		return resources.Level{}, fmt.Errorf("MapPlanからのLevel生成に失敗: %w", err)
	}

	// MapPlanからプレイヤー位置を取得
	plan, err := mapspawner.BuildPlan(chain)
	if err != nil {
		return resources.Level{}, fmt.Errorf("MapPlan生成に失敗: %w", err)
	}

	// プレイヤー位置を取得
	px, py, hasPlayerPos := plan.GetPlayerStartPosition()
	if !hasPlayerPos {
		return resources.Level{}, fmt.Errorf("MapPlanにプレイヤー開始位置が設定されていません")
	}
	playerX, playerY = px, py

	// プレイヤーを移動する
	if err := worldhelper.MovePlayerToPosition(world, playerX, playerY); err != nil {
		return resources.Level{}, fmt.Errorf("プレイヤー移動エラー: %w", err)
	}

	return level, nil
}

// createPlannerChain は指定されたプランナータイプに応じてプランナーチェーンを作成する
func createPlannerChain(plannerType mapplaner.PlannerType, width gc.Tile, height gc.Tile, seed uint64) *mapplaner.PlannerChain {
	// ランダム選択の場合は特別処理
	if plannerType.Name == mapplaner.PlannerTypeRandom.Name {
		return mapplaner.NewRandomPlanner(width, height, seed)
	}

	return plannerType.PlannerFunc(width, height, seed)
}
