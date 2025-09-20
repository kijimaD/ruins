package actions

import (
	"fmt"

	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
)

// WaitActivity は待機アクティビティの実装
type WaitActivity struct{}

func init() {
	// 待機アクティビティをレジストリに登録
	RegisterActivityActor(ActivityWait, &WaitActivity{})
}

// Validate は待機アクティビティの検証を行う
func (wa *WaitActivity) Validate(act *Activity, _ w.World) error {
	// 待機は基本的に常に実行可能
	// ただし、最低限のチェックは行う

	// アクターが存在するかチェック
	if act.Actor == 0 {
		return fmt.Errorf("待機するエンティティが指定されていません")
	}

	// 待機時間が妥当かチェック
	if act.TurnsTotal <= 0 {
		return fmt.Errorf("待機時間が無効です")
	}

	return nil
}

// Start は待機開始時の処理を実行する
func (wa *WaitActivity) Start(act *Activity, _ w.World) error {
	reason := "時間を過ごすため"
	act.Logger.Debug("待機開始", "actor", act.Actor, "reason", reason, "duration", act.TurnsLeft)
	return nil
}

// DoTurn は待機アクティビティの1ターン分の処理を実行する
func (wa *WaitActivity) DoTurn(act *Activity, world w.World) error {
	// 環境を観察
	wa.observeEnvironment(act, world)

	// 基本のターン処理
	if act.TurnsLeft <= 0 {
		act.Complete()
		return nil
	}

	// 1ターン進行
	act.TurnsLeft--
	act.Logger.Debug("待機進行",
		"turns_left", act.TurnsLeft,
		"progress", act.GetProgressPercent())

	// 完了チェック
	if act.TurnsLeft <= 0 {
		act.Complete()
		return nil
	}

	// メッセージ更新
	wa.updateMessage(act)
	return nil
}

// Finish は待機完了時の処理を実行する
func (wa *WaitActivity) Finish(act *Activity, world w.World) error {
	act.Logger.Debug("待機完了", "actor", act.Actor)

	// TODO: 1ターン待機の場合も出るのは微妙な感じがする
	// プレイヤーの場合のみ待機完了メッセージを表示
	if isPlayerActivity(act, world) {
		gamelog.New(gamelog.FieldLog).
			Append("待機を終了した").
			Log()
	}

	return nil
}

// Canceled は待機キャンセル時の処理を実行する
func (wa *WaitActivity) Canceled(act *Activity, _ w.World) error {
	act.Logger.Debug("待機キャンセル", "actor", act.Actor, "reason", act.CancelReason)
	return nil
}

// updateMessage は進行状況メッセージを更新する
func (wa *WaitActivity) updateMessage(act *Activity) {
	progress := act.GetProgressPercent()
	remainingTurns := act.TurnsLeft

	if progress < 25.0 {
		act.Message = "待機している..."
	} else if progress < 50.0 {
		act.Message = "時間を過ごしている..."
	} else if progress < 75.0 {
		act.Message = "のんびりと過ごしている..."
	} else if remainingTurns <= 1 {
		act.Message = "もうすぐ待機が終わりそうだ..."
	} else {
		act.Message = "引き続き待機している..."
	}
}

// observeEnvironment は環境観察処理を実行する
func (wa *WaitActivity) observeEnvironment(act *Activity, _ w.World) {
	// 待機中の環境観察（5ターン毎）
	if (act.TurnsTotal-act.TurnsLeft)%5 == 0 {
		// TODO: 環境観察の実装
		// - 周囲の敵の発見
		// - アイテムの発見
		// - 天候の変化など
		act.Logger.Debug("環境観察", "actor", act.Actor)
	}
}
