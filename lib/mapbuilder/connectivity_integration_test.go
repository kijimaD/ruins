package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

// TestBuilderChain_ValidateConnectivity はBuilderChainの接続性検証をテストする
func TestBuilderChain_ValidateConnectivity(t *testing.T) {
	t.Parallel()
	// 小さなテスト用マップを生成
	chain := NewSmallRoomBuilder(20, 20, 42) // 固定シードで再現可能
	chain.Build()

	// プレイヤーのスタート位置を部屋の中心付近に設定
	var playerStartX, playerStartY int
	if len(chain.BuildData.Rooms) > 0 {
		room := chain.BuildData.Rooms[0]
		playerStartX = int(room.X1+room.X2) / 2
		playerStartY = int(room.Y1+room.Y2) / 2
	} else {
		// 部屋がない場合はマップの中央付近を使用
		playerStartX = 10
		playerStartY = 10

		// プレイヤーのスタート位置を確実に床にする
		idx := chain.BuildData.Level.XYTileIndex(gc.Tile(playerStartX), gc.Tile(playerStartY))
		chain.BuildData.Tiles[idx] = TileFloor
	}

	// 接続性を検証
	result := chain.ValidateConnectivity(playerStartX, playerStartY)

	// プレイヤーのスタート位置は歩行可能である必要がある
	if !result.PlayerStartReachable {
		t.Error("Player start position should be reachable")
	}

	// 生成されたマップにはワープポータルや脱出ポータルはまだ配置されていないので
	// それらは0個である必要がある
	if len(result.WarpPortals) != 0 {
		t.Errorf("Expected 0 warp portals in basic room builder, got %d", len(result.WarpPortals))
	}

	if len(result.EscapePortals) != 0 {
		t.Errorf("Expected 0 escape portals in basic room builder, got %d", len(result.EscapePortals))
	}
}

// TestCaveBuilder_ValidateConnectivity は洞窟ビルダーの接続性をテストする
func TestCaveBuilder_ValidateConnectivity(t *testing.T) {
	t.Parallel()
	// 洞窟マップを生成
	chain := NewCaveBuilder(30, 30, 123)
	chain.Build()

	// 床タイルを見つけてプレイヤーのスタート位置とする
	var playerStartX, playerStartY int
	var foundFloor bool

	width := int(chain.BuildData.Level.TileWidth)
	height := int(chain.BuildData.Level.TileHeight)

	// 中央付近から床タイルを探す
	for x := width/2 - 5; x < width/2+5 && !foundFloor; x++ {
		for y := height/2 - 5; y < height/2+5 && !foundFloor; y++ {
			if x >= 0 && x < width && y >= 0 && y < height {
				idx := chain.BuildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				if chain.BuildData.Tiles[idx] == TileFloor {
					playerStartX = x
					playerStartY = y
					foundFloor = true
				}
			}
		}
	}

	if !foundFloor {
		t.Fatal("Could not find a floor tile for player start position")
	}

	// 接続性を検証
	result := chain.ValidateConnectivity(playerStartX, playerStartY)

	// プレイヤーのスタート位置は歩行可能である必要がある
	if !result.PlayerStartReachable {
		t.Error("Player start position should be reachable")
	}

	// 洞窟ビルダーもまだポータルを配置しないので0個である必要がある
	if len(result.WarpPortals) != 0 {
		t.Errorf("Expected 0 warp portals in cave builder, got %d", len(result.WarpPortals))
	}
}

// TestPathFinder_WithPortals は手動でポータルを配置して接続性をテストする
func TestPathFinder_WithPortals(t *testing.T) {
	t.Parallel()
	// テスト用の小さなマップを作成
	chain := NewBuilderChain(10, 10, 1)
	chain.StartWith(&TestRoomBuilder{})
	chain.Build()

	// プレイヤーのスタート位置（十字の中心）
	playerStartX, playerStartY := 5, 5

	// ワープポータルを配置（垂直通路上、到達可能な位置）
	warpIdx := chain.BuildData.Level.XYTileIndex(5, 2)
	chain.BuildData.Tiles[warpIdx] = TileWarpNext

	// 脱出ポータルを配置（到達不可能な位置：壁の中）
	escapeIdx := chain.BuildData.Level.XYTileIndex(1, 1)
	chain.BuildData.Tiles[escapeIdx] = TileWarpEscape

	// 接続性を検証
	result := chain.ValidateConnectivity(playerStartX, playerStartY)

	// デバッグ情報
	t.Logf("Player start reachable: %v", result.PlayerStartReachable)
	t.Logf("Warp portals: %d", len(result.WarpPortals))
	for i, portal := range result.WarpPortals {
		t.Logf("Warp portal %d: (%d, %d) reachable: %v", i, portal.X, portal.Y, portal.Reachable)
	}
	t.Logf("Escape portals: %d", len(result.EscapePortals))
	for i, portal := range result.EscapePortals {
		t.Logf("Escape portal %d: (%d, %d) reachable: %v", i, portal.X, portal.Y, portal.Reachable)
	}

	// ワープポータルは到達可能である必要がある
	if !result.HasReachableWarpPortal() {
		t.Error("Warp portal should be reachable")
	}

	// 脱出ポータルは到達不可能である必要がある
	if result.HasReachableEscapePortal() {
		t.Error("Escape portal should not be reachable")
	}

	// 完全には接続されていない
	if result.IsFullyConnected() {
		t.Error("Map should not be fully connected")
	}
}

// TestRoomBuilder はテスト用の簡単な部屋ビルダー
type TestRoomBuilder struct{}

func (t *TestRoomBuilder) BuildInitial(buildData *BuilderMap) {
	// 中央に簡単な十字型の部屋を作成
	width := int(buildData.Level.TileWidth)
	height := int(buildData.Level.TileHeight)

	// 全体を壁で埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	// 垂直方向の通路
	for y := 1; y < height-1; y++ {
		idx := buildData.Level.XYTileIndex(gc.Tile(width/2), gc.Tile(y))
		buildData.Tiles[idx] = TileFloor
	}

	// 水平方向の通路
	for x := 1; x < width-1; x++ {
		idx := buildData.Level.XYTileIndex(gc.Tile(x), gc.Tile(height/2))
		buildData.Tiles[idx] = TileFloor
	}
}
