package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
)

// FullRecoveryHP は非戦闘時のHP全回復エフェクト（ログ出力なし）
type FullRecoveryHP struct{}

// Apply は非戦闘時HP全回復エフェクトをターゲットに適用する
func (f FullRecoveryHP) Apply(world w.World, scope *Scope) error {
	if err := f.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)
		pools.HP.Current = pools.HP.Max
	}
	return nil
}

// Validate は非戦闘時HP全回復エフェクトの妥当性を検証する
func (f FullRecoveryHP) Validate(world w.World, scope *Scope) error {
	if len(scope.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (f FullRecoveryHP) String() string {
	return "FullRecoveryHP"
}

// FullRecoverySP は非戦闘時のSP全回復エフェクト（ログ出力なし）
type FullRecoverySP struct{}

// Apply は非戦闘時SP全回復エフェクトをターゲットに適用する
func (f FullRecoverySP) Apply(world w.World, scope *Scope) error {
	if err := f.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)
		pools.SP.Current = pools.SP.Max
	}
	return nil
}

// Validate は非戦闘時SP全回復エフェクトの妥当性を検証する
func (f FullRecoverySP) Validate(world w.World, scope *Scope) error {
	if len(scope.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (f FullRecoverySP) String() string {
	return "FullRecoverySP"
}


// RecoverySP は非戦闘時のSP回復エフェクト（ログ出力なし）
type RecoverySP struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

// Apply は非戦闘時SP部分回復エフェクトをターゲットに適用する
func (r RecoverySP) Apply(world w.World, scope *Scope) error {
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
			return fmt.Errorf("未対応の回復量タイプ: %T", amount)
		}
	}
	return nil
}

// Validate は非戦闘時SP部分回復エフェクトの妥当性を検証する
func (r RecoverySP) Validate(world w.World, scope *Scope) error {
	if r.Amount == nil {
		return errors.New("回復量が指定されていません")
	}
	if len(scope.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (r RecoverySP) String() string {
	return fmt.Sprintf("RecoverySP(%v)", r.Amount)
}
