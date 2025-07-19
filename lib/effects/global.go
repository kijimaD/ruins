package effects

import (
	"log"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// グローバルプロセッサー（既存コードとの互換性のため）
var globalProcessor = NewProcessor()

// 既存のエフェクトタイプ（下位互換性のため）
type Damage struct {
	Amount int
}

type Healing struct {
	Amount gc.Amounter
}

type ConsumptionStamina struct {
	Amount gc.Amounter
}

type RecoveryStamina struct {
	Amount gc.Amounter
}

type ItemUse struct {
	Item ecs.Entity
}

type WarpNext struct{}

type WarpEscape struct{}

// 既存のターゲットタイプ（下位互換性のため）
type Targets interface {
	isTarget()
}

type Party struct{}

func (Party) isTarget() {}

type Single struct {
	Target ecs.Entity
}

func (Single) isTarget() {}

type None struct{}

func (None) isTarget() {}

// 既存の関数（下位互換性のため）
func AddEffect(creator *ecs.Entity, effectType interface{}, targets Targets) {
	// エフェクトタイプを新しいシステムに変換
	var effect Effect

	switch e := effectType.(type) {
	case Damage:
		effect = CombatDamage{
			Amount: e.Amount,
			Source: DamageSourceWeapon,
		}
	case Healing:
		effect = CombatHealing{
			Amount: e.Amount,
		}
	case ConsumptionStamina:
		effect = ConsumeStamina{
			Amount: e.Amount,
		}
	case RecoveryStamina:
		effect = RestoreStamina{
			Amount: e.Amount,
		}
	case ItemUse:
		effect = UseItem{
			Item: e.Item,
		}
	case WarpNext:
		effect = MovementWarpNext{}
	case WarpEscape:
		effect = MovementWarpEscape{}
	default:
		log.Printf("未対応のエフェクトタイプ: %T", effectType)
		return
	}

	// ターゲットを変換してエフェクトを追加
	switch t := targets.(type) {
	case Single:
		if err := globalProcessor.AddEffect(effect, creator, t.Target); err != nil {
			log.Printf("エフェクト追加エラー: %v", err)
		}
	case Party:
		// パーティターゲットは後でRunEffectQueueで処理
		log.Printf("パーティターゲットは現在未対応")
	case None:
		if err := globalProcessor.AddEffect(effect, creator); err != nil {
			log.Printf("エフェクト追加エラー: %v", err)
		}
	}
}

func RunEffectQueue(world w.World) {
	if err := globalProcessor.Execute(world); err != nil {
		log.Printf("エフェクト実行エラー: %v", err)
	}
}

// ItemTrigger は既存のアイテムトリガー機能
func ItemTrigger(creator *ecs.Entity, item ecs.Entity, targets Targets, world w.World) {
	useItemEffect := UseItem{Item: item}

	switch t := targets.(type) {
	case Single:
		if err := globalProcessor.AddEffect(useItemEffect, creator, t.Target); err != nil {
			log.Printf("アイテムエフェクト追加エラー: %v", err)
		}
	case Party:
		selector := TargetParty{}
		if err := globalProcessor.AddTargetedEffect(useItemEffect, creator, selector, world); err != nil {
			log.Printf("アイテムエフェクト追加エラー: %v", err)
		}
	case None:
		if err := globalProcessor.AddEffect(useItemEffect, creator); err != nil {
			log.Printf("アイテムエフェクト追加エラー: %v", err)
		}
	}
}
