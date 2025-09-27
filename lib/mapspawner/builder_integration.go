package mapspawner

import (
	"fmt"
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	mapplanner "github.com/kijimaD/ruins/lib/mapplaner"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
)

// BuildPlanAndSpawn はPlannerChainを実行してMapPlanを生成し、Levelをスポーンする
func BuildPlanAndSpawn(world w.World, chain *mapplanner.PlannerChain, plannerType mapplanner.PlannerType) (resources.Level, error) {
	// プランナーチェーンを実行
	chain.Build()

	// PlanDataからMapPlanを構築
	plan, err := BuildPlanFromTiles(&chain.PlanData)
	if err != nil {
		return resources.Level{}, fmt.Errorf("MapPlan構築エラー: %w", err)
	}

	// プランナー設定に基づいてNPCとアイテムをMapPlanに追加
	if plannerType.SpawnEnemies {
		if err := addNPCsToPlan(world, chain, plan); err != nil {
			return resources.Level{}, fmt.Errorf("NPC追加エラー: %w", err)
		}
	}

	if plannerType.SpawnItems {
		if err := addItemsToPlan(world, chain, plan); err != nil {
			return resources.Level{}, fmt.Errorf("アイテム追加エラー: %w", err)
		}
	}

	// プランナータイプに応じてワープポータルを配置
	addWarpPortalsToPlan(world, chain, plan, plannerType)

	// プランナータイプに応じて固定Props配置を追加
	if err := addFixedPropsToPlan(chain, plan, plannerType); err != nil {
		return resources.Level{}, fmt.Errorf("固定Props追加エラー: %w", err)
	}

	// MapPlanからLevelをスポーン
	level, err := SpawnLevel(world, plan)
	if err != nil {
		return resources.Level{}, fmt.Errorf("level生成エラー: %w", err)
	}

	return level, nil
}

// BuildPlan はPlannerChainを実行してMapPlanを生成する
func BuildPlan(chain *mapplanner.PlannerChain) (*mapplanner.MapPlan, error) {
	// プランナーチェーンを実行
	chain.Build()

	// PlanDataからMapPlanを構築
	plan, err := BuildPlanFromTiles(&chain.PlanData)
	if err != nil {
		return nil, fmt.Errorf("MapPlan構築エラー: %w", err)
	}

	return plan, nil
}

// NPC/アイテム配置用の定数（levelgenから移動）
const (
	// NPC生成関連
	baseNPCCount    = 5   // NPC生成の基本数
	randomNPCCount  = 5   // NPC生成のランダム追加数（0-4の範囲）
	maxNPCFailCount = 200 // NPC生成の最大失敗回数

	// アイテム配置関連
	baseItemCount           = 2  // 通常アイテム配置の基本数
	randomItemCount         = 3  // 通常アイテム配置のランダム追加数（0-2の範囲）
	itemIncreaseDepth       = 5  // アイテム数増加の深度しきい値
	rareItemProbability     = 30 // レアアイテム配置確率（%）
	deepRareItemDepth       = 10 // 深い階層の判定深度
	deepRareItemProbability = 20 // 深い階層でのレアアイテム複数配置確率（%）

	// 配置処理関連
	maxPlacementAttempts = 200 // 配置処理の最大試行回数

	// ワープホール関連
	escapePortalInterval = 5 // 帰還ワープホール配置間隔（n階層ごと）
)

// addNPCsToPlan はMapPlanにNPCを追加する
func addNPCsToPlan(world w.World, chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan) error {
	failCount := 0
	total := baseNPCCount + chain.PlanData.RandomSource.Intn(randomNPCCount)
	successCount := 0

	for successCount < total && failCount <= maxNPCFailCount {
		tx := gc.Tile(chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileWidth)))
		ty := gc.Tile(chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileHeight)))

		if !chain.PlanData.IsSpawnableTile(world, tx, ty) {
			failCount++
			continue
		}

		// NPCタイプを選択（現在は固定、将来的にはテーブル化）
		npcType := "火の玉" // TODO: テーブルで選ぶ
		plan.AddNPC(int(tx), int(ty), npcType)

		successCount++
		failCount = 0
	}

	if failCount > maxNPCFailCount {
		return fmt.Errorf("NPC配置の試行回数が上限に達しました。配置数: %d/%d", successCount, total)
	}

	return nil
}

// addItemsToPlan はMapPlanにアイテムを追加する
func addItemsToPlan(world w.World, chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan) error {
	// 利用可能なアイテムリスト
	// TODO: テーブル化・レアリティ考慮する
	availableItems := []string{
		"回復薬",
		"回復スプレー",
		"手榴弾",
		"上級回復薬",
		"ルビー原石",
	}

	// レアアイテム
	rareItems := []string{
		"上級回復薬",
		"ルビー原石",
	}

	// 通常アイテムの配置数（階層の深度に応じて調整）
	normalItemCount := baseItemCount + chain.PlanData.RandomSource.Intn(randomItemCount)
	if world.Resources.Dungeon != nil && world.Resources.Dungeon.Depth > itemIncreaseDepth {
		normalItemCount++ // 深い階層ではアイテム数を増加
	}

	// レアアイテムの配置数（低確率）
	rareItemCount := 0
	if chain.PlanData.RandomSource.Intn(100) < rareItemProbability {
		rareItemCount = 1
		if world.Resources.Dungeon != nil && world.Resources.Dungeon.Depth > deepRareItemDepth && chain.PlanData.RandomSource.Intn(100) < deepRareItemProbability {
			rareItemCount = 2
		}
	}

	// 通常アイテムを配置
	if err := addItemsOfTypeToPlan(chain, plan, availableItems, normalItemCount); err != nil {
		return err
	}

	// レアアイテムを配置
	if rareItemCount > 0 {
		if err := addItemsOfTypeToPlan(chain, plan, rareItems, rareItemCount); err != nil {
			return err
		}
	}

	return nil
}

// addItemsOfTypeToPlan は指定された数のアイテムをMapPlanに追加する
func addItemsOfTypeToPlan(chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan, itemList []string, count int) error {
	failCount := 0
	successCount := 0

	for successCount < count {
		if failCount > maxPlacementAttempts {
			log.Printf("アイテム配置の試行回数が上限に達しました。配置数: %d/%d", successCount, count)
			break
		}

		// ランダムな位置を選択
		x := gc.Tile(chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileWidth)))
		y := gc.Tile(chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileHeight)))

		// スポーン可能な位置かチェック（world不要の代替チェック）
		if !isValidItemPosition(chain, x, y) {
			failCount++
			continue
		}

		// アイテム名をランダム選択
		itemName := itemList[chain.PlanData.RandomSource.Intn(len(itemList))]

		// MapPlanにアイテムを追加
		plan.AddItem(int(x), int(y), itemName)

		successCount++
		failCount = 0
	}
	return nil
}

// isValidItemPosition はアイテム配置に適した位置かチェックする（簡易版）
func isValidItemPosition(chain *mapplanner.PlannerChain, x, y gc.Tile) bool {
	tileIdx := chain.PlanData.Level.XYTileIndex(x, y)
	if int(tileIdx) >= len(chain.PlanData.Tiles) {
		return false
	}

	tile := chain.PlanData.Tiles[tileIdx]
	// 床、ワープタイルに配置可能
	return tile == mapplanner.TileFloor || tile == mapplanner.TileWarpNext || tile == mapplanner.TileWarpEscape
}

// addFixedPropsToPlan はプランナータイプに応じて固定Props配置をMapPlanに追加する
func addFixedPropsToPlan(chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan, plannerType mapplanner.PlannerType) error {
	// 町タイプの場合は固定Props配置を追加
	if plannerType.Name == mapplanner.PlannerTypeTown.Name {
		return addTownPropsToPlan(chain, plan)
	}

	// ダンジョンタイプの場合は固定Props配置を追加（必要に応じて実装）
	// TODO: 必要に応じて他のタイプも実装

	return nil
}

// addWarpPortalsToPlan はプランナータイプに応じてワープポータルをMapPlanに追加する
func addWarpPortalsToPlan(world w.World, chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan, plannerType mapplanner.PlannerType) {
	// プランナーが既にワープホールを配置済みかどうかを確認
	// StringMapPlannerベースのプランナーはentityMapでワープホールを完全に配置するため、
	// mapspawnerでの追加配置は不要
	existingWarpCount := 0
	for _, tile := range chain.PlanData.Tiles {
		if tile == mapplanner.TileWarpNext {
			existingWarpCount++
		}
	}

	// 進行ワープホールを配置
	if plannerType.UseFixedPortalPos {
		if existingWarpCount > 0 {
			// プランナーが既にワープホールを配置済みの場合は追加配置しない
			log.Printf("ワープポータル配置済み確認: 既存数=%d PlannerType:%s (mapspawner)",
				existingWarpCount, plannerType.Name)
		} else {
			// 街の公民館（下部の部屋）の中央にワープポータルを配置
			centerX := int(chain.PlanData.Level.TileWidth) / 2
			centerY := int(chain.PlanData.Level.TileHeight) / 2
			// 公民館の中央: Y1=centerY+10, Y2=centerY+22 の中央 = centerY+16
			warpX := centerX
			warpY := centerY + 16

			// 小さなマップの場合は範囲内に調整
			maxY := int(chain.PlanData.Level.TileHeight) - 1
			if warpY >= maxY {
				warpY = maxY - 1
			}

			plan.AddWarpNext(warpX, warpY)
		}
	} else {
		// ダンジョンの場合は通常のランダム配置
		for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
			x := chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileWidth))
			y := chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileHeight))

			if chain.PlanData.IsSpawnableTile(world, gc.Tile(x), gc.Tile(y)) {
				plan.AddWarpNext(x, y)
				break
			}
		}
	}

	// 帰還ワープホール配置（5階層ごと、またはデバッグ用）
	if world.Resources.Dungeon != nil && world.Resources.Dungeon.Depth%escapePortalInterval == 0 {
		if plannerType.UseFixedPortalPos {
			centerX := int(chain.PlanData.Level.TileWidth) / 2
			centerY := int(chain.PlanData.Level.TileHeight) / 2
			// 図書館（知識が集まる場所）に帰還ポータルを配置
			escapeX := centerX - 3
			escapeY := centerY - 15 // 図書館の中心
			plan.AddWarpEscape(escapeX, escapeY)
		} else {
			// ダンジョンの場合は通常のランダム配置
			for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
				x := chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileWidth))
				y := chain.PlanData.RandomSource.Intn(int(chain.PlanData.Level.TileHeight))

				if chain.PlanData.IsSpawnableTile(world, gc.Tile(x), gc.Tile(y)) {
					plan.AddWarpEscape(x, y)
					break
				}
			}
		}
	}
}

// addTownPropsToPlan は町用の固定Props配置をMapPlanに追加する
func addTownPropsToPlan(chain *mapplanner.PlannerChain, plan *mapplanner.MapPlan) error {
	centerX := int(chain.PlanData.Level.TileWidth) / 2
	centerY := int(chain.PlanData.Level.TileHeight) / 2

	// 図書館の家具配置
	libraryProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBookshelf, -8, -19}, // 北壁沿い
		{gc.PropTypeBookshelf, -6, -19}, // 北壁沿い
		{gc.PropTypeTable, -7, -17},     // 閲覧机
		{gc.PropTypeChair, -7, -16},     // 閲覧用椅子
		{gc.PropTypeTable, -4, -15},     // 学習机
	}

	for _, prop := range libraryProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 学校の家具配置
	schoolProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBookshelf, 6, -19},  // 北壁
		{gc.PropTypeBookshelf, 8, -19},  // 北壁
		{gc.PropTypeBookshelf, 10, -19}, // 北壁
		{gc.PropTypeBookshelf, 5, -16},  // 西壁
		{gc.PropTypeBookshelf, 11, -16}, // 東壁
		{gc.PropTypeTable, 8, -15},      // 教卓
		{gc.PropTypeChair, 8, -14},      // 教師用椅子
	}

	for _, prop := range schoolProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 住民の家1の家具配置
	house1Props := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 13, -7},   // 寝室
		{gc.PropTypeTable, 15, -5}, // 食事台
		{gc.PropTypeChair, 15, -4}, // 食事用椅子
		{gc.PropTypeChair, 16, -5}, // 食事用椅子
	}

	for _, prop := range house1Props {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 住民の家2の家具配置
	house2Props := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 14, 2},   // 寝室
		{gc.PropTypeTable, 16, 4}, // 食事台
		{gc.PropTypeChair, 16, 5}, // 食事用椅子
	}

	for _, prop := range house2Props {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 公民館の座席配置
	hallProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeChair, -6, 12}, // 集会用座席
		{gc.PropTypeChair, -4, 12}, // 集会用座席
		{gc.PropTypeChair, 4, 12},  // 集会用座席
		{gc.PropTypeChair, 6, 12},  // 集会用座席
	}

	for _, prop := range hallProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 事務所の家具配置
	officeProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeBed, 12, 13},       // 休憩用ベッド
		{gc.PropTypeTable, 14, 15},     // 事務机
		{gc.PropTypeChair, 14, 16},     // 事務用椅子
		{gc.PropTypeBookshelf, 18, 14}, // 書類棚
	}

	for _, prop := range officeProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	// 市場の露店（簡略化して一部のみ配置）
	marketProps := []struct {
		propType gc.PropType
		offsetX  int
		offsetY  int
	}{
		{gc.PropTypeTable, -12, 5}, // 露店1
		{gc.PropTypeTable, -9, 5},  // 露店2
		{gc.PropTypeTable, -6, 5},  // 露店3
	}

	for _, prop := range marketProps {
		x := centerX + prop.offsetX
		y := centerY + prop.offsetY
		if isValidPropPosition(chain, gc.Tile(x), gc.Tile(y)) {
			plan.AddProp(x, y, prop.propType)
		}
	}

	return nil
}

// isValidPropPosition はProp配置に適した位置かチェックする
func isValidPropPosition(chain *mapplanner.PlannerChain, x, y gc.Tile) bool {
	// 範囲チェック
	if x < 0 || x >= chain.PlanData.Level.TileWidth || y < 0 || y >= chain.PlanData.Level.TileHeight {
		return false
	}

	tileIdx := chain.PlanData.Level.XYTileIndex(x, y)
	if int(tileIdx) >= len(chain.PlanData.Tiles) {
		return false
	}

	tile := chain.PlanData.Tiles[tileIdx]
	// 床タイルにのみ配置可能
	return tile == mapplanner.TileFloor
}
