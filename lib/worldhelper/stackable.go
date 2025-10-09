package worldhelper

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// MergeStackableIntoInventory は既存のバックパック内Stackableアイテムと統合するか新規追加する
// Stackableコンポーネントを持つ場合は既存と数量統合、それ以外は個別アイテムとして追加
func MergeStackableIntoInventory(world w.World, newItemEntity ecs.Entity, itemName string) error {
	// Stackableコンポーネントがない場合は何もしない（個別アイテムとして扱う）
	if !newItemEntity.HasComponent(world.Components.Stackable) {
		return nil
	}

	// 既存の同名Stackableアイテムを探してマージ
	existingEntity, found := FindStackableInInventory(world, itemName)
	if found && existingEntity != newItemEntity {
		mergeStackables(world, existingEntity, newItemEntity)
	}

	return nil
}

// mergeStackables はStackableアイテムをマージする。数量を統合してnewItemエンティティは削除する
func mergeStackables(world w.World, existingItem, newItem ecs.Entity) {
	// 新しいアイテムの数量を既存のアイテムに追加
	existingStackable := world.Components.Stackable.Get(existingItem).(*gc.Stackable)
	newStackable := world.Components.Stackable.Get(newItem).(*gc.Stackable)

	// 数量を統合
	existingStackable.Count += newStackable.Count

	// 新しいアイテムエンティティを削除
	world.Manager.DeleteEntity(newItem)
}

// AddStackableCount は指定した名前のStackableアイテムの数量を増やす
// アイテムが存在しない場合は新規作成する
func AddStackableCount(world w.World, name string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive: %d", amount)
	}

	// 既存のアイテムを検索
	entity, found := FindStackableInInventory(world, name)
	if found {
		// 既存アイテムの数量を増やす
		stackable := world.Components.Stackable.Get(entity).(*gc.Stackable)
		stackable.Count += amount
		return nil
	}

	// 存在しない場合は新規作成
	_, err := SpawnStackable(world, name, amount, gc.ItemLocationInBackpack)
	return err
}

// RemoveStackableCount は指定した名前のStackableアイテムの数量を減らす
// 0個以下になった場合はエンティティを削除する
func RemoveStackableCount(world w.World, name string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive: %d", amount)
	}

	entity, found := FindStackableInInventory(world, name)
	if !found {
		return fmt.Errorf("stackable item not found: %s", name)
	}

	stackable := world.Components.Stackable.Get(entity).(*gc.Stackable)
	stackable.Count -= amount

	// 0個以下になったらエンティティを削除
	if stackable.Count <= 0 {
		world.Manager.DeleteEntity(entity)
	}

	return nil
}
