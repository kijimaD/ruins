package effects

import (
	"errors"
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/engine/entities"
	"github.com/kijimaD/ruins/lib/raw"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// UseItem はアイテムを使用するエフェクト
// アイテムのコンポーネントに基づいて他のエフェクトに分解される
type UseItem struct {
	Item ecs.Entity // 使用するアイテムのエンティティ
}

// Apply はアイテム使用エフェクトを適用する
func (u UseItem) Apply(world w.World, scope *Scope) error {
	if err := u.Validate(world, scope); err != nil {
		return err
	}

	// アイテムのコンポーネントを読み取って対応するエフェクトを直接適用

	// 回復効果があるかチェック
	if healing := world.Components.ProvidesHealing.Get(u.Item); healing != nil {
		healingComponent := healing.(*gc.ProvidesHealing)
		// アイテムによる回復は非戦闘時の回復として処理（ログ出力なし）
		healingEffect := RecoveryHP{Amount: healingComponent.Amount}

		// Apply内でValidateが呼ばれるため直接呼び出し
		if err := healingEffect.Apply(world, scope); err != nil {
			return fmt.Errorf("回復エフェクト適用失敗: %w", err)
		}
	}

	// ダメージ効果があるかチェック
	if damage := world.Components.InflictsDamage.Get(u.Item); damage != nil {
		damageComponent := damage.(*gc.InflictsDamage)
		damageEffect := CombatDamage{
			Amount: damageComponent.Amount,
			Source: DamageSourceItem,
		}

		// Apply内でValidateが呼ばれるため直接呼び出し
		if err := damageEffect.Apply(world, scope); err != nil {
			return fmt.Errorf("ダメージエフェクト適用失敗: %w", err)
		}
	}

	// 消費可能アイテムの場合は削除
	if consumable := world.Components.Consumable.Get(u.Item); consumable != nil {
		world.Manager.DeleteEntity(u.Item)
	}

	return nil
}

// Validate はアイテム使用エフェクトの妥当性を検証する
func (u UseItem) Validate(world w.World, scope *Scope) error {
	// アイテムエンティティにItemコンポーネントがあるかチェック
	if !u.Item.HasComponent(world.Components.Item) {
		return errors.New("無効なアイテムエンティティです")
	}

	// 何らかの効果があるかチェック
	hasEffect := world.Components.ProvidesHealing.Get(u.Item) != nil ||
		world.Components.InflictsDamage.Get(u.Item) != nil

	if !hasEffect {
		return errors.New("このアイテムには効果がありません")
	}

	// アイテム効果のターゲットのPoolsコンポーネント存在確認
	// 回復またはダメージ効果がある場合、ターゲットにPoolsコンポーネントが必要
	if world.Components.ProvidesHealing.Get(u.Item) != nil || world.Components.InflictsDamage.Get(u.Item) != nil {
		for _, target := range scope.Targets {
			if world.Components.Pools.Get(target) == nil {
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

// Apply はアイテム消費エフェクトを適用する
func (c ConsumeItem) Apply(world w.World, scope *Scope) error {
	if err := c.Validate(world, scope); err != nil {
		return err
	}

	world.Manager.DeleteEntity(c.Item)
	return nil
}

// Validate はアイテム消費エフェクトの妥当性を検証する
func (c ConsumeItem) Validate(world w.World, _ *Scope) error {
	// アイテムエンティティにItemコンポーネントがあるかチェック
	if !c.Item.HasComponent(world.Components.Item) {
		return errors.New("無効なアイテムエンティティです")
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

// Apply はアイテム生成エフェクトを適用する
func (c CreateItem) Apply(world w.World, scope *Scope) error {
	if err := c.Validate(world, scope); err != nil {
		return err
	}

	// RawMasterを直接使用してアイテムを生成（循環インポート回避）
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	for i := 0; i < c.Quantity; i++ {
		componentList := entities.ComponentList{}
		gameComponent, err := rawMaster.GenerateItem(c.ItemType, gc.ItemLocationInBackpack)
		if err != nil {
			return fmt.Errorf("アイテム生成失敗: %w", err)
		}
		componentList.Game = append(componentList.Game, gameComponent)
		entities.AddEntities(world, componentList)
	}
	return nil
}

// Validate はアイテム生成エフェクトの妥当性を検証する
func (c CreateItem) Validate(_ w.World, _ *Scope) error {
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
