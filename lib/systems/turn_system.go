package systems

import (
	"github.com/kijimaD/ruins/lib/aiinput"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
)

// TurnSystem はターン管理を行うシステム
func TurnSystem(world w.World) error {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	switch turnManager.TurnPhase {
	case turns.PlayerTurn:
		// プレイヤー入力処理はDungeonStateで実行される
	case turns.AITurn:
		// AIターン: 全AI・NPCを一括処理
		if err := processAITurn(world); err != nil {
			return err
		}
		turnManager.AdvanceToTurnEnd()
	case turns.TurnEnd:
		// ターン終了処理
		if err := processTurnEnd(world); err != nil {
			return err
		}
		turnManager.StartNewTurn()
	}
	return nil
}

// processAITurn はAIターンの処理を行う
func processAITurn(world w.World) error {
	logger := logger.New(logger.CategoryTurn)
	logger.Debug("AIターン処理開始")

	// AI・NPCエンティティを処理
	processor := aiinput.NewProcessor()
	if err := processor.ProcessAllEntities(world); err != nil {
		return err
	}

	logger.Debug("AIターン処理完了")
	return nil
}

// processTurnEnd はターン終了処理を行う
func processTurnEnd(world w.World) error {
	logger := logger.New(logger.CategoryTurn)
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	logger.Debug("ターン終了処理", "turn", turnManager.TurnNumber)

	// 全エンティティのアクションポイントを回復
	if err := turnManager.RestoreAllActionPoints(world); err != nil {
		return err
	}

	// Deadエンティティの削除処理
	if err := DeadCleanupSystem(world); err != nil {
		return err
	}

	// 自動相互作用の実行
	if err := AutoInteractionSystem(world); err != nil {
		return err
	}

	// TODO: ターン終了時の共通処理をここに追加
	// - エフェクトの持続時間減少
	// - 状態異常の更新
	// - 環境変化の処理
	// など
	return nil
}
