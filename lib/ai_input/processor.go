package ai_input

import (
	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Processor はAIエンティティの行動処理を管理する
type Processor struct {
	logger        *logger.Logger
	stateMachine  StateMachine
	actionPlanner ActionPlanner
	visionSystem  VisionSystem
}

// NewProcessor は新しいProcessorを作成する
func NewProcessor() *Processor {
	return &Processor{
		logger:        logger.New(logger.CategoryTurn),
		stateMachine:  NewStateMachine(),
		actionPlanner: NewActionPlanner(),
		visionSystem:  NewVisionSystem(),
	}
}

// ProcessAllEntities は全てのAIエンティティを処理する
func (p *Processor) ProcessAllEntities(world w.World) {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	p.logger.Debug("AI処理開始", "turn", turnManager.TurnNumber, "playerMoves", turnManager.PlayerMoves)

	executor := actions.NewExecutor()
	entityCount := 0

	// AIMoveFSMコンポーネントを持つ全エンティティを処理
	world.Manager.Join(
		world.Components.AIMoveFSM,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++
		p.logger.Debug("AIエンティティを処理中", "entity", entity, "count", entityCount)
		p.ProcessEntity(world, executor, entity)
	}))

	p.logger.Debug("AI処理完了", "処理されたエンティティ数", entityCount, "turn", turnManager.TurnNumber, "playerMoves", turnManager.PlayerMoves)
}

// ProcessEntity は個別のAIエンティティを処理する
func (p *Processor) ProcessEntity(world w.World, executor *actions.Executor, entity ecs.Entity) {
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	p.logger.Debug("AIエンティティ処理開始", "entity", entity)

	// 必要なコンポーネントを取得
	context, err := p.gatherEntityContext(world, entity)
	if err != nil {
		p.logger.Warn("AIエンティティコンテキスト取得失敗", "entity", entity, "error", err.Error())
		return
	}

	// プレイヤー検索
	playerEntity := p.findPlayer(world)
	if playerEntity == nil {
		p.logger.Warn("プレイヤーが見つからない", "entity", entity)
		return
	}

	// 視界チェック
	canSeePlayer := p.visionSystem.CanSeeTarget(world, entity, *playerEntity, context.Vision)
	p.logger.Debug("プレイヤー視界チェック", "entity", entity, "canSee", canSeePlayer)

	// 状態更新
	oldState := context.Roaming.SubState
	p.stateMachine.UpdateState(context.Roaming, canSeePlayer, turnManager.TurnNumber)
	if oldState != context.Roaming.SubState {
		p.logger.Debug("AI状態変化", "entity", entity, "from", oldState, "to", context.Roaming.SubState)
	}

	// 残りターン数を計算してログ出力
	elapsedTurns := turnManager.TurnNumber - context.Roaming.StartSubStateTurn
	remainingTurns := context.Roaming.DurationSubStateTurns - elapsedTurns
	if remainingTurns < 0 {
		remainingTurns = 0
	}
	p.logger.Debug("AIRoaming状態", "entity", entity, "subState", context.Roaming.SubState, "remainingTurns", remainingTurns)

	// アクション決定
	actionCtx, actionID := p.actionPlanner.PlanAction(world, entity, *playerEntity, context, canSeePlayer)

	// アクション実行
	p.logger.Debug("アクション決定", "entity", entity, "action", actionID.String(), "state", context.Roaming.SubState)
	if actionID != actions.ActionNull {
		result, err := executor.Execute(actionID, actionCtx, world)
		if err != nil {
			p.logger.Warn("AIアクション実行失敗", "entity", entity, "action", actionID.String(), "error", err.Error())
		} else {
			p.logger.Debug("AIアクション実行成功", "entity", entity, "action", actionID.String(), "success", result.Success, "state", context.Roaming.SubState, "message", result.Message)
		}
	} else {
		p.logger.Debug("アクション無し", "entity", entity)
	}
	p.logger.Debug("AIエンティティ処理完了", "entity", entity)
}

// EntityContext はAIエンティティの必要な情報をまとめる
type EntityContext struct {
	GridElement *gc.GridElement
	Vision      *gc.AIVision
	Roaming     *gc.AIRoaming
}

// gatherEntityContext はエンティティから必要なコンポーネントを収集する
func (p *Processor) gatherEntityContext(world w.World, entity ecs.Entity) (*EntityContext, error) {
	// GridElementコンポーネント取得
	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
	p.logger.Debug("AIエンティティ位置", "entity", entity, "x", gridElement.X, "y", gridElement.Y)

	// AIVisionコンポーネント確認
	aiVision := world.Components.AIVision.Get(entity)
	if aiVision == nil {
		return nil, &AIError{Type: "component_missing", Message: "AIVisionコンポーネントなし", Entity: entity}
	}
	vision := aiVision.(*gc.AIVision)
	p.logger.Debug("AIVision設定", "entity", entity, "viewDistance", vision.ViewDistance)

	// AIRoamingコンポーネント確認
	aiRoaming := world.Components.AIRoaming.Get(entity)
	if aiRoaming == nil {
		return nil, &AIError{Type: "component_missing", Message: "AIRoamingコンポーネントなし", Entity: entity}
	}
	roaming := aiRoaming.(*gc.AIRoaming)

	return &EntityContext{
		GridElement: gridElement,
		Vision:      vision,
		Roaming:     roaming,
	}, nil
}

// findPlayer はプレイヤーエンティティを探す
func (p *Processor) findPlayer(world w.World) *ecs.Entity {
	var playerEntity *ecs.Entity

	world.Manager.Join(world.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = &entity
	}))

	return playerEntity
}
