package components

// フィールド上に存在する
type Position struct {
	X int
	Y int
}

// フィールド上で通り抜けできない
type BlockPass struct{}

// フィールド上で視界を遮る
type BlockView struct{}

// フィールド上で描画できる
type Renderable struct{}
