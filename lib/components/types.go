package components

// 対象
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
	// 火属性などの属性耐性
	// 頑丈+1、貫通+2などのパッシブスキル
	// 「救護」「乱射」などのアクティブスキル
}
