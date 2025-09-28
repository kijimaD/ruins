package mapplanner

// FillAll は全体を指定したタイルで埋めるビルダー
type FillAll struct {
	TileName string
}

// NewFillAll は新しいFillAllビルダーを作成する
func NewFillAll(tileName string) FillAll {
	return FillAll{
		TileName: tileName,
	}
}

// PlanMeta はメタデータをビルドする
func (b FillAll) PlanMeta(planData *MetaPlan) {
	b.build(planData)
}

func (b FillAll) build(planData *MetaPlan) {
	// 全体を指定したタイルで埋める
	for i := range planData.Tiles {
		planData.Tiles[i] = planData.GenerateTile(b.TileName)
	}
}
