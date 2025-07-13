package worldhelper

import (
	"testing"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/game"
	"github.com/stretchr/testify/assert"
)

func TestEquipDisarm(t *testing.T) {
	world := game.InitWorld(960, 720)
	gameComponents := world.Components.Game.(*gc.Components)

	// アイテムエンティティを作成
	item := world.Manager.NewEntity()
	item.AddComponent(gameComponents.Item, &gc.Item{})
	item.AddComponent(gameComponents.ItemLocationInBackpack, &gc.ItemLocationInBackpack)

	// オーナーエンティティを作成
	owner := world.Manager.NewEntity()

	// 装備する
	Equip(world, item, owner, gc.EquipmentSlotNumber(0))

	// 装備されたことを確認
	assert.True(t, item.HasComponent(gameComponents.ItemLocationEquipped), "アイテムが装備されていない")
	assert.False(t, item.HasComponent(gameComponents.ItemLocationInBackpack), "アイテムがまだバックパックにある")
	assert.True(t, item.HasComponent(gameComponents.EquipmentChanged), "装備変更フラグが設定されていない")

	equipped := gameComponents.ItemLocationEquipped.Get(item).(*gc.LocationEquipped)
	assert.Equal(t, owner, equipped.Owner, "オーナーが正しく設定されていない")
	assert.Equal(t, gc.EquipmentSlotNumber(0), equipped.EquipmentSlot, "スロット番号が正しく設定されていない")

	// 装備を外す
	Disarm(world, item)

	// 装備が外されたことを確認
	assert.False(t, item.HasComponent(gameComponents.ItemLocationEquipped), "アイテムがまだ装備されている")
	assert.True(t, item.HasComponent(gameComponents.ItemLocationInBackpack), "アイテムがバックパックに戻っていない")
	assert.True(t, item.HasComponent(gameComponents.EquipmentChanged), "装備変更フラグが設定されていない")

	// クリーンアップ
	world.Manager.DeleteEntity(item)
	world.Manager.DeleteEntity(owner)
}
