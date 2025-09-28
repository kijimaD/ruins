package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

// TestPlannerChain_ValidateConnectivity はPlannerChainの接続性検証をテストする
func TestPlannerChain_ValidateConnectivity(t *testing.T) {
	t.Parallel()
	// 小さなテスト用マップを生成
	chain := NewSmallRoomPlanner(20, 20, 42) // 固定シードで再現可能
	chain.PlanData.RawMaster = createTestRawMaster()
	chain.Plan()

	// プレイヤーのスタート位置を部屋の中心付近に設定
	var playerStartX, playerStartY int
	if len(chain.PlanData.Rooms) > 0 {
		room := chain.PlanData.Rooms[0]
		playerStartX = int(room.X1+room.X2) / 2
		playerStartY = int(room.Y1+room.Y2) / 2
	} else {
		// 部屋がない場合はマップの中央付近を使用
		playerStartX = 10
		playerStartY = 10

		// プレイヤーのスタート位置を確実に床にする
		idx := chain.PlanData.Level.XYTileIndex(gc.Tile(playerStartX), gc.Tile(playerStartY))
		chain.PlanData.Tiles[idx] = chain.PlanData.GenerateTile("Floor")
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

// TestCavePlanner_ValidateConnectivity は洞窟ビルダーの接続性をテストする
func TestCavePlanner_ValidateConnectivity(t *testing.T) {
	t.Parallel()
	// 洞窟マップを生成
	chain := NewCavePlanner(30, 30, 123)
	chain.PlanData.RawMaster = createTestRawMaster()
	chain.Plan()

	// 床タイルを見つけてプレイヤーのスタート位置とする
	var playerStartX, playerStartY int
	var foundFloor bool

	width := int(chain.PlanData.Level.TileWidth)
	height := int(chain.PlanData.Level.TileHeight)

	// 中央付近から床タイルを探す
	for x := width/2 - 5; x < width/2+5 && !foundFloor; x++ {
		for y := height/2 - 5; y < height/2+5 && !foundFloor; y++ {
			if x >= 0 && x < width && y >= 0 && y < height {
				idx := chain.PlanData.Level.XYTileIndex(gc.Tile(x), gc.Tile(y))
				if chain.PlanData.Tiles[idx].Walkable {
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
	chain := NewPlannerChain(10, 10, 1)
	chain.PlanData.RawMaster = createTestRawMaster()
	chain.StartWith(&TestRoomPlanner{})
	chain.Plan()

	// プレイヤーのスタート位置（十字の中心）
	playerStartX, playerStartY := 5, 5

	// ワープポータルを配置（垂直通路上、到達可能な位置）
	warpIdx := chain.PlanData.Level.XYTileIndex(5, 2)
	chain.PlanData.Tiles[warpIdx] = chain.PlanData.GenerateTile("Floor")
	// ワープポータルエンティティを追加
	chain.PlanData.WarpPortals = append(chain.PlanData.WarpPortals, WarpPortal{
		X:    5,
		Y:    2,
		Type: WarpPortalNext,
	})

	// 脱出ポータルを配置（到達不可能な位置：壁の中）
	escapeIdx := chain.PlanData.Level.XYTileIndex(1, 1)
	chain.PlanData.Tiles[escapeIdx] = chain.PlanData.GenerateTile("Floor")
	// 脱出ポータルエンティティを追加
	chain.PlanData.WarpPortals = append(chain.PlanData.WarpPortals, WarpPortal{
		X:    1,
		Y:    1,
		Type: WarpPortalEscape,
	})

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

// TestRoomPlanner はテスト用の簡単な部屋ビルダー
type TestRoomPlanner struct{}

func (t *TestRoomPlanner) PlanInitial(planData *MetaPlan) error {
	// 中央に簡単な十字型の部屋を作成
	width := int(planData.Level.TileWidth)
	height := int(planData.Level.TileHeight)

	// 全体を壁で埋める
	for i := range planData.Tiles {
		planData.Tiles[i] = planData.GenerateTile("Wall")
	}

	// 垂直方向の通路
	for y := 1; y < height-1; y++ {
		idx := planData.Level.XYTileIndex(gc.Tile(width/2), gc.Tile(y))
		planData.Tiles[idx] = planData.GenerateTile("Floor")
	}

	// 水平方向の通路
	for x := 1; x < width-1; x++ {
		idx := planData.Level.XYTileIndex(gc.Tile(x), gc.Tile(height/2))
		planData.Tiles[idx] = planData.GenerateTile("Floor")
	}

	return nil
}
