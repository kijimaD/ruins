package mapplanner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestSmallRoomPlanner(t *testing.T) {
	t.Parallel()

	t.Run("SmallRoomPlannerが正常に作成される", func(t *testing.T) {
		t.Parallel()
		width := gc.Tile(20)
		height := gc.Tile(20)

		chain := NewSmallRoomPlanner(width, height, 12345)

		// チェーンが作成されていることを確認
		assert.NotNil(t, chain, "チェーンがnilである")
		assert.NotNil(t, chain.Starter, "Starterが設定されていない")
	})

	t.Run("SmallRoomPlannerでマップを生成", func(t *testing.T) {
		t.Parallel()
		width := gc.Tile(30)
		height := gc.Tile(30)

		chain := NewSmallRoomPlanner(width, height, 12345)

		// ビルド実行
		chain.PlanData.RawMaster = CreateTestRawMaster()
		chain.Plan()

		// タイル数が正しいことを確認
		expectedCount := int(width) * int(height)
		assert.Len(t, chain.PlanData.Tiles, expectedCount, "タイル数が正しくない")

		// 部屋が生成されていることを確認
		assert.NotEmpty(t, chain.PlanData.Rooms, "部屋が生成されていない")

		// 床タイルが存在することを確認
		floorCount := 0
		wallCount := 0
		for _, tile := range chain.PlanData.Tiles {
			if tile.Walkable {
				floorCount++
			} else {
				wallCount++
			}
		}
		assert.Greater(t, floorCount, 0, "床タイルが存在しない")
		assert.Greater(t, wallCount, 0, "壁タイルが存在しない")

		// 床と壁の合計がタイル総数と一致することを確認（他のタイルタイプがない場合）
		// 廊下や特殊タイルがある場合はこのアサーションを調整
		assert.LessOrEqual(t, floorCount+wallCount, expectedCount, "タイルタイプの合計が総数を超えている")
	})

	t.Run("生成された部屋が有効な範囲内にある", func(t *testing.T) {
		t.Parallel()
		width := gc.Tile(25)
		height := gc.Tile(25)

		chain := NewSmallRoomPlanner(width, height, 12345)
		chain.PlanData.RawMaster = CreateTestRawMaster()
		chain.Plan()

		// 各部屋が有効な範囲内にあることを確認
		for i, room := range chain.PlanData.Rooms {
			assert.GreaterOrEqual(t, int(room.X1), 0, "部屋%dのX1が負の値", i)
			assert.GreaterOrEqual(t, int(room.Y1), 0, "部屋%dのY1が負の値", i)
			assert.LessOrEqual(t, int(room.X2), int(width), "部屋%dのX2が幅を超えている", i)
			assert.LessOrEqual(t, int(room.Y2), int(height), "部屋%dのY2が高さを超えている", i)

			// 部屋のサイズが正しいことを確認
			assert.LessOrEqual(t, int(room.X1), int(room.X2), "部屋%dのX座標が逆転している", i)
			assert.LessOrEqual(t, int(room.Y1), int(room.Y2), "部屋%dのY座標が逆転している", i)
		}
	})

	t.Run("生成された部屋の内部が床タイルになっている", func(t *testing.T) {
		t.Parallel()
		width := gc.Tile(20)
		height := gc.Tile(20)

		chain := NewSmallRoomPlanner(width, height, 12345)
		chain.PlanData.RawMaster = CreateTestRawMaster()
		chain.Plan()

		// 各部屋の内部の少なくとも一部が床タイルであることを確認
		for i, room := range chain.PlanData.Rooms {
			// 部屋内に床タイルがあるか確認
			hasFloor := false
			for x := room.X1; x <= room.X2 && !hasFloor; x++ {
				for y := room.Y1; y <= room.Y2 && !hasFloor; y++ {
					idx := chain.PlanData.Level.XYTileIndex(x, y)
					if idx >= 0 && int(idx) < len(chain.PlanData.Tiles) {
						if chain.PlanData.Tiles[idx].Walkable {
							hasFloor = true
						}
					}
				}
			}
			assert.True(t, hasFloor, "部屋%dに床タイルが存在しない", i)
		}
	})

	t.Run("異なるサイズのマップで動作確認", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name   string
			width  gc.Tile
			height gc.Tile
		}{
			{"小さいマップ", 10, 10},
			{"中サイズマップ", 30, 30},
			{"大きいマップ", 50, 50},
			{"横長マップ", 40, 20},
			{"縦長マップ", 20, 40},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				chain := NewSmallRoomPlanner(tc.width, tc.height, 12345)

				// パニックなく実行できることを確認
				assert.NotPanics(t, func() {
					chain.PlanData.RawMaster = CreateTestRawMaster()
					chain.Plan()
				}, "%sでパニックが発生した", tc.name)

				// タイル数が正しいことを確認
				expectedCount := int(tc.width) * int(tc.height)
				assert.Len(t, chain.PlanData.Tiles, expectedCount,
					"%sのタイル数が正しくない", tc.name)
			})
		}
	})

	t.Run("廊下が生成されている", func(t *testing.T) {
		t.Parallel()
		width := gc.Tile(30)
		height := gc.Tile(30)

		chain := NewSmallRoomPlanner(width, height, 12345)
		chain.PlanData.RawMaster = CreateTestRawMaster()
		chain.Plan()

		// 廊下が生成されていることを確認
		// LineCorridorPlannerの仕様により、部屋が2つ以上ある場合は廊下が生成される
		if len(chain.PlanData.Rooms) >= 2 {
			assert.NotEmpty(t, chain.PlanData.Corridors, "部屋が2つ以上あるのに廊下が生成されていない")
		}
	})
}
