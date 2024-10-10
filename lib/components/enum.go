package components

import (
	"log"

	"github.com/pkg/errors"
)

var ErrInvalidEnumType = errors.New("enumに無効な値が指定された")

// ================

// ワープモード
type warpMode string

const (
	WarpModeNext   = warpMode("NEXT")
	WarpModeEscape = warpMode("ESCAPE")
)

// ================

// ターゲット数
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

// ターゲットの種別
type TargetFactionType string

const (
	TargetFactionAlly  = TargetFactionType("ALLY")  // 味方
	TargetFactionEnemy = TargetFactionType("ENEMY") // 敵
	TargetFactionCard  = TargetFactionType("CARD")  // カード
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

// 使えるシーン
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
// 種別によって適用する計算式が異なる
type AttackType string

const (
	AttackSword   = AttackType("SWORD")   // 刀剣
	AttackSpear   = AttackType("SPEAR")   // 長物
	AttackHandgun = AttackType("HANDGUN") // 拳銃
	AttackRifle   = AttackType("RIFLE")   // 小銃
	AttackFist    = AttackType("FIST")    // 格闘
)

func (enum AttackType) Valid() error {
	switch enum {
	case AttackSword, AttackSpear, AttackHandgun, AttackRifle, AttackFist:
		return nil
	}

	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

func (at AttackType) String() string {
	var result string
	switch at {
	case AttackSword:
		result = "刀剣"
	case AttackSpear:
		result = "長物"
	case AttackHandgun:
		result = "拳銃"
	case AttackRifle:
		result = "小銃"
	case AttackFist:
		result = "格闘"
	default:
		log.Fatal("invalid attack type")
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
type ElementType string

const (
	ElementTypeNone    ElementType = "NONE"
	ElementTypeFire    ElementType = "FIRE"
	ElementTypeThunder ElementType = "THUNDER"
	ElementTypeChill   ElementType = "CHILL"
	ElementTypePhoton  ElementType = "PHOTON"
)

func (enum ElementType) Valid() error {
	switch enum {
	case ElementTypeNone, ElementTypeFire, ElementTypeThunder, ElementTypeChill, ElementTypePhoton:
		return nil
	}
	return errors.Wrapf(ErrInvalidEnumType, "get %s", enum)
}

func (et ElementType) String() string {
	var result string
	switch et {
	case ElementTypeNone:
		result = "無"
	case ElementTypeFire:
		result = "火"
	case ElementTypeThunder:
		result = "電"
	case ElementTypeChill:
		result = "冷"
	case ElementTypePhoton:
		result = "光"
	default:
		log.Fatal("invalid element type")
	}
	return result
}

// ================

// 装備スロット番号。0始まり
type EquipmentSlotNumber int
