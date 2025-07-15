package components

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// ErrInvalidEnumType はenumに無効な値が指定された場合のエラー
var ErrInvalidEnumType = errors.New("enumに無効な値が指定された")

// ================

// ワープモード
type warpMode string

const (
	// WarpModeNext は次の階層へワープする
	WarpModeNext = warpMode("NEXT")
	// WarpModeEscape は脱出ワープする
	WarpModeEscape = warpMode("ESCAPE")
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
	// TargetGroupCard はカードグループ
	TargetGroupCard = TargetGroupType("CARD") // カード
	// TargetGroupNone はグループなし
	TargetGroupNone = TargetGroupType("NONE") // なし
)

// Valid はTargetGroupTypeの値が有効かを検証する
func (enum TargetGroupType) Valid() error {
	switch enum {
	case TargetGroupAlly, TargetGroupEnemy, TargetGroupCard, TargetGroupNone:
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

// AttackType は武器種別を表す。種別によって適用する計算式が異なる
type AttackType string

const (
	// AttackSword は刀剣
	AttackSword = AttackType("SWORD") // 刀剣
	// AttackSpear は長物
	AttackSpear = AttackType("SPEAR") // 長物
	// AttackHandgun は拳銃
	AttackHandgun = AttackType("HANDGUN") // 拳銃
	// AttackRifle は小銃
	AttackRifle = AttackType("RIFLE") // 小銃
	// AttackFist は格闘
	AttackFist = AttackType("FIST") // 格闘
	// AttackCanon は大砲
	AttackCanon = AttackType("CANON") // 大砲
)

// Valid はAttackTypeの値が有効かを検証する
func (enum AttackType) Valid() error {
	switch enum {
	case AttackSword, AttackSpear, AttackHandgun, AttackRifle, AttackFist, AttackCanon:
		return nil
	}

	return fmt.Errorf("get %s: %w", enum, ErrInvalidEnumType)
}

func (enum AttackType) String() string {
	var result string
	switch enum {
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
	case AttackCanon:
		result = "大砲"
	default:
		log.Fatal("invalid attack type")
	}

	return result
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
		log.Fatal("invalid equiment slot type")
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
		log.Fatal("invalid element type")
	}
	return result
}
