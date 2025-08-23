package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/resources"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestBuildData_GetWallType(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成（7x7）
	width, height := gc.Row(7), gc.Col(7)
	buildData := &BuilderMap{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			TileSize:   consts.TileSize,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]Tile, int(width)*int(height)),
		Rooms:     []Rect{},
		Corridors: [][]resources.TileIdx{},
	}

	// 全体を壁で埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	// テストケース1: WallTypeTop（下に床がある壁）
	// 座標系注意: XYTileIndex(tx Row, ty Col) → tx は X座標（横方向）、ty は Y座標（縦方向）
	// インデックス計算: ty * width + tx
	centerWallX, centerWallY := gc.Row(3), gc.Col(3)
	bottomFloorX, bottomFloorY := centerWallX, centerWallY+1 // 下の床（Y座標が大きくなる）

	centerWallIdx := buildData.Level.XYTileIndex(centerWallX, centerWallY)
	bottomFloorIdx := buildData.Level.XYTileIndex(bottomFloorX, bottomFloorY)

	buildData.Tiles[bottomFloorIdx] = TileFloor

	// デバッグ情報を追加
	upFloor := buildData.isFloorOrWarp(buildData.UpTile(centerWallIdx))
	downFloor := buildData.isFloorOrWarp(buildData.DownTile(centerWallIdx))
	leftFloor := buildData.isFloorOrWarp(buildData.LeftTile(centerWallIdx))
	rightFloor := buildData.isFloorOrWarp(buildData.RightTile(centerWallIdx))

	wallType := buildData.GetWallType(centerWallIdx)
	if wallType != WallTypeTop {
		t.Errorf("WallTypeTopの判定が間違っています。期待値: %s, 実際: %s\n上:%t, 下:%t, 左:%t, 右:%t",
			WallTypeTop.String(), wallType.String(), upFloor, downFloor, leftFloor, rightFloor)
	}

	// テストケース2: WallTypeRight（左に床がある壁）
	leftFloorX, leftFloorY := centerWallX-1, centerWallY // 左の床（X座標が小さくなる）
	leftFloorIdx := buildData.Level.XYTileIndex(leftFloorX, leftFloorY)

	buildData.Tiles[leftFloorIdx] = TileFloor
	buildData.Tiles[bottomFloorIdx] = TileWall // 前のテストケースをリセット

	wallType = buildData.GetWallType(centerWallIdx)
	if wallType != WallTypeRight {
		t.Errorf("WallTypeRightの判定が間違っています。期待値: %s, 実際: %s", WallTypeRight.String(), wallType.String())
	}

	// テストケース3: WallTypeTopLeft（右下に床がある角壁）
	rightFloorX, rightFloorY := centerWallX+1, centerWallY // 右の床（X座標が大きくなる）
	downFloorX, downFloorY := centerWallX, centerWallY+1   // 下の床（Y座標が大きくなる）

	rightFloorIdx := buildData.Level.XYTileIndex(rightFloorX, rightFloorY)
	downFloorIdx := buildData.Level.XYTileIndex(downFloorX, downFloorY)

	buildData.Tiles[rightFloorIdx] = TileFloor
	buildData.Tiles[downFloorIdx] = TileFloor
	buildData.Tiles[leftFloorIdx] = TileWall // リセット

	wallType = buildData.GetWallType(centerWallIdx)
	if wallType != WallTypeTopLeft {
		t.Errorf("WallTypeTopLeftの判定が間違っています。期待値: %s, 実際: %s", WallTypeTopLeft.String(), wallType.String())
	}

	// テストケース4: WallTypeGeneric（複雑なパターン）
	upFloorX, upFloorY := centerWallX, centerWallY-1 // 上の床（Y座標が小さくなる）
	upFloorIdx := buildData.Level.XYTileIndex(upFloorX, upFloorY)
	buildData.Tiles[upFloorIdx] = TileFloor

	wallType = buildData.GetWallType(centerWallIdx) // 今は上、右、下に床がある状態
	if wallType != WallTypeGeneric {
		t.Errorf("WallTypeGenericの判定が間違っています。期待値: %s, 実際: %s", WallTypeGeneric.String(), wallType.String())
	}
}

func TestBuildData_GetWallType_WithWarpTiles(t *testing.T) {
	t.Parallel()
	// テスト用のマップを作成
	width, height := gc.Row(5), gc.Col(5)
	buildData := &BuilderMap{
		Level: resources.Level{
			TileWidth:  width,
			TileHeight: height,
			TileSize:   consts.TileSize,
			Entities:   make([]ecs.Entity, int(width)*int(height)),
		},
		Tiles:     make([]Tile, int(width)*int(height)),
		Rooms:     []Rect{},
		Corridors: [][]resources.TileIdx{},
	}

	// 全体を壁で埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = TileWall
	}

	// ワープタイルを配置
	wallX, wallY := gc.Row(2), gc.Col(2)
	warpX, warpY := wallX, wallY+1 // 下にワープネクスト（Y座標が大きくなる）

	warpNextIdx := buildData.Level.XYTileIndex(warpX, warpY)
	wallIdx := buildData.Level.XYTileIndex(wallX, wallY)
	buildData.Tiles[warpNextIdx] = TileWarpNext

	wallType := buildData.GetWallType(wallIdx)
	if wallType != WallTypeTop {
		t.Errorf("ワープタイルに対するWallTypeTopの判定が間違っています。期待値: %s, 実際: %s", WallTypeTop.String(), wallType.String())
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
