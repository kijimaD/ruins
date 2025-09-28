package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestPlanData_GetWallType(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成（7x7）
	width, height := gc.Tile(7), gc.Tile(7)
	planData := &MetaPlan{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]raw.TileRaw, int(width)*int(height)),
		Rooms:     []gc.Rect{},
		Corridors: [][]resources.TileIdx{},
		RawMaster: createTestRawMaster(),
	}

	// 全体を壁で埋める
	for i := range planData.Tiles {
		planData.Tiles[i] = planData.GenerateTile("Wall")
	}

	// テストケース1: WallTypeTop（下に床がある壁）
	// 座標系注意: XYTileIndex(tx Row, ty Col) → tx は X座標（横方向）、ty は Y座標（縦方向）
	// インデックス計算: ty * width + tx
	centerWallX, centerWallY := gc.Tile(3), gc.Tile(3)
	bottomFloorX, bottomFloorY := centerWallX, centerWallY+1 // 下の床（Y座標が大きくなる）

	centerWallIdx := planData.Level.XYTileIndex(centerWallX, centerWallY)
	bottomFloorIdx := planData.Level.XYTileIndex(bottomFloorX, bottomFloorY)

	planData.Tiles[bottomFloorIdx] = planData.GenerateTile("Floor")

	// デバッグ情報を追加
	upFloor := planData.isFloorOrWarp(planData.UpTile(centerWallIdx))
	downFloor := planData.isFloorOrWarp(planData.DownTile(centerWallIdx))
	leftFloor := planData.isFloorOrWarp(planData.LeftTile(centerWallIdx))
	rightFloor := planData.isFloorOrWarp(planData.RightTile(centerWallIdx))

	wallType := planData.GetWallType(centerWallIdx)
	if wallType != WallTypeTop {
		t.Errorf("WallTypeTopの判定が間違っています。期待値: %s, 実際: %s\n上:%t, 下:%t, 左:%t, 右:%t",
			WallTypeTop.String(), wallType.String(), upFloor, downFloor, leftFloor, rightFloor)
	}

	// テストケース2: WallTypeRight（左に床がある壁）
	leftFloorX, leftFloorY := centerWallX-1, centerWallY // 左の床（X座標が小さくなる）
	leftFloorIdx := planData.Level.XYTileIndex(leftFloorX, leftFloorY)

	planData.Tiles[leftFloorIdx] = planData.GenerateTile("Floor")
	planData.Tiles[bottomFloorIdx] = planData.GenerateTile("Wall") // 前のテストケースをリセット

	wallType = planData.GetWallType(centerWallIdx)
	if wallType != WallTypeRight {
		t.Errorf("WallTypeRightの判定が間違っています。期待値: %s, 実際: %s", WallTypeRight.String(), wallType.String())
	}

	// テストケース3: WallTypeTopLeft（右下に床がある角壁）
	rightFloorX, rightFloorY := centerWallX+1, centerWallY // 右の床（X座標が大きくなる）
	downFloorX, downFloorY := centerWallX, centerWallY+1   // 下の床（Y座標が大きくなる）

	rightFloorIdx := planData.Level.XYTileIndex(rightFloorX, rightFloorY)
	downFloorIdx := planData.Level.XYTileIndex(downFloorX, downFloorY)

	planData.Tiles[rightFloorIdx] = planData.GenerateTile("Floor")
	planData.Tiles[downFloorIdx] = planData.GenerateTile("Floor")
	planData.Tiles[leftFloorIdx] = planData.GenerateTile("Wall") // リセット

	wallType = planData.GetWallType(centerWallIdx)
	if wallType != WallTypeTopLeft {
		t.Errorf("WallTypeTopLeftの判定が間違っています。期待値: %s, 実際: %s", WallTypeTopLeft.String(), wallType.String())
	}

	// テストケース4: WallTypeGeneric（複雑なパターン）
	upFloorX, upFloorY := centerWallX, centerWallY-1 // 上の床（Y座標が小さくなる）
	upFloorIdx := planData.Level.XYTileIndex(upFloorX, upFloorY)
	planData.Tiles[upFloorIdx] = planData.GenerateTile("Floor")

	wallType = planData.GetWallType(centerWallIdx) // 今は上、右、下に床がある状態
	if wallType != WallTypeGeneric {
		t.Errorf("WallTypeGenericの判定が間違っています。期待値: %s, 実際: %s", WallTypeGeneric.String(), wallType.String())
	}
}

func TestPlanData_GetWallType_WithWarpTiles(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成
	width, height := gc.Tile(5), gc.Tile(5)
	planData := &MetaPlan{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]raw.TileRaw, int(width)*int(height)),
		Rooms:     []gc.Rect{},
		Corridors: [][]resources.TileIdx{},
		RawMaster: createTestRawMaster(),
	}

	// 全体を壁で埋める
	for i := range planData.Tiles {
		planData.Tiles[i] = planData.GenerateTile("Wall")
	}

	// ワープポータルを配置
	wallX, wallY := gc.Tile(2), gc.Tile(2)
	warpX, warpY := wallX, wallY+1 // 下にワープネクスト（Y座標が大きくなる）

	warpNextIdx := planData.Level.XYTileIndex(warpX, warpY)
	wallIdx := planData.Level.XYTileIndex(wallX, wallY)
	planData.Tiles[warpNextIdx] = planData.GenerateTile("Floor")
	planData.WarpPortals = append(planData.WarpPortals, WarpPortal{
		X:    int(warpX),
		Y:    int(warpY),
		Type: WarpPortalNext,
	})

	wallType := planData.GetWallType(wallIdx)
	if wallType != WallTypeTop {
		t.Errorf("ワープポータル（床）に対するWallTypeTopの判定が間違っています。期待値: %s, 実際: %s", WallTypeTop.String(), wallType.String())
	}
}

func TestWallType_String(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		wallType WallType
		expected string
	}{
		{WallTypeTop, "Top"},
		{WallTypeBottom, "Bottom"},
		{WallTypeLeft, "Left"},
		{WallTypeRight, "Right"},
		{WallTypeTopLeft, "TopLeft"},
		{WallTypeTopRight, "TopRight"},
		{WallTypeBottomLeft, "BottomLeft"},
		{WallTypeBottomRight, "BottomRight"},
		{WallTypeGeneric, "Generic"},
	}

	for _, tc := range testCases {
		actual := tc.wallType.String()
		if actual != tc.expected {
			t.Errorf("WallType.String()の結果が間違っています。期待値: %s, 実際: %s", tc.expected, actual)
		}
	}
}
