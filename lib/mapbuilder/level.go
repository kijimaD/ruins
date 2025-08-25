package mapbuilder

import (
	"log"

	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"

	gc "github.com/kijimaD/ruins/lib/components"
)

// NewLevel は新規に階層を生成する。
// 階層を初期化するので、具体的なコードであり、その分参照を多く含んでいる。循環参照を防ぐためにこの関数はLevel構造体とは同じpackageに属していない。
func NewLevel(world w.World, width gc.Row, height gc.Col, seed uint64, builderType BuilderType) resources.Level {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)

	var chain *BuilderChain
	var playerX, playerY int

	// 接続性検証付きマップ生成（最大10回まで再試行）
	validMap := false
	for attempt := 0; attempt < 10 && !validMap; attempt++ {
		// シードを少しずつ変えて再生成
		currentSeed := seed + uint64(attempt)
		chain = createBuilderChain(builderType, width, height, currentSeed)
		chain.Build()

		// プレイヤーのスタート位置を見つける（最初にスポーン可能な位置）
		playerX, playerY = findPlayerStartPosition(&chain.BuildData, world)
		if playerX == -1 || playerY == -1 {
			continue // スタート位置が見つからない場合は再生成
		}

		// 接続性を検証（ポータル配置後）
		validMap = validateMapWithPortals(chain, world, gameResources, playerX, playerY)

		if !validMap && attempt < 9 {
			log.Printf("マップ生成試行 %d: 接続性検証失敗、再生成します", attempt+1)
		}
	}

	if !validMap {
		log.Printf("警告: %d回の試行後も完全接続マップを生成できませんでした。部分的接続マップを使用します", 10)
	}

	// ポータルは既にvalidateMapWithPortals内で配置済み
	// フィールドに操作対象キャラを配置する（事前に見つけた位置を使用）
	worldhelper.SpawnOperator(
		world,
		gc.Pixel(playerX*int(consts.TileSize)+int(consts.TileSize)/2),
		gc.Pixel(playerY*int(consts.TileSize)+int(consts.TileSize)/2),
	)
	// フィールドにNPCを生成する
	{
		failCount := 0
		total := 5 + chain.BuildData.RandomSource.Intn(5)
		successCount := 0
		for {
			if failCount > 200 {
				log.Fatal("NPCの生成に失敗した")
			}
			tx := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			ty := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
			if !chain.BuildData.IsSpawnableTile(world, tx, ty) {
				failCount++
				continue
			}
			worldhelper.SpawnNPC(
				world,
				gc.Pixel(int(tx)*int(consts.TileSize)+int(consts.TileSize/2)),
				gc.Pixel(int(ty)*int(consts.TileSize)+int(consts.TileSize/2)),
			)
			successCount++
			failCount = 0
			if successCount > total {
				break
			}
		}
	}

	// tilesを元にタイルエンティティを生成する
	for _i, t := range chain.BuildData.Tiles {
		i := resources.TileIdx(_i)
		x, y := chain.BuildData.Level.XYTileCoord(i)
		switch t {
		case TileFloor:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFloor(world, gc.Row(x), gc.Col(y))
		case TileWall:
			// 近傍8タイル（直交・斜め）にフロアがあるときだけ壁にする
			if chain.BuildData.AdjacentAnyFloor(i) {
				// 壁タイプを判定してスプライト番号を決定
				wallType := chain.BuildData.GetWallType(i)
				spriteNumber := getSpriteNumberForWallType(wallType)
				chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWallWithSprite(world, gc.Row(x), gc.Col(y), spriteNumber)
			}
		case TileWarpNext:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWarpNext(world, gc.Row(x), gc.Col(y))
		case TileWarpEscape:
			chain.BuildData.Level.Entities[i] = worldhelper.SpawnFieldWarpEscape(world, gc.Row(x), gc.Col(y))
		}
	}

	return chain.BuildData.Level
}

// getSpriteNumberForWallType は壁タイプに対応するスプライト番号を返す
func getSpriteNumberForWallType(wallType WallType) int {
	switch wallType {
	case WallTypeTop:
		return 10 // 上壁（下に床がある）
	case WallTypeBottom:
		return 11 // 下壁（上に床がある）
	case WallTypeLeft:
		return 12 // 左壁（右に床がある）
	case WallTypeRight:
		return 13 // 右壁（左に床がある）
	case WallTypeTopLeft:
		return 14 // 左上角（右下に床がある）
	case WallTypeTopRight:
		return 15 // 右上角（左下に床がある）
	case WallTypeBottomLeft:
		return 16 // 左下角（右上に床がある）
	case WallTypeBottomRight:
		return 17 // 右下角（左上に床がある）
	case WallTypeGeneric:
		return 1 // 汎用壁（従来の壁）
	default:
		return 1 // デフォルトは従来の壁
	}
}

// findPlayerStartPosition はプレイヤーのスタート位置を見つける
func findPlayerStartPosition(buildData *BuilderMap, world w.World) (int, int) {
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
		if buildData.IsSpawnableTile(world, gc.Row(pos.x), gc.Col(pos.y)) {
			return pos.x, pos.y
		}
	}

	// 全体をランダムに探索
	for attempt := 0; attempt < 200; attempt++ {
		x := buildData.RandomSource.Intn(width)
		y := buildData.RandomSource.Intn(height)
		if buildData.IsSpawnableTile(world, gc.Row(x), gc.Col(y)) {
			return x, y
		}
	}

	return -1, -1 // 見つからない場合
}

// validateMapWithPortals はポータルを配置してマップの接続性を検証する
func validateMapWithPortals(chain *BuilderChain, world w.World, gameResources *resources.Dungeon, playerX, playerY int) bool {
	// 進行ワープホールを配置
	warpNextPlaced := false
	for attempt := 0; attempt < 200; attempt++ {
		x := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
		y := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))
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
	escapePortalRequired := gameResources.Depth%5 == 0
	escapePortalPlaced := !escapePortalRequired

	if escapePortalRequired {
		for attempt := 0; attempt < 200; attempt++ {
			x := gc.Row(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileWidth)))
			y := gc.Col(chain.BuildData.RandomSource.Intn(int(chain.BuildData.Level.TileHeight)))

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
func createBuilderChain(builderType BuilderType, width gc.Row, height gc.Col, seed uint64) *BuilderChain {
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
