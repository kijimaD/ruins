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
