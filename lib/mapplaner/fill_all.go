package mapplaner

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
func (b FillAll) BuildMeta(buildData *PlannerMap) {
	b.build(buildData)
}

func (b FillAll) build(buildData *PlannerMap) {
	// 全体を指定したタイルで埋める
	for i := range buildData.Tiles {
		buildData.Tiles[i] = b.TileType
	}
}
