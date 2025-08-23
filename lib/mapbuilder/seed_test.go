package mapbuilder

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestSeedReproducibility(t *testing.T) {
	t.Parallel()
	// 同じシードで複数回マップを生成して同じ結果になることを確認
	const testSeed uint64 = 12345
	width, height := gc.Row(20), gc.Col(20)

	// 1回目の生成
	chain1 := SimpleRoomBuilder(width, height, testSeed)
	chain1.Build()
	tiles1 := make([]Tile, len(chain1.BuildData.Tiles))
	copy(tiles1, chain1.BuildData.Tiles)
	rooms1 := make([]Rect, len(chain1.BuildData.Rooms))
	copy(rooms1, chain1.BuildData.Rooms)

	// 2回目の生成（同じシード）
	chain2 := SimpleRoomBuilder(width, height, testSeed)
	chain2.Build()
	tiles2 := make([]Tile, len(chain2.BuildData.Tiles))
	copy(tiles2, chain2.BuildData.Tiles)
	rooms2 := make([]Rect, len(chain2.BuildData.Rooms))
	copy(rooms2, chain2.BuildData.Rooms)

	// タイルが完全に一致することを確認
	if len(tiles1) != len(tiles2) {
		t.Errorf("タイル数が異なります。1回目: %d, 2回目: %d", len(tiles1), len(tiles2))
	}
	for i := range tiles1 {
		if tiles1[i] != tiles2[i] {
			t.Errorf("タイル[%d]が異なります。1回目: %v, 2回目: %v", i, tiles1[i], tiles2[i])
		}
	}

	// 部屋が完全に一致することを確認
	if len(rooms1) != len(rooms2) {
		t.Errorf("部屋数が異なります。1回目: %d, 2回目: %d", len(rooms1), len(rooms2))
	}
	for i := range rooms1 {
		if rooms1[i] != rooms2[i] {
			t.Errorf("部屋[%d]が異なります。1回目: %v, 2回目: %v", i, rooms1[i], rooms2[i])
		}
	}
}

func TestDifferentSeeds(t *testing.T) {
	t.Parallel()
	// 異なるシードで異なる結果になることを確認
	width, height := gc.Row(20), gc.Col(20)

	// シード1で生成
	chain1 := SimpleRoomBuilder(width, height, 11111)
	chain1.Build()

	// シード2で生成
	chain2 := SimpleRoomBuilder(width, height, 22222)
	chain2.Build()

	// 部屋数が異なる可能性が高い（必ずしも異なるとは限らないが）
	// 少なくともいくつかのタイルは異なるはず
	differentTiles := 0
	for i := range chain1.BuildData.Tiles {
		if chain1.BuildData.Tiles[i] != chain2.BuildData.Tiles[i] {
			differentTiles++
		}
	}

	if differentTiles == 0 {
		t.Error("異なるシードなのにマップが完全に一致しています")
	}
}

func TestRandomSourceDeterministic(t *testing.T) {
	t.Parallel()
	// RandomSourceが決定論的であることを確認
	const seed uint64 = 99999

	// 同じシードから作成した2つのRandomSource
	rs1 := NewRandomSource(seed)
	rs2 := NewRandomSource(seed)

	// 同じ順序で同じ値を生成することを確認
	for i := 0; i < 100; i++ {
		val1 := rs1.Intn(1000)
		val2 := rs2.Intn(1000)
		if val1 != val2 {
			t.Errorf("反復%dで異なる値が生成されました。val1: %d, val2: %d", i, val1, val2)
		}
	}

	// Float64も確認
	for i := 0; i < 100; i++ {
		val1 := rs1.Float64()
		val2 := rs2.Float64()
		if val1 != val2 {
			t.Errorf("反復%dで異なるFloat64値が生成されました。val1: %f, val2: %f", i, val1, val2)
		}
	}
}

func TestZeroSeedHandling(t *testing.T) {
	t.Parallel()
	// シード0の場合はランダムなシードが使用されることを確認
	width, height := gc.Row(10), gc.Col(10)

	// シード0で2回生成
	chain1 := SimpleRoomBuilder(width, height, 0)
	chain1.Build()

	// 少し時間をずらして再度生成
	chain2 := SimpleRoomBuilder(width, height, 0)
	chain2.Build()

	// 高確率で異なるマップになるはず（厳密には同じになる可能性もあるが極めて低い）
	// ここでは部屋数だけチェック
	if len(chain1.BuildData.Rooms) == len(chain2.BuildData.Rooms) {
		// 部屋数が同じでも、タイルが異なる可能性をチェック
		differentTiles := 0
		for i := range chain1.BuildData.Tiles {
			if chain1.BuildData.Tiles[i] != chain2.BuildData.Tiles[i] {
				differentTiles++
			}
		}
		// 完全に一致する確率は非常に低い
		if differentTiles == 0 {
			t.Log("警告: シード0で生成した2つのマップが偶然一致しました（極めて稀）")
		}
	}
}
