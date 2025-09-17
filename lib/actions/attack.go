package actions

import (
	"fmt"
	"math"
	"math/rand"

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

// AttackActivity は攻撃アクティビティの実装
type AttackActivity struct{}

func init() {
	// 攻撃アクティビティをレジストリに登録
	RegisterActivityActor(ActivityAttack, &AttackActivity{})
}

// Validate は攻撃アクティビティの検証を行う
func (aa *AttackActivity) Validate(act *Activity, world w.World) error {
	// 攻撃対象の確認
	if act.Target == nil {
		return fmt.Errorf("攻撃対象が設定されていません")
	}

	// 攻撃対象が有効なエンティティか
	if *act.Target == 0 {
		return fmt.Errorf("攻撃対象が無効です")
	}

	// 攻撃者の生存確認
	if world.Components.Dead.Get(act.Actor) != nil {
		return fmt.Errorf("攻撃者が死亡しています")
	}

	// ターゲットの存在確認（GridElementの存在で判定）
	if world.Components.GridElement.Get(*act.Target) == nil {
		return fmt.Errorf("攻撃対象が存在しません")
	}

	// ターゲットの生存確認
	if world.Components.Dead.Get(*act.Target) != nil {
		return fmt.Errorf("攻撃対象が既に死亡しています")
	}

	// 射程チェック
	if !aa.isInRange(act.Actor, *act.Target, world) {
		return fmt.Errorf("攻撃対象が射程外です")
	}

	// 攻撃者の装備チェック（武器または素手攻撃可能か）
	if !aa.canPerformAttack(act.Actor, world) {
		return fmt.Errorf("攻撃手段がありません")
	}

	return nil
}

// Start は攻撃開始時の処理を実行する
func (aa *AttackActivity) Start(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃開始", "actor", act.Actor, "target", *act.Target)
	return nil
}

// DoTurn は攻撃アクティビティの1ターン分の処理を実行する
func (aa *AttackActivity) DoTurn(act *Activity, world w.World) error {
	// 攻撃対象チェック
	if act.Target == nil {
		act.Cancel("攻撃対象が設定されていません")
		return fmt.Errorf("攻撃対象が設定されていません")
	}

	// 攻撃可能性を再チェック
	if !aa.canAttack(act, world) {
		act.Cancel("攻撃できません")
		return fmt.Errorf("攻撃対象が無効です")
	}

	// 攻撃実行
	if err := aa.performAttack(act, world); err != nil {
		act.Cancel(fmt.Sprintf("攻撃エラー: %s", err.Error()))
		return err
	}

	// 攻撃処理完了
	act.Complete()
	return nil
}

// Finish は攻撃完了時の処理を実行する
func (aa *AttackActivity) Finish(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃アクティビティ完了",
		"actor", act.Actor,
		"target", *act.Target)

	return nil
}

// Canceled は攻撃キャンセル時の処理を実行する
func (aa *AttackActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("攻撃キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// performAttack は実際の攻撃処理を実行する
func (aa *AttackActivity) performAttack(act *Activity, world w.World) error {
	attacker := act.Actor
	target := *act.Target

	act.Logger.Debug("攻撃実行", "attacker", attacker, "target", target)

	// 命中判定
	hit, criticalHit := aa.rollHitCheck(attacker, target, world)
	if !hit {
		// 攻撃外れ
		aa.logAttackResult(attacker, target, world, false, false, 0)
		return nil
	}

	// ダメージ計算
	damage := aa.calculateDamage(attacker, target, world, criticalHit)
	if damage < 0 {
		damage = 0
	}

	// ダメージ適用
	if err := aa.applyDamage(target, damage, world); err != nil {
		return fmt.Errorf("ダメージ適用エラー: %w", err)
	}

	// 結果ログ
	aa.logAttackResult(attacker, target, world, true, criticalHit, damage)

	// 死亡判定
	aa.checkDeath(target, world)

	return nil
}

// canAttack は攻撃可能かをチェックする
func (aa *AttackActivity) canAttack(act *Activity, world w.World) bool {
	// 攻撃対象の確認
	if act.Target == nil {
		return false
	}

	// ターゲットの存在チェック
	if *act.Target == ecs.Entity(0) {
		return false
	}

	// より詳細なチェック（Validateと同様）
	if err := aa.Validate(act, world); err != nil {
		return false
	}

	return true
}

// isInRange は攻撃対象が射程内にいるかチェックする
func (aa *AttackActivity) isInRange(attacker, target ecs.Entity, world w.World) bool {
	// 攻撃者の位置を取得
	attackerGrid := world.Components.GridElement.Get(attacker)
	if attackerGrid == nil {
		return false
	}

	// ターゲットの位置を取得
	targetGrid := world.Components.GridElement.Get(target)
	if targetGrid == nil {
		return false
	}

	attackerPos := attackerGrid.(*gc.GridElement)
	targetPos := targetGrid.(*gc.GridElement)

	// タイル間の距離を計算
	dx := float64(attackerPos.X - targetPos.X)
	dy := float64(attackerPos.Y - targetPos.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	// 近接攻撃の場合は隣接チェック（斜めも考慮）
	// TODO: 遠距離武器の場合は射程を武器から取得
	return distance <= MeleeAttackRange
}

// canPerformAttack は攻撃手段があるかチェックする
func (aa *AttackActivity) canPerformAttack(attacker ecs.Entity, world w.World) bool {
	// TODO: 装備武器のチェック
	// 現在は素手攻撃を常に許可

	// 属性による攻撃可能チェック
	attrs := world.Components.Attributes.Get(attacker)
	if attrs == nil {
		return false // 属性がないエンティティは攻撃不可
	}

	return true
}

// rollHitCheck は命中判定を行う
func (aa *AttackActivity) rollHitCheck(attacker, target ecs.Entity, world w.World) (hit bool, critical bool) {
	// 攻撃者の器用度を取得
	attackerAttrs := world.Components.Attributes.Get(attacker).(*gc.Attributes)
	attackerDexterity := attackerAttrs.Dexterity.Total

	// ターゲットの敏捷度を取得
	targetAttrs := world.Components.Attributes.Get(target).(*gc.Attributes)
	targetAgility := targetAttrs.Agility.Total

	// 基本命中率計算
	baseHitRate := BaseHitRate + (attackerDexterity-targetAgility)*HitRatePerStatPoint

	// 命中率の範囲制限
	if baseHitRate > MaxHitRate {
		baseHitRate = MaxHitRate
	}
	if baseHitRate < MinHitRate {
		baseHitRate = MinHitRate
	}

	// ダイス振り
	roll := rand.IntN(DiceMax) + 1
	hit = roll <= baseHitRate

	// クリティカルヒット判定
	critical = roll <= CriticalHitThreshold

	return hit, critical
}

// calculateDamage はダメージを計算する
func (aa *AttackActivity) calculateDamage(attacker, target ecs.Entity, world w.World, critical bool) int {
	// 攻撃者の筋力を取得
	attackerAttrs := world.Components.Attributes.Get(attacker).(*gc.Attributes)
	attackerStrength := attackerAttrs.Strength.Total

	// ターゲットの防御力を取得
	targetAttrs := world.Components.Attributes.Get(target).(*gc.Attributes)
	targetDefense := targetAttrs.Defense.Total

	// 基本ダメージ = 筋力 + ランダム要素
	baseDamage := attackerStrength + rand.IntN(DamageRandomRange) + 1

	// TODO: 武器攻撃力の追加

	// クリティカルヒット時のダメージ倍率適用
	if critical {
		baseDamage = baseDamage * CriticalDamageMultiplier / CriticalDamageBase
	}

	// 防御力分を減算
	finalDamage := baseDamage - targetDefense
	if finalDamage < MinDamage {
		finalDamage = MinDamage
	}

	return finalDamage
}

// applyDamage はダメージをターゲットに適用する
func (aa *AttackActivity) applyDamage(target ecs.Entity, damage int, world w.World) error {
	// ターゲットのPoolsコンポーネントを取得
	pools := world.Components.Pools.Get(target)
	if pools == nil {
		return fmt.Errorf("ターゲットにPoolsコンポーネントがありません")
	}

	targetPools := pools.(*gc.Pools)

	// HPからダメージを減算
	targetPools.HP.Current -= damage
	if targetPools.HP.Current < 0 {
		targetPools.HP.Current = 0
	}

	return nil
}

// checkDeath は死亡判定を行う
func (aa *AttackActivity) checkDeath(target ecs.Entity, world w.World) {
	// ターゲットのHPをチェック
	pools := world.Components.Pools.Get(target)
	if pools == nil {
		return // Poolsがない場合は死亡判定しない
	}

	targetPools := pools.(*gc.Pools)
	if targetPools.HP.Current <= 0 {
		// 死亡メッセージをログ出力（プレイヤーまたは敵の場合）
		if target.HasComponent(world.Components.Player) || target.HasComponent(world.Components.AIMoveFSM) {
			targetName := aa.getEntityName(target, world)
			gamelog.New(gamelog.FieldLog).
				Build(func(l *gamelog.Logger) {
					aa.appendNameWithColor(l, target, targetName, world)
				}).
				Append(" は倒れた。").
				Log()
		}

		// エンティティを完全に削除
		world.Manager.DeleteEntity(target)
	}
}

// logAttackResult は攻撃結果をログに出力する
func (aa *AttackActivity) logAttackResult(attacker, target ecs.Entity, world w.World, hit bool, critical bool, damage int) {
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
		Append(" が ").
		Build(func(l *gamelog.Logger) {
			aa.appendNameWithColor(l, target, targetName, world)
		}).
		Build(func(l *gamelog.Logger) {
			if !hit {
				l.Append(" を攻撃したが外れた。")
			} else if critical {
				l.Append(" にクリティカルヒット。").Damage(damage).Append("ダメージ")
			} else {
				l.Append(" を攻撃した。").Damage(damage).Append("ダメージ")
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

	// Nameコンポーネントがない場合のフォールバック
	if entity.HasComponent(world.Components.Player) {
		return "プレイヤー"
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
