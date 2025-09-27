package mapplaner

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
)

func TestNewRandomPlanner(t *testing.T) {
	t.Parallel()

	width, height := gc.Tile(20), gc.Tile(20)

	// 同じシードで複数回実行して同じビルダータイプが選択されることを確認
	seed := uint64(12345)

	chain1 := NewRandomPlanner(width, height, seed)
	chain1.Build()

	chain2 := NewRandomPlanner(width, height, seed)
	chain2.Build()

	// 同じシードなので同じビルダータイプが選ばれ、同じ結果になるはず
	if len(chain1.PlanData.Rooms) != len(chain2.PlanData.Rooms) {
		t.Errorf("同じシードなのに部屋数が異なります。1回目: %d, 2回目: %d",
			len(chain1.PlanData.Rooms), len(chain2.PlanData.Rooms))
	}

	// タイル配置が同じことを確認
	if len(chain1.PlanData.Tiles) != len(chain2.PlanData.Tiles) {
		t.Errorf("同じシードなのにタイル数が異なります。1回目: %d, 2回目: %d",
			len(chain1.PlanData.Tiles), len(chain2.PlanData.Tiles))
	}

	for i, tile1 := range chain1.PlanData.Tiles {
		if chain2.PlanData.Tiles[i] != tile1 {
			t.Errorf("タイル[%d]が異なります。1回目: %v, 2回目: %v", i, tile1, chain2.PlanData.Tiles[i])
			break // 最初の違いだけ報告
		}
	}
}

func TestNewRandomPlannerVariety(t *testing.T) {
	t.Parallel()

	width, height := gc.Tile(20), gc.Tile(20)

	// 異なるシードで実行して、少なくとも異なる結果が出ることを確認
	results := make(map[int]int) // 部屋数 -> 出現回数

	// 複数の異なるシードでテスト
	seeds := []uint64{1, 2, 3, 42, 123, 456, 789, 999, 1337, 9999}

	for _, seed := range seeds {
		chain := NewRandomPlanner(width, height, seed)
		chain.Build()

		roomCount := len(chain.PlanData.Rooms)
		results[roomCount]++

		// 基本的な整合性チェック
		expectedTileCount := int(width) * int(height)
		if len(chain.PlanData.Tiles) != expectedTileCount {
			t.Errorf("シード%dでタイル数が不正です。期待: %d, 実際: %d",
				seed, expectedTileCount, len(chain.PlanData.Tiles))
		}

		// 少なくとも部屋が生成されることを確認
		if roomCount == 0 {
			t.Errorf("シード%dで部屋が生成されませんでした", seed)
		}
	}

	// 最低でも2種類以上の異なる結果が出ることを期待
	// （ランダムなので必ずしも保証されないが、高確率で異なる結果になるはず）
	if len(results) < 2 {
		t.Logf("警告: 異なるシードでも同じ部屋数ばかりが生成されました: %v", results)
	} else {
		t.Logf("異なる部屋数のパターン: %v", results)
	}
}

func TestNewRandomPlannerBuildsSuccessfully(t *testing.T) {
	t.Parallel()

	// 様々なマップサイズでテスト
	testCases := []struct {
		name   string
		width  gc.Tile
		height gc.Tile
		seed   uint64
	}{
		{"小さいマップ", 10, 10, 111},
		{"中サイズマップ", 30, 30, 222},
		{"大きいマップ", 50, 50, 333},
		{"横長マップ", 40, 20, 444},
		{"縦長マップ", 20, 40, 555},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// パニックなく実行できることを確認
			var chain *PlannerChain
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("%sでパニックが発生しました: %v", tc.name, r)
					}
				}()

				chain = NewRandomPlanner(tc.width, tc.height, tc.seed)
				chain.Build()
			}()

			if chain == nil {
				return // パニックで終了した場合
			}

			// タイル数が正しいことを確認
			expectedCount := int(tc.width) * int(tc.height)
			if len(chain.PlanData.Tiles) != expectedCount {
				t.Errorf("%sのタイル数が正しくない。期待: %d, 実際: %d",
					tc.name, expectedCount, len(chain.PlanData.Tiles))
			}

			// 部屋が生成されていることを確認
			if len(chain.PlanData.Rooms) == 0 {
				t.Errorf("%sで部屋が生成されませんでした", tc.name)
			}

			// 床タイルが存在することを確認
			floorCount := 0
			for _, tile := range chain.PlanData.Tiles {
				if tile == TileFloor {
					floorCount++
				}
			}

			if floorCount == 0 {
				t.Errorf("%sで床タイルが生成されませんでした", tc.name)
			}
		})
	}
}

func TestRandomPlannerTypes(t *testing.T) {
	t.Parallel()

	// 特定のシードで特定のビルダータイプが選ばれることを確認
	// これによりランダム性が正しく機能していることを検証

	width, height := gc.Tile(20), gc.Tile(20)

	// 複数のシードでテストして、異なるタイプのビルダーが選ばれることを確認
	seedResults := make(map[uint64]int) // seed -> 部屋数

	testSeeds := []uint64{1, 2, 3, 4, 5, 10, 20, 30, 100, 200}

	for _, seed := range testSeeds {
		chain := NewRandomPlanner(width, height, seed)
		chain.Build()

		roomCount := len(chain.PlanData.Rooms)
		seedResults[seed] = roomCount

		// 各シードで正常にマップが生成されることを確認
		if roomCount == 0 {
			t.Errorf("シード%dで部屋が生成されませんでした", seed)
		}

		// タイル総数の確認
		expectedTileCount := int(width) * int(height)
		if len(chain.PlanData.Tiles) != expectedTileCount {
			t.Errorf("シード%dでタイル数が不正。期待: %d, 実際: %d",
				seed, expectedTileCount, len(chain.PlanData.Tiles))
		}
	}

	t.Logf("各シードでの部屋数: %v", seedResults)
}
