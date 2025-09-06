package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/mathutil"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// DamageSource はダメージの発生源を示す
type DamageSource int

const (
	// DamageSourceWeapon は武器によるダメージを表す
	DamageSourceWeapon DamageSource = iota // 武器によるダメージ
	// DamageSourceItem はアイテムによるダメージを表す
	DamageSourceItem // アイテムによるダメージ
)

// Damage はダメージを与えるエフェクト
type Damage struct {
	Amount int          // ダメージ量
	Source DamageSource // ダメージの発生源
}

// Apply はダメージエフェクトをターゲットに適用する
func (d Damage) Apply(world w.World, scope *Scope) error {
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
			// 死亡状態コンポーネントを付与
			target.AddComponent(world.Components.Dead, &gc.Dead{})
			d.logDeath(world, target, scope.Logger)
		}
	}
	return nil
}

// Validate はダメージエフェクトの妥当性を検証する
func (d Damage) Validate(world w.World, scope *Scope) error {
	if d.Amount < 0 {
		return errors.New("ダメージは0以上である必要があります")
	}
	if len(scope.Targets) == 0 {
		return errors.New("ダメージ対象が指定されていません")
	}

	// ターゲットのPoolsコンポーネント存在確認
	for _, target := range scope.Targets {
		if !target.HasComponent(world.Components.Pools) {
			return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
		}
	}
	return nil
}

func (d Damage) String() string {
	return fmt.Sprintf("Damage(%d, %s)", d.Amount, d.sourceString())
}

func (d Damage) logDamage(world w.World, target ecs.Entity, amount int, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		entry := fmt.Sprintf("%sに%dのダメージ。", name.Name, amount)
		logger.Push(entry)
	}
}

func (d Damage) logDeath(world w.World, target ecs.Entity, logger GameLogAppender) {
	if logger == nil {
		return // ゲームログ出力先が指定されていない場合は何もしない
	}
	if nameComponent := world.Components.Name.Get(target); nameComponent != nil {
		name := nameComponent.(*gc.Name)
		logger.Push(fmt.Sprintf("%sは倒れた。", name.Name))
	}
}

func (d Damage) sourceString() string {
	switch d.Source {
	case DamageSourceWeapon:
		return "武器"
	case DamageSourceItem:
		return "アイテム"
	default:
		return "不明"
	}
}
