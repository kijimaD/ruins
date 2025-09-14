package turns

import (
	"github.com/kijimaD/ruins/lib/actions"
	"github.com/kijimaD/ruins/lib/logger"
)

// TurnPhase はターンの段階を表す
type TurnPhase int

const (
	// PlayerTurn プレイヤーのターン（入力待ち・アクション実行）
	PlayerTurn TurnPhase = iota
	// AITurn AI・NPCのターン（一括処理）
	AITurn
	// TurnEnd ターン終了処理
	TurnEnd
)

func (tp TurnPhase) String() string {
	switch tp {
	case PlayerTurn:
		return "PlayerTurn"
	case AITurn:
		return "AITurn"
	case TurnEnd:
		return "TurnEnd"
	default:
		return "Unknown"
	}
}

// TurnManager はターンの管理を行う
type TurnManager struct {
	PlayerMoves int       // プレイヤーの残り移動ポイント
	TurnPhase   TurnPhase // 現在のターンフェーズ
	TurnNumber  int       // ターン番号
	logger      *logger.Logger
}

// NewTurnManager は新しいTurnManagerを作成する
func NewTurnManager() *TurnManager {
	return &TurnManager{
		PlayerMoves: GetInitialPlayerMoves(),
		TurnPhase:   PlayerTurn,
		TurnNumber:  1,
		logger:      logger.New(logger.CategoryTurn),
	}
}

// GetInitialPlayerMoves は初期プレイヤー移動ポイントを返す
func GetInitialPlayerMoves() int {
	return 100
}

// CanPlayerAct はプレイヤーがアクションを実行可能かチェックする
func (tm *TurnManager) CanPlayerAct() bool {
	return tm.TurnPhase == PlayerTurn && tm.PlayerMoves > 0
}

// ConsumePlayerMoves はプレイヤーの移動ポイントを消費する
func (tm *TurnManager) ConsumePlayerMoves(actionID actions.ActionID) {
	cost := actions.GetActionInfo(actionID).MoveCost
	tm.PlayerMoves -= cost

	tm.logger.Debug("プレイヤー移動ポイント消費",
		"action", actionID.String(),
		"cost", cost,
		"remaining", tm.PlayerMoves)

	// 移動ポイントが尽きた場合はAIターンに移行
	if tm.PlayerMoves <= 0 {
		tm.TurnPhase = AITurn
		tm.logger.Debug("AIターンに移行", "turn", tm.TurnNumber)
	}
}

// AdvanceToAITurn はAIターンに強制移行する（待機など）
func (tm *TurnManager) AdvanceToAITurn() {
	tm.TurnPhase = AITurn
	tm.logger.Debug("AIターンに強制移行", "turn", tm.TurnNumber)
}

// AdvanceToTurnEnd はターン終了に移行する
func (tm *TurnManager) AdvanceToTurnEnd() {
	tm.TurnPhase = TurnEnd
	tm.logger.Debug("ターン終了処理", "turn", tm.TurnNumber)
}

// StartNewTurn は新しいターンを開始する
func (tm *TurnManager) StartNewTurn() {
	tm.TurnNumber++
	tm.PlayerMoves = GetInitialPlayerMoves()
	tm.TurnPhase = PlayerTurn

	tm.logger.Debug("新ターン開始",
		"turn", tm.TurnNumber,
		"player_moves", tm.PlayerMoves)
}

// GetTurnNumber は現在のターン番号を返す
func (tm *TurnManager) GetTurnNumber() int {
	return tm.TurnNumber
}

// IsPlayerTurn はプレイヤーのターンかチェックする
func (tm *TurnManager) IsPlayerTurn() bool {
	return tm.TurnPhase == PlayerTurn
}

// IsAITurn はAIのターンかチェックする
func (tm *TurnManager) IsAITurn() bool {
	return tm.TurnPhase == AITurn
}
