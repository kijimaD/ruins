package utils

import (
	"github.com/hajimehoshi/ebiten/v2"
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/engine/world"
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
	HPLabel        = "HP"
	SPLabel        = "SP"
	VitalityLabel  = "体力"
	StrengthLabel  = "筋力"
	SensationLabel = "感覚"
	DexterityLabel = "器用"
	AgilityLabel   = "敏捷"
	DefenseLabel   = "防御"

	AccuracyLabel         = "命中"
	DamageLabel           = "攻撃力"
	AttackCountLabel      = "回数"
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
func Clamp[T int | float64](value, minVal, max T) T {
	if value < minVal {
		return minVal
	}
	if value > max {
		return max
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
	gameComponents := world.Components.Game.(*gc.Components)
	var camera *gc.Camera
	var cPos *gc.Position
	world.Manager.Join(
		gameComponents.Camera,
		gameComponents.Position,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		camera = gameComponents.Camera.Get(entity).(*gc.Camera)
		cPos = gameComponents.Position.Get(entity).(*gc.Position)
	}))

	cx, cy := float64(world.Resources.ScreenDimensions.Width/2), float64(world.Resources.ScreenDimensions.Height/2)

	// カメラ位置
	op.GeoM.Translate(float64(-cPos.X), float64(-cPos.Y))
	op.GeoM.Scale(camera.Scale, camera.Scale)
	// 画面の中央
	op.GeoM.Translate(float64(cx), float64(cy))
}
