package systems

import (
	"testing"

	"github.com/kijimaD/ruins/lib/actions"
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/turns"
)

func TestAISystem(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーエンティティ作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(10), Y: gc.Tile(10)})

	// AIエンティティ作成
	aiEntity := world.Manager.NewEntity()
	aiEntity.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})
	aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})
	aiEntity.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(100), // 3タイル程度の視界
		TargetEntity: &player,
	})
	aiEntity.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,
		DurationSubStateTurns: 2,
	})

	// システム実行前の位置を記録
	initialGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	initialX, initialY := int(initialGrid.X), int(initialGrid.Y)

	// AISystem実行
	AISystem(world)

	// 実行後の位置を確認
	finalGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	finalX, finalY := int(finalGrid.X), int(finalGrid.Y)

	// 位置が変化したかチェック（移動またはそのまま）
	moved := (finalX != initialX) || (finalY != initialY)
	t.Logf("AI移動: (%d,%d) -> (%d,%d), moved: %v", initialX, initialY, finalX, finalY, moved)

	// AIエンティティがまだ存在することを確認
	if !aiEntity.HasComponent(world.Components.AIMoveFSM) {
		t.Error("AIエンティティのAIMoveFSMコンポーネントが失われた")
	}

	// AIRoamingの状態を確認
	roamingComp := world.Components.AIRoaming.Get(aiEntity).(*gc.AIRoaming)
	t.Logf("AI状態: %v", roamingComp.SubState)
}

func TestUpdateAIState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialState   gc.AIRoamingSubState
		canSeePlayer   bool
		currentTurn    int
		expectedStates []gc.AIRoamingSubState // 複数の可能な状態
	}{
		{
			name:           "待機状態でプレイヤー発見",
			initialState:   gc.AIRoamingWaiting,
			canSeePlayer:   true,
			currentTurn:    2,
			expectedStates: []gc.AIRoamingSubState{gc.AIRoamingChasing},
		},
		{
			name:           "待機状態でターン経過",
			initialState:   gc.AIRoamingWaiting,
			canSeePlayer:   false,
			currentTurn:    5, // 1ターン目から開始、2ターン持続なので3ターン目で終了
			expectedStates: []gc.AIRoamingSubState{gc.AIRoamingDriving},
		},
		{
			name:           "移動状態でプレイヤー発見",
			initialState:   gc.AIRoamingDriving,
			canSeePlayer:   true,
			currentTurn:    2,
			expectedStates: []gc.AIRoamingSubState{gc.AIRoamingChasing},
		},
		{
			name:           "追跡状態でプレイヤー継続視認",
			initialState:   gc.AIRoamingChasing,
			canSeePlayer:   true,
			currentTurn:    2,
			expectedStates: []gc.AIRoamingSubState{gc.AIRoamingChasing},
		},
		{
			name:           "追跡状態でプレイヤー見失い（短期間）",
			initialState:   gc.AIRoamingChasing,
			canSeePlayer:   false,
			currentTurn:    3, // 3ターン以内
			expectedStates: []gc.AIRoamingSubState{gc.AIRoamingChasing},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// テスト用のAIRoamingを作成
			roaming := &gc.AIRoaming{
				SubState:              tt.initialState,
				StartSubStateTurn:     1,
				DurationSubStateTurns: 2, // 2ターンのデフォルト持続時間
			}

			// 状態更新を実行
			updateAIState(roaming, tt.canSeePlayer, tt.currentTurn)

			// 結果を確認
			found := false
			for _, expectedState := range tt.expectedStates {
				if roaming.SubState == expectedState {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("updateAIState() = %v, want one of %v", roaming.SubState, tt.expectedStates)
			}

			t.Logf("状態遷移: %v -> %v", tt.initialState, roaming.SubState)
		})
	}
}

func TestFindPlayer(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーが存在しない場合
	player := findPlayer(world)
	if player != nil {
		t.Error("プレイヤーが存在しないはずなのにプレイヤーが見つかった")
	}

	// プレイヤーエンティティ作成
	playerEntity := world.Manager.NewEntity()
	playerEntity.AddComponent(world.Components.Player, gc.Player{})

	// プレイヤーが存在する場合
	player = findPlayer(world)
	if player == nil {
		t.Error("プレイヤーが見つからない")
	} else if *player != playerEntity {
		t.Errorf("間違ったプレイヤーエンティティが返された: got %v, want %v", *player, playerEntity)
	}
}

func TestCheckPlayerInSight(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーエンティティ作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(10), Y: gc.Tile(10)})

	// AIエンティティ作成
	aiEntity := world.Manager.NewEntity()

	tests := []struct {
		name         string
		aiX, aiY     int
		viewDistance gc.Pixel
		expected     bool
	}{
		{
			name: "視界内（近距離）",
			aiX:  9, aiY: 10,
			viewDistance: gc.Pixel(64), // 2タイル分
			expected:     true,
		},
		{
			name: "視界内（ぎりぎり）",
			aiX:  8, aiY: 10,
			viewDistance: gc.Pixel(96), // 3タイル分（距離2をカバーするため）
			expected:     true,
		},
		{
			name: "視界外",
			aiX:  5, aiY: 5,
			viewDistance: gc.Pixel(64), // 2タイル分
			expected:     false,
		},
		{
			name: "視界が広い場合",
			aiX:  5, aiY: 5,
			viewDistance: gc.Pixel(320), // 10タイル分
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// AIの位置設定
			aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{
				X: gc.Tile(tt.aiX),
				Y: gc.Tile(tt.aiY),
			})

			vision := &gc.AIVision{
				ViewDistance: tt.viewDistance,
				TargetEntity: &player,
			}

			result := checkPlayerInSight(world, aiEntity, player, vision)
			if result != tt.expected {
				t.Errorf("checkPlayerInSight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPlanChaseAction(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーエンティティ作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(10), Y: gc.Tile(10)})

	// AIエンティティ作成
	aiEntity := world.Manager.NewEntity()

	tests := []struct {
		name              string
		aiX, aiY          int
		expectedDirection string
		expectMove        bool
	}{
		{
			name: "右に移動",
			aiX:  8, aiY: 10,
			expectedDirection: "右",
			expectMove:        true,
		},
		{
			name: "左に移動",
			aiX:  12, aiY: 10,
			expectedDirection: "左",
			expectMove:        true,
		},
		{
			name: "上に移動",
			aiX:  10, aiY: 12,
			expectedDirection: "上",
			expectMove:        true,
		},
		{
			name: "下に移動",
			aiX:  10, aiY: 8,
			expectedDirection: "下",
			expectMove:        true,
		},
		{
			name: "斜め移動",
			aiX:  8, aiY: 8,
			expectedDirection: "右下",
			expectMove:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// AIの位置設定
			aiGrid := &gc.GridElement{X: gc.Tile(tt.aiX), Y: gc.Tile(tt.aiY)}
			aiEntity.AddComponent(world.Components.GridElement, aiGrid)

			ctx, actionID := planChaseAction(world, aiEntity, player, aiGrid)

			if !tt.expectMove {
				if actionID != actions.ActionWait {
					t.Errorf("期待されるアクション: ActionWait, 実際: %v", actionID)
				}
				return
			}

			if actionID != actions.ActionMove {
				t.Errorf("期待されるアクション: ActionMove, 実際: %v", actionID)
			}

			if ctx.Actor != aiEntity {
				t.Errorf("Actor = %v, want %v", ctx.Actor, aiEntity)
			}

			if ctx.Dest == nil {
				t.Fatal("Dest is nil")
			}

			// 移動方向が正しいかチェック
			destX, destY := int(ctx.Dest.X), int(ctx.Dest.Y)
			expectedX := tt.aiX
			expectedY := tt.aiY

			// プレイヤーに近づく方向に移動しているかチェック
			if tt.aiX < 10 {
				expectedX++
			} else if tt.aiX > 10 {
				expectedX--
			}

			if tt.aiY < 10 {
				expectedY++
			} else if tt.aiY > 10 {
				expectedY--
			}

			if destX != expectedX || destY != expectedY {
				t.Errorf("移動先 = (%d,%d), want (%d,%d)", destX, destY, expectedX, expectedY)
			}
		})
	}
}

func TestPlanRandomMoveAction(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// AIエンティティ作成
	aiEntity := world.Manager.NewEntity()
	aiGrid := &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)}

	// 複数回実行してランダム性を確認
	var waitCount, moveCount int
	totalTests := 100

	for i := 0; i < totalTests; i++ {
		ctx, actionID := planRandomMoveAction(world, aiEntity, aiGrid)

		if ctx.Actor != aiEntity {
			t.Errorf("Actor = %v, want %v", ctx.Actor, aiEntity)
		}

		switch actionID {
		case actions.ActionWait:
			waitCount++
		case actions.ActionMove:
			moveCount++
			if ctx.Dest == nil {
				t.Fatal("Dest is nil for ActionMove")
			}
			// 移動先が隣接タイルかチェック
			destX, destY := int(ctx.Dest.X), int(ctx.Dest.Y)
			diffX := destX - 5
			if diffX < 0 {
				diffX = -diffX
			}
			diffY := destY - 5
			if diffY < 0 {
				diffY = -diffY
			}
			if diffX > 1 || diffY > 1 {
				t.Errorf("移動先が隣接タイルではない: (%d,%d)", destX, destY)
			}
		default:
			t.Errorf("予期しないアクション: %v", actionID)
		}
	}

	// 待機と移動の両方が発生することを確認
	if waitCount == 0 {
		t.Error("待機アクションが発生しなかった")
	}
	if moveCount == 0 {
		t.Error("移動アクションが発生しなかった")
	}

	t.Logf("待機: %d回, 移動: %d回 (全%d回)", waitCount, moveCount, totalTests)
}

func TestCanMoveTo(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// AIエンティティ作成
	aiEntity := world.Manager.NewEntity()
	aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})

	// 壁エンティティ作成
	wallEntity := world.Manager.NewEntity()
	wallEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(6), Y: gc.Tile(5)})
	wallEntity.AddComponent(world.Components.BlockPass, gc.BlockPass{})

	tests := []struct {
		name     string
		destX    int
		destY    int
		expected bool
	}{
		{
			name:     "空いているタイル",
			destX:    4,
			destY:    5,
			expected: true,
		},
		{
			name:     "壁があるタイル",
			destX:    6,
			destY:    5,
			expected: false,
		},
		{
			name:     "自分の現在位置",
			destX:    5,
			destY:    5,
			expected: true, // 自分自身は除外される
		},
		{
			name:     "遠い空のタイル",
			destX:    10,
			destY:    10,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := canMoveTo(world, tt.destX, tt.destY, aiEntity)
			if result != tt.expected {
				t.Errorf("canMoveTo(%d, %d) = %v, want %v", tt.destX, tt.destY, result, tt.expected)
			}
		})
	}
}

func TestAIStateMachine(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーエンティティ作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(10), Y: gc.Tile(5)})

	// AIエンティティ作成（プレイヤーの視界外）
	aiEntity := world.Manager.NewEntity()
	aiEntity.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})
	aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})
	aiEntity.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(64), // 2タイル分の視界
		TargetEntity: &player,
	})
	aiEntity.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1, // 1ターン目から開始
		DurationSubStateTurns: 2, // 2ターンの持続時間
	})

	// TurnManagerのターンを進める（1->3ターンに進めて待機時間を超過させる）
	turnManager := world.Resources.TurnManager.(*turns.TurnManager)
	turnManager.TurnNumber = 5 // 1+2=3ターンを超過

	// 最初の状態を記録
	initialRoaming := world.Components.AIRoaming.Get(aiEntity).(*gc.AIRoaming)
	initialState := initialRoaming.SubState

	// AISystem実行
	AISystem(world)

	// 実行後の状態を確認
	finalRoaming := world.Components.AIRoaming.Get(aiEntity).(*gc.AIRoaming)
	finalState := finalRoaming.SubState

	t.Logf("状態遷移: %v -> %v (turn: %d)", initialState, finalState, turnManager.TurnNumber)

	// プレイヤーが視界外で待機時間が過ぎているので、移動状態になるはず
	if initialState == gc.AIRoamingWaiting && finalState != gc.AIRoamingDriving {
		t.Errorf("期待される状態遷移: %v -> %v, 実際: %v -> %v",
			gc.AIRoamingWaiting, gc.AIRoamingDriving, initialState, finalState)
	}
}

func TestAIChaseMode(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成（TurnManagerを含む）
	world := CreateTestWorldWithResources(t)

	// プレイヤーエンティティ作成
	player := world.Manager.NewEntity()
	player.AddComponent(world.Components.Player, gc.Player{})
	player.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(8), Y: gc.Tile(5)})

	// AIエンティティ作成（プレイヤーの視界内）
	aiEntity := world.Manager.NewEntity()
	aiEntity.AddComponent(world.Components.AIMoveFSM, &gc.AIMoveFSM{})
	aiEntity.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})
	aiEntity.AddComponent(world.Components.AIVision, &gc.AIVision{
		ViewDistance: gc.Pixel(200), // 十分な視界
		TargetEntity: &player,
	})
	aiEntity.AddComponent(world.Components.AIRoaming, &gc.AIRoaming{
		SubState:              gc.AIRoamingWaiting,
		StartSubStateTurn:     1,
		DurationSubStateTurns: 10,
	})

	// 初期位置を記録
	initialGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	initialX, initialY := int(initialGrid.X), int(initialGrid.Y)

	// AISystem実行
	AISystem(world)

	// 実行後の状態と位置を確認
	finalRoaming := world.Components.AIRoaming.Get(aiEntity).(*gc.AIRoaming)
	finalGrid := world.Components.GridElement.Get(aiEntity).(*gc.GridElement)
	finalX, finalY := int(finalGrid.X), int(finalGrid.Y)

	// プレイヤーを発見して追跡状態になるはず
	if finalRoaming.SubState != gc.AIRoamingChasing {
		t.Errorf("期待される状態: %v, 実際: %v", gc.AIRoamingChasing, finalRoaming.SubState)
	}

	// プレイヤーに向かって移動したかチェック
	playerGrid := world.Components.GridElement.Get(player).(*gc.GridElement)
	playerX, playerY := int(playerGrid.X), int(playerGrid.Y)

	// プレイヤーに近づいたかチェック
	initialDiffX := playerX - initialX
	if initialDiffX < 0 {
		initialDiffX = -initialDiffX
	}
	initialDiffY := playerY - initialY
	if initialDiffY < 0 {
		initialDiffY = -initialDiffY
	}
	initialDistance := initialDiffX + initialDiffY

	finalDiffX := playerX - finalX
	if finalDiffX < 0 {
		finalDiffX = -finalDiffX
	}
	finalDiffY := playerY - finalY
	if finalDiffY < 0 {
		finalDiffY = -finalDiffY
	}
	finalDistance := finalDiffX + finalDiffY

	if finalDistance > initialDistance {
		t.Errorf("AIがプレイヤーから遠ざかった: 距離 %d -> %d", initialDistance, finalDistance)
	}

	t.Logf("AI追跡: (%d,%d) -> (%d,%d), プレイヤー: (%d,%d), 距離: %d -> %d",
		initialX, initialY, finalX, finalY, playerX, playerY, initialDistance, finalDistance)
}
