package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/resources"
)

// createTestPlanData はテスト用のPlanDataを作成する
func createTestPlanData(width, height int) *MetaPlan {
	tileCount := width * height
	tiles := make([]raw.TileRaw, tileCount)

	// 一時的なMetaPlanインスタンスを作成
	tempPlan := &MetaPlan{
		Level: resources.Level{
			TileWidth:  gc.Tile(width),
			TileHeight: gc.Tile(height),
		},
		Tiles:     tiles,
		RawMaster: CreateTestRawMaster(),
	}

	// デフォルトで全て壁にする
	for i := range tiles {
		tiles[i] = tempPlan.GetTile("Wall")
	}

	return tempPlan
}

func TestPathFinder_IsWalkable(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(5, 5)
	pf := NewPathFinder(planData)

	// 境界外テスト
	if pf.IsWalkable(-1, 0) {
		t.Error("Expected (-1, 0) to be not walkable")
	}
	if pf.IsWalkable(0, -1) {
		t.Error("Expected (0, -1) to be not walkable")
	}
	if pf.IsWalkable(5, 0) {
		t.Error("Expected (5, 0) to be not walkable")
	}
	if pf.IsWalkable(0, 5) {
		t.Error("Expected (0, 5) to be not walkable")
	}

	// 壁タイルテスト（デフォルト）
	if pf.IsWalkable(1, 1) {
		t.Error("Expected wall tile to be not walkable")
	}

	// 床タイルに変更してテスト
	idx := planData.Level.XYTileIndex(1, 1)
	planData.Tiles[idx] = planData.GetTile("Floor")
	if !pf.IsWalkable(1, 1) {
		t.Error("Expected floor tile to be walkable")
	}

	// ワープタイルテスト
	idx = planData.Level.XYTileIndex(2, 2)
	planData.Tiles[idx] = planData.GetTile("Floor")
	if !pf.IsWalkable(2, 2) {
		t.Error("Expected warp next tile to be walkable")
	}

	// 脱出タイルテスト
	idx = planData.Level.XYTileIndex(3, 3)
	planData.Tiles[idx] = planData.GetTile("Floor")
	if !pf.IsWalkable(3, 3) {
		t.Error("Expected warp escape tile to be walkable")
	}
}

func TestPathFinder_FindPath_SimplePath(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(5, 5)
	pf := NewPathFinder(planData)

	// 簡単な一直線のパスを作成
	// (1,1) -> (1,2) -> (1,3)
	for y := 1; y <= 3; y++ {
		idx := planData.Level.XYTileIndex(1, gc.Tile(y))
		planData.Tiles[idx] = planData.GetTile("Floor")
	}

	path := pf.FindPath(1, 1, 1, 3)

	expectedLength := 3 // スタート、中間、ゴール
	if len(path) != expectedLength {
		t.Errorf("Expected path length %d, got %d", expectedLength, len(path))
	}

	// パスの内容を検証
	expected := []Position{{1, 1}, {1, 2}, {1, 3}}
	for i, pos := range expected {
		if i >= len(path) || path[i].X != pos.X || path[i].Y != pos.Y {
			t.Errorf("Expected position %d to be (%d, %d), got (%d, %d)",
				i, pos.X, pos.Y, path[i].X, path[i].Y)
		}
	}
}

func TestPathFinder_FindPath_NoPath(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(5, 5)
	pf := NewPathFinder(planData)

	// スタート地点のみ床にする（ゴールは壁のまま）
	idx := planData.Level.XYTileIndex(1, 1)
	planData.Tiles[idx] = planData.GetTile("Floor")

	path := pf.FindPath(1, 1, 3, 3)

	if len(path) != 0 {
		t.Errorf("Expected no path, got path of length %d", len(path))
	}
}

func TestPathFinder_FindPath_LShapedPath(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(5, 5)
	pf := NewPathFinder(planData)

	// L字型のパスを作成
	// (1,1) -> (1,2) -> (2,2) -> (3,2)
	positions := []Position{{1, 1}, {1, 2}, {2, 2}, {3, 2}}
	for _, pos := range positions {
		idx := planData.Level.XYTileIndex(gc.Tile(pos.X), gc.Tile(pos.Y))
		planData.Tiles[idx] = planData.GetTile("Floor")
	}

	path := pf.FindPath(1, 1, 3, 2)

	if len(path) != 4 {
		t.Errorf("Expected path length 4, got %d", len(path))
	}

	// スタートとゴールが正しいことを確認
	if path[0].X != 1 || path[0].Y != 1 {
		t.Errorf("Expected start at (1,1), got (%d,%d)", path[0].X, path[0].Y)
	}
	if path[len(path)-1].X != 3 || path[len(path)-1].Y != 2 {
		t.Errorf("Expected goal at (3,2), got (%d,%d)",
			path[len(path)-1].X, path[len(path)-1].Y)
	}
}

func TestPathFinder_IsReachable(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(5, 5)
	pf := NewPathFinder(planData)

	// パスを作成
	positions := []Position{{1, 1}, {1, 2}, {2, 2}}
	for _, pos := range positions {
		idx := planData.Level.XYTileIndex(gc.Tile(pos.X), gc.Tile(pos.Y))
		planData.Tiles[idx] = planData.GetTile("Floor")
	}

	// 到達可能なテスト
	if !pf.IsReachable(1, 1, 2, 2) {
		t.Error("Expected (1,1) to (2,2) to be reachable")
	}

	// 到達不可能なテスト
	if pf.IsReachable(1, 1, 3, 3) {
		t.Error("Expected (1,1) to (3,3) to be not reachable")
	}
}

func TestPathFinder_ValidateConnectivity(t *testing.T) {
	t.Parallel()
	planData := createTestPlanData(6, 6)
	pf := NewPathFinder(planData)

	// プレイヤースタート地点
	playerX, playerY := 1, 1
	idx := planData.Level.XYTileIndex(gc.Tile(playerX), gc.Tile(playerY))
	planData.Tiles[idx] = planData.GetTile("Floor")

	// 到達可能なワープポータル
	idx = planData.Level.XYTileIndex(1, 2)
	planData.Tiles[idx] = planData.GetTile("Floor")
	idx = planData.Level.XYTileIndex(1, 3)
	planData.Tiles[idx] = planData.GetTile("Floor")
	// ワープポータルエンティティを追加
	planData.WarpPortals = append(planData.WarpPortals, WarpPortal{
		X:    1,
		Y:    3,
		Type: WarpPortalNext,
	})

	// 接続性検証 - 到達可能なのでエラーなし
	if err := pf.ValidateConnectivity(playerX, playerY); err != nil {
		t.Errorf("Expected connectivity validation to pass, but got error: %v", err)
	}

	// エラーケース: プレイヤー開始位置が歩行不可能
	err := pf.ValidateConnectivity(0, 0)
	if err != ErrPlayerPlacement {
		t.Errorf("Expected ErrPlayerPlacement, got %v", err)
	}

	// エラーケース: ワープポータルに到達不可能
	// 到達不可能なワープポータル（孤立している）を追加
	idx = planData.Level.XYTileIndex(4, 4)
	planData.Tiles[idx] = planData.GetTile("Floor")
	planData.WarpPortals = []WarpPortal{{
		X:    4,
		Y:    4,
		Type: WarpPortalNext,
	}}

	err = pf.ValidateConnectivity(playerX, playerY)
	if err != ErrConnectivity {
		t.Errorf("Expected ErrConnectivity, got %v", err)
	}
}
