package systems

import (
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
)

// TurnSystem はターン管理を行うシステム
// CDDAスタイルのプレイヤー優先ターンシステムを実装
func TurnSystem(world w.World) {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	switch turnManager.TurnPhase {
	case turns.PlayerTurn:
		// プレイヤーターンでは入力システムが処理
		// 移動ポイントが尽きるまでプレイヤーが連続行動
		// TileInputSystemがアクション実行時にConsumePlayerMovesを呼ぶ

	case turns.AITurn:
		// AIターン: 全AI・NPCを一括処理
		processAITurn(world)
		turnManager.AdvanceToTurnEnd()

	case turns.TurnEnd:
		// ターン終了処理
		processTurnEnd(world)
		turnManager.StartNewTurn()
	}
}

// processAITurn はAIターンの処理を行う
func processAITurn(world w.World) {
	logger := logger.New(logger.CategoryTurn)
	logger.Debug("AIターン処理開始")

	// TODO: AI・NPCエンティティを処理
	// 現時点ではAIが存在しないため、処理をスキップ
	// 将来的にここでAISystemを呼び出す:
	// AISystem(world)

	logger.Debug("AIターン処理完了")
}

// processTurnEnd はターン終了処理を行う
func processTurnEnd(world w.World) {
	logger := logger.New(logger.CategoryTurn)
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	logger.Debug("ターン終了処理", "turn", turnManager.GetTurnNumber())

	// ターン終了時の共通処理をここに追加
	// - エフェクトの持続時間減少
	// - 状態異常の更新
	// - 環境変化の処理
	// など
}

// GetTurnPhase は現在のターンフェーズを取得する（デバッグ用）
func GetTurnPhase(world w.World) turns.TurnPhase {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	return turnManager.TurnPhase
}

// GetPlayerMoves は現在のプレイヤー移動ポイントを取得する（デバッグ用）
func GetPlayerMoves(world w.World) int {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	return turnManager.PlayerMoves
}
