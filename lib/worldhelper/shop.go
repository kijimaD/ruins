package worldhelper

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// 価格倍率
const (
	BuyPriceMultiplier  = 2.0 // 購入価格は価値の2倍
	SellPriceMultiplier = 0.5 // 売却価格は価値の半分
)

// CalculateBuyPrice は購入価格を計算する（価値の2倍）
func CalculateBuyPrice(baseValue int) int {
	return int(float64(baseValue) * BuyPriceMultiplier)
}

// CalculateSellPrice は売却価格を計算する（価値の半分）
func CalculateSellPrice(baseValue int) int {
	return int(float64(baseValue) * SellPriceMultiplier)
}

// GetItemValue はアイテムの基本価値を取得する
func GetItemValue(world w.World, entity ecs.Entity) int {
	if !entity.HasComponent(world.Components.Value) {
		return 0
	}
	value := world.Components.Value.Get(entity).(*gc.Value)
	return value.Value
}

// BuyItem はプレイヤーがアイテムを購入する
// 通貨が足りない場合や購入に失敗した場合はエラーを返す
func BuyItem(world w.World, playerEntity ecs.Entity, itemName string) error {
	// アイテムの価値を取得
	rawMaster := world.Resources.RawMaster.(*raw.Master)
	itemIdx, ok := rawMaster.ItemIndex[itemName]
	if !ok {
		return fmt.Errorf("アイテムが見つかりません: %s", itemName)
	}
	itemDef := rawMaster.Raws.Items[itemIdx]

	if itemDef.Value == nil {
		return fmt.Errorf("アイテムに価値が設定されていません: %s", itemName)
	}

	baseValue := *itemDef.Value
	price := CalculateBuyPrice(baseValue)

	// 所持金をチェック
	if !HasCurrency(world, playerEntity, price) {
		return fmt.Errorf("遺光が足りません（必要: %d、所持: %d）", price, GetCurrency(world, playerEntity))
	}

	// 通貨を消費
	if !ConsumeCurrency(world, playerEntity, price) {
		return fmt.Errorf("通貨の消費に失敗しました")
	}

	// アイテムがStackable対応かチェック
	isStackable := itemDef.Stackable != nil && *itemDef.Stackable

	if isStackable {
		// 既存のスタックがあるか確認して追加
		err := AddStackableCount(world, itemName, 1)
		if err != nil {
			// スタックが存在しない場合は新規作成
			_, err = SpawnStackable(world, itemName, 1, gc.ItemLocationInBackpack)
			if err != nil {
				// 購入失敗時は通貨を返金
				if refundErr := AddCurrency(world, playerEntity, price); refundErr != nil {
					return fmt.Errorf("アイテムの生成に失敗し、返金も失敗しました: %w (返金エラー: %v)", err, refundErr)
				}
				return fmt.Errorf("アイテムの生成に失敗しました: %w", err)
			}
		}
	} else {
		// 通常アイテムは新規作成
		_, err := SpawnItem(world, itemName, gc.ItemLocationInBackpack)
		if err != nil {
			// 購入失敗時は通貨を返金
			if refundErr := AddCurrency(world, playerEntity, price); refundErr != nil {
				return fmt.Errorf("アイテムの生成に失敗し、返金も失敗しました: %w (返金エラー: %v)", err, refundErr)
			}
			return fmt.Errorf("アイテムの生成に失敗しました: %w", err)
		}
	}

	return nil
}

// SellItem はプレイヤーがアイテムを売却する
func SellItem(world w.World, playerEntity ecs.Entity, itemEntity ecs.Entity) error {
	// アイテムの価値を取得
	baseValue := GetItemValue(world, itemEntity)
	if baseValue == 0 {
		return fmt.Errorf("このアイテムは売却できません")
	}

	price := CalculateSellPrice(baseValue)

	// アイテムがStackableの場合は1個だけ減らす
	if itemEntity.HasComponent(world.Components.Stackable) {
		stackable := world.Components.Stackable.Get(itemEntity).(*gc.Stackable)
		if stackable.Count > 1 {
			stackable.Count--
		} else {
			// 最後の1個の場合はエンティティを削除
			world.Manager.DeleteEntity(itemEntity)
		}
	} else {
		// 通常アイテムはエンティティを削除
		world.Manager.DeleteEntity(itemEntity)
	}

	// 通貨を追加
	if err := AddCurrency(world, playerEntity, price); err != nil {
		return fmt.Errorf("通貨の追加に失敗しました: %w", err)
	}

	return nil
}

// GetShopInventory は店の品揃えを返す（ハードコーディング）
func GetShopInventory() []string {
	return []string{
		"木刀",
		"ハンドガン",
		"西洋鎧",
		"作業用ヘルメット",
		"革のブーツ",
		"回復薬",
		"回復スプレー",
	}
}
