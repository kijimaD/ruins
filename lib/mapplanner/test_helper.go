package mapplanner

import "github.com/kijimaD/ruins/lib/raw"

// CreateTestRawMaster はテスト用の raw.Master インスタンスを作成する
func CreateTestRawMaster() *raw.Master {
	// テスト用の基本的なタイルデータを定義
	testTiles := []raw.TileRaw{
		{Name: "Wall", Walkable: false},
		{Name: "Floor", Walkable: true},
		{Name: "Empty", Walkable: false},
		{Name: "Dirt", Walkable: true},
	}

	// テスト用のアイテムテーブルを定義
	testItemTables := []raw.ItemTable{
		{
			Name: "通常",
			Entries: []raw.ItemTableEntry{
				{ItemName: "回復薬", Weight: 1.0},
				{ItemName: "回復スプレー", Weight: 0.8},
				{ItemName: "手榴弾", Weight: 0.5},
			},
		},
		{
			Name: "洞窟",
			Entries: []raw.ItemTableEntry{
				{ItemName: "回復薬", Weight: 1.0},
				{ItemName: "毒消し", Weight: 0.8},
				{ItemName: "黒曜石", Weight: 0.6},
			},
		},
		{
			Name: "森",
			Entries: []raw.ItemTableEntry{
				{ItemName: "回復薬", Weight: 1.0},
				{ItemName: "緑ハーブ", Weight: 1.2},
			},
		},
		{
			Name: "廃墟",
			Entries: []raw.ItemTableEntry{
				{ItemName: "回復薬", Weight: 1.0},
				{ItemName: "銀の欠片", Weight: 0.8},
			},
		},
	}

	// インデックスを作成
	tileIndex := make(map[string]int)
	for i, tile := range testTiles {
		tileIndex[tile.Name] = i
	}

	itemTableIndex := make(map[string]int)
	for i, table := range testItemTables {
		itemTableIndex[table.Name] = i
	}

	return &raw.Master{
		Raws: raw.Raws{
			Tiles:      testTiles,
			ItemTables: testItemTables,
		},
		TileIndex:      tileIndex,
		ItemTableIndex: itemTableIndex,
	}
}
