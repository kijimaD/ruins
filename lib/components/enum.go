package components

import (
	"log"

	"github.com/pkg/errors"
)

var ErrInvalidEnumType = errors.New("enumに無効な値が指定された")

type warpMode string

const (
	WarpModeNext   = warpMode("NEXT")
	WarpModeEscape = warpMode("ESCAPE")
)

// ================

type TargetNumType string

const (
	TargetSingle = TargetNumType("SINGLE")
	TargetAll    = TargetNumType("ALL")
)

func (enum TargetNumType) Valid() error {
	switch enum {
	case TargetSingle, TargetAll:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

// ================

type TargetFactionType string

const (
	TargetFactionAlly  = TargetFactionType("ALLY")  // 味方
	TargetFactionEnemy = TargetFactionType("ENEMY") //  敵
	TargetFactionNone  = TargetFactionType("NONE")  // なし
)

func (enum TargetFactionType) Valid() error {
	switch enum {
	case TargetFactionAlly, TargetFactionEnemy, TargetFactionNone:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

// ================

type UsableSceneType string

const (
	UsableSceneBattle = UsableSceneType("BATTLE") // 戦闘
	UsableSceneField  = UsableSceneType("FIELD")  // フィールド
	UsableSceneAny    = UsableSceneType("ANY")    // いつでも
)

func (enum UsableSceneType) Valid() error {
	switch enum {
	case UsableSceneBattle, UsableSceneField, UsableSceneAny:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

// ================
// 武器種別

type WeaponType string

const (
	WeaponSword   = WeaponType("SWORD")   // 刀剣
	WeaponSpear   = WeaponType("SPEAR")   // 長物
	WeaponHandgun = WeaponType("HANDGUN") // 拳銃
	WeaponRifle   = WeaponType("RIFLE")   // 小銃
)

func (enum WeaponType) Valid() error {
	switch enum {
	case WeaponSword, WeaponSpear, WeaponHandgun, WeaponRifle:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

func (wc WeaponType) String() string {
	var result string
	switch wc {
	case WeaponSword:
		result = "刀剣"
	case WeaponSpear:
		result = "長物"
	case WeaponHandgun:
		result = "拳銃"
	case WeaponRifle:
		result = "小銃"
	default:
		log.Fatal("invalid weapon type")
	}

	return result
}

// ================
// 装備品種別

type EquipmentType string

const (
	EquipmentHead    = EquipmentType("HEAD")    // 頭部
	EquipmentTorso   = EquipmentType("TORSO")   // 胴体
	EquipmentLegs    = EquipmentType("LEGS")    // 脚
	EquipmentJewelry = EquipmentType("JEWELRY") // アクセサリ
)

func (enum EquipmentType) Valid() error {
	switch enum {
	case EquipmentHead, EquipmentTorso, EquipmentLegs, EquipmentJewelry:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

func (es EquipmentType) String() string {
	var result string
	switch es {
	case EquipmentHead:
		result = "頭部"
	case EquipmentTorso:
		result = "胴体"
	case EquipmentLegs:
		result = "脚部"
	case EquipmentJewelry:
		result = "装飾"
	default:
		log.Fatal("invalid equiment slot type")
	}
	return result
}

// ================
// 攻撃属性

type DamageAttrType string

const (
	DamageAttrNone    DamageAttrType = "NONE"
	DamageAttrFire    DamageAttrType = "FIRE"
	DamageAttrThunder DamageAttrType = "THUNDER"
	DamageAttrChill   DamageAttrType = "CHILL"
	DamageAttrPhoton  DamageAttrType = "PHOTON"
)

func (enum DamageAttrType) Valid() error {
	switch enum {
	case DamageAttrNone, DamageAttrFire, DamageAttrThunder, DamageAttrChill, DamageAttrPhoton:
		return nil
	}
	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

func (da DamageAttrType) String() string {
	var result string
	switch da {
	case DamageAttrNone:
		result = "無"
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

// ================
// 装備スロット番号

type EquipmentSlotNumber int

const (
	EquipmentSlotZero EquipmentSlotNumber = iota
	EquipmentSlotOne
	EquipmentSlotTwo
	EquipmentSlotThree
)

// ================
// 値タイプ

type ValueType string

const (
	PercentageType ValueType = "PERCENTAGE"
	NumeralType              = "NUMERAL"
)

func (enum ValueType) Valid() error {
	switch enum {
	case PercentageType, NumeralType:
		return nil
	}
	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}
