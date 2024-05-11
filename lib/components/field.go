package components

// フィールド上に存在する
type Position struct {
	X     int
	Y     int
	Angle float64  // 角度(ラジアン)。この角度分スプライトを回転させる
	Depth DepthNum // 描画順。小さい順に先に(下に)描画する
}

// フィールド上で通過できない
type BlockPass struct{}

// フィールド上で視界を遮る
type BlockView struct{}

// フィールド上で描画できる
type Renderable struct{}
