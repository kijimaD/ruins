package worldhelper

import (
	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// InitDebugData はデバッグ用の初期データを設定する
// パーティメンバーが存在しない場合のみ実行される
func InitDebugData(world w.World) {
	// 既にパーティメンバーが存在するかチェック
	memberCount := 0
	world.Manager.Join(
		world.Components.FactionAlly,
		world.Components.InParty,
	).Visit(ecs.Visit(func(_ ecs.Entity) {
		memberCount++
	}))

	// 既にメンバーがいる場合は何もしない
	if memberCount > 0 {
		return
	}

	// デバッグ用アイテム生成
	SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
	card1 := SpawnItem(world, "木刀", gc.ItemLocationInBackpack)
	card2 := SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
	card3 := SpawnItem(world, "M72 LAW", gc.ItemLocationInBackpack)
	SpawnItem(world, "ハンドガン", gc.ItemLocationInBackpack)
	SpawnItem(world, "レイガン", gc.ItemLocationInBackpack)
	armor := SpawnItem(world, "西洋鎧", gc.ItemLocationInBackpack)
	SpawnItem(world, "作業用ヘルメット", gc.ItemLocationInBackpack)
	SpawnItem(world, "革のブーツ", gc.ItemLocationInBackpack)
	SpawnItem(world, "ルビー原石", gc.ItemLocationInBackpack)
	SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	SpawnItem(world, "回復薬", gc.ItemLocationInBackpack)
	SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	SpawnItem(world, "回復スプレー", gc.ItemLocationInBackpack)
	SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)
	SpawnItem(world, "手榴弾", gc.ItemLocationInBackpack)

	// デバッグ用メンバー生成
	ishihara := SpawnMember(world, "イシハラ", true)
	shirase := SpawnMember(world, "シラセ", true)
	SpawnMember(world, "タチバナ", true)
	SpawnMember(world, "ハンス", false)
	SpawnMember(world, "カイン", false)
	SpawnMember(world, "メイ", false)

	// デバッグ用マテリアルとレシピ
	SpawnAllMaterials(world)
	PlusAmount("鉄", 40, world)
	PlusAmount("鉄くず", 4, world)
	PlusAmount("緑ハーブ", 2, world)
	PlusAmount("フェライトコア", 30, world)
	SpawnAllRecipes(world)
	SpawnAllCards(world)

	// デバッグ用装備
	Equip(world, card1, ishihara, gc.EquipmentSlotNumber(0))
	Equip(world, card2, ishihara, gc.EquipmentSlotNumber(0))
	Equip(world, card3, shirase, gc.EquipmentSlotNumber(0))
	Equip(world, armor, ishihara, gc.EquipmentSlotNumber(0))
}
