package consts

import gc "github.com/kijimaD/ruins/lib/components"

const (
	// Game window dimensions
	MinGameWidth  = 960
	MinGameHeight = 720

	// Tile dimensions
	TileSize gc.Pixel = 32

	// UI Labels
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
