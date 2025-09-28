package mapplanner

import "github.com/kijimaD/ruins/lib/raw"

// createTestRawMaster はテスト用の raw.Master インスタンスを作成する
func createTestRawMaster() *raw.Master {
	// テスト用の基本的なタイルデータを定義
	testTiles := []raw.TileRaw{
		{Name: "Wall", Walkable: false},
		{Name: "Floor", Walkable: true},
		{Name: "Empty", Walkable: false},
	}

	// インデックスを作成
	tileIndex := make(map[string]int)
	for i, tile := range testTiles {
		tileIndex[tile.Name] = i
	}

	return &raw.Master{
		Raws: raw.Raws{
			Tiles: testTiles,
		},
		TileIndex: tileIndex,
	}
}
