package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// UseItem はアイテムを使用するエフェクト
// アイテムのコンポーネントに基づいて他のエフェクトに分解される
type UseItem struct {
	Item ecs.Entity // 使用するアイテムのエンティティ
}

func (u UseItem) Apply(ctx *Context) error {
	// アイテムのコンポーネントを読み取って対応するエフェクトを直接適用

	// 回復効果があるかチェック
	if healing := ctx.World.Components.ProvidesHealing.Get(u.Item); healing != nil {
		healingComponent := healing.(*gc.ProvidesHealing)
		// アイテムによる回復は非戦闘時の回復として処理（ログ出力なし）
		healingEffect := RecoveryHP{Amount: healingComponent.Amount}

		// Validateで事前検証
		if err := healingEffect.Validate(ctx); err != nil {
			return fmt.Errorf("回復エフェクト検証失敗: %w", err)
		}

		if err := healingEffect.Apply(ctx); err != nil {
			return fmt.Errorf("回復エフェクト適用失敗: %w", err)
		}
	}

	// ダメージ効果があるかチェック
	if damage := ctx.World.Components.InflictsDamage.Get(u.Item); damage != nil {
		damageComponent := damage.(*gc.InflictsDamage)
		damageEffect := CombatDamage{
			Amount: damageComponent.Amount,
			Source: DamageSourceItem,
		}

		// Validateで事前検証
		if err := damageEffect.Validate(ctx); err != nil {
			return fmt.Errorf("ダメージエフェクト検証失敗: %w", err)
		}

		if err := damageEffect.Apply(ctx); err != nil {
			return fmt.Errorf("ダメージエフェクト適用失敗: %w", err)
		}
	}

	// 消費可能アイテムの場合は削除
	if consumable := ctx.World.Components.Consumable.Get(u.Item); consumable != nil {
		ctx.World.Manager.DeleteEntity(u.Item)
	}

	return nil
}

func (u UseItem) Validate(ctx *Context) error {
	// アイテムエンティティが有効かチェック
	if u.Item == 0 {
		return errors.New("無効なアイテムエンティティです")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}

	// 何らかの効果があるかチェック
	hasEffect := false
	if ctx.World.Components.ProvidesHealing.Get(u.Item) != nil {
		hasEffect = true
	}
	if ctx.World.Components.InflictsDamage.Get(u.Item) != nil {
		hasEffect = true
	}

	if !hasEffect {
		return errors.New("このアイテムには効果がありません")
	}

	// アイテム効果のターゲットのPoolsコンポーネント存在確認
	// 回復またはダメージ効果がある場合、ターゲットにPoolsコンポーネントが必要
	if ctx.World.Components.ProvidesHealing.Get(u.Item) != nil || ctx.World.Components.InflictsDamage.Get(u.Item) != nil {
		for _, target := range ctx.Targets {
			if ctx.World.Components.Pools.Get(target) == nil {
				return fmt.Errorf("ターゲット %d にPoolsコンポーネントがありません", target)
			}
		}
	}

	return nil
}

func (u UseItem) String() string {
	return fmt.Sprintf("UseItem(%d)", u.Item)
}

// ConsumeItem は単純にアイテムを消費するエフェクト（効果なし）
type ConsumeItem struct {
	Item ecs.Entity // 消費するアイテム
}

func (c ConsumeItem) Apply(ctx *Context) error {
	if c.Item != 0 {
		ctx.World.Manager.DeleteEntity(c.Item)
	}
	return nil
}

func (c ConsumeItem) Validate(ctx *Context) error {
	if c.Item == 0 {
		return errors.New("無効なアイテムエンティティです")
	}
	if ctx.World.Manager == nil {
		return errors.New("Worldが設定されていません")
	}
	return nil
}

func (c ConsumeItem) String() string {
	return fmt.Sprintf("ConsumeItem(%d)", c.Item)
}

// CreateItem はアイテムを生成するエフェクト（将来拡張用）
type CreateItem struct {
	ItemType string // 生成するアイテムのタイプ
	Quantity int    // 生成数量
}

func (c CreateItem) Apply(ctx *Context) error {
	// TODO: アイテム生成機能を実装
	// 現在はプレースホルダー
	return nil
}

func (c CreateItem) Validate(ctx *Context) error {
	if c.ItemType == "" {
		return errors.New("アイテムタイプが指定されていません")
	}
	if c.Quantity <= 0 {
		return errors.New("生成数量は1以上である必要があります")
	}
	return nil
}

func (c CreateItem) String() string {
	return fmt.Sprintf("CreateItem(%s, %d)", c.ItemType, c.Quantity)
}
