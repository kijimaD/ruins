package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
)

// ConsumeStamina はスタミナを消費するエフェクト
type ConsumeStamina struct {
	Amount gc.Amounter // 消費量（固定値または割合）
}

// Apply はスタミナ消費エフェクトをターゲットに適用する
func (c ConsumeStamina) Apply(world w.World, scope *Scope) error {
	if err := c.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)

		switch amount := c.Amount.(type) {
		case gc.RatioAmount:
			consumeAmount := amount.Calc(pools.SP.Max)
			pools.SP.Current = mathutil.Max(0, pools.SP.Current-consumeAmount)
		case gc.NumeralAmount:
			consumeAmount := amount.Calc()
			pools.SP.Current = mathutil.Max(0, pools.SP.Current-consumeAmount)
		default:
			return fmt.Errorf("未対応のスタミナ消費量タイプ: %T", amount)
		}
	}
	return nil
}

// Validate はスタミナ消費エフェクトの妥当性を検証する
func (c ConsumeStamina) Validate(world w.World, scope *Scope) error {
	if c.Amount == nil {
		return errors.New("スタミナ消費量が指定されていません")
	}
	if len(scope.Targets) == 0 {
		return errors.New("スタミナ消費対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if !target.HasComponent(world.Components.Pools) {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (c ConsumeStamina) String() string {
	return fmt.Sprintf("ConsumeStamina(%v)", c.Amount)
}

// RestoreStamina はスタミナを回復するエフェクト
type RestoreStamina struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

// Apply はスタミナ回復エフェクトをターゲットに適用する
func (r RestoreStamina) Apply(world w.World, scope *Scope) error {
	if err := r.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)

		switch amount := r.Amount.(type) {
		case gc.RatioAmount:
			restoreAmount := amount.Calc(pools.SP.Max)
			pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+restoreAmount)
		case gc.NumeralAmount:
			restoreAmount := amount.Calc()
			pools.SP.Current = mathutil.Min(pools.SP.Max, pools.SP.Current+restoreAmount)
		default:
			return fmt.Errorf("未対応のスタミナ回復量タイプ: %T", amount)
		}
	}
	return nil
}

// Validate はスタミナ回復エフェクトの妥当性を検証する
func (r RestoreStamina) Validate(world w.World, scope *Scope) error {
	if r.Amount == nil {
		return errors.New("スタミナ回復量が指定されていません")
	}
	if len(scope.Targets) == 0 {
		return errors.New("スタミナ回復対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if !target.HasComponent(world.Components.Pools) {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (r RestoreStamina) String() string {
	return fmt.Sprintf("RestoreStamina(%v)", r.Amount)
}
