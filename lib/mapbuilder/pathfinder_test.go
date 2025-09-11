package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
)

// createTestBuildData はテスト用のBuildDataを作成する
func createTestBuildData(width, height int) *BuilderMap {
	tileCount := width * height
	tiles := make([]Tile, tileCount)

	// デフォルトで全て壁にする
	for i := range tiles {
		tiles[i] = TileWall
	}

	return &BuilderMap{
		Level: resources.Level{
			TileWidth:  gc.Tile(width),
			TileHeight: gc.Tile(height),
			TileSize:   consts.TileSize,
		},
		Tiles: tiles,
	}
}

func TestPathFinder_IsWalkable(t *testing.T) {
	t.Parallel()
	buildData := createTestBuildData(5, 5)
	pf := NewPathFinder(buildData)

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
	idx := buildData.Level.XYTileIndex(1, 1)
	buildData.Tiles[idx] = TileFloor
	if !pf.IsWalkable(1, 1) {
		t.Error("Expected floor tile to be walkable")
	}

	// ワープタイルテスト
	idx = buildData.Level.XYTileIndex(2, 2)
	buildData.Tiles[idx] = TileWarpNext
	if !pf.IsWalkable(2, 2) {
		t.Error("Expected warp next tile to be walkable")
	}

	// 脱出タイルテスト
	idx = buildData.Level.XYTileIndex(3, 3)
	buildData.Tiles[idx] = TileWarpEscape
	if !pf.IsWalkable(3, 3) {
		t.Error("Expected warp escape tile to be walkable")
	}
}

func TestPathFinder_FindPath_SimplePath(t *testing.T) {
	t.Parallel()
	buildData := createTestBuildData(5, 5)
	pf := NewPathFinder(buildData)

	// 簡単な一直線のパスを作成
	// (1,1) -> (1,2) -> (1,3)
	for y := 1; y <= 3; y++ {
		idx := buildData.Level.XYTileIndex(1, gc.Tile(y))
		buildData.Tiles[idx] = TileFloor
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
	buildData := createTestBuildData(5, 5)
	pf := NewPathFinder(buildData)

	// スタート地点のみ床にする（ゴールは壁のまま）
	idx := buildData.Level.XYTileIndex(1, 1)
	buildData.Tiles[idx] = TileFloor

	path := pf.FindPath(1, 1, 3, 3)

	if len(path) != 0 {
		t.Errorf("Expected no path, got path of length %d", len(path))
	}
}

func TestPathFinder_FindPath_LShapedPath(t *testing.T) {
	t.Parallel()
	buildData := createTestBuildData(5, 5)
	pf := NewPathFinder(buildData)

	// L字型のパスを作成
	// (1,1) -> (1,2) -> (2,2) -> (3,2)
	positions := []Position{{1, 1}, {1, 2}, {2, 2}, {3, 2}}
	for _, pos := range positions {
		idx := buildData.Level.XYTileIndex(gc.Tile(pos.X), gc.Tile(pos.Y))
		buildData.Tiles[idx] = TileFloor
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
	buildData := createTestBuildData(5, 5)
	pf := NewPathFinder(buildData)

	// パスを作成
	positions := []Position{{1, 1}, {1, 2}, {2, 2}}
	for _, pos := range positions {
		idx := buildData.Level.XYTileIndex(gc.Tile(pos.X), gc.Tile(pos.Y))
		buildData.Tiles[idx] = TileFloor
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

func TestPathFinder_ValidateMapConnectivity(t *testing.T) {
	t.Parallel()
	buildData := createTestBuildData(6, 6)
	pf := NewPathFinder(buildData)

	// プレイヤースタート地点
	playerX, playerY := 1, 1
	idx := buildData.Level.XYTileIndex(gc.Tile(playerX), gc.Tile(playerY))
	buildData.Tiles[idx] = TileFloor

	// 到達可能なワープポータル
	idx = buildData.Level.XYTileIndex(1, 2)
	buildData.Tiles[idx] = TileFloor
	idx = buildData.Level.XYTileIndex(1, 3)
	buildData.Tiles[idx] = TileWarpNext

	// 到達不可能な脱出ポータル（孤立している）
	idx = buildData.Level.XYTileIndex(4, 4)
	buildData.Tiles[idx] = TileWarpEscape

	result := pf.ValidateMapConnectivity(playerX, playerY)

	// プレイヤーのスタート位置は歩行可能である必要がある
	if !result.PlayerStartReachable {
		t.Error("Expected player start position to be reachable")
	}

	// ワープポータルが1つ見つかり、到達可能である必要がある
	if len(result.WarpPortals) != 1 {
		t.Errorf("Expected 1 warp portal, got %d", len(result.WarpPortals))
	}
	if !result.WarpPortals[0].Reachable {
		t.Error("Expected warp portal to be reachable")
	}

	// 脱出ポータルが1つ見つかり、到達不可能である必要がある
	if len(result.EscapePortals) != 1 {
		t.Errorf("Expected 1 escape portal, got %d", len(result.EscapePortals))
	}
	if result.EscapePortals[0].Reachable {
		t.Error("Expected escape portal to be not reachable")
	}

	// 到達可能なワープポータルがある
	if !result.HasReachableWarpPortal() {
		t.Error("Expected to have reachable warp portal")
	}

	// 到達可能な脱出ポータルがない
	if result.HasReachableEscapePortal() {
		t.Error("Expected to not have reachable escape portal")
	}

	// すべてが接続されていない（脱出ポータルに到達できない）
	if result.IsFullyConnected() {
		t.Error("Expected map to not be fully connected")
	}
}

func TestPathFinder_ValidateMapConnectivity_FullyConnected(t *testing.T) {
	t.Parallel()
	buildData := createTestBuildData(6, 6)
	pf := NewPathFinder(buildData)

	// プレイヤースタート地点から全ポータルへの直線パス
	playerX, playerY := 2, 2

	// 床を配置
	for y := 1; y <= 4; y++ {
		idx := buildData.Level.XYTileIndex(2, gc.Tile(y))
		buildData.Tiles[idx] = TileFloor
	}

	// ワープポータル（到達可能）
	idx := buildData.Level.XYTileIndex(2, 1)
	buildData.Tiles[idx] = TileWarpNext

	// 脱出ポータル（到達可能）
	idx = buildData.Level.XYTileIndex(2, 4)
	buildData.Tiles[idx] = TileWarpEscape

	result := pf.ValidateMapConnectivity(playerX, playerY)

	// 完全に接続されている必要がある
	if !result.IsFullyConnected() {
		t.Error("Expected map to be fully connected")
	}

	if !result.HasReachableWarpPortal() {
		t.Error("Expected to have reachable warp portal")
	}

	if !result.HasReachableEscapePortal() {
		t.Error("Expected to have reachable escape portal")
	}
}
