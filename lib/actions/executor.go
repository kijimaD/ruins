package actions

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/effects"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/movement"
	"github.com/kijimaD/ruins/lib/resources"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
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
func (e *Executor) Execute(actionID ActionID, ctx Context, world w.World) (*Result, error) {
	e.logger.Debug("アクション実行開始", "action", actionID.String(), "actor", ctx.Actor)
	// 検証フェーズ
	if err := e.validateAction(actionID, ctx, world); err != nil {
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
		result, err = e.executeMove(ctx, world)
	case ActionWait:
		result = e.executeWait(ctx, world)
	case ActionAttack:
		result, err = e.executeAttack(ctx, world)
	case ActionPickupItem:
		result = e.executePickupItem(ctx, world)
	case ActionWarp:
		result, err = e.executeWarp(ctx, world)
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
		// アクション成功時、統合AP管理システムでAP消費
		if result.Success && world.Resources.TurnManager != nil {
			if turnManager, ok := world.Resources.TurnManager.(*turns.TurnManager); ok {
				actionInfo := GetActionInfo(actionID)
				consumed := turnManager.ConsumeActionPoints(world, ctx.Actor, actionID.String(), actionInfo.MoveCost)
				if !consumed {
					e.logger.Debug("AP消費失敗（APコンポーネントなしまたはAP不足）", "entity", ctx.Actor, "action", actionID.String())
				}
			}
		}
	}
	return result, err
}

// validateAction はアクション実行前の検証を行う
func (e *Executor) validateAction(actionID ActionID, ctx Context, world w.World) error {
	// 必要に応じて後でより適切な検証を実装
	// アクション固有の検証
	switch actionID {
	case ActionMove:
		return e.validateMove(ctx, world)
	case ActionAttack:
		return e.validateAttack(ctx, world)
	case ActionPickupItem:
		return e.validatePickupItem(ctx, world)
	case ActionWarp:
		return e.validateWarp(ctx, world)
	case ActionWait:
		return nil // 待機は常に有効
	default:
		return fmt.Errorf("不明なアクション: %v", actionID)
	}
}

// executeMove は移動アクションを実行する
func (e *Executor) executeMove(ctx Context, world w.World) (*Result, error) {
	if ctx.Dest == nil {
		return &Result{Success: false, ActionID: ActionMove, Message: "移動先が指定されていません"},
			fmt.Errorf("移動先が指定されていません")
	}
	// GridElementを直接更新
	gridElement := world.Components.GridElement.Get(ctx.Actor).(*gc.GridElement)
	oldX, oldY := int(gridElement.X), int(gridElement.Y)
	gridElement.X = gc.Tile(int(ctx.Dest.X))
	gridElement.Y = gc.Tile(int(ctx.Dest.Y))
	message := fmt.Sprintf("(%d,%d)から(%d,%d)に移動した", oldX, oldY, int(ctx.Dest.X), int(ctx.Dest.Y))
	e.logger.Debug("移動完了", "from", fmt.Sprintf("(%d,%d)", oldX, oldY),
		"to", fmt.Sprintf("(%d,%d)", int(ctx.Dest.X), int(ctx.Dest.Y)))
	return &Result{
		Success:  true,
		ActionID: ActionMove,
		Message:  message,
	}, nil
}

// executeWait は待機アクションを実行する
func (e *Executor) executeWait(_ Context, _ w.World) *Result {
	return &Result{
		Success:  true,
		ActionID: ActionWait,
		Message:  "待機した",
	}
}

// executeAttack は攻撃アクションを実行する
func (e *Executor) executeAttack(ctx Context, _ w.World) (*Result, error) {
	if ctx.Target == nil {
		return &Result{Success: false, ActionID: ActionAttack, Message: "攻撃対象が指定されていません"},
			fmt.Errorf("攻撃対象が指定されていません")
	}
	// TODO: 実装する
	e.logger.Debug("攻撃実行", "attacker", ctx.Actor, "target", *ctx.Target)
	return &Result{
		Success:  true,
		ActionID: ActionAttack,
		Message:  "攻撃した",
	}, nil
}

// validateMove は移動アクションの検証を行う
func (e *Executor) validateMove(ctx Context, world w.World) error {
	if ctx.Dest == nil {
		return fmt.Errorf("移動先が指定されていません")
	}
	gridElement := world.Components.GridElement.Get(ctx.Actor)
	if gridElement == nil {
		return fmt.Errorf("移動可能なエンティティではありません")
	}
	if !movement.CanMoveTo(world, int(ctx.Dest.X), int(ctx.Dest.Y), ctx.Actor) {
		return fmt.Errorf("そこには移動できません")
	}
	return nil
}

// validateAttack は攻撃アクションの検証を行う
func (e *Executor) validateAttack(ctx Context, _ w.World) error {
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

// executePickupItem はアイテム拾得アクションを実行する
func (e *Executor) executePickupItem(ctx Context, world w.World) *Result {
	// プレイヤー位置を取得
	gridElement := world.Components.GridElement.Get(ctx.Actor).(*gc.GridElement)
	playerTileX := int(gridElement.X)
	playerTileY := int(gridElement.Y)
	// 収集されたアイテムを記録するリスト
	var itemsToCollect []ecs.Entity
	// 同じタイルのフィールドアイテムを検索
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationOnField,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(itemEntity ecs.Entity) {
		itemGrid := world.Components.GridElement.Get(itemEntity).(*gc.GridElement)
		// タイル単位の位置判定
		if int(itemGrid.X) == playerTileX && int(itemGrid.Y) == playerTileY {
			itemsToCollect = append(itemsToCollect, itemEntity)
		}
	}))
	if len(itemsToCollect) == 0 {
		return &Result{
			Success:  false,
			ActionID: ActionPickupItem,
			Message:  "拾えるアイテムがありません",
		}
	}
	// 収集されたアイテムを処理
	collectedCount := 0
	for _, itemEntity := range itemsToCollect {
		if err := e.collectFieldItem(world, itemEntity); err != nil {
			e.logger.Warn("アイテム拾得エラー", "item", itemEntity, "error", err.Error())
			continue
		}
		collectedCount++
	}
	message := fmt.Sprintf("%d個のアイテムを拾得した", collectedCount)
	e.logger.Debug("アイテム拾得完了", "count", collectedCount)
	return &Result{
		Success:  true,
		ActionID: ActionPickupItem,
		Message:  message,
	}
}

// executeWarp はワープアクションを実行する
func (e *Executor) executeWarp(_ Context, world w.World) (*Result, error) {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	if gameResources.PlayerTileState.CurrentWarp == nil {
		return &Result{
			Success:  false,
			ActionID: ActionWarp,
			Message:  "ワープホールが見つかりません",
		}, fmt.Errorf("ワープホールが見つかりません")
	}
	switch gameResources.PlayerTileState.CurrentWarp.Mode {
	case gc.WarpModeNext:
		gameResources.SetStateEvent(resources.StateEventWarpNext)
		e.logger.Debug("次の階へワープ")
		return &Result{
			Success:  true,
			ActionID: ActionWarp,
			Message:  "次の階へ移動した",
		}, nil
	case gc.WarpModeEscape:
		gameResources.SetStateEvent(resources.StateEventWarpEscape)
		e.logger.Debug("脱出ワープ")
		return &Result{
			Success:  true,
			ActionID: ActionWarp,
			Message:  "脱出した",
		}, nil
	default:
		return &Result{
			Success:  false,
			ActionID: ActionWarp,
			Message:  "不明なワープタイプです",
		}, fmt.Errorf("不明なワープタイプ: %v", gameResources.PlayerTileState.CurrentWarp.Mode)
	}
}

// collectFieldItem はフィールドアイテムを収集してバックパックに移動する
func (e *Executor) collectFieldItem(world w.World, itemEntity ecs.Entity) error {
	itemName := "Unknown Item"
	if nameComp := world.Components.Name.Get(itemEntity); nameComp != nil {
		name := nameComp.(*gc.Name)
		itemName = name.Name
	}
	// フィールドからバックパックに移動
	// ItemLocationOnFieldコンポーネントを削除
	itemEntity.RemoveComponent(world.Components.ItemLocationOnField)
	// ItemLocationInBackpackコンポーネントを追加
	itemEntity.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})
	// グリッド表示コンポーネントを削除（フィールドから消す）
	if itemEntity.HasComponent(world.Components.GridElement) {
		itemEntity.RemoveComponent(world.Components.GridElement)
	}
	// スプライト表示コンポーネントを削除（フィールドから消す）
	if itemEntity.HasComponent(world.Components.SpriteRender) {
		itemEntity.RemoveComponent(world.Components.SpriteRender)
	}
	// 既存のバックパック内の同じアイテムと統合する処理
	if err := worldhelper.MergeMaterialIntoInventory(world, itemEntity, itemName); err != nil {
		return fmt.Errorf("インベントリ統合エラー: %w", err)
	}
	// 色付きログ
	gamelog.New(gamelog.FieldLog).
		Append("プレイヤーが ").
		ItemName(itemName).
		Append(" を入手した。").
		Log()
	return nil
}

// validatePickupItem はアイテム拾得アクションの検証を行う
func (e *Executor) validatePickupItem(ctx Context, world w.World) error {
	// プレイヤー位置にアイテムがあるかチェック
	gridElementRaw := world.Components.GridElement.Get(ctx.Actor)
	if gridElementRaw == nil {
		return fmt.Errorf("位置情報がありません")
	}
	gridElement := gridElementRaw.(*gc.GridElement)
	playerTileX := int(gridElement.X)
	playerTileY := int(gridElement.Y)
	// 同じタイルにアイテムがあるかチェック
	hasItem := false
	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationOnField,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(itemEntity ecs.Entity) {
		itemGrid := world.Components.GridElement.Get(itemEntity).(*gc.GridElement)
		if int(itemGrid.X) == playerTileX && int(itemGrid.Y) == playerTileY {
			hasItem = true
		}
	}))
	if !hasItem {
		return fmt.Errorf("拾えるアイテムがありません")
	}
	return nil
}

// validateWarp はワープアクションの検証を行う
func (e *Executor) validateWarp(_ Context, world w.World) error {
	gameResources := world.Resources.Dungeon.(*resources.Dungeon)
	if gameResources.PlayerTileState.CurrentWarp == nil {
		return fmt.Errorf("ワープホールがありません")
	}
	return nil
}
