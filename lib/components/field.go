package components

// フィールド上に座標をもって存在する。移動体に対して使う
type Position struct {
	X     int
	Y     int
	Angle float64 // 角度(ラジアン)。この角度分スプライトを回転させる
}

// フィールド上にグリッドに沿って存在する。静的なステージオブジェクトに使う
type GridElement struct {
	Row Row
	Col Col
}

// タイルの横位置。ピクセル数ではない
type Row int

// タイルの縦位置。ピクセル数ではない
type Col int

// フィールド上で通過できない
type BlockPass struct{}

// フィールド上で視界を遮る
type BlockView struct{}

// フィールド上で描画できる
type Renderable struct{}
