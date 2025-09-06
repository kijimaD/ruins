package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Resurrection は死亡したキャラクターを蘇生するエフェクト
type Resurrection struct {
	Amount gc.Amounter // 蘇生時のHP回復量（固定値または割合）
}

// Apply は蘇生エフェクトをターゲットに適用する
func (r Resurrection) Apply(world w.World, scope *Scope) error {
	if err := r.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)

		// まず死亡状態を解除
		target.RemoveComponent(world.Components.Dead)

		// HPを回復
		switch amount := r.Amount.(type) {
		case gc.RatioAmount:
			healAmount := amount.Calc(pools.HP.Max)
			pools.HP.Current = mathutil.Min(pools.HP.Max, healAmount)
		case gc.NumeralAmount:
			healAmount := amount.Calc()
			pools.HP.Current = mathutil.Min(pools.HP.Max, healAmount)
		default:
			return fmt.Errorf("未対応の蘇生回復量タイプ: %T", amount)
		}

		// 最低でもHP 1は回復させる
		if pools.HP.Current == 0 {
			pools.HP.Current = 1
		}

		r.logResurrection(world, target, pools.HP.Current, scope.Logger)
	}
	return nil
}

// Validate は蘇生エフェクトの妥当性を検証する
func (r Resurrection) Validate(world w.World, scope *Scope) error {
	if r.Amount == nil {
		return errors.New("蘇生回復量が指定されていません")
	}
	if len(scope.Targets) == 0 {
		return errors.New("蘇生対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認と死亡状態チェック
	for _, target := range scope.Targets {
		if !target.HasComponent(world.Components.Pools) {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
		// 生存しているキャラクターには蘇生エフェクトは使用不可
		if !target.HasComponent(world.Components.Dead) {
			return fmt.Errorf("生存しているキャラクターには蘇生エフェクトは使用できません")
		}
	}
	return nil
}

func (r Resurrection) String() string {
	return fmt.Sprintf("Resurrection(%v)", r.Amount)
}

func (r Resurrection) logResurrection(world w.World, target ecs.Entity, finalHP int, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sが蘇生した。HP %d で復活。", name.Name, finalHP)
		logger.Push(entry)
	}
}
