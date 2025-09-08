package mapbuilder

import (
	"errors"
	"fmt"
	"log"

	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"

	gc "github.com/kijimaD/ruins/lib/components"
)

const (
	// マップ生成関連
	maxMapGenerationAttempts = 10  // マップ生成の最大試行回数
	maxPlacementAttempts     = 200 // 配置処理の最大試行回数

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

	// ワープホール関連
	escapePortalInterval = 5 // 帰還ワープホール配置間隔（n階層ごと）

	// 壁スプライト番号
	spriteWallTop         = 10 // 上壁
	spriteWallBottom      = 11 // 下壁
	spriteWallLeft        = 12 // 左壁
	spriteWallRight       = 13 // 右壁
	spriteWallTopLeft     = 14 // 左上角
	spriteWallTopRight    = 15 // 右上角
	spriteWallBottomLeft  = 16 // 左下角
	spriteWallBottomRight = 17 // 右下角
	spriteWallGeneric     = 1  // 汎用壁
)

// エラーメッセージ
var (
	ErrNPCGenerationFailed     = errors.New("NPCの生成に失敗しました")
	ErrPlayerStartNotFound     = errors.New("プレイヤーのスタート位置が見つかりません")
	ErrMapGenerationFailed     = errors.New("マップ生成に失敗しました")
	ErrWarpPortalPlaceFailed   = errors.New("ワープポータルの配置に失敗しました")
	ErrEscapePortalPlaceFailed = errors.New("脱出ポータルの配置に失敗しました")
)

// NewLevel は新規に階層を生成する。
// 階層を初期化するので、具体的なコードであり、その分参照を多く含んでいる。循環参照を防ぐためにこの関数はLevel構造体とは同じpackageに属していない。
func NewLevel(world w.World, width gc.Tile, height gc.Tile, seed uint64, builderType BuilderType) (resources.Level, error) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	var chain *BuilderChain
	var playerX, playerY int

	// 接続性検証付きマップ生成（最大試行回数まで再試行）
	validMap := false
	for attempt := 0; attempt < maxMapGenerationAttempts && !validMap; attempt++ {
		// シードを少しずつ変えて再生成
		currentSeed := seed + uint64(attempt)
		chain = createBuilderChain(builderType, width, height, currentSeed)
		chain.Build()

		// プレイヤーのスタート位置を見つける（最初にスポーン可能な位置）
		var err error
		playerX, playerY, err = findPlayerStartPosition(&chain.BuildData, world)
		if err != nil {
			continue // スタート位置が見つからない場合は再生成
		}

		// 接続性を検証（ポータル配置後）
		validMap = validateMapWithPortals(chain, world, gameResources, playerX, playerY)

		if !validMap && attempt < maxMapGenerationAttempts-1 {
			log.Printf("マップ生成試行 %d: 接続性検証失敗、再生成します", attempt+1)
		}
	}

	if !validMap {
		log.Printf("警告: %d回の試行後も完全接続マップを生成できませんでした。部分的接続マップを使用します", maxMapGenerationAttempts)
	}

	// ポータルは既にvalidateMapWithPortals内で配置済み
	// フィールドに操作対象キャラを配置する（事前に見つけた位置を使用）
	if err := worldhelper.SpawnOperator(world, playerX, playerY); err != nil {
		return resources.Level{}, fmt.Errorf("プレイヤー生成エラー: %w", err)
	}

	// フィールドにNPCを生成する
	if err := spawnNPCs(world, chain); err != nil {
		return resources.Level{}, err
	}

	// フィールドアイテムを生成する
	if err := spawnFieldItems(world, chain); err != nil {
		return resources.Level{}, err
	}

	// tilesを元にタイルエンティティを生成する
	for _i, t := range chain.BuildData.Tiles {
		i := resources.TileIdx(_i)
		x, y := chain.BuildData.Level.XYTileCoord(i)
		switch t {
		case TileFloor:
			entity, err := worldhelper.SpawnFloor(world, gc.Tile(x), gc.Tile(y))
			if err != nil {
				return resources.Level{}, fmt.Errorf("床の生成に失敗 (x=%d, y=%d): %w", int(x), int(y), err)
			}
			chain.BuildData.Level.Entities[i] = entity
		case TileWall:
			// 近傍8タイル（直交・斜め）にフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentAnyFloor(i) {
				// 壁タイプを判定してスプライト番号を決定
				wallType := chain.BuildData.GetWallType(i)
				spriteNumber := getSpriteNumberForWallType(wallType)
				entity, err := worldhelper.SpawnFieldWallWithSprite(world, gc.Tile(x), gc.Tile(y), spriteNumber)
				if err != nil {
					return resources.Level{}, fmt.Errorf("壁の生成に失敗 (x=%d, y=%d): %w", int(x), int(y), err)
				}
				chain.BuildData.Level.Entities[i] = entity
			}
		case TileWarpNext:
			entity, err := worldhelper.SpawnFieldWarpNext(world, gc.Tile(x), gc.Tile(y))
			if err != nil {
				return resources.Level{}, fmt.Errorf("進行ワープホールの生成に失敗 (x=%d, y=%d): %w", int(x), int(y), err)
			}
			chain.BuildData.Level.Entities[i] = entity
		case TileWarpEscape:
			entity, err := worldhelper.SpawnFieldWarpEscape(world, gc.Tile(x), gc.Tile(y))
			if err != nil {
				return resources.Level{}, fmt.Errorf("脱出ワープホールの生成に失敗 (x=%d, y=%d): %w", int(x), int(y), err)
			}
			chain.BuildData.Level.Entities[i] = entity
		}
	}

	return chain.BuildData.Level, nil
}

// spawnNPCs はフィールドにNPCを生成する
func spawnNPCs(world w.World, chain *BuilderChain) error {
	failCount := 0
	total := baseNPCCount + chain.BuildData.RandomSource.Intn(randomNPCCount)
	successCount := 0

	for {
		if failCount > maxNPCFailCount {
			return ErrNPCGenerationFailed
		}
		tx := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
		ty := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
		if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
			failCount++
			continue
		}
		if err := worldhelper.SpawnNPC(
			world,
			gc.Pixel(int(tx)*int(consts.TileSize)+int(consts.TileSize/2)),
			gc.Pixel(int(ty)*int(consts.TileSize)+int(consts.TileSize/2)),
		); err != nil {
			log.Printf("NPC生成に失敗: %v", err)
			failCount++
			continue
		}
		successCount++
		failCount = 0
		if successCount > total {
			break
		}
	}

	return nil
}

// spawnFieldItems はフィールドにアイテムを配置する
func spawnFieldItems(world w.World, chain *BuilderChain) error {
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

	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	// 通常アイテムの配置数（階層の深度に応じて調整）
	normalItemCount := baseItemCount + chain.BuildData.RandomSource.Intn(randomItemCount)
	if gameResources.Depth > itemIncreaseDepth {
		normalItemCount++ // 深い階層ではアイテム数を増加
	}

	// レアアイテムの配置数（低確率）
	rareItemCount := 0
	if chain.BuildData.RandomSource.Intn(100) < rareItemProbability {
		rareItemCount = 1
		if gameResources.Depth > deepRareItemDepth && chain.BuildData.RandomSource.Intn(100) < deepRareItemProbability {
			rareItemCount = 2
		}
	}

	// 通常アイテムを配置
	if err := spawnItems(world, chain, availableItems, normalItemCount); err != nil {
		return err
	}

	// レアアイテムを配置
	if rareItemCount > 0 {
		if err := spawnItems(world, chain, rareItems, rareItemCount); err != nil {
			return err
		}
	}

	return nil
}

// spawnItems は指定された数のアイテムを配置する
func spawnItems(world w.World, chain *BuilderChain, itemList []string, count int) error {
	failCount := 0
	successCount := 0

	for successCount < count {
		if failCount > maxPlacementAttempts {
			log.Printf("アイテム配置の試行回数が上限に達しました。配置数: %d/%d", successCount, count)
			break
		}

		// ランダムな位置を選択
		x := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
		y := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))

		// スポーン可能な位置かチェック
		if !chain.BuildData.IsSpawnableTile(world, x, y) {
			failCount++
			continue
		}

		// アイテム名をランダム選択
		itemName := itemList[chain.BuildData.RandomSource.Intn(len(itemList))]

		// アイテムを配置
		_, err := worldhelper.SpawnFieldItem(world, itemName, x, y)
		if err != nil {
			return fmt.Errorf("フィールドアイテム配置エラー: %w", err)
		}

		successCount++
		failCount = 0
	}
	return nil
}

// getSpriteNumberForWallType は壁タイプに対応するスプライト番号を返す
func getSpriteNumberForWallType(wallType WallType) int {
	switch wallType {
	case WallTypeTop:
		return spriteWallTop // 上壁（下に床がある）
	case WallTypeBottom:
		return spriteWallBottom // 下壁（上に床がある）
	case WallTypeLeft:
		return spriteWallLeft // 左壁（右に床がある）
	case WallTypeRight:
		return spriteWallRight // 右壁（左に床がある）
	case WallTypeTopLeft:
		return spriteWallTopLeft // 左上角（右下に床がある）
	case WallTypeTopRight:
		return spriteWallTopRight // 右上角（左下に床がある）
	case WallTypeBottomLeft:
		return spriteWallBottomLeft // 左下角（右上に床がある）
	case WallTypeBottomRight:
		return spriteWallBottomRight // 右下角（左上に床がある）
	case WallTypeGeneric:
		return spriteWallGeneric // 汎用壁（従来の壁）
	default:
		return 1 // デフォルトは従来の壁
	}
}

// findPlayerStartPosition はプレイヤーのスタート位置を見つける
func findPlayerStartPosition(buildData *BuilderMap, world w.World) (int, int, error) {
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 複数の候補位置を試す
	attempts := []struct{ x, y int }{
		{width / 2, height / 2},         // 中央
		{width / 4, height / 4},         // 左上寄り
		{3 * width / 4, height / 4},     // 右上寄り
		{width / 4, 3 * height / 4},     // 左下寄り
		{3 * width / 4, 3 * height / 4}, // 右下寄り
	}

	// 最適な位置を探す
	for _, pos := range attempts {
		if buildData.IsSpawnableTile(world, gc.Tile(pos.x), gc.Tile(pos.y)) {
			return pos.x, pos.y, nil
		}
	}

	// 全体をランダムに探索
	for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
		x := buildData.RandomSource.Intn(width)
		y := buildData.RandomSource.Intn(height)
		if buildData.IsSpawnableTile(world, gc.Tile(x), gc.Tile(y)) {
			return x, y, nil
		}
	}

	return -1, -1, ErrPlayerStartNotFound // 見つからない場合
}

// validateMapWithPortals はポータルを配置してマップの接続性を検証する
func validateMapWithPortals(chain *BuilderChain, world w.World, gameResources *resources.Dungeon, playerX, playerY int) bool {
	// 進行ワープホールを配置
	warpNextPlaced := false
	for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
		x := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
		y := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
		tileIdx := chain.BuildData.Level.XYTileIndex(x, y)

		if chain.BuildData.IsSpawnableTile(world, x, y) {
			chain.BuildData.Tiles[tileIdx] = TileWarpNext
			warpNextPlaced = true
			break
		}
	}

	if !warpNextPlaced {
		return false // ワープホール配置失敗
	}

	// 帰還ワープホールを配置（5階層ごと）
	escapePortalRequired := gameResources.Depth%escapePortalInterval == 0
	escapePortalPlaced := !escapePortalRequired

	if escapePortalRequired {
		for attempt := 0; attempt < maxPlacementAttempts; attempt++ {
			x := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			y := gc.Tile(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))

			if chain.BuildData.IsSpawnableTile(world, x, y) {
				tileIdx := chain.BuildData.Level.XYTileIndex(x, y)
				chain.BuildData.Tiles[tileIdx] = TileWarpEscape
				escapePortalPlaced = true
				break
			}
		}
	}

	if !escapePortalPlaced {
		return false // 脱出ポータル配置失敗
	}

	// 接続性を検証
	result := chain.ValidateConnectivity(playerX, playerY)

	// プレイヤーのスタート位置が歩行可能で、必要なポータルに到達可能かチェック
	if !result.PlayerStartReachable {
		return false
	}

	if !result.HasReachableWarpPortal() {
		return false // ワープポータルに到達できない
	}

	if escapePortalRequired && !result.HasReachableEscapePortal() {
		return false // 脱出ポータルに到達できない
	}

	return true
}

// createBuilderChain は指定されたビルダータイプに応じてビルダーチェーンを作成する
func createBuilderChain(builderType BuilderType, width gc.Tile, height gc.Tile, seed uint64) *BuilderChain {
	switch builderType {
	case BuilderTypeSmallRoom:
		return NewSmallRoomBuilder(width, height, seed)
	case BuilderTypeBigRoom:
		return NewBigRoomBuilder(width, height, seed)
	case BuilderTypeCave:
		return NewCaveBuilder(width, height, seed)
	case BuilderTypeForest:
		return NewForestBuilder(width, height, seed)
	case BuilderTypeRuins:
		return NewRuinsBuilder(width, height, seed)
	case BuilderTypeRandom:
		fallthrough
	default:
		// デフォルト（BuilderTypeRandomを含む）はランダムビルダー
		return NewRandomBuilder(width, height, seed)
	}
}
