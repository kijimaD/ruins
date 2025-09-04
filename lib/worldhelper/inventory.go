package worldhelper

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AddToInventory は既存のバックパック内アイテムと統合するか新規追加する
// materialの場合は数量管理、それ以外は単純にlocation切り替え
func AddToInventory(world w.World, newItemEntity ecs.Entity, itemName string) {
	// ItemコンポーネントまたはMaterialコンポーネントを持っているかチェック
	hasItem := newItemEntity.HasComponent(world.Components.Item)
	hasMaterial := newItemEntity.HasComponent(world.Components.Material)

	if !hasItem && !hasMaterial {
		panic(fmt.Sprintf("Entity %v does not have Item or Material component", newItemEntity))
	}

	// materialかどうかを確認
	isMaterial := hasMaterial

	if isMaterial {
		// materialの場合は既存の同名materialを探して数量を追加
		var existingMaterial ecs.Entity
		var found bool

		world.Manager.Join(
			world.Components.Material,
			world.Components.ItemLocationInBackpack,
			world.Components.Name,
		).Visit(ecs.Visit(func(entity ecs.Entity) {
			if found {
				return // 既に見つかっている場合はスキップ
			}

			name := world.Components.Name.Get(entity).(*gc.Name)
			if name.Name == itemName && entity != newItemEntity {
				existingMaterial = entity
				found = true
			}
		}))

		if found {
			// 既存のmaterialに数量を追加
			mergeMaterials(world, existingMaterial, newItemEntity)
		} else {
			// 見つからなかった場合は新規materialとして追加（既にLocationInBackpackが設定済み）
		}
	} else {
		// material以外の場合は単純にlocation切り替えのみ（統合しない）
		// 既にLocationInBackpackが設定されているのでそのまま
	}
}

// mergeMaterials はmaterial入手を処理する。materialの場合は数量管理なので、入手数量を足してアイテムエンティティは消す
func mergeMaterials(world w.World, existingMaterial, newMaterial ecs.Entity) {
	// 新しいmaterialの数量を既存のmaterialに追加
	existingMat := world.Components.Material.Get(existingMaterial).(*gc.Material)
	newMat := world.Components.Material.Get(newMaterial).(*gc.Material)

	// 数量を統合
	existingMat.Amount += newMat.Amount

	// 新しいmaterialエンティティを削除
	world.Manager.DeleteEntity(newMaterial)
}

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
	case gc.ItemLocationNone:
		if itemEntity.HasComponent(world.Components.ItemLocationNone) {
			itemEntity.RemoveComponent(world.Components.ItemLocationNone)
		}
	}

	// 新しい位置コンポーネントを追加
	switch toLocation {
	case gc.ItemLocationInBackpack:
		itemEntity.AddComponent(world.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
	case gc.ItemLocationOnField:
		itemEntity.AddComponent(world.Components.ItemLocationOnField, &gc.LocationOnField{})
	case gc.ItemLocationNone:
		itemEntity.AddComponent(world.Components.ItemLocationNone, &gc.LocationNone{})
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

// GetInventoryMaterials はバックパック内のマテリアル一覧を取得する
func GetInventoryMaterials(world w.World) []ecs.Entity {
	var materials []ecs.Entity

	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		materials = append(materials, entity)
	}))

	return materials
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

// FindMaterialInInventory は名前でバックパック内のマテリアルを検索する
func FindMaterialInInventory(world w.World, materialName string) (ecs.Entity, bool) {
	var foundEntity ecs.Entity
	var found bool

	world.Manager.Join(
		world.Components.Material,
		world.Components.ItemLocationInBackpack,
		world.Components.Name,
	).Visit(ecs.Visit(func(entity ecs.Entity) {
		if found {
			return
		}
		name := world.Components.Name.Get(entity).(*gc.Name)
		if name.Name == materialName {
			foundEntity = entity
			found = true
		}
	}))

	return foundEntity, found
}
