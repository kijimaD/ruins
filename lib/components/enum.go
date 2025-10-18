package components

import (
	"errors"
	"fmt"
)

// ErrInvalidEnumType はenumに無効な値が指定された場合のエラー
var ErrInvalidEnumType = errors.New("enumに無効な値が指定された")

// ================

// WarpMode はワープモードを表す
type WarpMode string

const (
	// WarpModeNext は次の階層へワープする
	WarpModeNext = WarpMode("NEXT")
	// WarpModeEscape は脱出ワープする
	WarpModeEscape = WarpMode("ESCAPE")
)

// ================

// TargetNumType はターゲット数を表す
type TargetNumType string

const (
	// TargetSingle は単体ターゲット
	TargetSingle = TargetNumType("SINGLE")
	// TargetAll は全体ターゲット
	TargetAll = TargetNumType("ALL")
)

// Valid はTargetNumTypeの値が有効かを検証する
func (enum TargetNumType) Valid() error {
	switch enum {
	case TargetSingle, TargetAll:
		return nil
	}

	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

// ================

// TargetGroupType は使用者から見たターゲットの種別。相対的な指定なので、所有者が敵グループだと対象グループは逆転する
type TargetGroupType string

const (
	// TargetGroupAlly は味方グループ
	TargetGroupAlly = TargetGroupType("ALLY") // 味方
	// TargetGroupEnemy は敵グループ
	TargetGroupEnemy = TargetGroupType("ENEMY") // 敵
	// TargetGroupWeapon は武器グループ
	TargetGroupWeapon = TargetGroupType("WEAPON") // 武器
	// TargetGroupNone はグループなし
	TargetGroupNone = TargetGroupType("NONE") // なし
)

// Valid はTargetGroupTypeの値が有効かを検証する
func (enum TargetGroupType) Valid() error {
	switch enum {
	case TargetGroupAlly, TargetGroupEnemy, TargetGroupWeapon, TargetGroupNone:
		return nil
	}

	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

// ================

// UsableSceneType は使えるシーンを表す
type UsableSceneType string

const (
	// UsableSceneBattle は戦闘シーン
	UsableSceneBattle = UsableSceneType("BATTLE") // 戦闘
	// UsableSceneField はフィールドシーン
	UsableSceneField = UsableSceneType("FIELD") // フィールド
	// UsableSceneAny はいつでも使えるシーン
	UsableSceneAny = UsableSceneType("ANY") // いつでも
)

// Valid はUsableSceneTypeの値が有効かを検証する
func (enum UsableSceneType) Valid() error {
	switch enum {
	case UsableSceneBattle, UsableSceneField, UsableSceneAny:
		return nil
	}

	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

// ================

// AttackRangeType は攻撃の射程タイプを表す
type AttackRangeType string

const (
	// AttackRangeMelee は近接攻撃
	AttackRangeMelee = AttackRangeType("MELEE")
	// AttackRangeRanged は遠距離攻撃
	AttackRangeRanged = AttackRangeType("RANGED")
)

// AttackType は武器種別を表す。種別によって適用する計算式が異なる
type AttackType struct {
	Type  string          // 武器種別の識別子
	Range AttackRangeType // 近接/遠距離の区分
}

var (
	// AttackSword は刀剣
	AttackSword = AttackType{Type: "SWORD", Range: AttackRangeMelee}
	// AttackSpear は長物
	AttackSpear = AttackType{Type: "SPEAR", Range: AttackRangeMelee}
	// AttackHandgun は拳銃
	AttackHandgun = AttackType{Type: "HANDGUN", Range: AttackRangeRanged}
	// AttackRifle は小銃
	AttackRifle = AttackType{Type: "RIFLE", Range: AttackRangeRanged}
	// AttackFist は格闘
	AttackFist = AttackType{Type: "FIST", Range: AttackRangeMelee}
	// AttackCanon は大砲
	AttackCanon = AttackType{Type: "CANON", Range: AttackRangeRanged}
)

// Valid はAttackTypeの値が有効かを検証する
func (at AttackType) Valid() error {
	switch at.Type {
	case AttackSword.Type, AttackSpear.Type, AttackHandgun.Type, AttackRifle.Type, AttackFist.Type, AttackCanon.Type:
		return nil
	}

	return fmt.Errorf("get %s: %w", at.Type, ErrInvalidEnumType)
}

// IsMelee は近接武器かどうかを返す
func (at AttackType) IsMelee() bool {
	return at.Range == AttackRangeMelee
}

// IsRanged は遠距離武器かどうかを返す
func (at AttackType) IsRanged() bool {
	return at.Range == AttackRangeRanged
}

func (at AttackType) String() string {
	var result string
	switch at.Type {
	case "SWORD":
		result = "刀剣"
	case "SPEAR":
		result = "長物"
	case "HANDGUN":
		result = "拳銃"
	case "RIFLE":
		result = "小銃"
	case "FIST":
		result = "格闘"
	case "CANON":
		result = "大砲"
	default:
		panic("invalid attack type")
	}

	return result
}

// ParseAttackType は文字列からAttackTypeを生成する
func ParseAttackType(s string) AttackType {
	switch s {
	case "SWORD":
		return AttackSword
	case "SPEAR":
		return AttackSpear
	case "HANDGUN":
		return AttackHandgun
	case "RIFLE":
		return AttackRifle
	case "FIST":
		return AttackFist
	case "CANON":
		return AttackCanon
	default:
		panic(fmt.Sprintf("invalid attack type: %s", s))
	}
}

// ================

// EquipmentType は装備品種別を表す
type EquipmentType string

const (
	// EquipmentHead は頭部装備
	EquipmentHead = EquipmentType("HEAD") // 頭部
	// EquipmentTorso は胴体装備
	EquipmentTorso = EquipmentType("TORSO") // 胴体
	// EquipmentLegs は脚装備
	EquipmentLegs = EquipmentType("LEGS") // 脚
	// EquipmentJewelry はアクセサリ装備
	EquipmentJewelry = EquipmentType("JEWELRY") // アクセサリ
)

// Valid はEquipmentTypeの値が有効かを検証する
func (enum EquipmentType) Valid() error {
	switch enum {
	case EquipmentHead, EquipmentTorso, EquipmentLegs, EquipmentJewelry:
		return nil
	}

	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

func (enum EquipmentType) String() string {
	var result string
	switch enum {
	case EquipmentHead:
		result = "頭部"
	case EquipmentTorso:
		result = "胴体"
	case EquipmentLegs:
		result = "脚部"
	case EquipmentJewelry:
		result = "装飾"
	default:
		panic("invalid equiment slot type")
	}
	return result
}

// ================

// ElementType は攻撃属性を表す
type ElementType string

const (
	// ElementTypeNone は属性なし
	ElementTypeNone ElementType = "NONE"
	// ElementTypeFire は火属性
	ElementTypeFire ElementType = "FIRE"
	// ElementTypeThunder は雷属性
	ElementTypeThunder ElementType = "THUNDER"
	// ElementTypeChill は氷属性
	ElementTypeChill ElementType = "CHILL"
	// ElementTypePhoton は光属性
	ElementTypePhoton ElementType = "PHOTON"
)

// Valid はElementTypeの値が有効かを検証する
func (enum ElementType) Valid() error {
	switch enum {
	case ElementTypeNone, ElementTypeFire, ElementTypeThunder, ElementTypeChill, ElementTypePhoton:
		return nil
	}
	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

func (enum ElementType) String() string {
	var result string
	switch enum {
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
		panic("invalid element type")
	}
	return result
}
