package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// RemoveFromInventory はインベントリからアイテムを削除する
func RemoveFromInventory(world w.World, itemEntity ecs.Entity) {
	if !itemEntity.HasComponent(world.Components.ItemLocationInBackpack) {
		return // バックパックにないアイテムは何もしない
	}

	world.Manager.DeleteEntity(itemEntity)
}

// TransferItem はアイテムの位置を変更する
func TransferItem(world w.World, itemEntity ecs.Entity, fromLocation, toLocation gc.ItemLocationType) {
	// 現在の位置コンポーネントを削除
	switch fromLocation {
	case gc.ItemLocationInBackpack:
		if itemEntity.HasComponent(world.Components.ItemLocationInBackpack) {
			itemEntity.RemoveComponent(world.Components.ItemLocationInBackpack)
		}
	case gc.ItemLocationOnField:
		if itemEntity.HasComponent(world.Components.ItemLocationOnField) {
			itemEntity.RemoveComponent(world.Components.ItemLocationOnField)
		}
	}

	// 新しい位置コンポーネントを追加
	switch toLocation {
	case gc.ItemLocationInBackpack:
		itemEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	case gc.ItemLocationOnField:
		itemEntity.AddComponent(world.Components.ItemLocationOnField, &gc.LocationOnField{})
	}
}

// GetInventoryItems はバックパック内のアイテム一覧を取得する
func GetInventoryItems(world w.World) []ecs.Entity {
	var items []ecs.Entity

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		items = append(items, entity)
	}))

	return items
}

// GetInventoryStackables はバックパック内のスタック可能アイテム一覧を取得する
func GetInventoryStackables(world w.World) []ecs.Entity {
	var stackables []ecs.Entity

	world.Manager.Join(
		world.Components.Stackable,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		stackables = append(stackables, entity)
	}))

	return stackables
}

// FindStackableInInventory は名前でバックパック内のStackableアイテムを検索する
func FindStackableInInventory(world w.World, name string) (ecs.Entity, bool) {
	var foundEntity ecs.Entity
	var found bool

	world.Manager.Join(
		world.Components.Stackable,
		world.Components.ItemLocationInBackpack,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if found {
			return
		}
		itemName := world.Components.Name.Get(entity).(*gc.Name)
		if itemName.Name == name {
			foundEntity = entity
			found = true
		}
	}))

	return foundEntity, found
}

// FindItemInInventory は名前でバックパック内のアイテムを検索する
func FindItemInInventory(world w.World, itemName string) (ecs.Entity, bool) {
	var foundEntity ecs.Entity
	var found bool

	world.Manager.Join(
		world.Components.Item,
		world.Components.ItemLocationInBackpack,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if found {
			return
		}
		name := world.Components.Name.Get(entity).(*gc.Name)
		if name.Name == itemName {
			foundEntity = entity
			found = true
		}
	}))

	return foundEntity, found
}
