package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Healing は体力を回復するエフェクト（戦闘・非戦闘共用）
type Healing struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

// Apply はHP回復エフェクトをターゲットに適用する
func (h Healing) Apply(world w.World, scope *Scope) error {
	if err := h.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)

		originalHP := pools.HP.Current

		switch amount := h.Amount.(type) {
		case gc.RatioAmount:
			healAmount := amount.Calc(pools.HP.Max)
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+healAmount)
		case gc.NumeralAmount:
			healAmount := amount.Calc()
			pools.HP.Current = mathutil.Min(pools.HP.Max, pools.HP.Current+healAmount)
		default:
			return fmt.Errorf("未対応の回復量タイプ: %T", amount)
		}

		actualHealing := pools.HP.Current - originalHP
		h.logHealing(world, target, actualHealing, scope.Logger)
	}
	return nil
}

// Validate はHP回復エフェクトの妥当性を検証する
func (h Healing) Validate(world w.World, scope *Scope) error {
	if h.Amount == nil {
		return errors.New("回復量が指定されていません")
	}
	if len(scope.Targets) == 0 {
		return errors.New("回復対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認と死亡状態チェック
	for _, target := range scope.Targets {
		if !target.HasComponent(world.Components.Pools) {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
		// 死亡状態のキャラクターには通常の回復エフェクトは使用不可
		if target.HasComponent(world.Components.Dead) {
			return fmt.Errorf("死亡しているキャラクターには回復エフェクトは使用できません")
		}
	}
	return nil
}

func (h Healing) String() string {
	return fmt.Sprintf("Healing(%v)", h.Amount)
}

func (h Healing) logHealing(world w.World, target ecs.Entity, amount int, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sが%d回復。", name.Name, amount)
		logger.Append(entry)
	}
}
