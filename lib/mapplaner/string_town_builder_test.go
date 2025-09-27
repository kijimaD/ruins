package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestNewStringTownPlanner(t *testing.T) {
	t.Parallel()
	// 新しい文字列ベースの街ビルダーをテスト
	width, height := gc.Tile(50), gc.Tile(50)
	chain := NewStringTownPlanner(width, height, 12345)

	// ビルダーチェーンを実行
	chain.Build()

	// サイズをチェック
	if chain.PlanData.Level.TileWidth != width {
		t.Errorf("幅が期待値と違います: 期待値 %d, 実際 %d", width, chain.PlanData.Level.TileWidth)
	}
	if chain.PlanData.Level.TileHeight != height {
		t.Errorf("高さが期待値と違います: 期待値 %d, 実際 %d", height, chain.PlanData.Level.TileHeight)
	}

	// タイル配置をチェック（50x50マップ用）
	testCases := []struct {
		x, y         int
		expectedTile Tile
		description  string
	}{
		{0, 0, TileWall, "左上角の壁"},
		{25, 26, TileFloor, "中央広場の床"},
		{11, 12, TileFloor, "住宅区域の床"},
		{25, 19, TileFloor, "商業区域の床"},
		{49, 49, TileWall, "右下角の壁"},
	}

	for _, tc := range testCases {
		tileIdx := chain.PlanData.Level.XYTileIndex(gc.Tile(tc.x), gc.Tile(tc.y))
		actualTile := chain.PlanData.Tiles[tileIdx]
		if actualTile != tc.expectedTile {
			t.Errorf("%s (位置: %d,%d): 期待値 %v, 実際 %v", tc.description, tc.x, tc.y, tc.expectedTile, actualTile)
		}
	}
}

func TestStringTownPlannerWithEntityPlan(t *testing.T) {
	t.Parallel()
	// 街レイアウトを取得
	tileMap, entityMap := GetTownLayout()

	// EntityPlanを直接生成
	plan, err := BuildEntityPlanFromStrings(tileMap, entityMap)
	if err != nil {
		t.Fatalf("EntityPlan生成に失敗: %v", err)
	}

	// サイズをチェック
	expectedWidth := 50
	expectedHeight := 50
	if plan.Width != expectedWidth {
		t.Errorf("幅が期待値と違います: 期待値 %d, 実際 %d", expectedWidth, plan.Width)
	}
	if plan.Height != expectedHeight {
		t.Errorf("高さが期待値と違います: 期待値 %d, 実際 %d", expectedHeight, plan.Height)
	}

	// Props配置をチェック
	propCount := 0
	propTypes := make(map[gc.PropType]int)

	for _, entity := range plan.Entities {
		if entity.EntityType == EntityTypeProp && entity.PropType != nil {
			propCount++
			propTypes[*entity.PropType]++
		}
	}

	t.Logf("総Props数: %d", propCount)

	// 主要なPropタイプが存在することを確認
	expectedTypes := []gc.PropType{
		gc.PropTypeChair, gc.PropTypeTable, gc.PropTypeBed,
		gc.PropTypeBookshelf, gc.PropTypeAltar, gc.PropTypeChest,
		gc.PropTypeBarrel, gc.PropTypeCrate, gc.PropTypeTorch,
	}

	for _, expectedType := range expectedTypes {
		if count, exists := propTypes[expectedType]; !exists || count == 0 {
			t.Errorf("PropType %v が配置されていません", expectedType)
		} else {
			t.Logf("PropType %v: %d個", expectedType, count)
		}
	}
}

func TestTownLayoutStructure(t *testing.T) {
	t.Parallel()
	// 街レイアウトの構造をチェック
	tileMap, entityMap := GetTownLayout()

	// 基本的な整合性チェック
	if len(tileMap) != 50 {
		t.Errorf("タイルマップの行数が期待値と違います: 期待値 50, 実際 %d", len(tileMap))
	}
	if len(entityMap) != 50 {
		t.Errorf("エンティティマップの行数が期待値と違います: 期待値 50, 実際 %d", len(entityMap))
	}

	// 各行の長さをチェック
	for i, row := range tileMap {
		if len(row) != 50 {
			t.Errorf("タイルマップの行 %d の長さが期待値と違います: 期待値 50, 実際 %d", i, len(row))
		}
	}

	for i, row := range entityMap {
		if len(row) != 50 {
			t.Errorf("エンティティマップの行 %d の長さが期待値と違います: 期待値 50, 実際 %d", i, len(row))
		}
	}

	// 境界が壁で囲まれていることをチェック
	for x := 0; x < 50; x++ {
		if tileMap[0][x] != '#' || tileMap[49][x] != '#' {
			t.Errorf("上下の境界が壁でありません: 位置 (%d, 0) または (%d, 49)", x, x)
		}
	}

	for y := 0; y < 50; y++ {
		if tileMap[y][0] != '#' || tileMap[y][49] != '#' {
			t.Errorf("左右の境界が壁でありません: 位置 (0, %d) または (49, %d)", y, y)
		}
	}

	// 中央部分が床で構成されていることをチェック
	floorExists := false
	for x := 1; x < 49; x++ {
		if tileMap[25][x] == 'f' || tileMap[25][x] == 'r' {
			floorExists = true
			break
		}
	}
	if !floorExists {
		t.Error("中央に床が存在しません")
	}

	// プレイヤー開始位置が存在することをチェック
	playerExists := false
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			if entityMap[y][x] == '@' {
				playerExists = true
				break
			}
		}
		if playerExists {
			break
		}
	}
	if !playerExists {
		t.Error("プレイヤー開始位置が配置されていません")
	}
}

func TestNewTownPlannerIntegration(t *testing.T) {
	t.Parallel()
	// 既存のNewTownPlanner関数が新しいStringTownPlannerを使用していることをテスト
	width, height := gc.Tile(30), gc.Tile(30)
	chain := NewTownPlanner(width, height, 54321)

	// ビルダーチェーンを実行
	chain.Build()

	// 基本的な構造チェック
	if chain.PlanData.Level.TileWidth != width {
		t.Errorf("幅が期待値と違います: 期待値 %d, 実際 %d", width, chain.PlanData.Level.TileWidth)
	}
	if chain.PlanData.Level.TileHeight != height {
		t.Errorf("高さが期待値と違います: 期待値 %d, 実際 %d", height, chain.PlanData.Level.TileHeight)
	}

	// 街の特徴的な要素が存在することをチェック
	hasFloor := false
	hasWall := false

	for _, tile := range chain.PlanData.Tiles {
		switch tile {
		case TileFloor:
			hasFloor = true
		case TileWall:
			hasWall = true
		}
	}

	if !hasFloor {
		t.Error("床タイルが存在しません")
	}
	if !hasWall {
		t.Error("壁タイルが存在しません")
	}
}
