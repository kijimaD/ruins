package components

import (
	"log"
)

type warpMode string

var (
	WarpModeNext   = warpMode("NEXT")
	WarpModeEscape = warpMode("ESCAPE")
)

type TargetNum string

var (
	TargetSingle = TargetNum("SINGLE")
	TargetAll    = TargetNum("ALL")
)

type TargetFactionType string

var (
	TargetFactionAlly  = TargetFactionType("ALLY")  // 味方
	TargetFactionEnemy = TargetFactionType("ENEMY") //  敵
	TargetFactionNone  = TargetFactionType("NONE")  // なし
)

type UsableSceneType string

var (
	UsableSceneBattle = UsableSceneType("BATTLE") // 戦闘
	UsableSceneField  = UsableSceneType("FIELD")  // フィールド
	UsableSceneAny    = UsableSceneType("ANY")    // いつでも
)

// ================
// 装備スロット

type EquipmentSlotType string

var (
	EquipmentSlotHead    = EquipmentSlotType("HEAD")    // 頭部
	EquipmentSlotTorso   = EquipmentSlotType("TORSO")   // 胴体
	EquipmentSlotLegs    = EquipmentSlotType("LEGS")    // 脚
	EquipmentSlotJewelry = EquipmentSlotType("JEWELRY") // アクセサリ
)

func StringToEquipmentSlotType(input string) EquipmentSlotType {
	var result EquipmentSlotType
	switch input {
	case string(EquipmentSlotHead):
		result = EquipmentSlotHead
	case string(EquipmentSlotTorso):
		result = EquipmentSlotTorso
	case string(EquipmentSlotLegs):
		result = EquipmentSlotLegs
	case string(EquipmentSlotJewelry):
		result = EquipmentSlotJewelry
	default:
		log.Fatal("invalid equiment slot type")
	}

	return result
}

func (es EquipmentSlotType) String() string {
	var result string
	switch es {
	case EquipmentSlotHead:
		result = "頭部"
	case EquipmentSlotTorso:
		result = "胴体"
	case EquipmentSlotLegs:
		result = "脚部"
	case EquipmentSlotJewelry:
		result = "装飾"
	default:
		log.Fatal("invalid equiment slot type")
	}
	return result
}

// ================
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
		result = "火"
	case DamageAttrThunder:
		result = "電"
	case DamageAttrChill:
		result = "冷"
	case DamageAttrPhoton:
		result = "光"
	default:
		log.Fatal("invalid damage attr type")
	}
	return result
}
