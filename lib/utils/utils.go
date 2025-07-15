//nolint:revive // utils package name is acceptable for utility functions
package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ========== 定数定義 ==========

const (
	// MinGameWidth はゲームウィンドウの最小幅
	MinGameWidth = 960
	// MinGameHeight はゲームウィンドウの最小高さ
	MinGameHeight = 720

	// TileSize はタイルの寸法
	TileSize gc.Pixel = 32

	// HPLabel はHP表示ラベル
	HPLabel = "HP"
	// SPLabel はSP表示ラベル
	SPLabel = "SP"
	// VitalityLabel は体力表示ラベル
	VitalityLabel = "体力"
	// StrengthLabel は筋力表示ラベル
	StrengthLabel = "筋力"
	// SensationLabel は感覚表示ラベル
	SensationLabel = "感覚"
	// DexterityLabel は器用表示ラベル
	DexterityLabel = "器用"
	// AgilityLabel は敏捷表示ラベル
	AgilityLabel = "敏捷"
	// DefenseLabel は防御表示ラベル
	DefenseLabel = "防御"

	// AccuracyLabel は命中表示ラベル
	AccuracyLabel = "命中"
	// DamageLabel は攻撃力表示ラベル
	DamageLabel = "攻撃力"
	// AttackCountLabel は攻撃回数表示ラベル
	AttackCountLabel = "回数"
	// EquimentCategoryLabel は装備部位表示ラベル
	EquimentCategoryLabel = "部位"
)

// AppVersion はビルド時に挿入されるアプリケーションバージョン
var AppVersion = "v0.0.0"

// ========== 汎用ユーティリティ ==========

// GetPtr は値のポインタを返す
func GetPtr[T any](x T) *T {
	return &x
}

// ========== 数学ユーティリティ ==========

// Min はxとyの小さい方を返す
func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Max はxとyの大きい方を返す
func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// Clamp はvalueを[min, max]の範囲に制限する
func Clamp[T int | float64](value, minVal, maxVal T) T {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

// Abs はxの絶対値を返す
func Abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// ========== カメラユーティリティ ==========

// SetTranslate はカメラを考慮した画像配置オプションをセットする
// TODO: ズーム率を追加する
func SetTranslate(world w.World, op *ebiten.DrawImageOptions) {
	var camera *gc.Camera
	var cPos *gc.Position
	world.Manager.Join(
		world.Components.Camera,
		world.Components.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		camera = world.Components.Camera.Get(entity).(*gc.Camera)
		cPos = world.Components.Position.Get(entity).(*gc.Position)
	}))

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置
	op.GeoM.Translate(float64(-cPos.X), float64(-cPos.Y))
	op.GeoM.Scale(camera.Scale, camera.Scale)
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
}
