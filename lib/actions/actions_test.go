package actions

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
)

func TestActionID_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		actionID ActionID
		expected string
	}{
		{"Move", ActionMove, "移動"},
		{"Wait", ActionWait, "待機"},
		{"Attack", ActionAttack, "攻撃"},
		{"Null", ActionNull, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.actionID.String()
			if result != tt.expected {
				t.Errorf("ActionID.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestContext(t *testing.T) {
	t.Parallel()
	world, err := w.InitWorld(&gc.Components{})
	if err != nil {
		t.Fatalf("InitWorld failed: %v", err)
	}

	actor := world.Manager.NewEntity()
	target := world.Manager.NewEntity()
	position := gc.Position{X: 5, Y: 10}

	ctx := Context{
		Actor:    actor,
		Target:   &target,
		Position: &position,
		World:    world,
	}

	// コンテキストの基本的な検証
	if ctx.Actor != actor {
		t.Errorf("Context.Actor = %v, want %v", ctx.Actor, actor)
	}
	if ctx.Target == nil || *ctx.Target != target {
		t.Errorf("Context.Target = %v, want %v", ctx.Target, &target)
	}
	if ctx.Position == nil || ctx.Position.X != 5 || ctx.Position.Y != 10 {
		t.Errorf("Context.Position = %v, want Position{X:5, Y:10}", ctx.Position)
	}
}

func TestResult(t *testing.T) {
	t.Parallel()
	result := Result{
		Success:  true,
		ActionID: ActionMove,
		Message:  "移動完了",
	}

	if !result.Success {
		t.Error("Result.Success should be true")
	}
	if result.ActionID != ActionMove {
		t.Errorf("Result.ActionID = %v, want %v", result.ActionID, ActionMove)
	}
	if result.Message != "移動完了" {
		t.Errorf("Result.Message = %v, want %v", result.Message, "移動完了")
	}
}
