package aiinput

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

	manager := actions.NewActivityManager(logger.New(logger.CategoryAction))
	entityCount := 0

	// AIMoveFSMコンポーネントを持つ全エンティティを処理
	world.Manager.Join(
		world.Components.AIMoveFSM,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++
		p.logger.Debug("AIエンティティを処理中", "entity", entity, "count", entityCount)
		p.ProcessEntity(world, manager, entity)
	}))

	p.logger.Debug("AI処理完了", "処理されたエンティティ数", entityCount, "turn", turnManager.TurnNumber, "playerMoves", turnManager.PlayerMoves)
}

// ProcessEntity は個別のAIエンティティを処理する
func (p *Processor) ProcessEntity(world w.World, manager *actions.ActivityManager, entity ecs.Entity) {
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
		return
	}

	// プレイヤーエンティティの有効性チェック
	if !playerEntity.HasComponent(world.Components.GridElement) {
		p.logger.Warn("プレイヤーエンティティが無効（GridElementなし）", "entity", entity, "player", *playerEntity)
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

	// APが残っている限り連続してアクションを実行
	actionsExecuted := 0
	maxActions := 10 // 無限ループを防ぐためのリミット

	for actionsExecuted < maxActions {
		// アクション決定
		activityType, actionParams := p.actionPlanner.PlanAction(world, entity, *playerEntity, context, canSeePlayer)

		// アクション実行
		p.logger.Debug("アクション決定", "entity", entity, "activity", activityType.String(), "state", context.Roaming.SubState, "actions", actionsExecuted)
		if activityType.String() == "" {
			p.logger.Debug("アクション無し", "entity", entity)
			break
		}

		// アクティビティタイプに応じたAPコストを計算
		actionCost, _ := actions.GetActivityCost(activityType)
		if !turnManager.CanEntityAct(world, entity, actionCost) {
			p.logger.Debug("AP不足でアクション実行不可", "entity", entity, "activity", activityType.String(), "cost", actionCost)
			break
		}

		result, err := manager.Execute(activityType, actionParams, world)
		if err != nil {
			p.logger.Warn("AIアクション実行失敗", "entity", entity, "activity", activityType.String(), "error", err.Error())
			break
		}

		p.logger.Debug("AIアクション実行成功", "entity", entity, "activity", activityType.String(), "success", result.Success, "state", context.Roaming.SubState, "message", result.Message)
		actionsExecuted++

		// アクション失敗時は停止
		if !result.Success {
			p.logger.Debug("アクション失敗により停止", "entity", entity, "activity", activityType.String())
			break
		}
	}

	p.logger.Debug("AIエンティティ処理完了", "entity", entity, "実行されたアクション数", actionsExecuted)
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

	world.Manager.Join(world.Components.Player, world.Components.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
		// 死亡しているエンティティは除外
		if entity.HasComponent(world.Components.Dead) {
			return
		}
		playerEntity = &entity
	}))

	return playerEntity
}
