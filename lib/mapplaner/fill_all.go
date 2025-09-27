package mapplanner

// FillAll は全体を指定したタイルで埋めるビルダー
type FillAll struct {
	TileType Tile
}

// NewFillAll は新しいFillAllビルダーを作成する
func NewFillAll(tileType Tile) FillAll {
	return FillAll{
		TileType: tileType,
	}
}

// BuildMeta はメタデータをビルドする
func (b FillAll) BuildMeta(buildData *MetaPlan) {
	b.build(buildData)
}

func (b FillAll) build(buildData *MetaPlan) {
	// 全体を指定したタイルで埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = b.TileType
	}
}
