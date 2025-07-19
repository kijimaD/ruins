package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// GameLogAppender はゲームログ出力のインターフェース
type GameLogAppender interface {
	Append(entry string)
}

// Scope はエフェクトの影響範囲を保持する
type Scope struct {
	Creator *ecs.Entity      // 効果の発動者（nilの場合もある）
	Targets []ecs.Entity     // 効果の対象エンティティ一覧
	Logger  GameLogAppender  // ゲームログ出力先（nilの場合はログ出力なし）
}

// Effect はゲーム内の効果を表す核心インターフェース
type Effect interface {
	// Apply は効果を実際に適用する
	Apply(world w.World, scope *Scope) error

	// Validate は効果の適用前に妥当性を検証する
	Validate(world w.World, scope *Scope) error

	// String は効果の文字列表現を返す（ログとデバッグ用）
	String() string
}

// DamageSource はダメージの発生源を示す
type DamageSource int

const (
	// DamageSourceWeapon は武器によるダメージを表す
	DamageSourceWeapon DamageSource = iota // 武器によるダメージ
	// DamageSourceItem はアイテムによるダメージを表す
	DamageSourceItem                       // アイテムによるダメージ
)

// CombatDamage はダメージを与えるエフェクト
type CombatDamage struct {
	Amount int          // ダメージ量
	Source DamageSource // ダメージの発生源
}

// Apply はダメージエフェクトをターゲットに適用する
func (d CombatDamage) Apply(world w.World, scope *Scope) error {
	if err := d.Validate(world, scope); err != nil {
		return err
	}

	for _, target := range scope.Targets {
		// Validateで事前確認済みのためnilチェック不要
		pools := world.Components.Pools.Get(target).(*gc.Pools)

		originalHP := pools.HP.Current
		pools.HP.Current = mathutil.Max(0, pools.HP.Current-d.Amount)
		actualDamage := originalHP - pools.HP.Current

		// ダメージログを記録
		d.logDamage(world, target, actualDamage, scope.Logger)

		// 死亡チェック
		if pools.HP.Current == 0 {
			d.logDeath(world, target, scope.Logger)
		}
	}
	return nil
}

// Validate はダメージエフェクトの妥当性を検証する
func (d CombatDamage) Validate(world w.World, scope *Scope) error {
	if d.Amount < 0 {
		return errors.New("ダメージは0以上である必要があります")
	}
	if len(scope.Targets) == 0 {
		return errors.New("ダメージ対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (d CombatDamage) String() string {
	return fmt.Sprintf("Damage(%d, %s)", d.Amount, d.sourceString())
}

func (d CombatDamage) logDamage(world w.World, target ecs.Entity, amount int, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sに%dのダメージ。", name.Name, amount)
		logger.Append(entry)
	}
}

func (d CombatDamage) logDeath(world w.World, target ecs.Entity, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		logger.Append(fmt.Sprintf("%sは倒れた。", name.Name))
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

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
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
		if world.Components.Pools.Get(target) == nil {
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
		if world.Components.Pools.Get(target) == nil {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (r RestoreStamina) String() string {
	return fmt.Sprintf("RestoreStamina(%v)", r.Amount)
}
