package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/logger"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Executor はアクションの実行を管理する
type Executor struct {
	processor *effects.Processor
	logger    *logger.Logger
}

// NewExecutor は新しいExecutorを作成する
func NewExecutor() *Executor {
	return &Executor{
		processor: effects.NewProcessor(),
		logger:    logger.New(logger.CategoryAction),
	}
}

// Execute は指定されたアクションを実行する
func (e *Executor) Execute(actionID ActionID, ctx Context) (*Result, error) {
	e.logger.Debug("アクション実行開始", "action", actionID.String(), "actor", ctx.Actor)

	// 検証フェーズ
	if err := e.validateAction(actionID, ctx); err != nil {
		e.logger.Warn("アクション検証失敗", "action", actionID.String(), "error", err.Error())
		return &Result{
			Success:  false,
			ActionID: actionID,
			Message:  err.Error(),
		}, err
	}

	// 実行フェーズ
	var result *Result
	var err error

	switch actionID {
	case ActionMove:
		result, err = e.executeMove(ctx)
	case ActionWait:
		result, err = e.executeWait(ctx)
	case ActionAttack:
		result, err = e.executeAttack(ctx)
	default:
		err = fmt.Errorf("未実装のアクション: %v", actionID)
		result = &Result{
			Success:  false,
			ActionID: actionID,
			Message:  err.Error(),
		}
	}

	if err != nil {
		e.logger.Error("アクション実行エラー", "action", actionID.String(), "error", err.Error())
	} else {
		e.logger.Debug("アクション実行完了", "action", actionID.String(), "success", result.Success)
	}

	return result, err
}

// validateAction はアクション実行前の検証を行う
func (e *Executor) validateAction(actionID ActionID, ctx Context) error {
	// 基本検証は削除（Entity(0)も有効な場合がある）
	// 必要に応じて後でより適切な検証を実装

	// アクション固有の検証
	switch actionID {
	case ActionMove:
		return e.validateMove(ctx)
	case ActionAttack:
		return e.validateAttack(ctx)
	case ActionWait:
		return nil // 待機は常に有効
	default:
		return fmt.Errorf("不明なアクション: %v", actionID)
	}
}

// executeMove は移動アクションを実行する
func (e *Executor) executeMove(ctx Context) (*Result, error) {
	if ctx.Position == nil {
		return &Result{Success: false, ActionID: ActionMove, Message: "移動先が指定されていません"},
			fmt.Errorf("移動先が指定されていません")
	}

	// GridElementを直接更新
	gridElement := ctx.World.Components.GridElement.Get(ctx.Actor).(*gc.GridElement)
	oldX, oldY := int(gridElement.X), int(gridElement.Y)

	gridElement.X = gc.Tile(int(ctx.Position.X))
	gridElement.Y = gc.Tile(int(ctx.Position.Y))

	message := fmt.Sprintf("(%d,%d)から(%d,%d)に移動した", oldX, oldY, int(ctx.Position.X), int(ctx.Position.Y))
	e.logger.Info("移動完了", "from", fmt.Sprintf("(%d,%d)", oldX, oldY),
		"to", fmt.Sprintf("(%d,%d)", int(ctx.Position.X), int(ctx.Position.Y)))

	return &Result{
		Success:  true,
		ActionID: ActionMove,
		Message:  message,
	}, nil
}

// executeWait は待機アクションを実行する
func (e *Executor) executeWait(ctx Context) (*Result, error) {
	e.logger.Info("待機実行", "actor", ctx.Actor)

	return &Result{
		Success:  true,
		ActionID: ActionWait,
		Message:  "待機した",
	}, nil
}

// executeAttack は攻撃アクションを実行する
func (e *Executor) executeAttack(ctx Context) (*Result, error) {
	if ctx.Target == nil {
		return &Result{Success: false, ActionID: ActionAttack, Message: "攻撃対象が指定されていません"},
			fmt.Errorf("攻撃対象が指定されていません")
	}

	// 基本的な攻撃処理（簡易実装）
	e.logger.Info("攻撃実行", "attacker", ctx.Actor, "target", *ctx.Target)

	return &Result{
		Success:  true,
		ActionID: ActionAttack,
		Message:  "攻撃した",
	}, nil
}

// validateMove は移動アクションの検証を行う
func (e *Executor) validateMove(ctx Context) error {
	if ctx.Position == nil {
		return fmt.Errorf("移動先が指定されていません")
	}

	// GridElementコンポーネントが存在するかチェック
	if !ctx.Actor.HasComponent(ctx.World.Components.GridElement) {
		return fmt.Errorf("移動可能なエンティティではありません")
	}

	// 移動可能性をチェック（簡易実装）
	if !e.canMoveTo(ctx.World, int(ctx.Position.X), int(ctx.Position.Y), ctx.Actor) {
		return fmt.Errorf("そこには移動できません")
	}

	return nil
}

// validateAttack は攻撃アクションの検証を行う
func (e *Executor) validateAttack(ctx Context) error {
	if ctx.Target == nil {
		return fmt.Errorf("攻撃対象が指定されていません")
	}

	// ターゲットが存在するかチェック
	// TODO: 実際の存在チェックロジックを実装
	if *ctx.Target == ecs.Entity(0) {
		return fmt.Errorf("攻撃対象が無効です")
	}

	return nil
}

// canMoveTo は移動可能かどうかを判定する
func (e *Executor) canMoveTo(_ w.World, x, y int, _ ecs.Entity) bool {
	// 基本的な境界チェック
	if x < 0 || y < 0 || x >= 100 || y >= 100 { // 適当な境界値
		return false
	}

	// 他のエンティティとの重複チェックは後で実装
	return true
}
