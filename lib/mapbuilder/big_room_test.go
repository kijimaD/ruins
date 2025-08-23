package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestBigRoomBuilder(t *testing.T) {
	t.Parallel()

	width, height := gc.Row(20), gc.Col(20)
	seed := uint64(12345)

	chain := NewBigRoomBuilder(width, height, seed)
	chain.Build()

	// 部屋が1つだけ生成されることを確認
	if len(chain.BuildData.Rooms) != 1 {
		t.Errorf("期待される部屋数: 1, 実際: %d", len(chain.BuildData.Rooms))
	}

	// 生成された部屋のサイズを確認
	room := chain.BuildData.Rooms[0]
	expectedMinWidth := int(width) - 6 // margin * 2 + 境界
	expectedMinHeight := int(height) - 6

	actualWidth := int(room.X2 - room.X1 + 1)
	actualHeight := int(room.Y2 - room.Y1 + 1)

	if actualWidth < expectedMinWidth {
		t.Errorf("部屋の幅が小さすぎます。期待最小値: %d, 実際: %d", expectedMinWidth, actualWidth)
	}

	if actualHeight < expectedMinHeight {
		t.Errorf("部屋の高さが小さすぎます。期待最小値: %d, 実際: %d", expectedMinHeight, actualHeight)
	}

	// 部屋の内部が床になっていることを確認
	floorCount := 0
	for x := room.X1; x <= room.X2; x++ {
		for y := room.Y1; y <= room.Y2; y++ {
			idx := chain.BuildData.Level.XYTileIndex(x, y)
			if chain.BuildData.Tiles[idx] == TileFloor {
				floorCount++
			}
		}
	}

	expectedFloorCount := actualWidth * actualHeight
	if floorCount != expectedFloorCount {
		t.Errorf("床タイル数が不正です。期待: %d, 実際: %d", expectedFloorCount, floorCount)
	}
}

func TestBigRoomWithPillarsBuilder(t *testing.T) {
	t.Parallel()

	width, height := gc.Row(20), gc.Col(20)
	seed := uint64(54321)
	pillarSpacing := 3

	chain := NewBigRoomWithPillarsBuilder(width, height, seed, pillarSpacing)
	chain.Build()

	// 部屋が1つだけ生成されることを確認
	if len(chain.BuildData.Rooms) != 1 {
		t.Errorf("期待される部屋数: 1, 実際: %d", len(chain.BuildData.Rooms))
	}

	// 柱が配置されていることを確認
	room := chain.BuildData.Rooms[0]
	pillarCount := 0

	// 柱の予想配置位置をチェック
	startX := int(room.X1) + pillarSpacing
	startY := int(room.Y1) + pillarSpacing

	for x := startX; x < int(room.X2); x += pillarSpacing + 1 {
		for y := startY; y < int(room.Y2); y += pillarSpacing + 1 {
			idx := chain.BuildData.Level.XYTileIndex(gc.Row(x), gc.Col(y))
			if chain.BuildData.Tiles[idx] == TileWall {
				pillarCount++
			}
		}
	}

	// 少なくとも1つの柱があることを確認
	if pillarCount == 0 {
		t.Error("柱が配置されていません")
	}

	t.Logf("配置された柱の数: %d", pillarCount)
}

func TestBigRoomBuilderReproducibility(t *testing.T) {
	t.Parallel()

	// 同じシードで複数回生成して同じ結果になることを確認
	width, height := gc.Row(15), gc.Col(15)
	seed := uint64(99999)

	// 1回目の生成
	chain1 := NewBigRoomBuilder(width, height, seed)
	chain1.Build()

	// 2回目の生成
	chain2 := NewBigRoomBuilder(width, height, seed)
	chain2.Build()

	// 部屋数が同じことを確認
	if len(chain1.BuildData.Rooms) != len(chain2.BuildData.Rooms) {
		t.Errorf("部屋数が異なります。1回目: %d, 2回目: %d",
			len(chain1.BuildData.Rooms), len(chain2.BuildData.Rooms))
	}

	// 部屋の位置とサイズが同じことを確認
	for i := range chain1.BuildData.Rooms {
		room1 := chain1.BuildData.Rooms[i]
		room2 := chain2.BuildData.Rooms[i]

		if room1 != room2 {
			t.Errorf("部屋[%d]が異なります。1回目: %+v, 2回目: %+v", i, room1, room2)
		}
	}

	// タイル配置が同じことを確認
	for i, tile1 := range chain1.BuildData.Tiles {
		tile2 := chain2.BuildData.Tiles[i]
		if tile1 != tile2 {
			t.Errorf("タイル[%d]が異なります。1回目: %v, 2回目: %v", i, tile1, tile2)
		}
	}
}

func TestBigRoomBuilderBoundaries(t *testing.T) {
	t.Parallel()

	// 境界の処理が正しいかを確認
	width, height := gc.Row(10), gc.Col(10)
	seed := uint64(11111)

	chain := NewBigRoomBuilder(width, height, seed)
	chain.Build()

	// マップの境界が壁になっていることを確認
	for x := 0; x < int(width); x++ {
		// 上端
		idx := chain.BuildData.Level.XYTileIndex(gc.Row(x), gc.Col(0))
		if chain.BuildData.Tiles[idx] != TileWall {
			t.Errorf("上端の境界[%d,0]が壁になっていません: %v", x, chain.BuildData.Tiles[idx])
		}

		// 下端
		idx = chain.BuildData.Level.XYTileIndex(gc.Row(x), height-1)
		if chain.BuildData.Tiles[idx] != TileWall {
			t.Errorf("下端の境界[%d,%d]が壁になっていません: %v", x, height-1, chain.BuildData.Tiles[idx])
		}
	}

	for y := 0; y < int(height); y++ {
		// 左端
		idx := chain.BuildData.Level.XYTileIndex(gc.Row(0), gc.Col(y))
		if chain.BuildData.Tiles[idx] != TileWall {
			t.Errorf("左端の境界[0,%d]が壁になっていません: %v", y, chain.BuildData.Tiles[idx])
		}

		// 右端
		idx = chain.BuildData.Level.XYTileIndex(width-1, gc.Col(y))
		if chain.BuildData.Tiles[idx] != TileWall {
			t.Errorf("右端の境界[%d,%d]が壁になっていません: %v", width-1, y, chain.BuildData.Tiles[idx])
		}
	}
}
