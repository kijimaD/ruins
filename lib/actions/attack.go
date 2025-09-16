package actions

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AttackActivity は攻撃アクティビティの実装
type AttackActivity struct{}

func init() {
	// 攻撃アクティビティをレジストリに登録
	RegisterActivityActor(ActivityAttack, &AttackActivity{})
}

// Validate は攻撃アクティビティの検証を行う
func (aa *AttackActivity) Validate(act *Activity, _ w.World) error {
	// 攻撃対象の確認
	if act.Target == nil {
		return fmt.Errorf("攻撃対象が設定されていません")
	}

	// 攻撃対象が有効なエンティティか
	if *act.Target == 0 {
		return fmt.Errorf("攻撃対象が無効です")
	}

	// TODO: より詳細な攻撃可能チェック
	// - ターゲットが存在するか
	// - 射程内にいるか
	// - 視界内にいるか
	// - 武器を装備しているか
	// - 攻撃可能な状態か（スタンなどでないか）

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
	act.Logger.Debug("攻撃実行",
		"actor", act.Actor,
		"target", *act.Target)

	// TODO: 実際の攻撃ロジック実装
	// - ダメージ計算
	// - 命中判定
	// - 武器による攻撃力修正
	// - ターゲットのHP減少
	// - 攻撃エフェクトの表示

	// プレイヤーの場合のみ攻撃メッセージを表示
	if isPlayerActivity(act, world) {
		gamelog.New(gamelog.FieldLog).
			Append("攻撃した").
			Log()
	}

	return nil
}

// canAttack は攻撃可能かをチェックする
func (aa *AttackActivity) canAttack(act *Activity, _ w.World) bool {
	// 攻撃対象の確認
	if act.Target == nil {
		return false
	}

	// ターゲットの存在チェック
	if *act.Target == ecs.Entity(0) {
		return false
	}

	// TODO: より詳細な攻撃可能チェック
	// - ターゲットが存在するか
	// - 射程内にいるか
	// - 視界内にいるか
	// - 武器を装備しているか

	return true
}
