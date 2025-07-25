package consts

import (
	gc "github.com/kijimaD/ruins/lib/components"
)

// ========== ウィンドウサイズ ==========

const (
	// MinGameWidth はゲームウィンドウの最小幅
	MinGameWidth = 960
	// MinGameHeight はゲームウィンドウの最小高さ
	MinGameHeight = 720
)

// ========== ゲーム定数 ==========

const (
	// TileSize はタイルの寸法
	TileSize gc.Pixel = 32
)

// ========== UI表示ラベル ==========

const (
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
