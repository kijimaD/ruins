package components

// 最大値と現在値を持つようなパラメータ
type Pool struct {
	Max     int // 計算式で算出される
	Current int
}

// 変動するパラメータ値
type Attribute struct {
	Base     int // 固有の値
	Modifier int // 装備などで変動する値
	Total    int // 足し合わせた現在値。算出される値のメモ
}

// 選択対象
type TargetType struct {
	TargetFaction TargetFactionType // 対象派閥
	TargetNum     TargetNumType
}

// 合成の元になる素材
type RecipeInput struct {
	Name   string
	Amount int
}

// 装備品のオプショナルな性能。武器・防具で共通する
type EquipBonus struct {
	Vitality  int
	Strength  int
	Sensation int
	Dexterity int
	Agility   int

	// 残り項目:
	// - 火属性などの属性耐性
	// - 頑丈+1、連射+2などのスキル
}

// 装備スロット番号。0始まり
type EquipmentSlotNumber int

// ================
type Amounter interface {
	Amount()
}

var _ Amounter = RatioAmount{}

// 倍率指定
type RatioAmount struct {
	Ratio float64
}

func (ra RatioAmount) Amount() {}
func (ra RatioAmount) Calc(base int) int {
	return int(float64(base) * ra.Ratio)
}

var _ Amounter = NumeralAmount{}

// 絶対量指定
type NumeralAmount struct {
	Numeral int
}

func (na NumeralAmount) Amount() {}
func (na NumeralAmount) Calc() int {
	return na.Numeral
}
