package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/gamelog"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// Context はエフェクト実行時のコンテキスト情報を保持する
type Context struct {
	World   w.World      // ECSワールド
	Creator *ecs.Entity  // 効果の発動者（nilの場合もある）
	Targets []ecs.Entity // 効果の対象エンティティ一覧
}

// Effect はゲーム内の効果を表す核心インターフェース
type Effect interface {
	// Apply は効果を実際に適用する
	Apply(ctx *Context) error

	// Validate は効果の適用前に妥当性を検証する
	Validate(ctx *Context) error

	// String は効果の文字列表現を返す（ログとデバッグ用）
	String() string
}

// DamageSource はダメージの発生源を示す
type DamageSource int

const (
	DamageSourceWeapon DamageSource = iota // 武器によるダメージ
	DamageSourceItem                       // アイテムによるダメージ
)

// CombatDamage はダメージを与えるエフェクト
type CombatDamage struct {
	Amount int          // ダメージ量
	Source DamageSource // ダメージの発生源
}

func (d CombatDamage) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)

		originalHP := pools.HP.Current
		pools.HP.Current = mathutil.Max(0, pools.HP.Current-d.Amount)
		actualDamage := originalHP - pools.HP.Current

		// ダメージログを記録
		d.logDamage(ctx, target, actualDamage)

		// 死亡チェック
		if pools.HP.Current == 0 {
			d.logDeath(ctx, target)
		}
	}
	return nil
}

func (d CombatDamage) Validate(ctx *Context) error {
	if d.Amount < 0 {
		return errors.New("ダメージは0以上である必要があります")
	}
	if len(ctx.Targets) == 0 {
		return errors.New("ダメージ対象が指定されていません")
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

func (d CombatDamage) String() string {
	return fmt.Sprintf("Damage(%d, %s)", d.Amount, d.sourceString())
}

func (d CombatDamage) logDamage(ctx *Context, target ecs.Entity, amount int) {
	if nameComponent := ctx.World.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sに%dのダメージ。", name.Name, amount)
		gamelog.BattleLog.Append(entry)
	}
}

func (d CombatDamage) logDeath(ctx *Context, target ecs.Entity) {
	if nameComponent := ctx.World.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		gamelog.BattleLog.Append(fmt.Sprintf("%sは倒れた。", name.Name))
	}
}

func (d CombatDamage) sourceString() string {
	switch d.Source {
	case DamageSourceWeapon:
		return "武器"
	case DamageSourceItem:
		return "アイテム"
	default:
		return "不明"
	}
}

// CombatHealing は体力を回復するエフェクト
type CombatHealing struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

func (h CombatHealing) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)

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
		h.logHealing(ctx, target, actualHealing)
	}
	return nil
}

func (h CombatHealing) Validate(ctx *Context) error {
	if h.Amount == nil {
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

func (h CombatHealing) String() string {
	return fmt.Sprintf("Healing(%v)", h.Amount)
}

func (h CombatHealing) logHealing(ctx *Context, target ecs.Entity, amount int) {
	if nameComponent := ctx.World.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sが%d回復。", name.Name, amount)
		gamelog.BattleLog.Append(entry)
	}
}

// ConsumeStamina はスタミナを消費するエフェクト
type ConsumeStamina struct {
	Amount gc.Amounter // 消費量（固定値または割合）
}

func (c ConsumeStamina) Apply(ctx *Context) error {
	for _, target := range ctx.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := ctx.World.Components.Pools.Get(target).(*gc.Pools)

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

func (c ConsumeStamina) Validate(ctx *Context) error {
	if c.Amount == nil {
		return errors.New("スタミナ消費量が指定されていません")
	}
	if len(ctx.Targets) == 0 {
		return errors.New("スタミナ消費対象が指定されていません")
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

func (c ConsumeStamina) String() string {
	return fmt.Sprintf("ConsumeStamina(%v)", c.Amount)
}

// RestoreStamina はスタミナを回復するエフェクト
type RestoreStamina struct {
	Amount gc.Amounter // 回復量（固定値または割合）
}

func (r RestoreStamina) Apply(ctx *Context) error {
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
			return fmt.Errorf("未対応のスタミナ回復量タイプ: %T", amount)
		}
	}
	return nil
}

func (r RestoreStamina) Validate(ctx *Context) error {
	if r.Amount == nil {
		return errors.New("スタミナ回復量が指定されていません")
	}
	if len(ctx.Targets) == 0 {
		return errors.New("スタミナ回復対象が指定されていません")
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

func (r RestoreStamina) String() string {
	return fmt.Sprintf("RestoreStamina(%v)", r.Amount)
}
