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

// PlanMeta はメタデータをビルドする
func (b FillAll) PlanMeta(planData *MetaPlan) {
	b.build(planData)
}

func (b FillAll) build(planData *MetaPlan) {
	// 全体を指定したタイルで埋める
	for i := range planData.Tiles {
		planData.Tiles[i] = b.TileType
	}
}
