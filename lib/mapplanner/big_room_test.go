package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestBigRoomPlanner(t *testing.T) {
	t.Parallel()

	width, height := gc.Tile(20), gc.Tile(20)
	seed := uint64(12345)

	chain := NewBigRoomPlanner(width, height, seed)
	chain.PlanData.RawMaster = createTestRawMaster()
	chain.Plan()

	// 部屋が1つだけ生成されることを確認
	if len(chain.PlanData.Rooms) != 1 {
		t.Errorf("期待される部屋数: 1, 実際: %d", len(chain.PlanData.Rooms))
	}

	// 生成された部屋のサイズを確認
	room := chain.PlanData.Rooms[0]
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

	// 床と壁の両方が存在することを確認
	floorCount := 0
	wallCount := 0
	for _, tile := range chain.PlanData.Tiles {
		if tile.Walkable {
			floorCount++
		} else {
			wallCount++
		}
	}

	if floorCount == 0 {
		t.Error("床タイルが存在しません")
	}
	if wallCount == 0 {
		t.Error("壁タイルが存在しません")
	}

	t.Logf("床タイル: %d, 壁タイル: %d", floorCount, wallCount)
}

func TestBigRoomVariations(t *testing.T) {
	t.Parallel()

	// 異なるシードで複数回テストして、各バリエーションが出ることを確認
	seeds := []uint64{1, 42, 123, 456, 789, 1024, 2048, 3333, 5000, 9999}

	variantCounts := make(map[string]int)

	for _, seed := range seeds {
		chain := NewBigRoomPlanner(20, 20, seed)
		chain.PlanData.RawMaster = createTestRawMaster()
		chain.Plan()

		// 部屋が1つ生成されることを確認
		if len(chain.PlanData.Rooms) != 1 {
			t.Errorf("Seed %d: Expected 1 room, got %d", seed, len(chain.PlanData.Rooms))
		}

		// タイル構成を分析してバリエーションを推測
		wallCount := 0
		floorCount := 0

		for _, tile := range chain.PlanData.Tiles {
			if tile.Walkable {
				floorCount++
			} else {
				wallCount++
			}
		}

		// 壁と床の比率から大まかなバリエーションを判定
		ratio := float64(wallCount) / float64(wallCount+floorCount)
		variantType := ""

		if ratio <= 0.36 {
			variantType = "basic"
		} else if ratio <= 0.45 {
			variantType = "pillars_obstacles_platform"
		} else {
			variantType = "maze"
		}

		variantCounts[variantType]++

		t.Logf("Seed %d: walls=%d, floors=%d, ratio=%.3f, type=%s",
			seed, wallCount, floorCount, ratio, variantType)
	}

	// 複数のバリエーションが生成されていることを確認
	if len(variantCounts) < 2 {
		t.Errorf("Expected multiple variants to be generated, got: %v", variantCounts)
	}

	t.Logf("Variant distribution: %v", variantCounts)
}

func TestBigRoomPlannerReproducibility(t *testing.T) {
	t.Parallel()

	// 同じシードで複数回生成して同じ結果になることを確認
	width, height := gc.Tile(15), gc.Tile(15)
	seed := uint64(99999)

	// 1回目の生成
	chain1 := NewBigRoomPlanner(width, height, seed)
	chain1.PlanData.RawMaster = createTestRawMaster()
	chain1.Plan()

	// 2回目の生成
	chain2 := NewBigRoomPlanner(width, height, seed)
	chain2.PlanData.RawMaster = createTestRawMaster()
	chain2.Plan()

	// 部屋数が同じことを確認
	if len(chain1.PlanData.Rooms) != len(chain2.PlanData.Rooms) {
		t.Errorf("部屋数が異なります。1回目: %d, 2回目: %d",
			len(chain1.PlanData.Rooms), len(chain2.PlanData.Rooms))
	}

	// 部屋の位置とサイズが同じことを確認
	for i := range chain1.PlanData.Rooms {
		room1 := chain1.PlanData.Rooms[i]
		room2 := chain2.PlanData.Rooms[i]

		if room1 != room2 {
			t.Errorf("部屋[%d]が異なります。1回目: %+v, 2回目: %+v", i, room1, room2)
		}
	}

	// タイル配置が同じことを確認
	for i, tile1 := range chain1.PlanData.Tiles {
		tile2 := chain2.PlanData.Tiles[i]
		if tile1 != tile2 {
			t.Errorf("タイル[%d]が異なります。1回目: %v, 2回目: %v", i, tile1, tile2)
		}
	}
}

func TestBigRoomPlannerBoundaries(t *testing.T) {
	t.Parallel()

	// 境界の処理が正しいかを確認
	width, height := gc.Tile(10), gc.Tile(10)
	seed := uint64(11111)

	chain := NewBigRoomPlanner(width, height, seed)
	chain.PlanData.RawMaster = createTestRawMaster()
	chain.Plan()

	// マップの境界が壁になっていることを確認
	for x := 0; x < int(width); x++ {
		// 上端
		idx := chain.PlanData.Level.XYTileIndex(gc.Tile(x), gc.Tile(0))
		if chain.PlanData.Tiles[idx] != chain.PlanData.GenerateTile("Wall") {
			t.Errorf("上端の境界[%d,0]が壁になっていません: %v", x, chain.PlanData.Tiles[idx])
		}

		// 下端
		idx = chain.PlanData.Level.XYTileIndex(gc.Tile(x), height-1)
		if chain.PlanData.Tiles[idx] != chain.PlanData.GenerateTile("Wall") {
			t.Errorf("下端の境界[%d,%d]が壁になっていません: %v", x, height-1, chain.PlanData.Tiles[idx])
		}
	}

	for y := 0; y < int(height); y++ {
		// 左端
		idx := chain.PlanData.Level.XYTileIndex(gc.Tile(0), gc.Tile(y))
		if chain.PlanData.Tiles[idx] != chain.PlanData.GenerateTile("Wall") {
			t.Errorf("左端の境界[0,%d]が壁になっていません: %v", y, chain.PlanData.Tiles[idx])
		}

		// 右端
		idx = chain.PlanData.Level.XYTileIndex(width-1, gc.Tile(y))
		if chain.PlanData.Tiles[idx] != chain.PlanData.GenerateTile("Wall") {
			t.Errorf("右端の境界[%d,%d]が壁になっていません: %v", width-1, y, chain.PlanData.Tiles[idx])
		}
	}
}
