package mapplanner

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

// アイテム配置用の定数
const (
	// アイテム配置関連
	baseItemCount           = 2  // 通常アイテム配置の基本数
	randomItemCount         = 3  // 通常アイテム配置のランダム追加数（0-2の範囲）
	itemIncreaseDepth       = 5  // アイテム数増加の深度しきい値
	rareItemProbability     = 30 // レアアイテム配置確率（%）
	deepRareItemDepth       = 10 // 深い階層の判定深度
	deepRareItemProbability = 20 // 深い階層でのレアアイテム複数配置確率（%）

	// 配置処理関連
	maxItemPlacementAttempts = 200 // アイテム配置処理の最大試行回数
)

// ItemSpec はアイテム配置仕様を表す
type ItemSpec struct {
	X        int    // X座標
	Y        int    // Y座標
	ItemName string // アイテム名
}

// ItemPlanner はアイテム配置を担当するプランナー
type ItemPlanner struct {
	world       w.World
	plannerType PlannerType
}

// NewItemPlanner はアイテムプランナーを作成する
func NewItemPlanner(world w.World, plannerType PlannerType) *ItemPlanner {
	return &ItemPlanner{
		world:       world,
		plannerType: plannerType,
	}
}

// BuildMeta はアイテム配置情報をMetaPlanに追加する
func (i *ItemPlanner) BuildMeta(buildData *MetaPlan) {
	if !i.plannerType.SpawnItems {
		return // アイテムをスポーンしない設定の場合は何もしない
	}

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

	// Itemsフィールドが存在しない場合は初期化
	if buildData.Items == nil {
		buildData.Items = []ItemSpec{}
	}

	// 通常アイテムの配置数（階層の深度に応じて調整）
	normalItemCount := baseItemCount + buildData.RandomSource.Intn(randomItemCount)
	if i.world.Resources.Dungeon != nil && i.world.Resources.Dungeon.Depth > itemIncreaseDepth {
		normalItemCount++ // 深い階層ではアイテム数を増加
	}

	// レアアイテムの配置数（低確率）
	rareItemCount := 0
	if buildData.RandomSource.Intn(100) < rareItemProbability {
		rareItemCount = 1
		if i.world.Resources.Dungeon != nil && i.world.Resources.Dungeon.Depth > deepRareItemDepth && buildData.RandomSource.Intn(100) < deepRareItemProbability {
			rareItemCount = 2
		}
	}

	// 通常アイテムを配置
	i.addItemsOfType(buildData, availableItems, normalItemCount)

	// レアアイテムを配置
	if rareItemCount > 0 {
		i.addItemsOfType(buildData, rareItems, rareItemCount)
	}
}

// addItemsOfType は指定された数のアイテムをMetaPlanに追加する
func (i *ItemPlanner) addItemsOfType(buildData *MetaPlan, itemList []string, count int) {
	failCount := 0
	successCount := 0

	for successCount < count {
		if failCount > maxItemPlacementAttempts {
			log.Printf("アイテム配置の試行回数が上限に達しました。配置数: %d/%d", successCount, count)
			break
		}

		// ランダムな位置を選択
		x := gc.Tile(buildData.RandomSource.Intn(int(buildData.Level.TileWidth)))
		y := gc.Tile(buildData.RandomSource.Intn(int(buildData.Level.TileHeight)))

		// スポーン可能な位置かチェック
		if !i.isValidItemPosition(buildData, x, y) {
			failCount++
			continue
		}

		// アイテム名をランダム選択
		itemName := itemList[buildData.RandomSource.Intn(len(itemList))]

		// MetaPlanにアイテムを追加
		buildData.Items = append(buildData.Items, ItemSpec{
			X:        int(x),
			Y:        int(y),
			ItemName: itemName,
		})

		successCount++
		failCount = 0
	}
}

// isValidItemPosition はアイテム配置に適した位置かチェックする
func (i *ItemPlanner) isValidItemPosition(buildData *MetaPlan, x, y gc.Tile) bool {
	tileIdx := buildData.Level.XYTileIndex(x, y)
	if int(tileIdx) >= len(buildData.Tiles) {
		return false
	}

	tile := buildData.Tiles[tileIdx]
	// 歩行可能なタイルに配置可能
	return tile.Walkable
}
