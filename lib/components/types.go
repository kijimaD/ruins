package components

// 最大値と現在値を持つようなパラメータ
type Pool struct {
	Max     int // 計算式で算出される
	Current int // 現在値
}

// 変動するパラメータ値
type Attribute struct {
	Base     int // 固有の値
	Modifier int // 装備などで変動する値
	Total    int // 足し合わせた現在値。算出される値のメモ
}

// 選択対象
type TargetType struct {
	TargetGroup TargetGroupType // 対象グループ（味方、敵など）
	TargetNum   TargetNumType   // 対象数（単体、複数、全体）
}

// 合成の元になる素材
type RecipeInput struct {
	Name   string // 素材名
	Amount int    // 必要量
}

// 装備品のオプショナルな性能。武器・防具で共通する
type EquipBonus struct {
	Vitality  int // 体力ボーナス
	Strength  int // 筋力ボーナス
	Sensation int // 感覚ボーナス
	Dexterity int // 器用ボーナス
	Agility   int // 敏捷ボーナス

	// 残り項目:
	// - 火属性などの属性耐性
	// - 頑丈+1、連射+2などのスキル
}

// 装備スロット番号。0始まり
type EquipmentSlotNumber int

// ================
// 量を計算するためのインターフェース
type Amounter interface {
	Amount() // 量計算を識別するマーカーメソッド
}

var _ Amounter = RatioAmount{}

// 倍率指定
type RatioAmount struct {
	Ratio float64 // 倍率
}

// Amounterインターフェースの実装
func (ra RatioAmount) Amount() {}

// 倍率と基準値から実際の量を計算する
func (ra RatioAmount) Calc(base int) int {
	return int(float64(base) * ra.Ratio)
}

var _ Amounter = NumeralAmount{}

// 絶対量指定
type NumeralAmount struct {
	Numeral int // 絶対量
}

// Amounterインターフェースの実装
func (na NumeralAmount) Amount() {}

// 固定の数値量を返す
func (na NumeralAmount) Calc() int {
	return na.Numeral
}
