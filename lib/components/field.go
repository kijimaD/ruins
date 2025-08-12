package components

// Velocity は移動属性
type Velocity struct {
	// 角度(度数法)。この角度分スプライトを回転させる
	Angle float64
	// 現在の速度
	Speed float64
	// スロットルモード
	ThrottleMode ThrottleMode
	// 最高速度
	MaxSpeed float64
}

// ThrottleMode はスロットルモード
type ThrottleMode string

const (
	// ThrottleModeFront は前進スロットルモード
	ThrottleModeFront = ThrottleMode("FRONT")
	// ThrottleModeBack は後退スロットルモード
	ThrottleModeBack = ThrottleMode("BACK")
	// ThrottleModeNope はスロットルなし
	ThrottleModeNope = ThrottleMode("NOPE")
)

// Position はフィールド上に座標をもって存在する。移動体に対して使う
// スプライトはこの位置に中心を合わせて配置する
// -----
// |   |
// | + |
// |   |
// -----
type Position struct {
	X Pixel
	Y Pixel
}

// Pixel はピクセル単位。計算用にfloat64
type Pixel float64

// GridElement はフィールド上にグリッドに沿って存在する。静的なステージオブジェクトに使う
// TODO: 縦横の型を共通にする。タイル単位だということがわかればよい。TileUnitとか
// スプライトはグリッドに沿って配置する
// +----
// |   |
// |   |
// |   |
// -----
type GridElement struct {
	Row Row
	Col Col
}

// Row はタイルの横位置。ピクセル数ではない
type Row int

// Col はタイルの縦位置。ピクセル数ではない
type Col int

// BlockPass はフィールド上で通過できない
type BlockPass struct{}

// BlockView はフィールド上で視界を遮る
type BlockView struct{}

// Renderable はフィールド上で描画できる
type Renderable struct{}
