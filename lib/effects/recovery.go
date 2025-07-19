package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
)

// FullRecoveryHP は非戦闘時のHP全回復エフェクト（ログ出力なし）
type FullRecoveryHP struct{}

func (f FullRecoveryHP) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)
		pools.HP.Current = pools.HP.Max
	}
	return nil
}

func (f FullRecoveryHP) Validate(ctx *Context) error {
	if len(ctx.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	
	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range ctx.Targets {
		if ctx.World.Components.Pools.Get(target) == nil {
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

func (f FullRecoverySP) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)
		pools.SP.Current = pools.SP.Max
	}
	return nil
}

func (f FullRecoverySP) Validate(ctx *Context) error {
	if len(ctx.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	
	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range ctx.Targets {
		if ctx.World.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (f FullRecoverySP) String() string {
	return "FullRecoverySP"
}

// RecoveryHP は非戦闘時のHP回復エフェクト（ログ出力なし）
type RecoveryHP struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

func (r RecoveryHP) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)

		switch amount := r.Amount.(type) {
		case gc.RatioAmount:
			healAmount := amount.Calc(pools.HP.Max)
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+healAmount)
		case gc.NumeralAmount:
			healAmount := amount.Calc()
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+healAmount)
		default:
			return fmt.Errorf("未対応の回復量タイプ: %T", amount)
		}
	}
	return nil
}

func (r RecoveryHP) Validate(ctx *Context) error {
	if r.Amount == nil {
		return errors.New("回復量が指定されていません")
	}
	if len(ctx.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	
	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range ctx.Targets {
		if ctx.World.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (r RecoveryHP) String() string {
	return fmt.Sprintf("RecoveryHP(%v)", r.Amount)
}

// RecoverySP は非戦闘時のSP回復エフェクト（ログ出力なし）
type RecoverySP struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

func (r RecoverySP) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)

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

func (r RecoverySP) Validate(ctx *Context) error {
	if r.Amount == nil {
		return errors.New("回復量が指定されていません")
	}
	if len(ctx.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	
	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range ctx.Targets {
		if ctx.World.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (r RecoverySP) String() string {
	return fmt.Sprintf("RecoverySP(%v)", r.Amount)
}