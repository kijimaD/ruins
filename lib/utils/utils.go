package utils

import gc "github.com/kijimaD/ruins/lib/components"

// ========== 定数定義 ==========

const (
	// ゲームウィンドウの寸法
	MinGameWidth  = 960
	MinGameHeight = 720

	// タイルの寸法
	TileSize gc.Pixel = 32

	// UIラベル
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

// ビルド時に挿入する
var AppVersion = "v0.0.0"

// ========== 汎用ユーティリティ ==========

func GetPtr[T any](x T) *T {
	return &x
}

// ========== 数学ユーティリティ ==========

// xとyの小さい方を返す
func Min[T int | float64](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// xとyの大きい方を返す
func Max[T int | float64](x, y T) T {
	if x > y {
		return x
	}
	return y
}

// valueを[min, max]の範囲に制限する
func Clamp[T int | float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// xの絶対値を返す
func Abs[T int | float64](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
