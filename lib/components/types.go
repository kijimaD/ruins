package components

import "log"

// 対象
type TargetType struct {
	TargetFaction TargetFactionType // 対象派閥
	TargetNum     TargetNum
}

// 合成の元になる素材
type RecipeInput struct {
	Name   string
	Amount int
}

// 攻撃属性
type DamageAttrType string

var (
	DamageAttrNone    DamageAttrType = "NONE"
	DamageAttrFire    DamageAttrType = "FIRE"
	DamageAttrThunder DamageAttrType = "THUNDER"
	DamageAttrChill   DamageAttrType = "CHILL"
	DamageAttrPhoton  DamageAttrType = "PHOTON"
)

func StringToDamangeAttrType(input string) DamageAttrType {
	var result DamageAttrType
	switch input {
	case string(DamageAttrNone):
		result = DamageAttrNone
	case string(DamageAttrFire):
		result = DamageAttrFire
	case string(DamageAttrThunder):
		result = DamageAttrThunder
	case string(DamageAttrChill):
		result = DamageAttrChill
	case string(DamageAttrPhoton):
		result = DamageAttrPhoton
	default:
		log.Fatal("invalid damage attr type")
	}
	return result
}

func (da DamageAttrType) String() string {
	var result string
	switch da {
	case DamageAttrNone:
		result = "無属性"
	case DamageAttrFire:
		result = "火炎"
	case DamageAttrThunder:
		result = "電撃"
	case DamageAttrChill:
		result = "冷気"
	case DamageAttrPhoton:
		result = "光"
	}
	return result
}
