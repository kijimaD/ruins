package systems

import (
	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/consts"
	"github.com/kijimaD/ruins/lib/gamelog"
	w "github.com/kijimaD/ruins/lib/world"
	"github.com/kijimaD/ruins/lib/worldhelper"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ItemCollectionSystem はプレイヤーとフィールドアイテムの衝突を検出し、アイテムを収集する
func ItemCollectionSystem(world w.World) {
	// 収集されたアイテムを記録するリスト
	var itemsToCollect []ecs.Entity

	// プレイヤー（Operatorコンポーネントを持つエンティティ）とフィールドアイテムの衝突をチェック
	world.Manager.Join(
		world.Components.Position,
		world.Components.Operator,
	).Visit(ecs.Visit(func(playerEntity ecs.Entity) {
		playerPos := world.Components.Position.Get(playerEntity).(*gc.Position)

		// フィールドアイテムとの衝突をチェック
		world.Manager.Join(
			world.Components.Item,
			world.Components.ItemLocationOnField,
			world.Components.GridElement,
		).Visit(ecs.Visit(func(itemEntity ecs.Entity) {
			// グリッド位置からピクセル位置を計算
			gridElement := world.Components.GridElement.Get(itemEntity).(*gc.GridElement)
			itemPixelPos := &gc.Position{
				X: gc.Pixel(int(gridElement.Row)*int(consts.TileSize) + int(consts.TileSize)/2), // タイル中央
				Y: gc.Pixel(int(gridElement.Col)*int(consts.TileSize) + int(consts.TileSize)/2), // タイル中央
			}

			// 衝突判定（共通関数を使用）
			if checkCollisionSimple(world, playerEntity, itemEntity, playerPos, itemPixelPos) {
				itemsToCollect = append(itemsToCollect, itemEntity)
			}
		}))
	}))

	// 収集されたアイテムを処理（バックパックに移動）
	for _, itemEntity := range itemsToCollect {
		collectFieldItem(world, itemEntity)
	}
}

// collectFieldItem はフィールドアイテムを収集してバックパックに移動する
func collectFieldItem(world w.World, itemEntity ecs.Entity) {
	// アイテム名を取得（ログ用）
	itemName := "Unknown Item"
	if nameComp := world.Components.Name.Get(itemEntity); nameComp != nil {
		name := nameComp.(*gc.Name)
		itemName = name.Name
	}

	// フィールドからバックパックに移動
	// ItemLocationOnFieldコンポーネントを削除
	itemEntity.RemoveComponent(world.Components.ItemLocationOnField)

	// ItemLocationInBackpackコンポーネントを追加
	itemEntity.AddComponent(world.Components.ItemLocationInBackpack, gc.LocationInBackpack{})

	// グリッド表示コンポーネントを削除（フィールドから消す）
	if itemEntity.HasComponent(world.Components.GridElement) {
		itemEntity.RemoveComponent(world.Components.GridElement)
	}

	// スプライト表示コンポーネントを削除（フィールドから消す）
	if itemEntity.HasComponent(world.Components.SpriteRender) {
		itemEntity.RemoveComponent(world.Components.SpriteRender)
	}

	// 既存のバックパック内の同じアイテムと統合する処理
	if err := worldhelper.MergeMaterialIntoInventory(world, itemEntity, itemName); err != nil {
		panic(err)
	}

	gamelog.FieldLog.Append(itemName + "を入手した。")
}
