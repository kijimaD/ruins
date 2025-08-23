package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/stretchr/testify/assert"
)

func TestSmallRoomBuilder(t *testing.T) {
	t.Parallel()

	t.Run("SmallRoomBuilderが正しいビルダーチェーンを作成", func(t *testing.T) {
		t.Parallel()
		width := gc.Row(20)
		height := gc.Col(20)

		chain := NewSmallRoomBuilder(width, height, 12345)

		// チェーンが作成されていることを確認
		assert.NotNil(t, chain, "チェーンがnilである")
		assert.NotNil(t, chain.Starter, "Starterが設定されていない")

		// 期待されるビルダーが正しい順序で追加されているか確認
		assert.Len(t, chain.Builders, 4, "ビルダーの数が正しくない")

		// ビルダーの型を確認
		_, ok0 := chain.Builders[0].(FillAll)
		assert.True(t, ok0, "1番目のビルダーがFillAllでない")

		_, ok1 := chain.Builders[1].(RoomDraw)
		assert.True(t, ok1, "2番目のビルダーがRoomDrawでない")

		_, ok2 := chain.Builders[2].(LineCorridorBuilder)
		assert.True(t, ok2, "3番目のビルダーがLineCorridorBuilderでない")

		_, ok3 := chain.Builders[3].(BoundaryWall)
		assert.True(t, ok3, "4番目のビルダーがBoundaryWallでない")
	})

	t.Run("SmallRoomBuilderでマップを生成", func(t *testing.T) {
		t.Parallel()
		width := gc.Row(30)
		height := gc.Col(30)

		chain := NewSmallRoomBuilder(width, height, 12345)

		// ビルド実行
		chain.Build()

		// タイル数が正しいことを確認
		expectedCount := int(width) * int(height)
		assert.Len(t, chain.BuildData.Tiles, expectedCount, "タイル数が正しくない")

		// 部屋が生成されていることを確認
		assert.NotEmpty(t, chain.BuildData.Rooms, "部屋が生成されていない")

		// 床タイルが存在することを確認
		floorCount := 0
		wallCount := 0
		for _, tile := range chain.BuildData.Tiles {
			switch tile {
			case TileFloor:
				floorCount++
			case TileWall:
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
		width := gc.Row(25)
		height := gc.Col(25)

		chain := NewSmallRoomBuilder(width, height, 12345)
		chain.Build()

		// 各部屋が有効な範囲内にあることを確認
		for i, room := range chain.BuildData.Rooms {
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
		width := gc.Row(20)
		height := gc.Col(20)

		chain := NewSmallRoomBuilder(width, height, 12345)
		chain.Build()

		// 各部屋の内部の少なくとも一部が床タイルであることを確認
		for i, room := range chain.BuildData.Rooms {
			// 部屋内に床タイルがあるか確認
			hasFloor := false
			for x := room.X1; x <= room.X2 && !hasFloor; x++ {
				for y := room.Y1; y <= room.Y2 && !hasFloor; y++ {
					idx := chain.BuildData.Level.XYTileIndex(x, y)
					if idx >= 0 && int(idx) < len(chain.BuildData.Tiles) {
						if chain.BuildData.Tiles[idx] == TileFloor {
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
			width  gc.Row
			height gc.Col
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
				chain := NewSmallRoomBuilder(tc.width, tc.height, 12345)

				// パニックなく実行できることを確認
				assert.NotPanics(t, func() {
					chain.Build()
				}, "%sでパニックが発生した", tc.name)

				// タイル数が正しいことを確認
				expectedCount := int(tc.width) * int(tc.height)
				assert.Len(t, chain.BuildData.Tiles, expectedCount,
					"%sのタイル数が正しくない", tc.name)
			})
		}
	})

	t.Run("廊下が生成されている", func(t *testing.T) {
		t.Parallel()
		width := gc.Row(30)
		height := gc.Col(30)

		chain := NewSmallRoomBuilder(width, height, 12345)
		chain.Build()

		// 廊下が生成されていることを確認
		// LineCorridorBuilderの仕様により、部屋が2つ以上ある場合は廊下が生成される
		if len(chain.BuildData.Rooms) >= 2 {
			assert.NotEmpty(t, chain.BuildData.Corridors, "部屋が2つ以上あるのに廊下が生成されていない")
		}
	})
}
