package actions

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

func TestExecutor_ExecuteMove(t *testing.T) {
	t.Parallel()

	// テスト用のワールド作成
	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}

	// テスト用のエンティティ作成
	actor := world.Manager.NewEntity()
	actor.AddComponent(world.Components.GridElement, &gc.GridElement{X: gc.Tile(5), Y: gc.Tile(5)})

	executor := NewExecutor()

	// テストケース
	tests := []struct {
		name          string
		ctx           Context
		expectSuccess bool
		expectError   bool
	}{
		{
			name: "正常な移動",
			ctx: Context{
				Actor:    actor,
				Position: &gc.Position{X: gc.Pixel(6), Y: gc.Pixel(6)},
				World:    world,
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "移動先未指定",
			ctx: Context{
				Actor:    actor,
				Position: nil,
				World:    world,
			},
			expectSuccess: false,
			expectError:   true,
		},
		{
			name: "範囲外移動",
			ctx: Context{
				Actor:    actor,
				Position: &gc.Position{X: gc.Pixel(-1), Y: gc.Pixel(0)},
				World:    world,
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := executor.Execute(ActionMove, tt.ctx)

			// エラーのチェック
			if (err != nil) != tt.expectError {
				t.Errorf("Execute() error = %v, expectError %v", err, tt.expectError)
				return
			}

			// 結果のチェック
			if result == nil {
				t.Fatal("Execute() returned nil result")
			}

			if result.Success != tt.expectSuccess {
				t.Errorf("Execute() result.Success = %v, expect %v", result.Success, tt.expectSuccess)
			}

			if result.ActionID != ActionMove {
				t.Errorf("Execute() result.ActionID = %v, expect %v", result.ActionID, ActionMove)
			}

			// 正常な移動の場合、実際に位置が変更されたかチェック
			if tt.expectSuccess && tt.ctx.Position != nil {
				gridElement := world.Components.GridElement.Get(actor).(*gc.GridElement)
				if int(gridElement.X) != int(tt.ctx.Position.X) || int(gridElement.Y) != int(tt.ctx.Position.Y) {
					t.Errorf("Position not updated correctly. Got (%d,%d), expect (%d,%d)",
						gridElement.X, gridElement.Y, int(tt.ctx.Position.X), int(tt.ctx.Position.Y))
				}
			}
		})
	}
}

func TestExecutor_ExecuteWait(t *testing.T) {
	t.Parallel()

	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}

	actor := world.Manager.NewEntity()
	executor := NewExecutor()

	ctx := Context{
		Actor: actor,
		World: world,
	}

	result, err := executor.Execute(ActionWait, ctx)

	if err != nil {
		t.Errorf("Execute(ActionWait) unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Execute(ActionWait) returned nil result")
	}

	if !result.Success {
		t.Errorf("Execute(ActionWait) result.Success = false, expect true")
	}

	if result.ActionID != ActionWait {
		t.Errorf("Execute(ActionWait) result.ActionID = %v, expect %v", result.ActionID, ActionWait)
	}
}

func TestExecutor_ExecuteAttack(t *testing.T) {
	t.Parallel()

	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}

	actor := world.Manager.NewEntity()
	target := world.Manager.NewEntity()
	executor := NewExecutor()

	tests := []struct {
		name          string
		ctx           Context
		expectSuccess bool
		expectError   bool
	}{
		{
			name: "正常な攻撃",
			ctx: Context{
				Actor:  actor,
				Target: &target,
				World:  world,
			},
			expectSuccess: true,
			expectError:   false,
		},
		{
			name: "ターゲット未指定",
			ctx: Context{
				Actor:  actor,
				Target: nil,
				World:  world,
			},
			expectSuccess: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := executor.Execute(ActionAttack, tt.ctx)

			if (err != nil) != tt.expectError {
				t.Errorf("Execute() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if result == nil {
				t.Fatal("Execute() returned nil result")
			}

			if result.Success != tt.expectSuccess {
				t.Errorf("Execute() result.Success = %v, expect %v", result.Success, tt.expectSuccess)
			}

			if result.ActionID != ActionAttack {
				t.Errorf("Execute() result.ActionID = %v, expect %v", result.ActionID, ActionAttack)
			}
		})
	}
}

func TestExecutor_ValidateAction(t *testing.T) {
	t.Parallel()

	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}

	executor := NewExecutor()

	tests := []struct {
		name      string
		actionID  ActionID
		ctx       Context
		expectErr bool
	}{
		{
			name:     "Actorが0の場合",
			actionID: ActionMove,
			ctx: Context{
				Actor: ecs.Entity(0),
				World: world,
			},
			expectErr: true,
		},
		{
			name:     "不明なアクション",
			actionID: ActionID(999),
			ctx: Context{
				Actor: ecs.Entity(1),
				World: world,
			},
			expectErr: true,
		},
		{
			name:     "待機は常に有効",
			actionID: ActionWait,
			ctx: Context{
				Actor: ecs.Entity(1),
				World: world,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := executor.validateAction(tt.actionID, tt.ctx)
			if (err != nil) != tt.expectErr {
				t.Errorf("validateAction() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestNewExecutor(t *testing.T) {
	t.Parallel()

	executor := NewExecutor()

	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if executor.processor == nil {
		t.Error("NewExecutor() processor is nil")
	}

	if executor.logger == nil {
		t.Error("NewExecutor() logger is nil")
	}
}
