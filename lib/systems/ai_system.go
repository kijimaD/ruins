package systems

import (
	"math"
	"math/rand/v2"

	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/logger"
	"github.com/kijimaD/ruins/lib/turns"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AISystem はAIエンティティの行動を処理するシステム
func AISystem(world w.World) {
	aiLogger := logger.New(logger.CategoryTurn)
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)

	aiLogger.Debug("AISystem開始", "turn", turnManager.TurnNumber, "playerMoves", turnManager.PlayerMoves)

	executor := actions.NewExecutor()

	// AIエンティティの数をカウント
	entityCount := 0
	// AIMoveFSMコンポーネントを持つ全エンティティを処理
	world.Manager.Join(
		world.Components.AIMoveFSM,
		world.Components.GridElement,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		entityCount++
		aiLogger.Debug("AIエンティティを処理中", "entity", entity, "count", entityCount)
		processAIEntity(world, executor, entity)
	}))

	aiLogger.Debug("AISystem完了", "処理されたエンティティ数", entityCount, "turn", turnManager.TurnNumber, "playerMoves", turnManager.PlayerMoves)
}

// processAIEntity は個別のAIエンティティを処理する
func processAIEntity(world w.World, executor *actions.Executor, entity ecs.Entity) {
	aiLogger := logger.New(logger.CategoryTurn)
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	aiLogger.Debug("processAIEntity開始", "entity", entity)

	// 基本コンポーネント取得
	gridElement := world.Components.GridElement.Get(entity).(*gc.GridElement)
	aiLogger.Debug("AIエンティティ位置", "entity", entity, "x", gridElement.X, "y", gridElement.Y)

	// AIVisionコンポーネントを確認
	aiVision := world.Components.AIVision.Get(entity)
	if aiVision == nil {
		aiLogger.Warn("AIVisionコンポーネントなし", "ai", entity)
		return
	}
	vision := aiVision.(*gc.AIVision)
	aiLogger.Debug("AIVision設定", "entity", entity, "viewDistance", vision.ViewDistance)

	// AIRoamingコンポーネントを確認
	aiRoaming := world.Components.AIRoaming.Get(entity)
	if aiRoaming == nil {
		aiLogger.Warn("AIRoamingコンポーネントなし", "ai", entity)
		return
	}
	roaming := aiRoaming.(*gc.AIRoaming)

	// 現在の状態での残りターン数を計算
	elapsedTurns := turnManager.TurnNumber - roaming.StartSubStateTurn
	remainingTurns := roaming.DurationSubStateTurns - elapsedTurns
	if remainingTurns < 0 {
		remainingTurns = 0
	}

	aiLogger.Debug("AIRoaming状態", "entity", entity, "subState", roaming.SubState, "remainingTurns", remainingTurns)

	// プレイヤーを探す
	playerEntity := findPlayer(world)
	if playerEntity == nil {
		aiLogger.Warn("プレイヤーが見つからない", "ai", entity)
		return
	}
	aiLogger.Debug("プレイヤー発見", "entity", entity, "player", *playerEntity)

	// 視界チェック
	canSeePlayer := checkPlayerInSight(world, entity, *playerEntity, vision)
	aiLogger.Debug("プレイヤー視界チェック", "entity", entity, "canSee", canSeePlayer)

	// 状態更新
	oldState := roaming.SubState
	updateAIState(roaming, canSeePlayer, turnManager.TurnNumber)
	if oldState != roaming.SubState {
		aiLogger.Debug("AI状態変化", "entity", entity, "from", oldState, "to", roaming.SubState)
	}

	// 現在の状態に基づいてアクション決定
	var actionCtx actions.Context
	var actionID actions.ActionID

	switch roaming.SubState {
	case gc.AIRoamingChasing:
		// 追跡モード：プレイヤーに向かって移動
		actionCtx, actionID = planChaseAction(world, entity, *playerEntity, gridElement)
		aiLogger.Debug("追跡モード", "ai", entity, "player", *playerEntity)

	case gc.AIRoamingDriving:
		// 移動モード：ランダム移動
		actionCtx, actionID = planRandomMoveAction(world, entity, gridElement)
		aiLogger.Debug("移動モード", "ai", entity)

	case gc.AIRoamingWaiting:
		// 待機モード：何もしない
		actionCtx, actionID = actions.Context{Actor: entity}, actions.ActionWait
		aiLogger.Debug("待機モード", "ai", entity)

	default:
		// 不明な状態：待機
		actionCtx, actionID = actions.Context{Actor: entity}, actions.ActionWait
		aiLogger.Warn("不明なAI状態", "ai", entity, "state", roaming.SubState)
	}

	// アクション実行
	aiLogger.Debug("アクション決定", "entity", entity, "action", actionID.String(), "state", roaming.SubState)
	if actionID != actions.ActionNull {
		result, err := executor.Execute(actionID, actionCtx, world)
		if err != nil {
			aiLogger.Warn("AIアクション実行失敗", "ai", entity, "action", actionID.String(), "error", err.Error())
		} else {
			aiLogger.Debug("AIアクション実行成功", "ai", entity, "action", actionID.String(), "success", result.Success, "state", roaming.SubState, "message", result.Message)
		}
	} else {
		aiLogger.Debug("アクション無し", "entity", entity)
	}
	aiLogger.Debug("processAIEntity完了", "entity", entity)
}

// updateAIState はAIの状態を更新する有限状態機械
func updateAIState(roaming *gc.AIRoaming, canSeePlayer bool, currentTurn int) {
	elapsedTurns := currentTurn - roaming.StartSubStateTurn

	// 現在の状態によって遷移ロジックを決定
	switch roaming.SubState {
	case gc.AIRoamingWaiting:
		// 待機状態からの遷移
		if canSeePlayer {
			// プレイヤー発見 → 追跡状態へ
			roaming.SubState = gc.AIRoamingChasing
			roaming.StartSubStateTurn = currentTurn
			roaming.DurationSubStateTurns = 10 + rand.IntN(5) // 10-14ターン追跡
		} else if elapsedTurns >= roaming.DurationSubStateTurns {
			// 待機ターン終了 → 移動状態へ
			roaming.SubState = gc.AIRoamingDriving
			roaming.StartSubStateTurn = currentTurn
			roaming.DurationSubStateTurns = 3 + rand.IntN(7) // 3-9ターン移動
		}

	case gc.AIRoamingDriving:
		// 移動状態からの遷移
		if canSeePlayer {
			// プレイヤー発見 → 追跡状態へ
			roaming.SubState = gc.AIRoamingChasing
			roaming.StartSubStateTurn = currentTurn
			roaming.DurationSubStateTurns = 10 + rand.IntN(5) // 10-14ターン追跡
		} else if elapsedTurns >= roaming.DurationSubStateTurns {
			// 移動ターン終了 → 待機状態へ
			roaming.SubState = gc.AIRoamingWaiting
			roaming.StartSubStateTurn = currentTurn
			roaming.DurationSubStateTurns = 2 + rand.IntN(4) // 2-5ターン待機
		}

	case gc.AIRoamingChasing:
		// 追跡状態からの遷移
		if !canSeePlayer {
			// プレイヤーを見失った場合
			if elapsedTurns >= 3 {
				// 3ターン以上見失った → 移動状態へ
				roaming.SubState = gc.AIRoamingDriving
				roaming.StartSubStateTurn = currentTurn
				roaming.DurationSubStateTurns = 5 + rand.IntN(5) // 5-9ターン移動
			}
			// 3ターン以内なら追跡継続
		} else if elapsedTurns >= roaming.DurationSubStateTurns {
			// 追跡ターン終了 → 待機状態へ
			roaming.SubState = gc.AIRoamingWaiting
			roaming.StartSubStateTurn = currentTurn
			roaming.DurationSubStateTurns = 3 + rand.IntN(4) // 3-6ターン待機
		} else {
			// プレイヤー視認中：追跡継続、ターンリセット
			roaming.StartSubStateTurn = currentTurn
		}

	default:
		// 不明な状態：待機状態に初期化
		roaming.SubState = gc.AIRoamingWaiting
		roaming.StartSubStateTurn = currentTurn
		roaming.DurationSubStateTurns = 2 + rand.IntN(3) // 2-4ターン待機
	}
}

// findPlayer はプレイヤーエンティティを探す
func findPlayer(world w.World) *ecs.Entity {
	var playerEntity *ecs.Entity

	world.Manager.Join(world.Components.Player).Visit(ecs.Visit(func(entity ecs.Entity) {
		playerEntity = &entity
	}))

	return playerEntity
}

// checkPlayerInSight はプレイヤーが視界内にいるかチェック
func checkPlayerInSight(world w.World, aiEntity, playerEntity ecs.Entity, vision *gc.AIVision) bool {
	aiGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	playerGrid := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)

	// 距離計算（タイル単位）
	dx := float64(int(playerGrid.X) - int(aiGrid.X))
	dy := float64(int(playerGrid.Y) - int(aiGrid.Y))
	distance := math.Sqrt(dx*dx + dy*dy)

	// 視界距離内かチェック（タイル単位で計算）
	viewDistanceInTiles := float64(vision.ViewDistance) / 32.0 // 仮にタイル1つ=32ピクセル

	return distance <= viewDistanceInTiles
}

// planChaseAction はプレイヤー追跡アクションを計画
func planChaseAction(world w.World, aiEntity, playerEntity ecs.Entity, aiGrid *gc.GridElement) (actions.Context, actions.ActionID) {
	playerGrid := world.Components.GridElement.Get(playerEntity).(*gc.GridElement)

	// プレイヤーに向かう方向を計算
	dx := int(playerGrid.X) - int(aiGrid.X)
	dy := int(playerGrid.Y) - int(aiGrid.Y)

	// 移動候補を優先度順で試行
	var moveCandidates []struct{ x, y int }

	// 直接的な方向を最優先
	if dx != 0 && dy != 0 {
		// 斜め移動
		moveX := 0
		moveY := 0
		if dx > 0 {
			moveX = 1
		} else {
			moveX = -1
		}
		if dy > 0 {
			moveY = 1
		} else {
			moveY = -1
		}
		moveCandidates = append(moveCandidates, struct{ x, y int }{moveX, moveY})

		// 代替案として軸に沿った移動
		absDx := dx
		if absDx < 0 {
			absDx = -absDx
		}
		absDy := dy
		if absDy < 0 {
			absDy = -absDy
		}
		if absDx > absDy {
			moveCandidates = append(moveCandidates, struct{ x, y int }{moveX, 0})
			moveCandidates = append(moveCandidates, struct{ x, y int }{0, moveY})
		} else {
			moveCandidates = append(moveCandidates, struct{ x, y int }{0, moveY})
			moveCandidates = append(moveCandidates, struct{ x, y int }{moveX, 0})
		}
	} else if dx != 0 {
		// 水平移動のみ
		moveX := 1
		if dx < 0 {
			moveX = -1
		}
		moveCandidates = append(moveCandidates, struct{ x, y int }{moveX, 0})
		// 代替案として垂直移動
		moveCandidates = append(moveCandidates, struct{ x, y int }{0, 1})
		moveCandidates = append(moveCandidates, struct{ x, y int }{0, -1})
	} else if dy != 0 {
		// 垂直移動のみ
		moveY := 1
		if dy < 0 {
			moveY = -1
		}
		moveCandidates = append(moveCandidates, struct{ x, y int }{0, moveY})
		// 代替案として水平移動
		moveCandidates = append(moveCandidates, struct{ x, y int }{1, 0})
		moveCandidates = append(moveCandidates, struct{ x, y int }{-1, 0})
	}

	// 移動可能な候補を探す
	for _, candidate := range moveCandidates {
		destX := int(aiGrid.X) + candidate.x
		destY := int(aiGrid.Y) + candidate.y

		if canMoveTo(world, destX, destY, aiEntity) {
			return actions.Context{
				Actor: aiEntity,
				Dest:  &gc.Position{X: gc.Pixel(destX), Y: gc.Pixel(destY)},
			}, actions.ActionMove
		}
	}

	// どこにも移動できない場合は待機
	return actions.Context{Actor: aiEntity}, actions.ActionWait
}

// planRandomMoveAction はランダム移動アクションを計画
func planRandomMoveAction(world w.World, aiEntity ecs.Entity, aiGrid *gc.GridElement) (actions.Context, actions.ActionID) {
	// 30%の確率で待機
	if rand.Float64() < 0.3 {
		return actions.Context{Actor: aiEntity}, actions.ActionWait
	}

	// ランダムに隣接する8方向から選択
	directions := []struct{ x, y int }{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	// 移動可能な方向をシャッフルして試行
	shuffledDirections := make([]struct{ x, y int }, len(directions))
	copy(shuffledDirections, directions)

	// Fisher-Yatesアルゴリズムでシャッフル
	for i := len(shuffledDirections) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		shuffledDirections[i], shuffledDirections[j] = shuffledDirections[j], shuffledDirections[i]
	}

	// 移動可能な方向を探す
	for _, direction := range shuffledDirections {
		destX := int(aiGrid.X) + direction.x
		destY := int(aiGrid.Y) + direction.y

		if canMoveTo(world, destX, destY, aiEntity) {
			return actions.Context{
				Actor: aiEntity,
				Dest:  &gc.Position{X: gc.Pixel(destX), Y: gc.Pixel(destY)},
			}, actions.ActionMove
		}
	}

	// どこにも移動できない場合は待機
	return actions.Context{Actor: aiEntity}, actions.ActionWait
}
