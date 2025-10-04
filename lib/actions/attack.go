package actions

import (
	"fmt"
	"math"
	"math/rand/v2"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 攻撃システムの定数
const (
	// 射程・距離関連
	MeleeAttackRange = 1.5 // 近接攻撃の最大射程（斜めも考慮）

	// 命中率関連
	BaseHitRate          = 80 // 基本命中率（%）
	HitRatePerStatPoint  = 2  // 器用度と敏捷度の差1点あたりの命中率変化（%）
	MaxHitRate           = 95 // 最大命中率（%）
	MinHitRate           = 5  // 最小命中率（%）
	CriticalHitThreshold = 5  // クリティカルヒット判定しきい値（%以下）

	// ダメージ関連
	DamageRandomRange        = 6 // ダメージのランダム要素（1-6）
	CriticalDamageMultiplier = 3 // クリティカルダメージ倍率の分子
	CriticalDamageBase       = 2 // クリティカルダメージ倍率の分母（3/2 = 1.5倍）
	MinDamage                = 1 // 最低保証ダメージ

	// 確率計算関連
	DiceMax = 100 // ダイス最大値（1-100）
)

// AttackActivity はActivityInterfaceの実装
type AttackActivity struct{}

// Info はActivityInterfaceの実装
func (aa *AttackActivity) Info() ActivityInfo {
	return ActivityInfo{
		Name:            "攻撃",
		Description:     "敵を攻撃する",
		Interruptible:   false,
		Resumable:       false,
		ActionPointCost: 100,
		TotalRequiredAP: 100,
	}
}

// String はActivityInterfaceの実装
func (aa *AttackActivity) String() string {
	return "Attack"
}

// Validate はActivityInterfaceの実装
func (aa *AttackActivity) Validate(act *Activity, world w.World) error {
	if act.Target == nil {
		return ErrAttackTargetNotSet
	}

	if *act.Target == 0 {
		return ErrAttackTargetInvalid
	}

	if world.Components.Dead.Get(act.Actor) != nil {
		return ErrAttackerDead
	}

	if world.Components.GridElement.Get(*act.Target) == nil {
		return ErrAttackTargetNotExists
	}

	if world.Components.Dead.Get(*act.Target) != nil {
		return ErrAttackTargetDead
	}

	if !aa.isInRange(act.Actor, *act.Target, world) {
		return ErrAttackOutOfRange
	}

	if !aa.canPerformAttack(act.Actor, world) {
		return ErrAttackNoWeapon
	}

	return nil
}

// Start はActivityInterfaceの実装
func (aa *AttackActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃開始", "actor", act.Actor, "target", *act.Target)
	return nil
}

// DoTurn はActivityInterfaceの実装
func (aa *AttackActivity) DoTurn(act *Activity, world w.World) error {
	if act.Target == nil {
		act.Cancel("攻撃対象が設定されていません")
		return ErrAttackTargetNotSet
	}

	if !aa.canAttack(act, world) {
		act.Cancel("攻撃できません")
		return ErrAttackTargetInvalid
	}

	if err := aa.performAttack(act, world); err != nil {
		act.Cancel(fmt.Sprintf("攻撃エラー: %s", err.Error()))
		return err
	}

	act.Complete()
	return nil
}

// Finish はActivityInterfaceの実装
func (aa *AttackActivity) Finish(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃アクティビティ完了",
		"actor", act.Actor,
		"target", *act.Target)

	return nil
}

// Canceled はActivityInterfaceの実装
func (aa *AttackActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

func (aa *AttackActivity) performAttack(act *Activity, world w.World) error {
	attacker := act.Actor
	target := *act.Target

	act.Logger.Debug("攻撃実行", "attacker", attacker, "target", target)

	hit, criticalHit := aa.rollHitCheck(attacker, target, world)
	if !hit {
		aa.logAttackResult(attacker, target, world, false, false, 0)
		return nil
	}

	damage := aa.calculateDamage(attacker, target, world, criticalHit)
	if damage < 0 {
		damage = 0
	}

	// 1. 攻撃試行ログ
	aa.logAttackResult(attacker, target, world, true, criticalHit, damage)

	// 2. ダメージを直接適用（ダメージログと死亡ログはapplyDamage内で出力）
	aa.applyDamage(world, target, damage, attacker)

	return nil
}

func (aa *AttackActivity) canAttack(act *Activity, world w.World) bool {
	if act.Target == nil {
		return false
	}

	if *act.Target == ecs.Entity(0) {
		return false
	}

	if err := aa.Validate(act, world); err != nil {
		return false
	}

	return true
}

func (aa *AttackActivity) isInRange(attacker, target ecs.Entity, world w.World) bool {
	attackerGrid := world.Components.GridElement.Get(attacker)
	if attackerGrid == nil {
		return false
	}

	targetGrid := world.Components.GridElement.Get(target)
	if targetGrid == nil {
		return false
	}

	attackerPos := attackerGrid.(*gc.GridElement)
	targetPos := targetGrid.(*gc.GridElement)

	dx := float64(attackerPos.X - targetPos.X)
	dy := float64(attackerPos.Y - targetPos.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	// TODO: 遠距離武器の場合は射程を武器から取得
	return distance <= MeleeAttackRange
}

func (aa *AttackActivity) canPerformAttack(attacker ecs.Entity, world w.World) bool {
	// TODO: 装備武器のチェック
	attrs := world.Components.Attributes.Get(attacker)
	return attrs != nil
}

func (aa *AttackActivity) rollHitCheck(attacker, target ecs.Entity, world w.World) (hit bool, critical bool) {
	attackerAttrs := world.Components.Attributes.Get(attacker).(*gc.Attributes)
	attackerDexterity := attackerAttrs.Dexterity.Total

	targetAttrs := world.Components.Attributes.Get(target).(*gc.Attributes)
	targetAgility := targetAttrs.Agility.Total

	baseHitRate := BaseHitRate + (attackerDexterity-targetAgility)*HitRatePerStatPoint

	if baseHitRate > MaxHitRate {
		baseHitRate = MaxHitRate
	}
	if baseHitRate < MinHitRate {
		baseHitRate = MinHitRate
	}

	roll := rand.IntN(DiceMax) + 1
	hit = roll <= baseHitRate
	critical = roll <= CriticalHitThreshold

	return hit, critical
}

func (aa *AttackActivity) calculateDamage(attacker, target ecs.Entity, world w.World, critical bool) int {
	attackerAttrs := world.Components.Attributes.Get(attacker).(*gc.Attributes)
	attackerStrength := attackerAttrs.Strength.Total

	targetAttrs := world.Components.Attributes.Get(target).(*gc.Attributes)
	targetDefense := targetAttrs.Defense.Total

	baseDamage := attackerStrength + rand.IntN(DamageRandomRange) + 1

	// TODO: 武器攻撃力の追加

	if critical {
		baseDamage = baseDamage * CriticalDamageMultiplier / CriticalDamageBase
	}

	finalDamage := baseDamage - targetDefense
	if finalDamage < MinDamage {
		finalDamage = MinDamage
	}

	return finalDamage
}

// applyDamage はダメージを適用する
func (aa *AttackActivity) applyDamage(world w.World, target ecs.Entity, damage int, source ecs.Entity) {
	pools := world.Components.Pools.Get(target).(*gc.Pools)

	beforeHP := pools.HP.Current
	pools.HP.Current -= damage
	if pools.HP.Current < 0 {
		pools.HP.Current = 0
	}

	// ダメージログ出力（プレイヤー関連の場合のみ）
	if isPlayerActivity(&Activity{Actor: source}, world) || isPlayerActivity(&Activity{Actor: target}, world) {
		sourceName := aa.getEntityName(source, world)
		targetName := aa.getEntityName(target, world)

		logger := gamelog.New(gamelog.FieldLog)
		logger.Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, source, sourceName, world)
		}).Append(" は ").Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, target, targetName, world)
		}).Append(fmt.Sprintf(" に %d のダメージを与えた。", damage)).Log()
	}

	// 死亡チェック
	if pools.HP.Current <= 0 && beforeHP > 0 {
		target.AddComponent(world.Components.Dead, &gc.Dead{})
		aa.logDeath(world, target, source)
	}
}

// logDeath は死亡ログを出力する
func (aa *AttackActivity) logDeath(world w.World, target ecs.Entity, source ecs.Entity) {
	// プレイヤー関連の場合のみログ出力
	if !isPlayerActivity(&Activity{Actor: source}, world) && !isPlayerActivity(&Activity{Actor: target}, world) {
		return
	}

	targetName := aa.getEntityName(target, world)

	gamelog.New(gamelog.FieldLog).
		Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, target, targetName, world)
		}).
		Append(" は倒れた。").
		Log()
}

// logAttackResult は攻撃結果をログに出力する
// ダメージと死亡は別途ログ出力されるため、ここでは攻撃の成否のみを出力
func (aa *AttackActivity) logAttackResult(attacker, target ecs.Entity, world w.World, hit bool, critical bool, _ int) {
	// プレイヤーが関わる攻撃のみログ出力
	if !isPlayerActivity(&Activity{Actor: attacker}, world) && !isPlayerActivity(&Activity{Actor: target}, world) {
		return
	}

	// 攻撃者名とターゲット名を取得
	attackerName := aa.getEntityName(attacker, world)
	targetName := aa.getEntityName(target, world)

	gamelog.New(gamelog.FieldLog).
		Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, attacker, attackerName, world)
		}).
		Append(" は ").
		Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, target, targetName, world)
		}).
		Build(func(l *gamelog.Logger) {
			if !hit {
				l.Append(" を攻撃したが外れた。")
			} else if critical {
				l.Append(" にクリティカルヒット！")
			} else {
				l.Append(" を攻撃した。")
			}
		}).
		Log()
}

// getEntityName はエンティティの名前を取得する
func (aa *AttackActivity) getEntityName(entity ecs.Entity, world w.World) string {
	// Nameコンポーネントから名前を取得
	name := world.Components.Name.Get(entity)
	if name != nil {
		return name.(*gc.Name).Name
	}

	return "Unknown"
}

// appendNameWithColor はエンティティの種類に応じて色付きで名前を追加する
func (aa *AttackActivity) appendNameWithColor(logger *gamelog.Logger, entity ecs.Entity, name string, world w.World) {
	if entity.HasComponent(world.Components.Player) {
		logger.PlayerName(name)
	} else if entity.HasComponent(world.Components.AIMoveFSM) {
		logger.NPCName(name)
	} else {
		logger.Append(name)
	}
}
